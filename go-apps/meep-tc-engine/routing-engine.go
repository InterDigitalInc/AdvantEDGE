/*
 * Copyright (c) 2019  InterDigital Communications, Inc
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

package main

import (
	"encoding/json"
	"strconv"

	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mgModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mg-manager-model"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
)

const typeLb string = "lb"
const typeMeSvc string = "ME-SVC"
const typeIngressSvc string = "INGRESS-SVC"
const typeEgressSvc string = "EGRESS-SVC"

const DEFAULT_LB_RULES_DB = 0

// LbRulesStore -
type LbRulesStore struct {
	baseKey string
	rc      *redis.Connector
}

type RoutingEngine struct {
	name         string
	sandboxName  string
	lbRulesStore *LbRulesStore
}

func NewRoutingEngine(name string, sandboxName string) (re *RoutingEngine, err error) {
	// Create new Routing Engine instance
	re = new(RoutingEngine)
	re.name = name
	re.sandboxName = sandboxName

	// Open Load Balancing Rules Store
	re.lbRulesStore = new(LbRulesStore)
	re.lbRulesStore.baseKey = dkm.GetKeyRoot(tce.sandboxName) + mgManagerKey
	re.lbRulesStore.rc, err = redis.NewConnector(redisAddr, DEFAULT_LB_RULES_DB)
	if err != nil {
		log.Error("Failed connection to LB Rules Store Redis DB. Error: ", err)
		return nil, err
	}
	log.Info("Connected to LB Rules Store redis DB")

	log.Info("Successfully create Routing Engine")
	return re, nil
}

// RefreshLbRules - Fetch & apply latest MG Manager LB rules
func (re *RoutingEngine) RefreshLbRules() {

	// Retrieve LB rules from DB
	jsonNetElemList, err := re.lbRulesStore.rc.JSONGetEntry(re.lbRulesStore.baseKey+typeLb, ".")
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Unmarshal MG Service Maps
	var netElemList mgModel.NetworkElementList
	err = json.Unmarshal([]byte(jsonNetElemList), &netElemList)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Update pod MG service mappings
	for _, netElem := range netElemList.NetworkElements {
		podInfo := podInfoMap[netElem.Name]
		if podInfo == nil {
			log.Error("Failed to find network element: ", netElem.Name)
			continue
		}

		// Set load balanced MG Service instance
		for _, svcMap := range netElem.ServiceMaps {
			if svcInfo, found := svcInfoMap[svcMap.LbSvcName]; found {
				podInfo.MgSvcMap[svcMap.MgSvcName] = svcInfo
			} else {
				log.Error("failed to find service instance: ", svcMap.LbSvcName)
			}
		}
	}

	// Apply new MG Service mapping rules
	re.applyLbRules()

	// Inform sidecars of LB rule updates
	re.publishLbRulesUpdate()
}

// publishLbRulesUpdate - Inform sidecars of LB rules update
func (re *RoutingEngine) publishLbRulesUpdate() {

	// Send TC LB Rules update message to TC Sidecars for enforcement
	msg := tce.mqLocal.CreateMsg(mq.MsgTcLbRulesUpdate, moduleTcSidecar, tce.sandboxName)
	log.Debug("TX MSG: ", mq.PrintMsg(msg))
	err := tce.mqLocal.SendMsg(msg)
	if err != nil {
		log.Error("Failed to send message. Error: ", err.Error())
	}
}

// Generate & store rules based on mapping
func (re *RoutingEngine) applyLbRules() {
	log.Debug("applyLbRules")

	keys := map[string]bool{}

	// For each pod, add MG, ingress & egress Service LB rules
	for _, podInfo := range podInfoMap {
		// MG Service LB rules
		for _, svcInfo := range podInfo.MgSvcMap {
			// Add one rule per port
			for _, portInfo := range svcInfo.Ports {
				// Populate rule fields
				fields := make(map[string]interface{})
				fields[fieldSvcType] = typeMeSvc
				fields[fieldSvcName] = svcInfo.MgSvc.Name
				fields[fieldSvcIp] = tce.ipManager.GetSvcIp(svcInfo.MgSvc.Name)
				fields[fieldSvcProtocol] = portInfo.Protocol
				fields[fieldSvcPort] = portInfo.Port
				fields[fieldLbSvcName] = svcInfo.Name
				fields[fieldLbSvcIp] = tce.ipManager.GetSvcIp(svcInfo.Name)
				fields[fieldLbSvcPort] = portInfo.Port

				// Make unique key
				key := tce.netCharStore.baseKey + typeLb + ":" + podInfo.Name + ":" +
					svcInfo.MgSvc.Name + ":" + strconv.Itoa(int(portInfo.Port))
				keys[key] = true

				// Set rule information in DB
				_ = tce.netCharStore.rc.SetEntry(key, fields)
			}
		}

		// Ingress Service rules
		for _, svcMap := range podInfo.IngressSvcMapList {
			// Get Service info from exposed service name
			// Check if MG Service first
			var svcInfo *ServiceInfo
			var found bool
			if svcInfo, found = podInfo.MgSvcMap[svcMap.SvcName]; !found {
				// If not found, must be unique service
				if svcInfo, found = svcInfoMap[svcMap.SvcName]; !found {
					log.Warn("Failed to find service instance: ", svcMap.SvcName)
					continue
				}
			}

			// Populate rule fields
			fields := make(map[string]interface{})
			fields[fieldSvcType] = typeIngressSvc
			fields[fieldSvcName] = svcMap.SvcName
			fields[fieldSvcIp] = "0.0.0.0/0"
			fields[fieldSvcProtocol] = svcMap.Protocol
			fields[fieldSvcPort] = svcMap.NodePort
			fields[fieldLbSvcName] = svcInfo.Name
			fields[fieldLbSvcIp] = tce.ipManager.GetSvcIp(svcInfo.Name)
			fields[fieldLbSvcPort] = svcMap.SvcPort

			// Make unique key
			key := tce.netCharStore.baseKey + typeLb + ":" + podInfo.Name + ":" +
				svcMap.SvcName + ":" + strconv.Itoa(int(svcMap.NodePort))
			keys[key] = true

			// Set rule information in DB
			_ = tce.netCharStore.rc.SetEntry(key, fields)
		}

		// Egress Service rules
		for _, svcMap := range podInfo.EgressSvcMapList {
			// Populate rule fields
			fields := make(map[string]interface{})
			fields[fieldSvcType] = typeEgressSvc
			fields[fieldSvcName] = svcMap.SvcName
			fields[fieldSvcIp] = "0.0.0.0/0"
			fields[fieldSvcProtocol] = svcMap.Protocol
			fields[fieldSvcPort] = svcMap.SvcPort
			fields[fieldLbSvcName] = svcMap.SvcName
			fields[fieldLbSvcIp] = svcMap.SvcIp
			fields[fieldLbSvcPort] = svcMap.SvcPort

			// Make unique key
			key := tce.netCharStore.baseKey + typeLb + ":" + podInfo.Name + ":" +
				svcMap.SvcName + ":" + strconv.Itoa(int(svcMap.SvcPort))
			keys[key] = true

			// Set rule information in DB
			_ = tce.netCharStore.rc.SetEntry(key, fields)
		}
	}

	// Remove stale DB entries
	keyName := tce.netCharStore.baseKey + typeLb + ":*"
	err := tce.netCharStore.rc.ForEachEntry(keyName, removeLbEntryHandler, &keys)
	if err != nil {
		log.Error("Failed to remove old entries with err: ", err)
		return
	}
}

func removeLbEntryHandler(key string, fields map[string]string, userData interface{}) error {
	keys := userData.(*map[string]bool)

	if _, found := (*keys)[key]; !found {
		_ = tce.netCharStore.rc.DelEntry(key)
	}
	return nil
}
