/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 *
 * NOTICE: File content based on https://github.com/morvencao/kube-mutating-webhook-tutorial (Apache 2.0)
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	ceModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"

	"github.com/ghodss/yaml"
	"k8s.io/api/admission/v1beta1"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

const moduleCtrlEngine string = "ctrl-engine"
const typeActive string = "active"
const channelCtrlActive string = moduleCtrlEngine + "-" + typeActive
const meepOrigin = "scenario"

// Active scenarion name
var activeScenarioName string

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()
)

type WebhookServer struct {
	sidecarConfig *Config
	server        *http.Server
}

// Webhook Server parameters
type WhSvrParameters struct {
	port           int    // webhook server port
	certFile       string // path to the x509 certificate for https
	keyFile        string // path to the x509 private key matching `CertFile`
	sidecarCfgFile string // path to sidecar injector configuration file
}

type Config struct {
	Containers []corev1.Container `yaml:"containers"`
	Volumes    []corev1.Volume    `yaml:"volumes"`
}

type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

func init() {
	_ = corev1.AddToScheme(runtimeScheme)
	_ = admissionregistrationv1beta1.AddToScheme(runtimeScheme)
}

func activeDBConnect() (err error) {
	// Connect to Active DB
	err = DBConnect()
	if err != nil {
		log.Error("Failed connection to Active DB. Error: ", err)
		return err
	}
	log.Info("Connected to Active DB")

	// Subscribe to Pub-Sub events for MEEP Controller
	// NOTE: Current implementation is RedisDB Pub-Sub
	err = Subscribe(channelCtrlActive)
	if err != nil {
		log.Error("Failed to subscribe to Pub/Sub events. Error: ", err)
		return
	}
	log.Info("Subscribed to Pub/Sub events")

	// Initialize using current active scenario
	processActiveScenarioUpdate()

	return nil
}

func activeDBListen() {
	// Listen for subscribed events. Provide event handler method.
	_ = Listen(eventHandler)
}

func eventHandler(channel string, payload string) {
	// Handle Message according to Rx Channel
	switch channel {

	// MEEP Ctrl Engine active scenario update Channel
	case channelCtrlActive:
		log.Debug("Event received on channel: ", channelCtrlActive)
		processActiveScenarioUpdate()

	default:
		log.Warn("Unsupported channel")
	}
}

func processActiveScenarioUpdate() {
	// Retrieve active scenario from DB
	jsonScenario, err := DBJsonGetEntry(moduleCtrlEngine+":"+typeActive, ".")
	if err != nil {
		log.Error(err.Error())
		clearScenario()
		return
	}

	// Unmarshal Active scenario
	var scenario ceModel.Scenario
	err = json.Unmarshal([]byte(jsonScenario), &scenario)
	if err != nil {
		log.Error(err.Error())
		clearScenario()
		return
	}

	// Parse scenario
	parseScenario(scenario)
}

func clearScenario() {
	log.Debug("clearScenario() -- Resetting all variables")
	activeScenarioName = ""
}

func parseScenario(scenario ceModel.Scenario) {
	log.Debug("parseScenario")

	// Update active scenatio name
	activeScenarioName = scenario.Name
	log.Info("Active scenario name set to: ", activeScenarioName)
}

func loadConfig(configFile string) (*Config, error) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Retrieve App Name from provided network element name string, if any
func getAppName(name string) string {
	names := bytes.Split([]byte(name), []byte(activeScenarioName+"-"))
	if len(names) != 2 {
		return ""
	}
	return string(names[1])
}

func getSidecarPatch(template corev1.PodTemplateSpec, sidecarConfig *Config, meepAppName string) (patch []byte, err error) {

	// Apply labels
	newLabels := make(map[string]string)
	newLabels["meepApp"] = meepAppName
	newLabels["meepOrigin"] = meepOrigin
	newLabels["meepScenario"] = activeScenarioName
	newLabels["processId"] = meepAppName

	// Add environment variables to sidecar containers
	var envVars []corev1.EnvVar
	var envVar corev1.EnvVar
	envVar.Name = "MEEP_POD_NAME"
	envVar.Value = meepAppName
	envVars = append(envVars, envVar)

	var sidecarContainers []corev1.Container
	for _, container := range sidecarConfig.Containers {
		container.Env = envVars
		sidecarContainers = append(sidecarContainers, container)
	}

	// Create patch operations
	var patchOps []patchOperation
	patchOps = append(patchOps, addContainer(template.Spec.Containers, sidecarContainers, "/spec/template/spec/containers")...)
	patchOps = append(patchOps, addVolume(template.Spec.Volumes, sidecarConfig.Volumes, "/spec/template/spec/volumes")...)
	patchOps = append(patchOps, updateLabels(template.ObjectMeta.Labels, newLabels, "/spec/template/metadata/labels")...)

	// Serialize patch
	patch, err = json.Marshal(patchOps)
	if err != nil {
		return nil, err
	}

	return patch, nil
}

