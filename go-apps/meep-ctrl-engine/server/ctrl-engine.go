/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/flimzy/kivik"
	_ "github.com/go-kivik/couchdb"
	"github.com/gorilla/mux"

	log "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-ctrl-engine/log"
	tce "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-tc-engine-client"
	ve "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-virt-engine-client"
)

const scenarioDBName = "scenarios"
const activeScenarioName = "active"
const moduleCtrlEngine string = "ctrl-engine"
const moduleMonEngine string = "mon-engine"

const typeActive string = "active"
const channelCtrlActive string = moduleCtrlEngine + "-" + typeActive

const ALLUP = "0"
const ATLEASTONENOTUP = "1"
const NOUP = "2"

const NB_CORE_PODS = 9 //although virt-engine is not a pod yet... it is considered as one as is appended to the list of pods

var virtEngine *ve.APIClient
var tcEngine *tce.APIClient

var db *kivik.DB

var clientServiceMapList []ClientServiceMap

func getCorePodsList() map[string]bool {

	innerMap := map[string]bool{
		"couchdb-couchdb-0": false,
		"meep-ctrl-engine":  false,
		"meep-initializer":  false,
		"meep-mg-manager":   false,
		"meep-mon-engine":   false,
		"meep-tc-engine":    false,
		"metricbeat":        false,
		"virt-engine":       false,
	}
	return innerMap
}

