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

	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
	uuid "github.com/google/uuid"
	"github.com/gorilla/mux"
)

const appEnablementModule = "meep-app-enablement"
const appEnablementKey = "app-enablement"
const ACTIVE = "ACTIVE"
const INACTIVE = "INACTIVE"
const APP_ENABLEMENT_DB = 0

// App Info fields
const fieldAppInstanceId = "id"
const fieldAppName = "name"
const fieldAppType = "type"
const fieldMepName = "mep"
const fieldState = "state"
const fieldVersion = "version"

// MQ payload fields
const mqFieldAppInstanceId = "id"
const mqFieldMepName = "mep"

type AppCtrl struct {
	sandboxName string
	rc          *redis.Connector
	mqLocal     *mq.MsgQueue
	baseKey     string
}

type ApplicationData struct {
	AppInfoList      []dataModel.ApplicationInfo
	FilterParameters *FilterParameters
}

type FilterParameters struct {
	appId   string
	appName string
	appType string
	state   string
	mepName string
}

// App Controller
var appCtrl *AppCtrl

// Initialize App Controller
func appCtrlInit(sandboxName string, mqLocal *mq.MsgQueue) (err error) {

	// Create new App Controller
	appCtrl = new(AppCtrl)
	appCtrl.sandboxName = sandboxName
	appCtrl.mqLocal = mqLocal

	// Set base storage key
	appCtrl.baseKey = dkm.GetKeyRoot(sandboxName) + appEnablementKey

	// Connect to Redis DB
	appCtrl.rc, err = redis.NewConnector(redisDBAddr, APP_ENABLEMENT_DB)
	if err != nil {
		log.Error("Failed connection to Redis DB. Error: ", err)
		return err
	}
	_ = appCtrl.rc.DBFlush(appCtrl.baseKey)
	log.Info("Connected to Redis DB")

	return nil
}

// Start App Controller
func appCtrlRun() (err error) {
	return nil
}

// Stop App Controller
func appCtrlStop() (err error) {
	return nil
}

func applicationsPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Info("applicationsPOST")

	// Retrieve & validate Application info from request body
	var appInfo dataModel.ApplicationInfo
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&appInfo)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = validateAppInfo(&appInfo, "")
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Obtain a new App Instance ID if none provided
	if appInfo.Id == "" {
		appInstanceId, err := getNewInstanceId()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		appInfo.Id = appInstanceId
	}

	// Make sure App instance does not exist
	if appInstanceExists(appInfo.MepName, appInfo.Id) {
		errStr := "AppInstanceId already exists"
		log.Error(errStr)
		http.Error(w, errStr, http.StatusBadRequest)
		return
	}

	// create entry in DB
	err = setAppInfo(&appInfo)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	jsonResponse := convertApplicationInfoToJson(&appInfo)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func applicationsAppInstanceIdPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Info("applicationsAppInstanceIdPUT")
	vars := mux.Vars(r)
	appInstanceId := vars["appInstanceId"]

	// Retrieve & Validate Application info from request body
	var appInfo dataModel.ApplicationInfo
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&appInfo)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = validateAppInfo(&appInfo, appInstanceId)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// Make sure App instance exists
	if !appInstanceExists(appInfo.MepName, appInfo.Id) {
		errStr := "AppInstanceId does not exist"
		log.Error(errStr)
		http.Error(w, errStr, http.StatusNotFound)
		return
	}

	// override entry in DB
	err = setAppInfo(&appInfo)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	jsonResponse := convertApplicationInfoToJson(&appInfo)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func applicationsAppInstanceIdGET(w http.ResponseWriter, r *http.Request) {
	log.Info("applicationsByIdGET")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	appInstanceId := vars["appInstanceId"]

	// Get App info for requested App instance ID
	var appData ApplicationData
	var filterParams FilterParameters
	filterParams.appId = appInstanceId
	appData.FilterParameters = &filterParams
	key := appCtrl.baseKey + ":mep:*:app:" + appInstanceId + ":info"
	err := appCtrl.rc.ForEachEntry(key, populateAppInfoList, &appData)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if App instance was found
	if len(appData.AppInfoList) != 1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Send response
	jsonResponse := convertApplicationInfoToJson(&appData.AppInfoList[0])
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func applicationsAppInstanceIdDELETE(w http.ResponseWriter, r *http.Request) {
	log.Info("applicationsByIdDELETE")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	appInstanceId := vars["appInstanceId"]

	// Get App info for requested App instance ID
	var appData ApplicationData
	var filterParams FilterParameters
	filterParams.appId = appInstanceId
	appData.FilterParameters = &filterParams
	key := appCtrl.baseKey + ":mep:*:app:" + appInstanceId + ":info"
	err := appCtrl.rc.ForEachEntry(key, populateAppInfoList, &appData)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if App instance was found
	if len(appData.AppInfoList) != 1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	appInfo := appData.AppInfoList[0]

	// Inform MEP instance to terminate App instance
	msg := appCtrl.mqLocal.CreateMsg(mq.MsgAppTerminate, appEnablementModule, appCtrl.sandboxName)
	msg.Payload[mqFieldAppInstanceId] = appInfo.Id
	msg.Payload[mqFieldMepName] = appInfo.MepName
	log.Debug("TX MSG: ", mq.PrintMsg(msg))
	err = appCtrl.mqLocal.SendMsg(msg)
	if err != nil {
		log.Error("Failed to send message. Error: ", err.Error())

		// TODO -- [Graceful Terminate Failure] Update App instance Service availability + Flush App Instance data

		// Flush App instance data
		key = appCtrl.baseKey + ":mep:" + appInfo.MepName + ":app:" + appInfo.Id
		_ = appCtrl.rc.DBFlush(key)
	}

	// Send response
	w.WriteHeader(http.StatusNoContent)
}

func applicationsGET(w http.ResponseWriter, r *http.Request) {
	log.Info("applicationsGET")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Validate & retrieve query parameters
	u, _ := url.Parse(r.URL.String())
	q := u.Query()
	validParams := []string{"app", "type", "state", "mep"}
	err := validateQueryParams(q, validParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	appName := q.Get("app")
	appType := q.Get("type")
	appState := q.Get("state")
	mepName := q.Get("mep")

	// Get filter App info list
	var appData ApplicationData
	var filterParameters FilterParameters
	filterParameters.appName = appName
	filterParameters.appType = appType
	filterParameters.state = appState
	filterParameters.mepName = mepName
	appData.FilterParameters = &filterParameters

	key := appCtrl.baseKey + ":mep:*:app:*:info"
	err = appCtrl.rc.ForEachEntry(key, populateAppInfoList, &appData)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	jsonResponse, err := json.Marshal(appData.AppInfoList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

// Get a new unique app instance
func getNewInstanceId() (string, error) {
	//allow 3 tries, if not return an error
	maxNbRetries := 3
	for try := maxNbRetries; try > 0; try-- {
		appInstanceId := uuid.New().String()
		if !appInstanceExists("", appInstanceId) {
			return appInstanceId, nil
		}
		try--
	}
	return "", errors.New("Can't allocate a unique instance Id")
}

// Validate that App Instance exists
func appInstanceExists(mepName string, appInstanceId string) bool {
	if mepName == "" {
		// Get list of running App instances
		var appData ApplicationData
		key := appCtrl.baseKey + ":mep:*:app:" + appInstanceId + ":info"
		err := appCtrl.rc.ForEachEntry(key, populateAppInfoList, &appData)
		if err != nil {
			log.Error("Failed to retrieve App information")
			return false
		}

		// Check if App instance exists
		for _, appInfo := range appData.AppInfoList {
			if appInfo.Id == appInstanceId {
				return true
			}
		}
	} else {
		// Find
		key := appCtrl.baseKey + ":mep:" + mepName + ":app:" + appInstanceId + ":info"
		return appCtrl.rc.EntryExists(key)
	}
	return false
}

// Validate query params
func validateQueryParams(params url.Values, validParams []string) error {
	for param := range params {
		found := false
		for _, validParam := range validParams {
			if param == validParam {
				found = true
				break
			}
		}
		if !found {
			err := errors.New("Invalid query param: " + param)
			log.Error(err.Error())
			return err
		}
	}
	return nil
}

// Validate App Info params
func validateAppInfo(appInfo *dataModel.ApplicationInfo, appInstanceId string) error {
	// Validate content
	if appInstanceId != "" && appInfo.Id != appInstanceId {
		return errors.New("Mandatory Application Instance Id parameter and body content not matching")
	}
	if appInfo.Name == "" {
		return errors.New("Mandatory Name not present")
	}
	if appInfo.MepName == "" {
		return errors.New("Mandatory MEP Name not present")
	}
	if appInfo.State == nil {
		return errors.New("Mandatory State not present")
	}
	switch *appInfo.State {
	case dataModel.READY_ApplicationState, dataModel.INITIALIZED_ApplicationState:
	default:
		return errors.New("Mandatory State value not valid")
	}

	// Initialize default App Type if not provided
	if appInfo.Type_ == nil || *appInfo.Type_ == "" {
		appType := dataModel.USER_ApplicationType
		appInfo.Type_ = &appType
	}
	return nil
}

// Set Application Information in DB
func setAppInfo(appInfo *dataModel.ApplicationInfo) error {
	fields := make(map[string]interface{})
	fields[fieldAppInstanceId] = appInfo.Id
	fields[fieldAppName] = appInfo.Name
	fields[fieldAppType] = string(*appInfo.Type_)
	fields[fieldMepName] = appInfo.MepName
	fields[fieldState] = string(*appInfo.State)
	fields[fieldVersion] = appInfo.Version

	key := appCtrl.baseKey + ":mep:" + appInfo.MepName + ":app:" + appInfo.Id + ":info"
	err := appCtrl.rc.SetEntry(key, fields)
	if err != nil {
		return err
	}
	return nil
}

func populateAppInfoList(key string, fields map[string]string, appData interface{}) error {
	// Get query params & userlist from user data
	data := appData.(*ApplicationData)
	if data == nil {
		return errors.New("AppInfoList not found in user data")
	}

	// Retrieve user info from DB
	var appInfo dataModel.ApplicationInfo
	appInfo.Id = fields[fieldAppInstanceId]
	appInfo.Name = fields[fieldAppName]
	appInfo.MepName = fields[fieldMepName]
	appInfo.Version = fields[fieldVersion]
	appType := dataModel.ApplicationType(fields[fieldAppType])
	appInfo.Type_ = &appType
	appState := dataModel.ApplicationState(fields[fieldState])
	appInfo.State = &appState

	// Filter Apps
	if data.FilterParameters != nil {
		// App Instance ID
		if data.FilterParameters.appId != "" && data.FilterParameters.appId != appInfo.Id {
			return nil
		}
		// App Name
		if data.FilterParameters.appName != "" && data.FilterParameters.appName != appInfo.Name {
			return nil
		}
		// App Type
		if data.FilterParameters.appType != "" && data.FilterParameters.appType != string(*appInfo.Type_) {
			return nil
		}
		// App state
		if data.FilterParameters.state != "" && data.FilterParameters.state != string(*appInfo.State) {
			return nil
		}
		// MEP name
		if data.FilterParameters.mepName != "" && data.FilterParameters.mepName != appInfo.MepName {
			return nil
		}
	}

	// Add App info to list
	data.AppInfoList = append(data.AppInfoList, appInfo)
	return nil
}

func convertApplicationInfoToJson(appInfo *dataModel.ApplicationInfo) string {
	jsonInfo, err := json.Marshal(*appInfo)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(jsonInfo)
}
