package server

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	
	log "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-virt-engine/log"
	"github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-virt-engine/helm"
)

func VirtEngineInit() {
	log.Debug("Initializing MEEP Virtualization Engine")
}

func readAndPrintRequest(r *http.Request) {

	// Read the Body content
	var bodyBytes []byte
	bodyBytes, _ = ioutil.ReadAll(r.Body)
	log.Info(bodyBytes)

	// Restore the io.ReadCloser to its original state
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

}

func populateScenario(r *http.Request) (Scenario, error) {

	log.Debug("populateScenario")

	var scenario Scenario

	//readAndPrintRequest(r)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&scenario)

	if err != nil {
		log.Error(err.Error())
		return scenario, err
	}

	return scenario, nil
}

func veActivateScenario(w http.ResponseWriter, r *http.Request) {
	scenario, err := populateScenario(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = CreateYamlScenarioFile(scenario)
	if err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func veGetActiveScenario(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func veSendEvent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func veTerminateScenario(w http.ResponseWriter, r *http.Request) {
	// Read Parameters
	vars := mux.Vars(r)
	name := vars["name"]

	// Retrieve list of releases
	rels, err := helm.GetReleasesName()
	var toDelete []helm.Chart
	for _, rel := range rels {
		if strings.Contains(rel.Name, name) {
			// just keep releases related to the current scenario
			var c helm.Chart
			c.ReleaseName = rel.Name
			toDelete = append(toDelete, c)
		}
	}

	// Delete releases
	if len(toDelete) > 0 {
		err = helm.DeleteReleases(toDelete)
		log.Debug(err)
	}

	// Then delete charts
	homePath := os.Getenv("HOME")
	path := homePath + "/.meep/active/" + name
	if _, err := os.Stat(path); err == nil {
		log.Debug("Removing charts ", path)
		os.RemoveAll(path)
	}

	w.WriteHeader(http.StatusOK)
}
