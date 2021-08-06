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
const mqFieldApiList = "apilist"

type SwaggerApiList struct {
	Apis     []SwaggerApi
	UserApis []SwaggerApi
}

type SwaggerApi struct {
	Name string
	Url  string
}

type SwaggerApiMgr struct {
	moduleName  string
	sandboxName string
	svcUrl      string
	msgQueue    *mq.MsgQueue
	handlerId   int
}

// NewSwaggerApiMgr - Creates and initialize a Swagger API Manager instance
func NewSwaggerApiMgr(moduleName string, sandboxName string, mepName string, msgQueue *mq.MsgQueue) (sam *SwaggerApiMgr, err error) {
	// Create new Swagger API Manager instance
	sam = new(SwaggerApiMgr)
	sam.sandboxName = sandboxName
	sam.msgQueue = msgQueue
	sam.handlerId = -1

	// Get Host Url
	hostUrl, err := url.Parse(strings.TrimSpace(os.Getenv("MEEP_HOST_URL")))
	if err != nil {
		hostUrl = new(url.URL)
	}
	log.Info("MEEP_HOST_URL: ", hostUrl)

	// Get Service Path
	svcPath := strings.TrimSpace(os.Getenv("MEEP_SVC_PATH"))
	log.Info("MEEP_SVC_PATH: ", svcPath)

	// Module name
	sam.moduleName = moduleName
	if mepName != "" {
		sam.moduleName = mepName + "-" + moduleName
	}
	log.Info("Module name: ", sam.moduleName)

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
func (sam *SwaggerApiMgr) Start() error {
	var err error

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

// AddApis - Send message to inform listeners of added APIs
func (sam *SwaggerApiMgr) AddApis() error {
	// Get service APIs from filesystem
	apiFiles, err := ioutil.ReadDir(apiDir)
	if err != nil {
		log.Error("Failed to read API dir with error: ", err.Error())
		return err
	}
	userApiFiles, err := ioutil.ReadDir(userApiDir)
	if err != nil {
		log.Error("Failed to read API dir with error: ", err.Error())
		return err
	}

	// Populate API list
	var apiList SwaggerApiList
	for _, file := range apiFiles {
		var api SwaggerApi
		api.Name = file.Name()
		api.Url = sam.svcUrl + "/api/" + api.Name
		apiList.Apis = append(apiList.Apis, api)
	}
	for _, file := range userApiFiles {
		var api SwaggerApi
		api.Name = file.Name()
		api.Url = sam.svcUrl + "/user-api/" + api.Name
		apiList.Apis = append(apiList.Apis, api)
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
	msg.Payload[mqFieldApiList] = apiListStr
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
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		moduleName := msg.Payload[mqFieldModule]
		apiListStr := msg.Payload[mqFieldApiList]
		sam.processApiUpdate(moduleName, convertJsonToSwaggerApiList(apiListStr))
	default:
	}
}

func (sam *SwaggerApiMgr) processApiUpdate(moduleName string, apiList *SwaggerApiList) {

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
