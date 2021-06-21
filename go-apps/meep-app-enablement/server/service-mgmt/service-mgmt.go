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
	"strconv"
	"strings"
	"sync"
	"time"

	//        sbi "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-wais/sbi"
	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	//httpLog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-http-logger"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"

	"github.com/gorilla/mux"
)

const msmgmtBasePath = "/mec_service_mgmt/v2/"
const msmgmtKey = "msmgmt"
const appEnablementKey = "app-enablement"

//const logModuleMSMgmt = "meep-app-enablement"
//const serviceName = "MEC Service Management"

var redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"

//var influxAddr string = "http://meep-influxdb.default.svc.cluster.local:8086"

var MSMGMT_DB = 5

var rc *redis.Connector
var hostUrl *url.URL
var sandboxName string
var selfName string
var basePath string
var baseKey string
var appEnablementBaseKey string

var mutex sync.Mutex

var expiryTicker *time.Ticker

var nextSubscriptionIdAvailable int
var nextServiceRegistrationIdAvailable int

type ServiceInfoList struct {
	ServiceInfos     []ServiceInfo
	filterParameters *FilterParameters
}

type FilterParameters struct {
	serInstanceId     []string
	serName           []string
	serCategoryId     string
	consumedLocalOnly string
	isLocal           string
	scopeOfLocality   string
}

/*func notImplemented(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusNotImplemented)
}
*/
func Init() (err error) {
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
	/*	if selfName == "" {
			err = errors.New("MEEP_SELF_NAME env variable not set")
			log.Error(err.Error())
			return err
		}
	*/
	//HARDCODE SOMETHING
	selfName = "mep1"
	log.Info("MEEP_SELF_NAME: ", selfName)

	// Set base path
	basePath = "/" + sandboxName + msmgmtBasePath
	// Get base store key
	baseKey = dkm.GetKeyRoot(sandboxName) + msmgmtKey
	appEnablementBaseKey = dkm.GetKeyRoot(sandboxName) + selfName + ":" + appEnablementKey

	// Connect to Redis DB
	rc, err = redis.NewConnector(redisAddr, MSMGMT_DB)
	if err != nil {
		log.Error("Failed connection to Redis DB. Error: ", err)
		return err
	}

	_ = rc.DBFlush(baseKey)
	_ = rc.DBFlush(appEnablementBaseKey)

	log.Info("Connected to Redis DB")

	reInit()

	expiryTicker = time.NewTicker(time.Second)
	go func() {
		for range expiryTicker.C {
			//checkForExpiredSubscriptions()
		}
	}()
	/*
	   // Initialize SBI
	   sbiCfg := sbi.SbiCfg{
	           SandboxName:    sandboxName,
	           RedisAddr:      redisAddr,
	           InfluxAddr:     influxAddr,
	           StaInfoCb:      updateStaInfo,
	           ApInfoCb:       updateApInfo,
	           ScenarioNameCb: updateStoreName,
	           CleanUpCb:      cleanUp,
	   }

	   err = sbi.Init(sbiCfg)
	   if err != nil {
	           log.Error("Failed initialize SBI. Error: ", err)
	           return err
	   }
	   log.Info("SBI Initialized")
	*/
	return nil
}

// reInit - finds the value already in the DB to repopulate local stored info
func reInit() {
	//next available subsId will be overrriden if subscriptions already existed
	nextSubscriptionIdAvailable = 1
	nextServiceRegistrationIdAvailable = 1
}

// Run - Start WAIS
func Run() (err error) {
	return nil //sbi.Run()
}

// Stop - Stop WAIS
func Stop() (err error) {
	return nil //sbi.Stop()
}

func appServicesGET(w http.ResponseWriter, r *http.Request) {
	log.Info("appServicesGET")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)

	appServicesGETByAppInstanceId(w, r, vars["appInstanceId"])
}