func addContainer(target, added []corev1.Container, basePath string) (patch []patchOperation) {
	first := len(target) == 0
	var value interface{}
	for _, add := range added {
		value = add
		path := basePath
		if first {
			first = false
			value = []corev1.Container{add}
		} else {
			path = path + "/-"
		}

		patch = append(patch, patchOperation{
			Op:    "add",
			Path:  path,
			Value: value,
		})
	}
	return patch
}

func addVolume(target, added []corev1.Volume, basePath string) (patch []patchOperation) {
	first := len(target) == 0
	var value interface{}
	for _, add := range added {
		value = add
		path := basePath
		if first {
			first = false
			value = []corev1.Volume{add}
		} else {
			path = path + "/-"
		}
		patch = append(patch, patchOperation{
			Op:    "add",
			Path:  path,
			Value: value,
		})
	}
	return patch
}

func updateLabels(target map[string]string, added map[string]string, basePath string) (patch []patchOperation) {
	for key, value := range added {
		path := basePath + "/" + key
		op := "add"
		if target != nil && target[key] != "" {
			op = "replace"
		}
		patch = append(patch, patchOperation{
			Op:    op,
			Path:  path,
			Value: value,
		})
	}
	return patch
}

// main mutation process
func (whsvr *WebhookServer) mutate(ar *v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	req := ar.Request
	log.Info("Mutate request Name[", req.Name, "] Kind[", req.Kind, "] Namespace[", req.Namespace, "]")

	// Ignore if no active scenario
	if activeScenarioName == "" {
		log.Info("No active scenario. Ignoring request...")
		return &v1beta1.AdmissionResponse{
			Allowed: true,
		}
	}

	// Retrieve resource-specific information
	var resourceName string
	var template corev1.PodTemplateSpec

	switch req.Kind.Kind {
	case "Deployment":
		// Unmarshal Deployment
		var deployment appsv1.Deployment
		if err := json.Unmarshal(req.Object.Raw, &deployment); err != nil {
			log.Error("Could not unmarshal raw object: ", err.Error())
			return &v1beta1.AdmissionResponse{
				Result: &metav1.Status{
					Message: err.Error(),
				},
			}
		}
		log.Info("Deployment Name: ", deployment.Name)
		resourceName = deployment.Name
		template = deployment.Spec.Template

	case "StatefulSet":
		// Unmarshal StatefulSet
		var statefulset appsv1.StatefulSet
		if err := json.Unmarshal(req.Object.Raw, &statefulset); err != nil {
			log.Error("Could not unmarshal raw object: ", err.Error())
			return &v1beta1.AdmissionResponse{
				Result: &metav1.Status{
					Message: err.Error(),
				},
			}
		}
		log.Info("StatefulSet Name: ", statefulset.Name)
		resourceName = statefulset.Name
		template = statefulset.Spec.Template

	default:
		log.Info("Unsupported admission request Kind[", req.Kind.Kind, "]")
		return &v1beta1.AdmissionResponse{
			Allowed: true,
		}
	}

	// Retrieve App Name from resource name
	meepAppName := getAppName(resourceName)
	if meepAppName == "" {
		log.Info("Resource not part of active scenario. Ignoring request...")
		return &v1beta1.AdmissionResponse{
			Allowed: true,
		}
	}
	log.Info("MEEP App Name: ", meepAppName)

	// Get sidecar patch
	patch, err := getSidecarPatch(template, whsvr.sidecarConfig, meepAppName)
	if err != nil {
		return &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	log.Debug("AdmissionResponse: patch=", string(patch))
	return &v1beta1.AdmissionResponse{
		Allowed: true,
		Patch:   patch,
		PatchType: func() *v1beta1.PatchType {
			pt := v1beta1.PatchTypeJSONPatch
			return &pt
		}(),
	}
}

// Serve method for webhook server
func (whsvr *WebhookServer) serve(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	if len(body) == 0 {
		log.Error("empty body")
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		log.Error("Content-Type=", contentType, ", expect application/json")
		http.Error(w, "invalid Content-Type, expect `application/json`", http.StatusUnsupportedMediaType)
		return
	}

	var admissionResponse *v1beta1.AdmissionResponse
	ar := v1beta1.AdmissionReview{}
	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		log.Error("Can't decode body: ", err.Error())
		admissionResponse = &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	} else {
		admissionResponse = whsvr.mutate(&ar)
	}

	admissionReview := v1beta1.AdmissionReview{}
	if admissionResponse != nil {
		admissionReview.Response = admissionResponse
		if ar.Request != nil {
			admissionReview.Response.UID = ar.Request.UID
		}
	}

	resp, err := json.Marshal(admissionReview)
	if err != nil {
		log.Error("Can't encode response: ", err.Error())
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
	}
	log.Info("Ready to write reponse ...")
	if _, err := w.Write(resp); err != nil {
		log.Error("Can't write response: ", err.Error())
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}
}