// Establish new couchDB connection
func connectDb(dbName string) (*kivik.DB, error) {

	// Connect to Couch DB
	log.Debug("Establish new couchDB connection")
	dbClient, err := kivik.New(context.TODO(), "couch", "http://couchdb-svc-couchdb:5984/")
	if err != nil {
		return nil, err
	}

	// Create Scenario DB if id does not exist
	log.Debug("Check if scenario DB exists: " + dbName)
	debExists, err := dbClient.DBExists(context.TODO(), dbName)
	if err != nil {
		return nil, err
	}
	if debExists == false {
		log.Debug("Create new DB: " + dbName)
		err = dbClient.CreateDB(context.TODO(), dbName)
		if err != nil {
			return nil, err
		}
	}

	// Open scenario DB
	log.Debug("Open scenario DB: " + dbName)
	db, err := dbClient.DB(context.TODO(), dbName)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Get scenario from DB
func getScenario(returnNilOnNotFound bool, db *kivik.DB, scenarioName string, scenario interface{}) error {

	// Get scenario from DB
	log.Debug("Get scenario from DB: " + scenarioName)
	row, err := db.Get(context.TODO(), scenarioName)
	if err != nil {
		//that's a call to the couch DB.. in order not to return nil, we override it
		if returnNilOnNotFound == true {
			if err.Error() == "Not Found: deleted" {
				//specifically for the case where there is nothing.. so the scenario object will be empty
				return nil
			}
		}
		return err
	}
	// Decode JSON-encoded document
	return row.ScanDoc(scenario)
}

// Get scenario list from DB
func getScenarioList(db *kivik.DB, scenarioList *ScenarioList) error {

	// Retrieve all scenarios from DB
	log.Debug("Get all scenarios from DB")
	rows, err := db.AllDocs(context.TODO())
	if err != nil {
		return err
	}

	// Loop through scenarios and populate scenario list to return
	log.Debug("Loop through scenarios")
	for rows.Next() {
		var scenario Scenario
		if rows.ID() != activeScenarioName {
			err = getScenario(false, db, rows.ID(), &scenario)
			if err == nil {
				// Append scenario to list
				scenarioList.Scenarios = append(scenarioList.Scenarios, scenario)
			}
		}
	}

	return nil
}

// Add scenario to DB
func addScenario(db *kivik.DB, scenarioName string, scenario interface{}) (string, error) {

	// Add scenario to DB
	log.Debug("Add new scenario to DB: " + scenarioName)
	rev, err := db.Put(context.TODO(), scenarioName, scenario)
	if err != nil {
		return "", err
	}

	// Add active scenario to Redis DB & publish update
	//   - Marshal object to JSON string
	//   - Store in Redis DB as REJSON object
	//   - Publish active scenario update event
	if scenarioName == activeScenarioName {
		jsonScenario, err := json.Marshal(scenario)
		if err != nil {
			log.Error(err.Error())
			return "", err
		}
		err = RedisDBJsonSetEntry(moduleCtrlEngine+":"+scenarioName, ".", string(jsonScenario))
		if err != nil {
			log.Error(err.Error())
			return "", err
		}
		err = RedisDBPublish(channelCtrlActive, "")
		if err != nil {
			log.Error(err.Error())
		}
	}

	return rev, nil
}

// Update scenario in DB
func setScenario(db *kivik.DB, scenarioName string, scenario Scenario) (string, error) {

	// Remove previous version
	err := removeScenario(db, scenarioName)
	if err != nil {
		return "", err
	}

	// Add updated version
	rev, err := addScenario(db, scenarioName, scenario)
	if err != nil {
		return "", err
	}

	return rev, nil
}

// Remove scenario from DB
func removeScenario(db *kivik.DB, scenarioName string) error {

	// Get latest Rev of stored scenario to remove
	rev, err := db.Rev(context.TODO(), scenarioName)
	if err != nil {
		return err
	}

	// Remove scenario from DB
	log.Debug("Remove scenario from DB: " + scenarioName)
	rev, err = db.Delete(context.TODO(), scenarioName, rev)
	if err != nil {
		return err
	}

	// Remove active scenario from Redis DB
	// NOTE: Update not published here because remove is also called on updates
	// TODO: Don't remove on update...
	if scenarioName == activeScenarioName {
		err = RedisDBJsonDelEntry(moduleCtrlEngine+":"+scenarioName, ".")
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}

	return nil
}

// Remove all scenarios from DB
func removeAllScenarios(db *kivik.DB) error {

	// Retrieve all scenarios from DB
	log.Debug("Get all scenarios from DB")
	rows, err := db.AllDocs(context.TODO())
	if err != nil {
		return err
	}

	// Loop through scenarios and remove each one
	log.Debug("Loop through scenarios")
	for rows.Next() {
		removeScenario(db, rows.ID())
	}

	return nil
}

// ceGetActiveScenario retrieves the deployed scenario status
func monitorActiveDeployment(deployment *Deployment) error {
	log.Debug("monitorActiveDeployment")

	// Connect to K8s API Server
	clientset, err := connectToAPISvr()
	if err != nil {
		return err
	}

	// Retrieve all pods
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{LabelSelector: "vertical=edge-ar-vr"})
	if err != nil {
		return err
	}
	if len(pods.Items) == 0 {
		log.Debug("No pods found in the cluster")
		return nil
	}

	// Loop through all pods and populate scneario
	log.Debug("Found ", len(pods.Items), " pods in the cluster")
	for _, pod := range pods.Items {
		podLabels := pod.ObjectMeta.Labels

		// Retrieve pod information
		domainID := podLabels["domainId"]
		domainName := podLabels["domainName"]
		domainType := podLabels["domainType"]
		zoneID := podLabels["zoneId"]
		zoneName := podLabels["zoneName"]
		zoneType := podLabels["zoneType"]
		networkLocationID := podLabels["networkLocationId"]
		networkLocationName := podLabels["networkLocationName"]
		networkLocationType := podLabels["networkLocationType"]
		physicalLocationID := podLabels["physicalLocationId"]
		physicalLocationName := podLabels["physicalLocationName"]
		physicalLocationType := podLabels["physicalLocationType"]
		physicalLocationIsExternal := (podLabels["physicalLocationIsExternal"] == "true")
		processID := podLabels["processId"]
		processName := podLabels["processName"]
		processType := podLabels["processType"]
		processIsExternal := (podLabels["processIsExternal"] == "true")

		log.Debug("domainID[", domainID, "]",
			"domainName[", domainName, "]",
			"domainType[", domainType, "]",
			"zoneID[", zoneID, "]",
			"zoneName[", zoneName, "]",
			"zoneType[", zoneType, "]",
			"networkLocationID[", networkLocationID, "]",
			"networkLocationName[", networkLocationName, "]",
			"networkLocationType[", networkLocationType, "]",
			"physicalLocationID[", physicalLocationID, "]",
			"physicalLocationName[", physicalLocationName, "]",
			"physicalLocationType[", physicalLocationType, "]",
			"physicalLocationIsExternal[", physicalLocationIsExternal, "]",
			"processID[", processID, "]",
			"processName[", processName, "]",
			"processType[", processType, "]",
			"processIsExternal[", processIsExternal, "]")

		// Get domain
		var domain *Domain
		for i, d := range deployment.Domains {
			if d.Id == domainID {
				domain = &deployment.Domains[i]
				break
			}
		}
		if domain == nil {
			var newDomain Domain
			newDomain.Id = domainID
			newDomain.Name = domainName
			newDomain.Type_ = domainType
			deployment.Domains = append(deployment.Domains, newDomain)
			domain = &deployment.Domains[len(deployment.Domains)-1]
		}

		// Get zone
		var zone *Zone
		for i, z := range domain.Zones {
			if z.Id == zoneID {
				zone = &domain.Zones[i]
				break
			}
		}
		if zone == nil {
			var newZone Zone
			newZone.Id = zoneID
			newZone.Name = zoneName
			newZone.Type_ = zoneType
			domain.Zones = append(domain.Zones, newZone)
			zone = &domain.Zones[len(domain.Zones)-1]
		}

		// Get networkLocation
		var networkLocation *NetworkLocation
		for i, nl := range zone.NetworkLocations {
			if nl.Id == networkLocationID {
				networkLocation = &zone.NetworkLocations[i]
				break
			}
		}
		if networkLocation == nil {
			var newNetworkLocation NetworkLocation
			newNetworkLocation.Id = networkLocationID
			newNetworkLocation.Name = networkLocationName
			newNetworkLocation.Type_ = networkLocationType
			zone.NetworkLocations = append(zone.NetworkLocations, newNetworkLocation)
			networkLocation = &zone.NetworkLocations[len(zone.NetworkLocations)-1]
		}

		// Get physicalLocation
		var physicalLocation *PhysicalLocation
		for i, nl := range networkLocation.PhysicalLocations {
			if nl.Id == physicalLocationID {
				physicalLocation = &networkLocation.PhysicalLocations[i]
				break
			}
		}
		if physicalLocation == nil {
			var newPhysicalLocation PhysicalLocation
			newPhysicalLocation.Id = physicalLocationID
			newPhysicalLocation.Name = physicalLocationName
			newPhysicalLocation.Type_ = physicalLocationType
			newPhysicalLocation.IsExternal = physicalLocationIsExternal
			networkLocation.PhysicalLocations = append(networkLocation.PhysicalLocations, newPhysicalLocation)
			physicalLocation = &networkLocation.PhysicalLocations[len(networkLocation.PhysicalLocations)-1]
		}

		// Get Process
		var process *Process
		for i, p := range physicalLocation.Processes {
			if p.Id == physicalLocationID {
				process = &physicalLocation.Processes[i]
				break
			}
		}
		if process != nil {
			log.Debug("Process[", process.Id, "] already exists in domain[", domain.Id, "] zone[", zone.Id, "] networkLocation[", networkLocation.Id, "] physicalLocation[", physicalLocation.Id, "]")
		} else {
			var newProcess Process
			newProcess.Id = processID
			newProcess.Name = processName
			newProcess.Type_ = processType
			newProcess.IsExternal = processIsExternal
			physicalLocation.Processes = append(physicalLocation.Processes, newProcess)
		}
	}

	return nil
}