func appServicesGETByAppInstanceId(w http.ResponseWriter, r *http.Request, appInstanceId string) {

	u, _ := url.Parse(r.URL.String())
	q := u.Query()
	serInstanceId := q["ser_instance_id"]
	serName := q["ser_name"]
	serCategoryId := q.Get("ser_category_id")
	consumedLocalOnly := q.Get("consumed_local_only")
	isLocal := q.Get("is_local")
	scopeOfLocality := q.Get("scope_of_locality")

	validQueryParams := []string{"ser_instance_id", "ser_name", "ser_category_id", "consumed_local_only", "is_local", "scope_of_locality"}

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

	//if specific appInstanceId is queried ("*" is wildcard)
	if appInstanceId != "*" {
		err := validateAppInstanceId(appInstanceId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	var sInfoList ServiceInfoList
	var filterParameters FilterParameters
	filterParameters.serInstanceId = serInstanceId
	filterParameters.serName = serName
	filterParameters.serCategoryId = serCategoryId
	filterParameters.consumedLocalOnly = consumedLocalOnly
	filterParameters.isLocal = isLocal
	filterParameters.scopeOfLocality = scopeOfLocality
	sInfoList.filterParameters = &filterParameters

	keyName := appEnablementBaseKey + ":apps:" + appInstanceId + ":svcs:*"
	err := rc.ForEachJSONEntry(keyName, populateServiceInfoList, &sInfoList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(sInfoList.ServiceInfos)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func validateAppInstanceId(appInstanceId string) error {
	//check if entry exist for the application in the DB
	key := appEnablementBaseKey + ":apps:" + appInstanceId
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

func populateServiceInfoList(key string, jsonInfo string, sInfoList interface{}) error {
	// Get query params & userlist from user data
	data := sInfoList.(*ServiceInfoList)

	if data == nil {
		return errors.New("ServiceInfos not found in serviceInfoList")
	}

	// Retrieve user info from DB
	var sInfo ServiceInfo
	err := json.Unmarshal([]byte(jsonInfo), &sInfo)
	if err != nil {
		return err
	}

	match := true

	if data.filterParameters != nil {
		//compare with filter to return the service info or not
		match = false
		if len(data.filterParameters.serInstanceId) > 0 {
			for _, value := range data.filterParameters.serInstanceId {
				if sInfo.SerInstanceId == value {
					match = true
					break
				}
			}
		} else {
			match = true
		}

		if match {
			match = false
			if len(data.filterParameters.serName) > 0 {
				for _, value := range data.filterParameters.serName {
					if sInfo.SerName == value {
						match = true
						break
					}
				}
			} else {
				match = true
			}
		}

		if match {
			if data.filterParameters.serCategoryId != "" {
				//comparing with either the category name or id, spec is not clear
				match = (data.filterParameters.serCategoryId == sInfo.SerCategory.Name) || (data.filterParameters.serCategoryId == sInfo.SerCategory.Id)
			}
		}

		if match {
			if data.filterParameters.consumedLocalOnly != "" {
				consumedLocalOnlyStr := strconv.FormatBool(sInfo.ConsumedLocalOnly)
				match = (data.filterParameters.consumedLocalOnly == consumedLocalOnlyStr)
			}
		}

		if match {
			if data.filterParameters.isLocal != "" {
				isLocalStr := strconv.FormatBool(sInfo.IsLocal)
				match = (data.filterParameters.isLocal == isLocalStr)
			}
		}

		if match {
			if data.filterParameters.scopeOfLocality != "" {
				match = (data.filterParameters.scopeOfLocality == string(*sInfo.ScopeOfLocality))
			}
		}
	}

	if match {
		data.ServiceInfos = append(data.ServiceInfos, sInfo)
	}
	return nil
}

func appServicesPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Info("appServicesPOST")

	//var response InlineCircleNotificationSubscription
	vars := mux.Vars(r)
	appInstanceId := vars["appInstanceId"]

	err := validateAppInstanceId(appInstanceId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var sInfoPost ServiceInfoPost
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&sInfoPost)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//checking for mandatory properties
	if sInfoPost.SerName == "" {
		log.Error("Mandatory Service Name parameter not present")
		http.Error(w, "Mandatory Service Name parameter not present", http.StatusBadRequest)
		return
	}
	if sInfoPost.Version == "" {
		log.Error("Mandatory Service Version parameter not present")
		http.Error(w, "Mandatory Service Version parameter not present", http.StatusBadRequest)
		return
	}
	if sInfoPost.State == nil {
		log.Error("Mandatory Service State parameter not present")
		http.Error(w, "Mandatory Service State parameter not present", http.StatusBadRequest)
		return
	}
	if sInfoPost.Serializer == nil {
		log.Error("Mandatory Serializer parameter not present")
		http.Error(w, "Mandatory Serializer parameter not present", http.StatusBadRequest)
		return
	}

	sInfo, err := registerService(appInstanceId, &sInfoPost)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	jsonResponse, err := json.Marshal(sInfo)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func deregisterService(appInstanceId string, serviceId string) error {

	return nil
}

func registerService(appInstanceId string, sInfoPost *ServiceInfoPost) (*ServiceInfo, error) {

	mutex.Lock()
	defer mutex.Unlock()

	newServiceId := nextServiceRegistrationIdAvailable
	nextServiceRegistrationIdAvailable++
	serviceId := strconv.Itoa(newServiceId)

	sInfo, err := updateService(appInstanceId, sInfoPost, serviceId, false)
	if err != nil {
		nextServiceRegistrationIdAvailable--
		return nil, err
	}

	return sInfo, nil
}

func updateService(appInstanceId string, sInfoPost *ServiceInfoPost, serviceId string, needMutex bool) (*ServiceInfo, error) {

	if sInfoPost == nil {
		return nil, errors.New("Service Info is null")
	}

	if needMutex {
		mutex.Lock()
		defer mutex.Unlock()
	}

	sInfo, err := createSInfoFromSInfoPost(sInfoPost, serviceId)
	if err == nil {
		err = rc.JSONSetEntry(appEnablementBaseKey+":apps:"+appInstanceId+":svcs:"+serviceId, ".", ConvertServiceInfoToJson(sInfo))
	}

	return sInfo, err
}

func createSInfoFromSInfoPost(sInfoPost *ServiceInfoPost, serviceId string) (*ServiceInfo, error) {
	var sInfo ServiceInfo
	if sInfoPost.SerInstanceId == "" {
		sInfo.SerInstanceId = serviceId
	} else if sInfoPost.SerInstanceId != serviceId {
		return nil, errors.New("ServiceId not matching serviceId provided in the body of the request")
	} else {
		sInfo.SerInstanceId = sInfoPost.SerInstanceId
	}

	sInfo.SerName = sInfoPost.SerName
	sInfo.SerCategory = sInfoPost.SerCategory
	sInfo.Version = sInfoPost.Version
	sInfo.State = sInfoPost.State
	sInfo.TransportInfo = sInfoPost.TransportInfo
	sInfo.Serializer = sInfoPost.Serializer
	sInfo.ScopeOfLocality = sInfoPost.ScopeOfLocality
	sInfo.ConsumedLocalOnly = sInfoPost.ConsumedLocalOnly
	sInfo.IsLocal = sInfoPost.IsLocal
	return &sInfo, nil
}

func appServicesByIdDELETE(w http.ResponseWriter, r *http.Request) {
	log.Info("appServicesByIdDELETE")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	appInstanceId := vars["appInstanceId"]
	serviceId := vars["serviceId"]

	err := validateAppInstanceId(appInstanceId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	present, _ := rc.JSONGetEntry(appEnablementBaseKey+":apps:"+appInstanceId+":svcs:"+serviceId, ".")
	if present == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = rc.JSONDelEntry(appEnablementBaseKey+":apps:"+appInstanceId+":svcs:"+serviceId, ".")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = deregisterService(appInstanceId, serviceId)

	w.WriteHeader(http.StatusNoContent)

}

func appServicesByIdGET(w http.ResponseWriter, r *http.Request) {
	log.Info("appServicesByIdGET")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	serviceId := vars["serviceId"]
	appInstanceId := vars["appInstanceId"]

	err := validateAppInstanceId(appInstanceId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse, _ := rc.JSONGetEntry(appEnablementBaseKey+":apps:"+appInstanceId+":svcs:"+serviceId, ".")
	if jsonResponse == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func appServicesByIdPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Info("appServicesPUT")

	//var response InlineCircleNotificationSubscription
	vars := mux.Vars(r)
	appInstanceId := vars["appInstanceId"]
	serviceId := vars["serviceId"]

	err := validateAppInstanceId(appInstanceId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var sInfoPost ServiceInfoPost
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&sInfoPost)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//checking for mandatory properties
	if sInfoPost.SerInstanceId == "" {
		log.Error("Mandatory Service Instance Id parameter not present")
		http.Error(w, "Mandatory Service Instance Id parameter not present", http.StatusBadRequest)
		return
	}
	if sInfoPost.SerInstanceId != serviceId {
		log.Error("Mandatory Service Instance Id parameter and body content not matching")
		http.Error(w, "Mandatory Service Instance Id parameter and body content not matching", http.StatusBadRequest)
		return
	}
	if sInfoPost.SerName == "" {
		log.Error("Mandatory Service Name parameter not present")
		http.Error(w, "Mandatory Service Name parameter not present", http.StatusBadRequest)
		return
	}
	if sInfoPost.Version == "" {
		log.Error("Mandatory Service Version parameter not present")
		http.Error(w, "Mandatory Service Version parameter not present", http.StatusBadRequest)
		return
	}
	if sInfoPost.State == nil {
		log.Error("Mandatory Service State parameter not present")
		http.Error(w, "Mandatory Service State parameter not present", http.StatusBadRequest)
		return
	}
	if sInfoPost.Serializer == nil {
		log.Error("Mandatory Serializer parameter not present")
		http.Error(w, "Mandatory Serializer parameter not present", http.StatusBadRequest)
		return
	}

	sInfo, err := updateService(appInstanceId, &sInfoPost, serviceId, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse, err := json.Marshal(sInfo)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))

}

func servicesByIdGET(w http.ResponseWriter, r *http.Request) {
	log.Info("servicesByIdGET")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	serviceId := vars["serviceId"]

	var sInfoList ServiceInfoList

	keyName := appEnablementBaseKey + ":apps:*:svcs:" + serviceId
	err := rc.ForEachJSONEntry(keyName, populateServiceInfoList, &sInfoList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//can only be 1 if successful, so reject 0 or more than 1 result
	if len(sInfoList.ServiceInfos) != 1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonResponse, err := json.Marshal(sInfoList.ServiceInfos[0])
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func servicesGET(w http.ResponseWriter, r *http.Request) {
	log.Info("servicesGET")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	//use appInstanceId wildcard
	appServicesGETByAppInstanceId(w, r, "*")
}
