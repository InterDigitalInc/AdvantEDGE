/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
*/
package main

import (
	"fmt"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-mon-engine/log"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

const moduleMonEngine string = "mon-engine"

//index in array
const EVENT_POD_ADDED = 0
const EVENT_POD_MODIFIED = 1
const EVENT_POD_DELETED = 2

var pod_event_str = [3]string{"pod added", "pod modified", "pod deleted"}

type MonEngineInfo struct {
	PodName              string
	Namespace            string
	MeepApp              string
	MeepOrigin           string
	MeepScenario         string
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

// Init - Mon Engine initialization
func Init() (err error) {

	// Connect to Redis DB
	err = DBConnect()
	if err != nil {
		log.Error("Failed connection to Active DB. Error: ", err)
		return err
	}
	log.Info("Connected to Active DB")

	// Empty DB
	DBFlush(moduleMonEngine)

	return nil
}

// Run - Mon Engine main loop
func Run() (err error) {

	// Watch k8s pods (main loop)
	err = k8sConnect()
	if err != nil {
		log.Error("Failed to watch k8s pods")
		return err
	}

	return nil
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
		"pod name : ", monEngineInfo.PodName,
		"namespace : ", monEngineInfo.Namespace,
		"meepApp : ", monEngineInfo.MeepApp,
		"meepOrigin : ", monEngineInfo.MeepOrigin,
		"meepScenario : ", monEngineInfo.MeepScenario,
		"phase : ", monEngineInfo.Phase,
		"podInitialized : ", monEngineInfo.PodInitialized,
		"podUnschedulable : ", monEngineInfo.PodUnschedulable,
		"podScheduled : ", monEngineInfo.PodScheduled,
		"podReady : ", monEngineInfo.PodReady,
		"podConditionError : ", monEngineInfo.PodConditionError,
		"ContainerStatusesMsg : ", monEngineInfo.ContainerStatusesMsg,
		"NbOkContainers : ", monEngineInfo.NbOkContainers,
		"NbTotalContainers : ", monEngineInfo.NbTotalContainers,
		"NbPodRestart : ", monEngineInfo.NbPodRestart,
		"LogicalState : ", monEngineInfo.LogicalState,
		"StartTime : ", monEngineInfo.StartTime)

}

func processEvent(obj interface{}, reason int) {
	if pod, ok := obj.(*v1.Pod); ok {

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
				if pod.Status.ContainerStatuses[i].Ready == true {
					okContainers++
				} else {
					if pod.Status.ContainerStatuses[i].State.Waiting != nil {
						reasonFailureStr = pod.Status.ContainerStatuses[i].State.Waiting.Reason
					} else if pod.Status.ContainerStatuses[i].State.Terminated != nil {
						if reasonFailureStr != "" {
							reasonFailureStr = pod.Status.ContainerStatuses[i].State.Terminated.Reason
						}
					}
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
		monEngineInfo.MeepOrigin = pod.Labels["meepOrigin"]
		monEngineInfo.MeepScenario = pod.Labels["meepScenario"]
		if pod.Status.StartTime != nil {
			monEngineInfo.StartTime = pod.Status.StartTime.String()
		}

		monEngineInfo.LogicalState = monEngineInfo.Phase

		//Phase is Running but might not really be because of some other attributes
		//start of override section of the LogicalState by specific conditions

		if pod.GetObjectMeta().GetDeletionTimestamp() != nil {
			monEngineInfo.LogicalState = "Terminating"
		} else {
			if monEngineInfo.PodReady != "True" {
				monEngineInfo.LogicalState = "Pending"
			} else {
				if monEngineInfo.NbOkContainers < monEngineInfo.NbTotalContainers {
					monEngineInfo.LogicalState = "Failed"
				}
			}
		}
		//end of override section

		printfMonEngineInfo(monEngineInfo, reason)

		if reason == EVENT_POD_DELETED {
			deleteEntryInDB(monEngineInfo)
		} else {
			addOrUpdateEntryInDB(monEngineInfo)
		}
	}
}

func addOrUpdateEntryInDB(monEngineInfo MonEngineInfo) {
	// Populate rule fields
	fields := make(map[string]interface{})
	fields["name"] = monEngineInfo.PodName
	fields["namespace"] = monEngineInfo.Namespace
	fields["meepApp"] = monEngineInfo.MeepApp
	fields["meepOrigin"] = monEngineInfo.MeepOrigin
	fields["meepScenario"] = monEngineInfo.MeepScenario
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
	key := moduleMonEngine + ":MO-" + monEngineInfo.MeepOrigin + ":MS-" + monEngineInfo.MeepScenario + ":MA-" + monEngineInfo.MeepApp + ":" + monEngineInfo.PodName

	// Set rule information in DB
	DBSetEntry(key, fields)

}

func deleteEntryInDB(monEngineInfo MonEngineInfo) {

	// Make unique key
	key := moduleMonEngine + ":MO-" + monEngineInfo.MeepOrigin + ":MS-" + monEngineInfo.MeepScenario + ":MA-" + monEngineInfo.MeepApp + ":" + monEngineInfo.PodName

	// Set rule information in DB
	DBRemoveEntry(key)

}

func k8sConnect() (err error) {

	// Connect to K8s API Server
	clientset, err := connectToAPISvr()
	if err != nil {
		log.Error("Failed to connect with k8s API Server. Error: ", err)
		return err
	}

	//scenarioName := "latency-demo"
	meepOrigin := "core"

	// Retrieve pods from k8s api with scenario label
	pods, err := clientset.CoreV1().Pods("").List(
		metav1.ListOptions{LabelSelector: fmt.Sprintf("meepOrigin=%s", meepOrigin)})
	if err != nil {
		log.Error("Failed to retrieve services from k8s API Server. Error: ", err)
		return err
	}

	// Store service IPs
	for _, pod := range pods.Items {
		podName := pod.ObjectMeta.Name
		podPhase := pod.Status.Phase
		log.Debug("podName: ", podName, " podPhase: ", podPhase)
	}

	watchlist := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", v1.NamespaceAll, fields.Everything())

	_, controller := cache.NewInformer( // also take a look at NewSharedIndexInformer
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

	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(stop)
	for {
		time.Sleep(time.Second)
	}
}
