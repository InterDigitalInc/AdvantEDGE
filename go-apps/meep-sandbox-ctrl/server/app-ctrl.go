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
	"strconv"

	apps "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-applications"
	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	uuid "github.com/google/uuid"
	"github.com/gorilla/mux"
)

// MQ payload fields
const mqFieldAppInstanceId = "id"
const mqFieldPersist = "persist"

type AppCtrl struct {
	sandboxName string
	appStore    *apps.ApplicationStore
	mqLocal     *mq.MsgQueue
}

// App Controller
var appCtrl *AppCtrl

// Initialize App Controller
func appCtrlInit(sandboxName string, mqLocal *mq.MsgQueue) (err error) {

	// Create new App Controller
	appCtrl = new(AppCtrl)
	appCtrl.sandboxName = sandboxName
	appCtrl.mqLocal = mqLocal

	// Create Application Store
	cfg := &apps.ApplicationStoreCfg{
		Name:      moduleName,
		Namespace: sandboxName,
		UpdateCb:  appStoreUpdateCb,
		RedisAddr: redisDBAddr,
	}
	appCtrl.appStore, err = apps.NewApplicationStore(cfg)
	if err != nil {
		log.Error("Failed to connect to Application Store. Error: ", err)
		return err
	}
	log.Info("Connected to Application Store")

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

// Flush App instances
func appCtrlFlushAppInstances() (err error) {
	appCtrl.appStore.FlushNonPersistent()
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
	} else {
		// Make sure App instance does not exist
		if appInstanceExists(appInfo.Id) {
			errStr := "AppInstanceId already exists"
			log.Error(errStr)
			http.Error(w, errStr, http.StatusBadRequest)
			return
		}
	}

	// Create new App instance
	err = appCtrl.appStore.Set(convertApplicationInfoToApp(&appInfo))
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
	if !appInstanceExists(appInfo.Id) {
		errStr := "AppInstanceId does not exist"
		log.Error(errStr)
		http.Error(w, errStr, http.StatusNotFound)
		return
	}

	// Override entry in DB
	err = appCtrl.appStore.Set(convertApplicationInfoToApp(&appInfo))
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
	log.Info("applicationsAppInstanceIdGET")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	appInstanceId := vars["appInstanceId"]

	// Get App info for requested App instance ID
	app, err := appCtrl.appStore.Get(appInstanceId)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	appInfo := convertAppToApplicationInfo(app)

	// Send response
	jsonResponse := convertApplicationInfoToJson(appInfo)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func applicationsAppInstanceIdDELETE(w http.ResponseWriter, r *http.Request) {
	log.Info("applicationsAppInstanceIdDELETE")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	appId := vars["appInstanceId"]

	// Flush App instance data
	err := appCtrl.appStore.Del(appId)
	if err != nil {
		log.Error(err.Error())
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
	validParams := []string{"app", "mep", "type"}
	err := validateQueryParams(q, validParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	appName := q.Get("app")
	mepName := q.Get("mep")
	appType := q.Get("type")

	// Get application list
	appList, err := appCtrl.appStore.GetAll()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare AppInfo list
	appInfoList := make([]dataModel.ApplicationInfo, 0)
	for _, app := range appList {
		// Filter using query params
		if (appName != "" && app.Name != appName) ||
			(mepName != "" && app.Mep != mepName) ||
			(appType != "" && app.Type != appType) {
			continue
		}
		// Append appInfo
		appInfo := convertAppToApplicationInfo(app)
		appInfoList = append(appInfoList, *appInfo)
	}

	// Send response
	jsonResponse, err := json.Marshal(appInfoList)
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
		if !appInstanceExists(appInstanceId) {
			return appInstanceId, nil
		}
		try--
	}
	return "", errors.New("Can't allocate a unique instance Id")
}

// Validate that App Instance exists
func appInstanceExists(appInstanceId string) bool {
	_, err := appCtrl.appStore.Get(appInstanceId)
	return err == nil
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

	// Set default App type if missing
	if appInfo.Type_ == "" {
		appInfo.Type_ = apps.TypeUser
	}
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

func convertAppToApplicationInfo(app *apps.Application) *dataModel.ApplicationInfo {
	appInfo := &dataModel.ApplicationInfo{
		Id:      app.Id,
		Name:    app.Name,
		MepName: app.Mep,
		Type_:   app.Type,
		Persist: app.Persist,
	}
	return appInfo
}

func convertApplicationInfoToApp(appInfo *dataModel.ApplicationInfo) *apps.Application {
	app := &apps.Application{
		Id:      appInfo.Id,
		Name:    appInfo.Name,
		Mep:     appInfo.MepName,
		Type:    appInfo.Type_,
		Persist: appInfo.Persist,
	}
	return app
}

func appStoreUpdateCb(eventType string, eventData interface{}) {
	var msg *mq.Msg

	// Create message to send on MQ
	switch eventType {
	case apps.EventAdd:
		msg = appCtrl.mqLocal.CreateMsg(mq.MsgAppUpdate, mq.TargetAll, appCtrl.sandboxName)
		msg.Payload[mqFieldAppInstanceId] = eventData.(string)
	case apps.EventRemove:
		msg = appCtrl.mqLocal.CreateMsg(mq.MsgAppRemove, mq.TargetAll, appCtrl.sandboxName)
		msg.Payload[mqFieldAppInstanceId] = eventData.(string)
	case apps.EventFlush:
		msg = appCtrl.mqLocal.CreateMsg(mq.MsgAppFlush, mq.TargetAll, appCtrl.sandboxName)
		msg.Payload[mqFieldPersist] = strconv.FormatBool(eventData.(bool))
	default:
		return
	}

	// Send message to inform other modules of app store changes
	log.Debug("TX MSG: ", mq.PrintMsg(msg))
	err := appCtrl.mqLocal.SendMsg(msg)
	if err != nil {
		log.Error("Failed to send message. Error: ", err.Error())
		return
	}
}
