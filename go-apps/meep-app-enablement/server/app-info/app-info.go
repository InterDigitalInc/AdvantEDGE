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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
	uuid "github.com/google/uuid"
	"github.com/gorilla/mux"
)

const appInfoBasePath = "app_info/v1/"
const appEnablementKey = "app-enablement"
const defaultMepName = "global"
const ACTIVE = "ACTIVE"
const INACTIVE = "INACTIVE"

var redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"

var APP_ENABLEMENT_DB = 0

var mutex *sync.Mutex
var rc *redis.Connector
var hostUrl *url.URL
var sandboxName string
var mepName string = defaultMepName
var basePath string
var baseKey string

var expiryTicker *time.Ticker

var nextAppInstanceIdAvailable int

type ApplicationInfoList struct {
	ApplicationInfos []ApplicationInfo
	filterParameters *FilterParameters
}

type FilterParameters struct {
	appName  string
	appState string
}

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

	// hostUrl is the url of the node serving the resourceURL
	// Retrieve public url address where service is reachable, if not present, use Host URL environment variable
	hostUrl, err = url.Parse(strings.TrimSpace(os.Getenv("MEEP_PUBLIC_URL")))
	if err != nil || hostUrl == nil || hostUrl.String() == "" {
		hostUrl, err = url.Parse(strings.TrimSpace(os.Getenv("MEEP_HOST_URL")))
		if err != nil {
			hostUrl = new(url.URL)
		}
	}
	log.Info("MEEP_HOST_URL: ", hostUrl)

	// Get MEP name
	mepNameEnv := strings.TrimSpace(os.Getenv("MEEP_MEP_NAME"))
	if mepNameEnv != "" {
		mepName = mepNameEnv
	}
	log.Info("MEEP_MEP_NAME: ", mepName)

	// Set base path
	if mepName == defaultMepName {
		basePath = "/" + sandboxName + "/" + appInfoBasePath
	} else {
		basePath = "/" + sandboxName + "/" + mepName + "/" + appInfoBasePath
	}

	// Set base storage key
	baseKey = dkm.GetKeyRoot(sandboxName) + appEnablementKey + ":mep:" + mepName

	// Connect to Redis DB
	rc, err = redis.NewConnector(redisAddr, APP_ENABLEMENT_DB)
	if err != nil {
		log.Error("Failed connection to Redis DB. Error: ", err)
		return err
	}
	_ = rc.DBFlush(baseKey)
	log.Info("Connected to Redis DB")

	reInit()

	expiryTicker = time.NewTicker(time.Second)
	go func() {
		for range expiryTicker.C {
			//checkForExpiredSubscriptions()
		}
	}()
	return nil
}

// reInit - finds the value already in the DB to repopulate local stored info
func reInit() {
	nextAppInstanceIdAvailable = 1
}

// Run - Start
func Run() (err error) {
	return nil
}

// Stop - Stop
func Stop() (err error) {
	return nil
}

func getNewInstanceId() (string, error) {
	/*	appInstanceId := strconv.Itoa(nextAppInstanceIdAvailable)
		nextAppInstanceIdAvailable++
		return appInstanceId
	*/
	//allow 3 tries, if not return an error
	maxNbRetries := 3
	for try := maxNbRetries; try > 0; try-- {
		appInstanceId := uuid.New().String()
		err := validateAppInstanceId(appInstanceId)
		//if there is an error, it means the instance id already exist, get a new one
		if err != nil {
			return appInstanceId, nil
		}
		try--
	}
	return "", errors.New("Can't allocate a unique instance Id")
}

func validateAppInstanceId(appInstanceId string) error {
	//check if entry exist for the application in the DB
	key := baseKey + ":app:" + appInstanceId + ":info"
	fields, err := rc.GetEntry(key)
	if err != nil || len(fields) == 0 {
		if err != nil {
			//problem when reading the DB
			log.Error(err.Error())
			return err
		} else {
			newError := errors.New("AppInstanceId does not exist, app is not running")
			log.Error(newError.Error())
			return newError
		}
	}
	return nil
}

func applicationsPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Info("applicationsPOST")

	var appInfo ApplicationInfo
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&appInfo)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	appInstanceId, err := getNewInstanceId()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	appInfo.AppInstanceId = appInstanceId

	//check if entry exist for the application in the DB
	key := baseKey + ":app:" + appInstanceId + ":info"
	fields, err := rc.GetEntry(key)
	if err == nil && len(fields) > 0 {
		log.Error("AppInstanceId already exists")
		http.Error(w, "AppInstanceId already exists", http.StatusBadRequest)
		return
	}

	//checking for mandatory properties
	if appInfo.AppName == "" {
		log.Error("Mandatory Name not present")
		http.Error(w, "Mandatory Name not present", http.StatusBadRequest)
		return
	}
	if appInfo.State == nil {
		log.Error("Mandatory State not present")
		http.Error(w, "Mandatory State not present", http.StatusBadRequest)
		return
	}
	switch *appInfo.State {
	case APPLICATION_STATE_ACTIVE, APPLICATION_STATE_INACTIVE:
	default:
		log.Error("Mandatory State value not valid")
		http.Error(w, "Mandatory State value not valid", http.StatusBadRequest)
		return
	}

	// create entry in DB
	newAppInfoFields := make(map[string]interface{})
	newAppInfoFields["appInstanceId"] = appInfo.AppInstanceId
	newAppInfoFields["appName"] = appInfo.AppName
	newAppInfoFields["state"] = string(*appInfo.State)
	newAppInfoFields["version"] = appInfo.Version

	err = rc.SetEntry(key, newAppInfoFields)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse := ConvertApplicationInfoToJson(&appInfo)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func applicationsAppInstanceIdPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Info("applicationsAppInstanceIdPUT")

	vars := mux.Vars(r)
	appInstanceId := vars["appInstanceId"]

	var appInfo ApplicationInfo
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&appInfo)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	key := baseKey + ":app:" + appInstanceId + ":info"
	fields, err := rc.GetEntry(key)
	if err != nil || len(fields) == 0 {
		log.Error("AppInstanceId does not already exist")
		http.Error(w, "AppInstanceId does not already exist", http.StatusNotFound)
		return
	}

	//checking for mandatory properties
	if appInfo.AppName == "" {
		log.Error("Mandatory Name not present")
		http.Error(w, "Mandatory Name not present", http.StatusBadRequest)
		return
	}
	if appInfo.State == nil {
		log.Error("Mandatory State not present")
		http.Error(w, "Mandatory State not present", http.StatusBadRequest)
		return
	}
	switch *appInfo.State {
	case APPLICATION_STATE_ACTIVE, APPLICATION_STATE_INACTIVE:
	default:
		log.Error("Mandatory State value not valid")
		http.Error(w, "Mandatory State value not valid", http.StatusBadRequest)
		return
	}

	if appInfo.AppInstanceId != appInstanceId {
		log.Error("Mandatory Application Instance Id parameter and body content not matching")
		http.Error(w, "Mandatory Application Instance Id parameter and body content not matching", http.StatusBadRequest)
		return
	}

	// override entry in DB
	newAppInfoFields := make(map[string]interface{})
	newAppInfoFields["appName"] = appInfo.AppName
	newAppInfoFields["state"] = string(*appInfo.State)
	newAppInfoFields["version"] = appInfo.Version

	err = rc.SetEntry(key, newAppInfoFields)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse := ConvertApplicationInfoToJson(&appInfo)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))

}

func applicationsAppInstanceIdGET(w http.ResponseWriter, r *http.Request) {
	log.Info("applicationsByIdGET")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	appInstanceId := vars["appInstanceId"]

	fields, err := rc.GetEntry(baseKey + ":app:" + appInstanceId + ":info")
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if len(fields) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var appInfo ApplicationInfo
	appInfo.AppName = fields["appName"]
	appInfo.AppInstanceId = fields["appInstanceId"]
	appInfo.Version = fields["version"]
	appInfoState := ApplicationState(fields["state"])
	appInfo.State = &appInfoState
	jsonResponse := ConvertApplicationInfoToJson(&appInfo)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func applicationsAppInstanceIdDELETE(w http.ResponseWriter, r *http.Request) {
	log.Info("applicationsByIdDELETE")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	appInstanceId := vars["appInstanceId"]

	mutex.Lock()
	defer mutex.Unlock()

	// Check if App instance exists
	fields, err := rc.GetEntry(baseKey + ":app:" + appInstanceId + ":info")
	if err != nil || len(fields) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Flush App instance
	_ = rc.DBFlush(baseKey + ":app:" + appInstanceId)
	w.WriteHeader(http.StatusNoContent)
}

func applicationsGET(w http.ResponseWriter, r *http.Request) {
	log.Info("applicationsGET")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	u, _ := url.Parse(r.URL.String())
	q := u.Query()
	appName := q.Get("app_name")
	appState := q.Get("app_state")

	validQueryParams := []string{"app_name", "app_state"}

	found := false
	for queryParam := range q {
		found = false
		for _, validQueryParam := range validQueryParams {
			if queryParam == validQueryParam {
				found = true
				break
			}
		}
		if !found {
			log.Error("Query param not valid: ", queryParam)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	var applicationInfoList ApplicationInfoList
	var filterParameters FilterParameters
	filterParameters.appName = appName
	filterParameters.appState = appState
	applicationInfoList.filterParameters = &filterParameters

	keyName := baseKey + ":app:*:info"

	err := rc.ForEachEntry(keyName, populateApplicationInfoList, &applicationInfoList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse, err := json.Marshal(applicationInfoList.ApplicationInfos)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func populateApplicationInfoList(key string, fields map[string]string, applicationInfoList interface{}) error {
	// Get query params & userlist from user data
	data := applicationInfoList.(*ApplicationInfoList)

	if data == nil {
		return errors.New("ApplicationInfos not found in applicationInfoList")
	}

	// Retrieve user info from DB
	var appInfo ApplicationInfo
	appInfo.AppName = fields["appName"]
	appInfo.AppInstanceId = fields["appInstanceId"]
	appInfoState := ApplicationState(fields["state"])
	appInfo.State = &appInfoState
	appInfo.Version = fields["version"]

	match := true

	if data.filterParameters != nil {
		//compare with filter to return the service info or not
		//checking for wildcard or identical values
		match = (data.filterParameters.appName == "" || data.filterParameters.appName == appInfo.AppName)

		if match {
			match = (data.filterParameters.appState == "" || data.filterParameters.appState == string(*appInfo.State))
		}
	}
	if match {
		data.ApplicationInfos = append(data.ApplicationInfos, appInfo)
	}
	return nil
}
