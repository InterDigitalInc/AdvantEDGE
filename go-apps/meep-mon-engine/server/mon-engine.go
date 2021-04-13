/*
 * Copyright (c) 2019  InterDigital Communications, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
	sbs "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-store"
	v1 "k8s.io/api/core/v1"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type UserData struct {
	AllPodsStatus PodsStatus
	ExpectedPods  map[string]*PodStatus
}

type MonEngineInfo struct {
	PodType              string
	PodName              string
	Namespace            string
	MeepApp              string
	MeepOrigin           string
	MeepScenario         string
	Release              string
	Phase                string
	PodInitialized       string
	PodReady             string
	PodScheduled         string
	PodUnschedulable     string
	PodConditionError    string
	ContainerStatusesMsg string
	NbOkContainers       int
	NbTotalContainers    int
	NbPodRestart         int
	LogicalState         string
	StartTime            string
}

type Sandbox struct {
	Releases  map[string]bool
	StartTime time.Time
	Running   bool
}

const serviceName = "Monitoring Engine"
const moduleName = "meep-mon-engine"
const moduleNamespace = "default"
const notFoundStr = "na"
const monEngineKey = "mon-engine:"

// MQ payload fields
const fieldSandboxName = "sandbox-name"

// index in array
const EVENT_POD_ADDED = 0
const EVENT_POD_MODIFIED = 1
const EVENT_POD_DELETED = 2

// Metrics
var (
	metricSboxCreateDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "mon_engine_sbox_create_duration",
		Help:    "A histogram of sandbox creation durations",
		Buckets: prometheus.LinearBuckets(10, 5, 5),
	})
)

var pod_event_str = [3]string{"pod added", "pod modified", "pod deleted"}
var rc *redis.Connector
var redisAddr = "meep-redis-master:6379"
var baseKey string = dkm.GetKeyRootGlobal() + monEngineKey
var stopChan = make(chan struct{})
var mqGlobal *mq.MsgQueue
var handlerId int
var sandboxStore *sbs.SandboxStore
var mutex sync.Mutex

var sandboxes map[string]*Sandbox

var depPodsList []string
var corePodsList []string
var sboxPodsList []string

var expectedDepPods map[string]*PodStatus
var expectedCorePods map[string]*PodStatus
var expectedSboxPods map[string]*PodStatus

// Init - Mon Engine initialization
func Init() (err error) {

	// Retrieve dependency pod list from environment variable
	expectedDepPods = make(map[string]*PodStatus)
	depPodsStr := strings.TrimSpace(os.Getenv("MEEP_DEPENDENCY_PODS"))
	log.Info("MEEP_DEPENDENCY_PODS: ", depPodsStr)
	if depPodsStr != "" {
		depPodsList = strings.Split(depPodsStr, ",")
		for _, pod := range depPodsList {
			podStatus := new(PodStatus)
			podStatus.PodType = "core"
			podStatus.Sandbox = "default"
			podStatus.Name = pod
			podStatus.LogicalState = "NotAvailable"
			expectedDepPods[pod] = podStatus
		}
	}

	// Retrieve core pod list from environment variable
	expectedCorePods = make(map[string]*PodStatus)
	corePodsStr := strings.TrimSpace(os.Getenv("MEEP_CORE_PODS"))
	log.Info("MEEP_CORE_PODS: ", corePodsStr)
	if corePodsStr != "" {
		corePodsList = strings.Split(corePodsStr, ",")
		for _, pod := range corePodsList {
			podStatus := new(PodStatus)
			podStatus.PodType = "core"
			podStatus.Sandbox = "default"
			podStatus.Name = pod
			podStatus.LogicalState = "NotAvailable"
			expectedCorePods[pod] = podStatus
		}
	}

	// Retrieve sandbox pod list from environment variable
	expectedSboxPods = make(map[string]*PodStatus)
	sboxPodsStr := strings.TrimSpace(os.Getenv("MEEP_SANDBOX_PODS"))
	log.Info("MEEP_SANDBOX_PODS: ", sboxPodsStr)
	if sboxPodsStr != "" {
		sboxPodsList = strings.Split(sboxPodsStr, ",")
	}

	// Create message queue
	mqGlobal, err = mq.NewMsgQueue(mq.GetGlobalName(), moduleName, moduleNamespace, redisAddr)
	if err != nil {
		log.Error("Failed to create Message Queue with error: ", err)
		return err
	}
	log.Info("Message Queue created")

	// Connect to Redis DB
	rc, err = redis.NewConnector(redisAddr, 0)
	if err != nil {
		log.Error("Failed connection to Redis: ", err)
		return err
	}
	log.Info("Connected to Mon Engine DB")

	// Empty DB
	_ = rc.DBFlush(baseKey)

	// Connect to Sandbox Store
	sandboxStore, err = sbs.NewSandboxStore(redisAddr)
	if err != nil {
		log.Error("Failed connection to Sandbox Store: ", err.Error())
		return err
	}
	log.Info("Connected to Sandbox Store")

	// Initialize sandbox map
	sandboxes = make(map[string]*Sandbox)

	return nil
}

// Run - Mon Engine monitoring thread
func Run() (err error) {

	// Initialize expected pods for existing sandboxes
	if sboxMap, err := sandboxStore.GetAll(); err == nil {
		for _, sbox := range sboxMap {
			addExpectedPods(sbox.Name)
			createSandbox(sbox.Name)
		}
	}

	// Register Message Queue handler
	handler := mq.MsgHandler{Handler: msgHandler, UserData: nil}
	handlerId, err = mqGlobal.RegisterHandler(handler)
	if err != nil {
		log.Error("Failed to register MsgQueue handler: ", err.Error())
		return err
	}

	// Start thread to watch k8s pods
	err = k8sConnect()
	if err != nil {
		log.Error("Failed to watch k8s pods")
		return err
	}

	return nil
}

func Stop() {
	close(stopChan)
}

// Message Queue handler
func msgHandler(msg *mq.Msg, userData interface{}) {
	switch msg.Message {
	case mq.MsgSandboxCreate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		addExpectedPods(msg.Payload[fieldSandboxName])
		createSandbox(msg.Payload[fieldSandboxName])
	case mq.MsgSandboxDestroy:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		removeExpectedPods(msg.Payload[fieldSandboxName])
		deleteSandbox(msg.Payload[fieldSandboxName])
	default:
		log.Trace("Ignoring unsupported message: ", mq.PrintMsg(msg))
	}
}

func connectToAPISvr() (*kubernetes.Clientset, error) {

	// Create the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return clientset, nil
}

func printfMonEngineInfo(monEngineInfo MonEngineInfo, reason int) {

	log.Debug("Monitoring Engine info *** ", pod_event_str[reason], " *** ",
		" podType : ", monEngineInfo.PodType,
		" podName : ", monEngineInfo.PodName,
		" namespace : ", monEngineInfo.Namespace,
		" meepApp : ", monEngineInfo.MeepApp,
		" meepOrigin : ", monEngineInfo.MeepOrigin,
		" meepScenario : ", monEngineInfo.MeepScenario,
		" release : ", monEngineInfo.Release,
		" phase : ", monEngineInfo.Phase,
		" podInitialized : ", monEngineInfo.PodInitialized,
		" podUnschedulable : ", monEngineInfo.PodUnschedulable,
		" podScheduled : ", monEngineInfo.PodScheduled,
		" podReady : ", monEngineInfo.PodReady,
		" podConditionError : ", monEngineInfo.PodConditionError,
		" ContainerStatusesMsg : ", monEngineInfo.ContainerStatusesMsg,
		" NbOkContainers : ", monEngineInfo.NbOkContainers,
		" NbTotalContainers : ", monEngineInfo.NbTotalContainers,
		" NbPodRestart : ", monEngineInfo.NbPodRestart,
		" LogicalState : ", monEngineInfo.LogicalState,
		" StartTime : ", monEngineInfo.StartTime)
}

func processEvent(obj interface{}, reason int) {
	var ok bool
	var pod *v1.Pod

	// Validate object type is Pod
	if pod, ok = obj.(*v1.Pod); !ok {
		return
	}

	var monEngineInfo MonEngineInfo

	if reason != EVENT_POD_DELETED {
		podConditionMsg := ""
		podScheduled := "False"
		podReady := "False"
		podInitialized := "False"
		podUnschedulable := "False"
		nbConditions := len(pod.Status.Conditions)
		for i := 0; i < nbConditions; i++ {
			switch pod.Status.Conditions[i].Type {
			case "PodScheduled":
				podScheduled = string(pod.Status.Conditions[i].Status)
			case "Ready":
				podReady = string(pod.Status.Conditions[i].Status)
				if podReady == "False" {
					podConditionMsg = string(pod.Status.Conditions[i].Message)
				}
			case "Initialized":
				podInitialized = string(pod.Status.Conditions[i].Status)
			case "Unschedulable":
				podUnschedulable = string(pod.Status.Conditions[i].Status)
			}
		}

		nbContainers := len(pod.Status.ContainerStatuses)
		okContainers := 0
		restartCount := 0
		reasonFailureStr := ""
		for i := 0; i < nbContainers; i++ {
			if pod.Status.ContainerStatuses[i].Ready {
				okContainers++
			} else if pod.Status.ContainerStatuses[i].State.Waiting != nil {
				reasonFailureStr = pod.Status.ContainerStatuses[i].State.Waiting.Reason
			} else if pod.Status.ContainerStatuses[i].State.Terminated != nil && reasonFailureStr != "" {
				reasonFailureStr = pod.Status.ContainerStatuses[i].State.Terminated.Reason
			}
			//only update if the value is greater than 0, and we keep it
			if restartCount == 0 {
				restartCount = int(pod.Status.ContainerStatuses[i].RestartCount)
			}
		}

		monEngineInfo.PodInitialized = podInitialized
		monEngineInfo.PodUnschedulable = podUnschedulable
		monEngineInfo.PodScheduled = podScheduled
		monEngineInfo.PodReady = podReady
		monEngineInfo.PodConditionError = podConditionMsg
		monEngineInfo.ContainerStatusesMsg = reasonFailureStr
		monEngineInfo.NbOkContainers = okContainers
		monEngineInfo.NbTotalContainers = nbContainers
		monEngineInfo.NbPodRestart = restartCount
	}

	//common for both the add, update and delete
	monEngineInfo.Phase = string(pod.Status.Phase)
	monEngineInfo.PodName = pod.Name
	monEngineInfo.Namespace = pod.Namespace
	monEngineInfo.MeepApp = pod.Labels["meepApp"]
	monEngineInfo.MeepScenario = pod.Labels["meepScenario"]
	if monEngineInfo.Release, ok = pod.Labels["release"]; !ok {
		monEngineInfo.Release = notFoundStr
	}
	if monEngineInfo.MeepApp, ok = pod.Labels["meepApp"]; !ok {
		monEngineInfo.MeepApp = notFoundStr
	}
	if monEngineInfo.MeepOrigin, ok = pod.Labels["meepOrigin"]; !ok {
		monEngineInfo.MeepOrigin = notFoundStr
	}
	if monEngineInfo.MeepScenario, ok = pod.Labels["meepScenario"]; !ok {
		monEngineInfo.MeepScenario = notFoundStr
	}
	monEngineInfo.LogicalState = monEngineInfo.Phase
	monEngineInfo.PodType = getPodType(monEngineInfo.MeepOrigin, monEngineInfo.Release)

	//Phase is Running but might not really be because of some other attributes
	//start of override section of the LogicalState by specific conditions

	if pod.GetObjectMeta().GetDeletionTimestamp() != nil {
		monEngineInfo.LogicalState = "Terminating"
	} else if monEngineInfo.PodReady != "True" {
		monEngineInfo.LogicalState = "Pending"
	} else if monEngineInfo.NbOkContainers < monEngineInfo.NbTotalContainers {
		monEngineInfo.LogicalState = "Failed"
	}
	//end of override section

	printfMonEngineInfo(monEngineInfo, reason)

	// Add, update or remove entry in DB only if core or scenario pod
	if monEngineInfo.PodType != notFoundStr {
		if reason == EVENT_POD_DELETED {
			deleteEntryInDB(&monEngineInfo)
		} else {
			addOrUpdateEntryInDB(&monEngineInfo)
			monitorSboxCreation(&monEngineInfo)
		}
	} else {
		log.Debug("Ignoring non-AdvantEDGE pod: ", monEngineInfo.PodName)
	}
}

func addOrUpdateEntryInDB(monEngineInfo *MonEngineInfo) {
	// Populate rule fields
	fields := make(map[string]interface{})
	fields["type"] = monEngineInfo.PodType
	fields["name"] = monEngineInfo.PodName
	fields["namespace"] = monEngineInfo.Namespace
	fields["meepApp"] = monEngineInfo.MeepApp
	fields["meepOrigin"] = monEngineInfo.MeepOrigin
	fields["meepScenario"] = monEngineInfo.MeepScenario
	fields["release"] = monEngineInfo.Release
	fields["phase"] = monEngineInfo.Phase
	fields["initialised"] = monEngineInfo.PodInitialized
	fields["scheduled"] = monEngineInfo.PodScheduled
	fields["ready"] = monEngineInfo.PodReady
	fields["unschedulable"] = monEngineInfo.PodUnschedulable
	fields["condition-error"] = monEngineInfo.PodConditionError
	fields["nbOkContainers"] = monEngineInfo.NbOkContainers
	fields["nbTotalContainers"] = monEngineInfo.NbTotalContainers
	fields["nbPodRestart"] = monEngineInfo.NbPodRestart
	fields["logicalState"] = monEngineInfo.LogicalState
	fields["startTime"] = monEngineInfo.StartTime

	// Make unique key
	key := baseKey + monEngineInfo.Namespace + ":" + monEngineInfo.PodType + ":" + monEngineInfo.PodName

	// Set rule information in DB
	err := rc.SetEntry(key, fields)
	if err != nil {
		log.Error("Entry could not be updated in DB for ", monEngineInfo.MeepApp, ": ", err)
	}
}

func deleteEntryInDB(monEngineInfo *MonEngineInfo) {

	// Make unique key
	key := baseKey + monEngineInfo.Namespace + ":" + monEngineInfo.PodType + ":" + monEngineInfo.PodName

	// Set rule information in DB
	err := rc.DelEntry(key)
	if err != nil {
		log.Error("Entry could not be deleted in DB for ", monEngineInfo.MeepApp, ": ", err)
	}
}

func k8sConnect() (err error) {

	// Connect to K8s API Server
	clientset, err := connectToAPISvr()
	if err != nil {
		log.Error("Failed to connect with k8s API Server. Error: ", err)
		return err
	}

	watchlist := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", v1.NamespaceAll, fields.Everything())

	// also take a look at NewSharedIndexInformer
	_, controller := cache.NewInformer(
		watchlist,
		&v1.Pod{},
		0, //Duration is int64
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				processEvent(obj, EVENT_POD_ADDED)
			},
			DeleteFunc: func(obj interface{}) {
				processEvent(obj, EVENT_POD_DELETED)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				processEvent(newObj, EVENT_POD_MODIFIED)
			},
		},
	)

	go controller.Run(stopChan)
	return nil
}

func monitorSboxCreation(monEngineInfo *MonEngineInfo) {
	mutex.Lock()
	defer mutex.Unlock()

	// Find matching sandbox entry
	sboxName := monEngineInfo.Namespace
	if sbox, found := sandboxes[sboxName]; found {

		// Ignore if sbox already running
		if !sbox.Running {

			// Set release running state
			if _, exists := sbox.Releases[monEngineInfo.Release]; exists {

				sbox.Releases[monEngineInfo.Release] = (monEngineInfo.LogicalState == "Running")

				// Check if sandbox is running
				sboxRunning := true
				for _, running := range sbox.Releases {
					if !running {
						sboxRunning = false
						break
					}
				}

				// If all releases are running, log sandbox creation time metric
				if sboxRunning {
					sbox.Running = true
					creationTime := float64(time.Since(sbox.StartTime).Milliseconds()) / 1000.0
					log.Info("Sbox: ", sboxName, " creationTime: ", creationTime)
					metricSboxCreateDuration.Observe(creationTime)
				}
			}
		}
	}
}

// Retrieve POD states
// GET /states
func meGetStates(w http.ResponseWriter, r *http.Request) {
	var err error
	var data UserData

	// Retrieve query parameters
	query := r.URL.Query()
	queryType := query.Get("type")
	querySandbox := query.Get("sandbox")
	queryLong := query.Get("long")

	// Get expected pods list
	data.ExpectedPods = make(map[string]*PodStatus)
	if queryType != "scenario" {
		if querySandbox == "" || querySandbox == "all" {
			for k, v := range expectedCorePods {
				data.ExpectedPods[k] = v
			}
			for k, v := range expectedDepPods {
				data.ExpectedPods[k] = v
			}
		}
		if querySandbox != "" || querySandbox == "all" {
			for _, v := range expectedSboxPods {
				if v.Sandbox == querySandbox || querySandbox == "all" {
					data.ExpectedPods[v.Name] = v
				}
			}
		}
	}

	// Create DB key using query filters
	sandboxKey := ""
	if querySandbox == "" {
		sandboxKey = "default:"
	} else if querySandbox == "all" {
		sandboxKey = "*:"
	} else {
		sandboxKey = querySandbox + ":"
	}

	typeKey := ""
	if queryType != "" {
		typeKey = queryType + ":"
	} else {
		typeKey = "*"
	}

	keyName := baseKey + sandboxKey + typeKey + "*"

	// Retrieve pod status information from DB
	if queryLong == "true" {
		err = rc.ForEachEntry(keyName, getPodDetails, &data)
	} else {
		err = rc.ForEachEntry(keyName, getPodStatesOnly, &data)
	}
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add missing pods status
	for _, podStatus := range data.ExpectedPods {
		data.AllPodsStatus.PodStatus = append(data.AllPodsStatus.PodStatus, *podStatus)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Format response
	jsonResponse, err := json.Marshal(data.AllPodsStatus)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

func getPodDetails(key string, fields map[string]string, userData interface{}) error {
	data := userData.(*UserData)

	// Append pod status
	var podStatus PodStatus
	podStatus.PodType = fields["type"]
	podStatus.Sandbox = fields["namespace"]
	podStatus.Name = getPodName(fields["meepApp"], fields["name"])
	podStatus.Namespace = fields["namespace"]
	podStatus.MeepApp = fields["meepApp"]
	podStatus.MeepOrigin = fields["meepOrigin"]
	podStatus.MeepScenario = fields["meepScenario"]
	podStatus.Phase = fields["phase"]
	podStatus.PodInitialized = fields["initialised"]
	podStatus.PodScheduled = fields["scheduled"]
	podStatus.PodReady = fields["ready"]
	podStatus.PodUnschedulable = fields["unschedulable"]
	podStatus.PodConditionError = fields["condition-error"]
	podStatus.NbOkContainers = fields["nbOkContainers"]
	podStatus.NbTotalContainers = fields["nbTotalContainers"]
	podStatus.NbPodRestart = fields["nbPodRestart"]
	podStatus.LogicalState = fields["logicalState"]
	podStatus.StartTime = fields["startTime"]
	data.AllPodsStatus.PodStatus = append(data.AllPodsStatus.PodStatus, podStatus)

	// Remove from expected pods
	delete(data.ExpectedPods, fields["release"])

	return nil
}

func getPodStatesOnly(key string, fields map[string]string, userData interface{}) error {
	data := userData.(*UserData)

	// Append pod status
	var podStatus PodStatus
	podStatus.PodType = fields["type"]
	podStatus.Sandbox = fields["namespace"]
	podStatus.Name = getPodName(fields["meepApp"], fields["name"])
	podStatus.LogicalState = fields["logicalState"]
	data.AllPodsStatus.PodStatus = append(data.AllPodsStatus.PodStatus, podStatus)

	// Remove from expected pods
	delete(data.ExpectedPods, fields["release"])

	return nil
}

func getPodType(origin string, release string) string {
	podType := notFoundStr
	if origin == "core" || origin == "scenario" {
		podType = origin
	} else if release != notFoundStr {
		if _, ok := expectedDepPods[release]; ok {
			podType = "core"
		} else if _, ok := expectedCorePods[release]; ok {
			podType = "core"
		}
	}
	return podType
}

func getPodName(app string, name string) string {
	var podName string
	if app != notFoundStr {
		podName = app
	} else {
		podName = name
	}
	return podName
}

func addExpectedPods(sandboxName string) {
	for _, pod := range sboxPodsList {
		// Get sandbox-specific pod name
		var podName, podKeyName string
		prefix := "meep-"
		sandboxPrefix := prefix + sandboxName + "-"
		if strings.HasPrefix(pod, prefix) {
			podName = pod
			podKeyName = sandboxPrefix + pod[len(prefix):]
		} else {
			podName = prefix + pod
			podKeyName = sandboxPrefix + pod
		}

		// Add to expected sandbox pods list
		podStatus := new(PodStatus)
		podStatus.PodType = "core"
		podStatus.Sandbox = sandboxName
		podStatus.Name = podName
		podStatus.LogicalState = "NotAvailable"
		expectedSboxPods[podKeyName] = podStatus
	}
}

func removeExpectedPods(sandboxName string) {
	for _, pod := range sboxPodsList {
		// Get sandbox-specific pod name
		var podName string
		prefix := "meep-"
		sandboxPrefix := prefix + sandboxName + "-"
		if strings.HasPrefix(pod, prefix) {
			podName = sandboxPrefix + pod[len(prefix):]
		} else {
			podName = sandboxPrefix + pod
		}

		// Delete from expected list
		delete(expectedSboxPods, podName)
	}
}

// Create new sandbox to monitor
func createSandbox(sandboxName string) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := sandboxes[sandboxName]; !exists {
		sbox := new(Sandbox)
		sbox.Running = false
		sbox.Releases = make(map[string]bool)
		sbox.StartTime = time.Now()
		for _, pod := range sboxPodsList {
			sbox.Releases[pod] = false
		}
		sandboxes[sandboxName] = sbox
		log.Info("Created new sandbox to monitor: ", sandboxName)
	}
}

// Delete monitored sandbox
func deleteSandbox(sandboxName string) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := sandboxes[sandboxName]; exists {
		delete(sandboxes, sandboxName)
		log.Info("Removed sandbox to monitor: ", sandboxName)
	}
}
