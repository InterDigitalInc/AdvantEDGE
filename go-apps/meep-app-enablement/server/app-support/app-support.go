/*
 * Copyright (c) 2020  InterDigital Communications, Inc
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
	//	"time"

	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	msmgmt "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-app-enablement/server/service-mgmt"
	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"

	"github.com/gorilla/mux"
)

const mappsupportBasePath = "/mec_app_support/v1/"
const mappsupportKey = "mec-app-support"
const appEnablementKey = "app-enablement"
const ACTIVE = "ACTIVE"
const INACTIVE = "INACTIVE"

var mutex *sync.Mutex

//const logModuleMSMgmt = "meep-app-enablement"
//const serviceName = "MEC Service Management"

var redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"

//var influxAddr string = "http://meep-influxdb.default.svc.cluster.local:8086"

var APP_ENABLEMENT_DB = 5

var rc *redis.Connector
var hostUrl *url.URL
var sandboxName string
var selfName string
var basePath string
var baseKey string
var appEnablementBaseKey string

//var expiryTicker *time.Ticker

var nextSubscriptionIdAvailable int

/*func notImplemented(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusNotImplemented)
}
*/
func Init(globalMutex *sync.Mutex) (err error) {
	mutex = globalMutex
	// Retrieve Sandbox name from environment variable
	sandboxNameEnv := strings.TrimSpace(os.Getenv("MEEP_SANDBOX_NAME"))
	if sandboxNameEnv != "" {
		sandboxName = sandboxNameEnv
	}
	if sandboxName == "" {
		err = errors.New("MEEP_SANDBOX_NAME env variable not set")
		log.Error(err.Error())
		return err
	}
	log.Info("MEEP_SANDBOX_NAME: ", sandboxName)

	// hostUrl is the url of the node serving the resourceURL
	// Retrieve public url address where service is reachable, if not present, use Host URL environment variable
	hostUrl, err = url.Parse(strings.TrimSpace(os.Getenv("MEEP_PUBLIC_URL")))
	if err != nil || hostUrl == nil || hostUrl.String() == "" {
		hostUrl, err = url.Parse(strings.TrimSpace(os.Getenv("MEEP_HOST_URL")))
		if err != nil {
			hostUrl = new(url.URL)
		}
	}
	log.Info("resource URL: ", hostUrl)

	selfNameEnv := strings.TrimSpace(os.Getenv("MEEP_SELF_NAME"))
	if selfNameEnv != "" {
		selfName = selfNameEnv
	}
	//TODO
	/*
		if selfName == "" {
			err = errors.New("MEEP_SELF_NAME env variable not set")
			log.Error(err.Error())
			return err
		}
	*/
	//HARDCODE SOMETHING
	selfName = "mep1"
	log.Info("MEEP_SELF_NAME: ", selfName)

	// Set base path
	basePath = "/" + sandboxName + mappsupportBasePath
	// Get base store key
	baseKey = dkm.GetKeyRoot(sandboxName) + mappsupportKey
	appEnablementBaseKey = dkm.GetKeyRoot(sandboxName) + selfName + ":" + appEnablementKey

	// Connect to Redis DB
	rc, err = redis.NewConnector(redisAddr, APP_ENABLEMENT_DB)
	if err != nil {
		log.Error("Failed connection to Redis DB. Error: ", err)
		return err
	}

	_ = rc.DBFlush(baseKey)
	_ = rc.DBFlush(appEnablementBaseKey)

	log.Info("Connected to Redis DB")

	reInit()
	/*
		expiryTicker = time.NewTicker(time.Second)
		go func() {
			for range expiryTicker.C {
				//checkForExpiredSubscriptions()
			}
		}()
	*/
	return nil
}

// reInit - finds the value already in the DB to repopulate local stored info
func reInit() {
	//next available subsId will be overrriden if subscriptions already existed
	nextSubscriptionIdAvailable = 1
}

// Run - Start APP support
func Run() (err error) {
	return nil
}

// Stop - Stop APP support
func Stop() (err error) {
	return nil
}

func applicationsConfirmReadyPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Info("applicationsConfirmReadyPOST")

	vars := mux.Vars(r)
	appInstanceId := vars["appInstanceId"]

	var confirmation AppReadyConfirmation
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&confirmation)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//check if entry exist for the application in the DB
	key := appEnablementBaseKey + ":apps:" + appInstanceId
	fields, err := rc.GetEntry(key)
	if err != nil || len(fields) == 0 {
		log.Error("AppInstanceId does not exist, app is not running")
		http.Error(w, "AppInstanceId does not exist, app is not running", http.StatusBadRequest)
		return
	}

	//checking for mandatory properties
	if confirmation.Indication == nil {
		log.Error("Mandatory Indication not present")
		http.Error(w, "Mandatory Indication not present", http.StatusBadRequest)
		return
	}
	switch *confirmation.Indication {
	case READY:
	default:
		log.Error("Mandatory OperationAction value not valid")
		http.Error(w, "Mandatory OperationAction value not valid", http.StatusBadRequest)
		return
	}

	err = updateAllServices(appInstanceId, msmgmt.ACTIVE)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update entry in DB
	updatedFields := make(map[string]interface{})
	updatedFields["state"] = ACTIVE
	err = rc.SetEntry(key, updatedFields)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func applicationsConfirmTerminationPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Info("applicationsConfirmTerminationPOST")

	vars := mux.Vars(r)
	appInstanceId := vars["appInstanceId"]

	//check if entry exist for the application in the DB
	key := appEnablementBaseKey + ":apps:" + appInstanceId
	fields, err := rc.GetEntry(key)
	if err != nil || len(fields) == 0 {
		log.Error("AppInstanceId does not exist, app is not running")
		http.Error(w, "AppInstanceId does not exist, app is not running", http.StatusBadRequest)
		return
	}

	var confirmation AppTerminationConfirmation
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&confirmation)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//checking for mandatory properties
	if confirmation.OperationAction == nil {
		log.Error("Mandatory OperationAction not present")
		http.Error(w, "Mandatory OperationAction not present", http.StatusBadRequest)
		return
	}
	switch *confirmation.OperationAction {
	case STOPPING, TERMINATING:
	default:
		log.Error("Mandatory OperationAction value not valid")
		http.Error(w, "Mandatory OperationAction value not valid", http.StatusBadRequest)
		return
	}

	//do nothing if state is STOPPING, spec says : retention of state
	if *confirmation.OperationAction == TERMINATING {
		err = updateAllServices(appInstanceId, msmgmt.INACTIVE)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Update entry in DB
	updatedFields := make(map[string]interface{})
	updatedFields["state"] = INACTIVE
	err = rc.SetEntry(key, updatedFields)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func updateAllServices(appInstanceId string, state msmgmt.ServiceState) error {
	var sInfoList msmgmt.ServiceInfoList

	keyName := appEnablementBaseKey + ":apps:" + appInstanceId + ":svcs:*"
	mutex.Lock()
	defer mutex.Unlock()
	err := rc.ForEachJSONEntry(keyName, populateServiceInfoList, &sInfoList)
	if err != nil {
		return err
	}
	for _, sInfo := range sInfoList.ServiceInfos {
		serviceId := sInfo.SerInstanceId
		sInfo.State = &state
		err = rc.JSONSetEntry(appEnablementBaseKey+":apps:"+appInstanceId+":svcs:"+serviceId, ".", msmgmt.ConvertServiceInfoToJson(&sInfo))
		if err != nil {
			return err
		}

	}
	return nil
}

func populateServiceInfoList(key string, jsonInfo string, sInfoList interface{}) error {
	// Get query params & userlist from user data
	data := sInfoList.(*msmgmt.ServiceInfoList)

	if data == nil {
		return errors.New("ServiceInfos not found in serviceInfoList")
	}

	// Retrieve user info from DB
	var sInfo msmgmt.ServiceInfo
	err := json.Unmarshal([]byte(jsonInfo), &sInfo)
	if err != nil {
		return err
	}
	data.ServiceInfos = append(data.ServiceInfos, sInfo)
	return nil
}