func populateClientServiceMap(activeScenario *Scenario) {

	// Clear client service mapping if there is no active scenario
	if activeScenario == nil {
		clientServiceMapList = nil
		return
	}

	// Parse through scenario and fill external client service mappings
	for _, domain := range activeScenario.Deployment.Domains {
		for _, zone := range domain.Zones {
			for _, nl := range zone.NetworkLocations {
				for _, pl := range nl.PhysicalLocations {
					for _, proc := range pl.Processes {
						if proc.IsExternal {
							// Create new client service map
							var clientServiceMap ClientServiceMap
							clientServiceMap.Client = proc.Name
							for _, serviceMap := range proc.ExternalConfig.IngressServiceMap {
								clientServiceMap.ServiceMap = append(clientServiceMap.ServiceMap, serviceMap)
							}

							// Add new map to list
							clientServiceMapList = append(clientServiceMapList, clientServiceMap)
						}
					}
				}
			}
		}
	}
}

func connectToAPISvr() (*kubernetes.Clientset, error) {
	// Create the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return clientset, nil
}

// CtrlEngineInit Initializes the Controller Engine
func CtrlEngineInit() (err error) {
	log.Debug("CtrlEngineInit")

	// Make Scenario DB connection
	db, err = connectDb(scenarioDBName)
	if err != nil {
		log.Error("Failed connection to Scenario DB. Error: ", err)
		return err
	}
	log.Info("Connected to Scenario DB")

	// Connect to Redis DB
	err = RedisDBConnect()
	if err != nil {
		log.Error("Failed connection to Active DB. Error: ", err)
		return err
	}
	log.Info("Connected to Active DB")

	// Create client for Virtualization Engine API
	veCfg := ve.NewConfiguration()
	veCfg.BasePath = "http://meep-virt-engine/v1"
	virtEngine = ve.NewAPIClient(veCfg)
	if virtEngine == nil {
		log.Debug("Cannot find the Virtualization Engine API")
		return err
	}
	log.Info("Created Virt Engine client")

	// Create client for TC Controller API
	tcCfg := tce.NewConfiguration()
	tcCfg.BasePath = "http://meep-tc-engine/v1"
	tcEngine = tce.NewAPIClient(tcCfg)
	if tcEngine == nil {
		log.Debug("Cannot find the TC Engine API")
		return err
	}
	log.Info("Created TC Engine client")

	return nil
}

