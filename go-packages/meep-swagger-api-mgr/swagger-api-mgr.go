/*
 * Copyright (c) 2021  InterDigital Communications, Inc
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

package giscache

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
	"os"
	"sort"
	"strings"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
)

const apiDir = "/api"
const userApiDir = "/user-api"
const modulePlatformCtrl = "meep-platform-ctrl"
const moduleSandboxCtrl = "meep-sandbox-ctrl"

// MQ payload fields
const mqFieldModule = "module"
const mqFieldMepName = "mep"
const mqFieldApiList = "apilist"

type SwaggerApiList struct {
	isMep    bool
	Apis     []SwaggerApi
	UserApis []SwaggerApi
}

type SwaggerApi struct {
	Name string
	Url  string
}

type SwaggerApiMgr struct {
	moduleName   string
	sandboxName  string
	mepName      string
	isProvider   bool
	isAggregator bool
	svcUrl       string
	msgQueue     *mq.MsgQueue
	handlerId    int
	apiMap       map[string]*SwaggerApiList
}

// NewSwaggerApiMgr - Creates and initialize a Swagger API Manager instance
func NewSwaggerApiMgr(moduleName string, sandboxName string, mepName string, msgQueue *mq.MsgQueue) (sam *SwaggerApiMgr, err error) {
	// Create new Swagger API Manager instance
	sam = new(SwaggerApiMgr)
	sam.moduleName = moduleName
	sam.sandboxName = sandboxName
	sam.mepName = mepName
	sam.msgQueue = msgQueue
	sam.handlerId = -1
	sam.apiMap = make(map[string]*SwaggerApiList)

	// Get Host Url
	hostUrl, err := url.Parse(strings.TrimSpace(os.Getenv("MEEP_HOST_URL")))
	if err != nil {
		hostUrl = new(url.URL)
	}
	log.Info("MEEP_HOST_URL: ", hostUrl)

	// Get Service Path
	svcPath := strings.TrimSpace(os.Getenv("MEEP_SVC_PATH"))
	log.Info("MEEP_SVC_PATH: ", svcPath)

	// Create full service url
	sam.svcUrl = hostUrl.String()
	if sandboxName != "" {
		sam.svcUrl += "/" + sandboxName
	}
	if mepName != "" {
		sam.svcUrl += "/" + mepName
	}
	sam.svcUrl += svcPath
	log.Info("Service URL: ", sam.svcUrl)

	log.Info("Created Swagger API Manager")
	return sam, nil
}

// Start - Start Swagger API message Handler
func (sam *SwaggerApiMgr) Start(isProvider bool, isAggregator bool) error {
	var err error
	sam.isProvider = isProvider
	sam.isAggregator = isAggregator

	// Make sure handler is not running
	if sam.handlerId != -1 {
		return errors.New("Handler already running")
	}

	// Register Message Queue handler
	handler := mq.MsgHandler{Handler: sam.msgHandler, UserData: nil}
	sam.handlerId, err = sam.msgQueue.RegisterHandler(handler)
	if err != nil {
		log.Error("Failed to listen for sandbox updates: ", err.Error())
		return err
	}

	// If aggregator, send API update Request to providers
	if sam.isAggregator {
		_ = sam.sendApiRequest()
	}
	return nil
}

// Stop - Stop Swagger API message Handler
func (sam *SwaggerApiMgr) Stop() error {
	// Make sure handler is running
	if sam.handlerId == -1 {
		return errors.New("Handler not running")
	}

	// Unregister Message Queue handler
	sam.msgQueue.UnregisterHandler(sam.handlerId)
	sam.handlerId = -1
	return nil
}

// FlushMepApis - Flush all MEP APIs
func (sam *SwaggerApiMgr) FlushMepApis() error {

	// Remove MEP APIs from API map
	for key, apiList := range sam.apiMap {
		if apiList.isMep {
			delete(sam.apiMap, key)
		}
	}

	// Refresh API lists
	sam.refreshApiLists()

	return nil
}

// AddApis - Send message to inform listeners of added APIs
func (sam *SwaggerApiMgr) AddApis() error {
	// Get service APIs from filesystem
	apiFiles, err := ioutil.ReadDir(apiDir)
	if err != nil {
		log.Error("Failed to read API dir with error: ", err.Error())
		// Should never happen. For UT, return nil.
		return nil
	}
	userApiFiles, err := ioutil.ReadDir(userApiDir)
	if err != nil {
		log.Error("Failed to read API dir with error: ", err.Error())
		// Should never happen. For UT, return nil.
		return nil
	}

	// Get MEP prefix
	mepPrefix := ""
	if sam.mepName != "" {
		mepPrefix = sam.mepName + " - "
	}

	// Populate API lists
	var apiList SwaggerApiList
	for _, file := range apiFiles {
		var api SwaggerApi
		api.Name = mepPrefix + file.Name()
		api.Url = sam.svcUrl + "/api/" + file.Name()
		apiList.Apis = append(apiList.Apis, api)
	}
	for _, file := range userApiFiles {
		var api SwaggerApi
		api.Name = mepPrefix + file.Name()
		api.Url = sam.svcUrl + "/user-api/" + file.Name()
		apiList.UserApis = append(apiList.UserApis, api)
	}

	// Publish new API list
	return sam.publishApiUpdate(&apiList)
}

// RemoveApis - Send message to inform listeners of removed APIs
func (sam *SwaggerApiMgr) RemoveApis() error {
	// Publish empty API list
	return sam.publishApiUpdate(nil)
}

func (sam *SwaggerApiMgr) publishApiUpdate(apiList *SwaggerApiList) error {
	// Populate API list
	apiListStr := ""
	if apiList != nil {
		apiListStr = convertSwaggerApiListToJson(apiList)
	}

	// Send message to inform listeners of API update
	var msg *mq.Msg
	if sam.sandboxName == "" {
		msg = sam.msgQueue.CreateMsg(mq.MsgApiUpdate, modulePlatformCtrl, mq.TargetAll)
	} else {
		msg = sam.msgQueue.CreateMsg(mq.MsgApiUpdate, moduleSandboxCtrl, sam.sandboxName)
	}
	msg.Payload[mqFieldModule] = sam.moduleName
	msg.Payload[mqFieldMepName] = sam.mepName
	msg.Payload[mqFieldApiList] = apiListStr
	log.Debug("TX MSG: ", mq.PrintMsg(msg))
	err := sam.msgQueue.SendMsg(msg)
	if err != nil {
		log.Error("Failed to send message. Error: ", err.Error())
		return err
	}
	return nil
}

func (sam *SwaggerApiMgr) sendApiRequest() error {
	var msg *mq.Msg
	if sam.sandboxName == "" {
		msg = sam.msgQueue.CreateMsg(mq.MsgApiRequest, mq.TargetAll, mq.TargetAll)
	} else {
		msg = sam.msgQueue.CreateMsg(mq.MsgApiRequest, mq.TargetAll, sam.sandboxName)
	}
	log.Debug("TX MSG: ", mq.PrintMsg(msg))
	err := sam.msgQueue.SendMsg(msg)
	if err != nil {
		log.Error("Failed to send message. Error: ", err.Error())
		return err
	}
	return nil
}

// Message Queue handler
func (sam *SwaggerApiMgr) msgHandler(msg *mq.Msg, userData interface{}) {
	switch msg.Message {
	case mq.MsgApiUpdate:
		if sam.isAggregator {
			log.Debug("RX MSG: ", mq.PrintMsg(msg))
			moduleName := msg.Payload[mqFieldModule]
			mepName := msg.Payload[mqFieldMepName]
			apiListStr := msg.Payload[mqFieldApiList]
			apiList := convertJsonToSwaggerApiList(apiListStr)
			if apiList == nil {
				apiList = new(SwaggerApiList)
			}
			sam.processApiUpdate(moduleName, mepName, apiList)
		}
	case mq.MsgApiRequest:
		if sam.isProvider {
			log.Debug("RX MSG: ", mq.PrintMsg(msg))
			sam.processApiRequest()
		}
	default:
	}
}

func (sam *SwaggerApiMgr) processApiRequest() {
	// Retieve & send APIs
	_ = sam.AddApis()
}

func (sam *SwaggerApiMgr) processApiUpdate(moduleName string, mepName string, apiList *SwaggerApiList) {

	// Validate params
	if moduleName == "" {
		log.Error("Invalid module name")
		return
	}
	if apiList == nil {
		log.Error("apiList == nil")
		return
	}

	// Update module API list
	mepPrefix := ""
	apiList.isMep = false
	if mepName != "" {
		mepPrefix = mepName + "-"
		apiList.isMep = true
	}
	sam.apiMap[mepPrefix+moduleName] = apiList

	// Refresh API lists
	sam.refreshApiLists()
}

func (sam *SwaggerApiMgr) refreshApiLists() {

	// Get ordered api lists
	var orderedApiList []SwaggerApi
	var orderedUserApiList []SwaggerApi
	for _, apiList := range sam.apiMap {
		orderedApiList = append(orderedApiList, apiList.Apis...)
		orderedUserApiList = append(orderedUserApiList, apiList.UserApis...)
	}
	sortApiList(orderedApiList)
	sortApiList(orderedUserApiList)

	// Update index.html file API lists
	_ = updateApiListInFile("/swagger/index.html", orderedApiList)
	_ = updateApiListInFile("/user-swagger/index.html", orderedUserApiList)
}

func updateApiListInFile(filename string, apiList []SwaggerApi) error {
	matchPrefix := "        url"
	apiListNone := "        url: \"\","
	apiListPrefix := "        urls: ["
	apiListPostfix := " ],"

	// Generate API list string to replace in file
	apiListStr := apiListNone
	if len(apiList) > 0 {
		apiListStr = apiListPrefix
		for _, api := range apiList {
			// Replace http/https prefix from URL with browser location protocol
			url := strings.TrimPrefix(api.Url, "https:")
			url = strings.TrimPrefix(url, "http:")
			apiListStr += "{\"name\": \"" + api.Name + "\", \"url\": location.protocol + \"" + url + "\"},"
		}
		apiListStr += apiListPostfix
	}
	log.Debug("apiListStr: ", apiListStr)

	// Read input file
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error("Failed to read file with err: ", err.Error())
		return err
	}

	// Replace API list line
	lines := strings.Split(string(input), "\n")
	apiListUpdated := false
	for i, line := range lines {
		if strings.HasPrefix(line, matchPrefix) {
			lines[i] = apiListStr
			apiListUpdated = true
			break
		}
	}

	// Write file if updated
	if apiListUpdated {
		output := strings.Join(lines, "\n")
		err = ioutil.WriteFile(filename, []byte(output), 0644)
		if err != nil {
			log.Error("Failed to write file with err: ", err.Error())
			return err
		}
	}
	return nil
}

func convertSwaggerApiListToJson(apiList *SwaggerApiList) string {
	jsonInfo, err := json.Marshal(*apiList)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(jsonInfo)
}

func convertJsonToSwaggerApiList(jsonInfo string) *SwaggerApiList {
	var obj SwaggerApiList
	err := json.Unmarshal([]byte(jsonInfo), &obj)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	return &obj
}

func sortApiList(apiList []SwaggerApi) {
	sort.Slice(apiList, func(i, j int) bool {
		return apiList[i].Name < apiList[j].Name
	})
}