// Create a new scenario in store
func ceCreateScenario(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceCreateScenario")

	// Get scenario name from request parameters
	vars := mux.Vars(r)
	scenarioName := vars["name"]
	log.Debug("Scenario name: ", scenarioName)

	// Retrieve scenario from request body
	var scenario Scenario
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&scenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Add new scenario to DB
	rev, err := addScenario(db, scenarioName, scenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Debug("Scenario added with rev: ", rev)

	// OK
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// Delete scenario from store
func ceDeleteScenario(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceDeleteScenario")

	// Get scenario name from request parameters
	vars := mux.Vars(r)
	scenarioName := vars["name"]
	log.Debug("Scenario name: ", scenarioName)

	// Remove scenario from DB
	err := removeScenario(db, scenarioName)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// OK
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// Remove all scenarios from DB
func ceDeleteScenarioList(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceDeleteScenarioList")

	// Remove all scenario from DB
	err := removeAllScenarios(db)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// Retrieve the requested scenario
func ceGetScenario(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceGetScenario")

	// Get scenario name from request parameters
	vars := mux.Vars(r)
	scenarioName := vars["name"]
	log.Debug("Scenario name: ", scenarioName)

	// Validate scenario name
	if scenarioName == "" {
		log.Debug("Invalid scenario name")
		http.Error(w, "Invalid scenario name", http.StatusBadRequest)
		return
	}

	// Retrieve scenario from DB
	var scenario Scenario
	err := getScenario(false, db, scenarioName, &scenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Format response
	jsonResponse, err := json.Marshal(scenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

func ceGetScenarioList(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceGetScenarioList")

	// Retrieve scenario list from DB
	var scenarioList ScenarioList
	err := getScenarioList(db, &scenarioList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Format response
	jsonResponse, err := json.Marshal(scenarioList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

// Update stored scenario
func ceSetScenario(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceSetScenario")

	// Get scenario name from request parameters
	vars := mux.Vars(r)
	scenarioName := vars["name"]
	log.Debug("Scenario name: ", scenarioName)

	// Retrieve scenario from request body
	var scenario Scenario
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&scenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update scenario in DB
	rev, err := setScenario(db, scenarioName, scenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	log.Debug("Scenario updated with rev: ", rev)

	// OK
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// Activate/Deploy scenario
func ceActivateScenario(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceActivateScenario")

	// Get scenario name from request parameters
	vars := mux.Vars(r)
	scenarioName := vars["name"]
	log.Debug("Scenario name: ", scenarioName)

	// Make sure scenario is not already deployed
	var activeScenario Scenario
	err := getScenario(false, db, activeScenarioName, &activeScenario)
	if err == nil {
		log.Error("Scenario already active")
		http.Error(w, "Scenario already active", http.StatusBadRequest)
		return
	}

	// Retrieve scenario to activate from DB

	// !!!!! IMPORTANT NOTE !!!!!
	// Scenario stored in DB is unmarshalled into a VE Scenario object
	var veScenario ve.Scenario
	err = getScenario(false, db, scenarioName, &veScenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Set active scenario in DB
	rev, err := addScenario(db, activeScenarioName, veScenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Debug("Active scenario set with rev: ", rev)

	// Activate scenario in virtualization Engine
	resp, err := virtEngine.ScenarioDeploymentApi.ActivateScenario(nil, veScenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		removeScenario(db, activeScenarioName)
		return
	}

	// Apply network characteristics on active scenario

	// !!!!! IMPORTANT NOTE !!!!!
	// Active scenario stored in DB is unmarshalled into a TC Engine Scenario object
	var tceScenario tce.Scenario
	err = getScenario(false, db, activeScenarioName, &tceScenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		removeScenario(db, activeScenarioName)
		return
	}

	// Activate scenario in TC Controller
	resp, err = tcEngine.ScenarioDeploymentApi.ActivateScenario(nil, tceScenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		removeScenario(db, activeScenarioName)
		return
	}

	// Retrieve active scenario stored in DB
	err = getScenario(false, db, activeScenarioName, &activeScenario)
	if err != nil {
		log.Error("Scenario not active")
		http.Error(w, "Scenario not active", http.StatusBadRequest)
		removeScenario(db, activeScenarioName)
		return
	}

	// Populate active external client service map
	populateClientServiceMap(&activeScenario)

	// Return response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if resp != nil {
		w.WriteHeader(resp.StatusCode)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// ceGetActiveScenario retrieves the deployed scenario status
func ceGetActiveScenario(w http.ResponseWriter, r *http.Request) {
	log.Debug("CEGetActiveScenario")

	// Retrieve active scenario
	var scenario Scenario
	err := getScenario(true, db, activeScenarioName, &scenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// NOTE: For now, return full scenario without status information.
	// Eventually, we will need to fetch latest status information from DB or k8s.

	// // Create Scenario object
	// var deployment Deployment
	// var scenario Scenario
	// scenario.Name = "Edge-Enabled 5G Video"
	// scenario.Deployment = &deployment

	// err := monitorActiveDeployment(&deployment)
	// if err != nil {
	// 	log.Error(err.Error())
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// Format response
	jsonResponse, err := json.Marshal(scenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

// ceGetActiveClientServiceMaps retrieves the deployed scenario external client service mappings
// NOTE: query parameters 'client' and 'service' may be specified to filter results
func ceGetActiveClientServiceMaps(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceGetActiveClientServiceMaps")
	var filteredList *[]ClientServiceMap

	// Retrieve client ID & service name from query parameters
	query := r.URL.Query()
	client := query.Get("client")
	service := query.Get("service")

	// Filter only requested service mappings from client service map list
	if client == "" && service == "" {
		// Any client & service
		filteredList = &clientServiceMapList
	} else {
		filteredList = new([]ClientServiceMap)
		if service == "" {
			// Any service for requested client
			for _, clientServiceMap := range clientServiceMapList {
				if clientServiceMap.Client == client {
					*filteredList = append(*filteredList, clientServiceMap)
					break
				}
			}
		} else if client == "" {
			// Any client for requested service
			for _, clientServiceMap := range clientServiceMapList {
				var svcMap ClientServiceMap
				svcMap.Client = clientServiceMap.Client
				for _, serviceMap := range clientServiceMap.ServiceMap {
					if serviceMap.Name == service {
						svcMap.ServiceMap = append(svcMap.ServiceMap, serviceMap)
						break
					}
				}

				// Only append if at least one match found
				if len(svcMap.ServiceMap) > 0 {
					*filteredList = append(*filteredList, svcMap)
				}
			}
		} else {
			// Requested client and service
			for _, clientServiceMap := range clientServiceMapList {
				if clientServiceMap.Client == client {
					for _, serviceMap := range clientServiceMap.ServiceMap {
						if serviceMap.Name == service {
							var svcMap ClientServiceMap
							svcMap.Client = clientServiceMap.Client
							svcMap.ServiceMap = append(svcMap.ServiceMap, serviceMap)

							*filteredList = append(*filteredList, svcMap)
							break
						}
					}
					break
				}
			}
		}
	}

	// Format response
	jsonResponse, err := json.Marshal(*filteredList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

// Terminate the active/deployed scenario
func ceTerminateScenario(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceTerminateScenario")

	// Clear active external client service map
	populateClientServiceMap(nil)

	// Retrieve active scenario from DB
	var scenario Scenario
	err := getScenario(false, db, activeScenarioName, &scenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Remove active scenario from DB
	err = removeScenario(db, activeScenarioName)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Publish active scenario update event
	err = RedisDBPublish(channelCtrlActive, "")
	if err != nil {
		log.Error(err.Error())
	}

	// Terminate scenario in TC Controller
	resp, err := tcEngine.ScenarioDeploymentApi.DeleteNetworkCharacteristicsTable(nil)
	if err != nil {
		log.Error(err.Error())
	}

	// Terminate scenario in virtualization Engine
	resp, err = virtEngine.ScenarioDeploymentApi.TerminateScenario(nil, scenario.Name)
	if err != nil {
		log.Error(err.Error())
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if resp != nil {
		w.WriteHeader(resp.StatusCode)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func ceGetEventList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func sendEventNetworkCharacteristics(event Event) (string, int) {

	// Retrieve active scenario
	var scenario Scenario
	err := getScenario(false, db, activeScenarioName, &scenario)
	if err != nil {
		return err.Error(), http.StatusNotFound
	}

	elementFound := false
	netChar := event.EventNetworkCharacteristicsUpdate

	// Retrieve element name & type
	elementName := netChar.ElementName
	elementType := strings.ToUpper(netChar.ElementType)

	// Find the element
	if elementType == "SCENARIO" {
		scenario.Deployment.InterDomainLatency = netChar.Latency
		scenario.Deployment.InterDomainLatencyVariation = netChar.LatencyVariation
		scenario.Deployment.InterDomainThroughput = netChar.Throughput
		scenario.Deployment.InterDomainPacketLoss = netChar.PacketLoss
		elementFound = true
	}

	for index_d, d := range scenario.Deployment.Domains {
		if elementFound == true {
			break
		} else if elementType == "OPERATOR" && elementName == d.Name {
			domain := &scenario.Deployment.Domains[index_d]
			domain.InterZoneLatency = netChar.Latency
			domain.InterZoneLatencyVariation = netChar.LatencyVariation
			domain.InterZoneThroughput = netChar.Throughput
			domain.InterZonePacketLoss = netChar.PacketLoss
			elementFound = true
			break
		}

		// Parse zones
		for index_z, z := range d.Zones {
			if elementFound == true {
				break
			} else if elementType == "ZONE-INTER-EDGE" && elementName == z.Name {
				zone := &scenario.Deployment.Domains[index_d].Zones[index_z]
				zone.InterEdgeLatency = netChar.Latency
				zone.InterEdgeLatencyVariation = netChar.LatencyVariation
				zone.InterEdgeThroughput = netChar.Throughput
				zone.InterEdgePacketLoss = netChar.PacketLoss
				elementFound = true
				break
			} else if elementType == "ZONE-INTER-FOG" && elementName == z.Name {
				zone := &scenario.Deployment.Domains[index_d].Zones[index_z]
				zone.InterFogLatency = netChar.Latency
				zone.InterFogLatencyVariation = netChar.LatencyVariation
				zone.InterFogThroughput = netChar.Throughput
				zone.InterFogPacketLoss = netChar.PacketLoss
				elementFound = true
				break
			} else if elementType == "ZONE-EDGE-FOG" && elementName == z.Name {
				zone := &scenario.Deployment.Domains[index_d].Zones[index_z]
				zone.EdgeFogLatency = netChar.Latency
				zone.EdgeFogLatencyVariation = netChar.LatencyVariation
				zone.EdgeFogThroughput = netChar.Throughput
				zone.EdgeFogPacketLoss = netChar.PacketLoss
				elementFound = true
				break
			}

			// Parse Network Locations
			for index_nl, nl := range z.NetworkLocations {
				if elementType == "POA" && elementName == nl.Name {
					netloc := &scenario.Deployment.Domains[index_d].Zones[index_z].NetworkLocations[index_nl]
					netloc.TerminalLinkLatency = netChar.Latency
					netloc.TerminalLinkLatencyVariation = netChar.LatencyVariation
					netloc.TerminalLinkThroughput = netChar.Throughput
					netloc.TerminalLinkPacketLoss = netChar.PacketLoss
					elementFound = true
					break
				}
			}
		}
	}

	if elementFound == true {
		log.Debug("element was found and updates should be applied")
	} else {
		return "Element not found in the scenario", http.StatusNotFound
	}

	// Store updated active scenario in DB
	rev, err := setScenario(db, activeScenarioName, scenario)
	if err != nil {
		return err.Error(), http.StatusNotFound
	}
	log.Debug("Active scenario updated with rev: ", rev)

	// TODO in Execution Engine:
	//    - Update any deployed location services
	//    - Inform monitoring engine?

	// Inform TC Controller of updated scenario

	// !!!!! IMPORTANT NOTE !!!!!
	// Active scenario stored in DB is unmarshalled into a TC Scenario object
	var tceScenario tce.Scenario
	err = getScenario(false, db, activeScenarioName, &tceScenario)
	if err != nil {
		return err.Error(), http.StatusNotFound
	}

	// Activate scenario in TC Controller
	_, err = tcEngine.ScenarioDeploymentApi.ActivateScenario(nil, tceScenario)
	if err != nil {
		return err.Error(), http.StatusNotFound
	}

	return "", -1
}

func sendEventUeMobility(event Event) (string, int) {

	// Retrieve active scenario
	var scenario Scenario
	err := getScenario(false, db, activeScenarioName, &scenario)
	if err != nil {
		return err.Error(), http.StatusNotFound
	}

	// Retrieve UE name and destination PoA name
	ueName := event.EventUeMobility.Ue
	poaName := event.EventUeMobility.Dest

	var oldNL *NetworkLocation
	var newNL *NetworkLocation
	var ue *PhysicalLocation
	var ueIndex int

	// Find UE & destination PoA
	log.Debug("Searching for UE and destination PoA in active scenario")
	for i := range scenario.Deployment.Domains {
		domain := &scenario.Deployment.Domains[i]

		for j := range domain.Zones {
			zone := &domain.Zones[j]

			for k := range zone.NetworkLocations {
				nl := &zone.NetworkLocations[k]

				// Destination PoA
				if nl.Name == poaName {
					newNL = nl
				}

				for l := range nl.PhysicalLocations {
					pl := &nl.PhysicalLocations[l]

					// UE to move
					if pl.Type_ == "UE" && pl.Name == ueName {
						oldNL = nl
						ue = pl
						ueIndex = l
					}
				}
			}
		}
	}

	// Update UE location if necessary
	if ue != nil && oldNL != nil && newNL != nil && oldNL != newNL {
		log.Debug("Found UE and destination PoA. Updating UE location.")

		// Add UE to new location
		newNL.PhysicalLocations = append(newNL.PhysicalLocations, *ue)

		// Remove UE from old location
		oldNL.PhysicalLocations[ueIndex] = oldNL.PhysicalLocations[len(oldNL.PhysicalLocations)-1]
		oldNL.PhysicalLocations = oldNL.PhysicalLocations[:len(oldNL.PhysicalLocations)-1]

		// Store updated active scenario in DB
		rev, err := setScenario(db, activeScenarioName, scenario)
		if err != nil {
			return err.Error(), http.StatusNotFound
		}
		log.Debug("Active scenario updated with rev: ", rev)

		// TODO in Execution Engine:
		//    - Update any deployed location services
		//    - Inform monitoring engine?

		// Inform TC Controller of updated scenario

		// !!!!! IMPORTANT NOTE !!!!!
		// Active scenario stored in DB is unmarshalled into a TC Scenario object
		var tceScenario tce.Scenario
		err = getScenario(false, db, activeScenarioName, &tceScenario)
		if err != nil {
			return err.Error(), http.StatusNotFound
		}

		// Activate scenario in TC Controller
		_, err = tcEngine.ScenarioDeploymentApi.ActivateScenario(nil, tceScenario)
		if err != nil {
			return err.Error(), http.StatusNotFound
		}

	} else {
		err := "Failed to find UE or destination PoA"
		return err, http.StatusNotFound
	}
	return "", -1
}

func sendEventPoasInRange(event Event) (string, int) {
	var ue *PhysicalLocation

	// Retrieve active scenario
	var scenario Scenario
	err := getScenario(false, db, activeScenarioName, &scenario)
	if err != nil {
		return err.Error(), http.StatusNotFound
	}

	// Retrieve UE name
	ueName := event.EventPoasInRange.Ue

	// Retrieve list of visible POAs and sort them
	poasInRange := event.EventPoasInRange.PoasInRange
	sort.Strings(poasInRange)

	// Find UE
	log.Debug("Searching for UE in active scenario")
	for i := range scenario.Deployment.Domains {
		domain := &scenario.Deployment.Domains[i]

		for j := range domain.Zones {
			zone := &domain.Zones[j]

			for k := range zone.NetworkLocations {
				nl := &zone.NetworkLocations[k]

				for l := range nl.PhysicalLocations {
					pl := &nl.PhysicalLocations[l]

					// UE to update
					if pl.Type_ == "UE" && pl.Name == ueName {
						ue = pl
						break
					}
				}
				if ue != nil {
					break
				}
			}
			if ue != nil {
				break
			}
		}
		if ue != nil {
			break
		}
	}

	// Update POAS in range if necessary
	if ue != nil {
		log.Debug("UE Found. Checking for update to visible POAs")

		// Compare new list of poas with current UE list and update if necessary
		if Equal(poasInRange, ue.NetworkLocationsInRange) == false {
			log.Debug("Updating POAs in range for UE: " + ue.Name)
			ue.NetworkLocationsInRange = poasInRange

			// Store updated active scenario in DB
			rev, err := setScenario(db, activeScenarioName, scenario)
			if err != nil {
				return err.Error(), http.StatusNotFound
			}
			log.Debug("Active scenario updated with rev: ", rev)
		} else {
			log.Debug("POA list unchanged. Ignoring.")
		}
	} else {
		err := "Failed to find UE"
		return err, http.StatusNotFound
	}
	return "", -1
}

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func ceSendEvent(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceSendEvent")

	// Get event type from request parameters
	vars := mux.Vars(r)
	eventType := vars["type"]
	log.Debug("Event Type: ", eventType)

	// Retrieve event from request body
	var event Event
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&event)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Process Event
	var httpStatus int
	var error string
	switch eventType {
	case "UE-MOBILITY":
		error, httpStatus = sendEventUeMobility(event)
	case "NETWORK-CHARACTERISTICS-UPDATE":
		error, httpStatus = sendEventNetworkCharacteristics(event)
	case "POAS-IN-RANGE":
		error, httpStatus = sendEventPoasInRange(event)
	default:
		error = "Unsupported event type"
		httpStatus = http.StatusBadRequest
	}

	if error != "" {
		log.Error(error)
		http.Error(w, error, httpStatus)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ceGetMeepSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ceSetMeepSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func getPodDetails(key string, fields map[string]string, userData interface{}) error {

	podsStatus := userData.(*PodsStatus)
	var podStatus PodStatus
	if fields["meepApp"] != "" {
		podStatus.Name = fields["meepApp"]
	} else {
		podStatus.Name = fields["name"]
	}

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

	podsStatus.PodStatus = append(podsStatus.PodStatus, podStatus)
	return nil
}

func getPodStatesOnly(key string, fields map[string]string, userData interface{}) error {
	podsStatus := userData.(*PodsStatus)
	var podStatus PodStatus
	if fields["meepApp"] != "" {
		podStatus.Name = fields["meepApp"]
	} else {
		podStatus.Name = fields["name"]
	}
	podStatus.LogicalState = fields["logicalState"]

	podsStatus.PodStatus = append(podsStatus.PodStatus, podStatus)

	return nil
}

func ceGetStates(w http.ResponseWriter, r *http.Request) {

	subKey := ""
	var podsStatus PodsStatus
	// Retrieve client ID & service name from query parameters
	query := r.URL.Query()
	longParam := query.Get("long")
	typeParam := query.Get("type")

	detailed := false
	if longParam == "true" {
		detailed = true
	}

	if typeParam == "" {
		subKey = "MO-scenario:"
	} else {
		subKey = "MO-" + typeParam + ":"
	}

	//values for pod name
	keyName := moduleMonEngine + "*" + subKey + "*"

	//get will be unique... but reusing the generic function
	var err error
	if detailed == true {
		err = RedisDBForEachEntry(keyName, getPodDetails, &podsStatus)
	} else {
		err = RedisDBForEachEntry(keyName, getPodStatesOnly, &podsStatus)
	}

	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if typeParam == "core" {
		// ***** virt-engine is not a pod yet, but we need to make sure it is started to have a functional system
		var podStatus PodStatus
		podStatus.Name = "virt-engine"
		//we do not care about the content of the answer, simply that there is one
		_, resp, _ := virtEngine.ScenarioDeploymentApi.GetActiveScenario(nil, "dummy")

		if resp != nil {
			if resp.StatusCode == http.StatusOK {
				podStatus.LogicalState = "Running"
			} else {
				podStatus.LogicalState = "InternalError"
			}
		} else {
			podStatus.LogicalState = "NotRunning"
		}
		podsStatus.PodStatus = append(podsStatus.PodStatus, podStatus)
		// ***** virt-engine running or not code END

		//if some are missing... its because its coming up and as such... we cannot return a success yet... adding one entry that will be false

		var corePods map[string]bool
		corePods = getCorePodsList()

		//loop through each of them by name
		for _, statusPod := range podsStatus.PodStatus {
			for corePod := range corePods {
				if strings.Contains(statusPod.Name, corePod) {
					corePods[corePod] = true
					break
				}
			}
		}

		//loop through the list of pods to see which one might be missing
		for corePod := range corePods {
			if !corePods[corePod] {
				var podStatus PodStatus
				podStatus.Name = corePod
				podStatus.LogicalState = "NotAvailable"
				podsStatus.PodStatus = append(podsStatus.PodStatus, podStatus)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Format response
	jsonResponse, err := json.Marshal(podsStatus)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}
