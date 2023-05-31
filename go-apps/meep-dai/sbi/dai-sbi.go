/*
 * Copyright (c) 2022  The AdvantEDGE Authors
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
package sbi

import (
	"errors"
	"strconv"
	"sync"

	//"time"

	//dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	tm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-dai-mgr"
	gc "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-gis-cache"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	met "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	sam "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-swagger-api-mgr"
)

const moduleName string = "meep-dai-sbi"

var metricStore *met.MetricStore
var redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"
var influxAddr string = "http://meep-influxdb.default.svc.cluster.local:8086"

const postgresUser = "postgres"
const postgresPwd = "pwd"

var notifyAppContextDeletion func(string, string)

type SbiCfg struct {
	ModuleName                     string
	SandboxName                    string
	HostUrl                        string
	MepName                        string
	RedisAddr                      string
	InfluxAddr                     string
	PostgisHost                    string
	PostgisPort                    string
	OnboardedMecApplicationsFolder string
	Locality                       []string
	ScenarioNameCb                 func(string)
	AppInfoList                    func(tm.AppInfoList)
	NotifyAppContextDeletion       func(string, string)
	CleanUpCb                      func()
}

type DaiSbi struct {
	moduleName      string
	sandboxName     string
	hostUrl         string
	mepName         string
	scenarioName    string
	localityEnabled bool
	locality        map[string]bool
	mqLocal         *mq.MsgQueue
	handlerId       int
	apiMgr          *sam.SwaggerApiMgr
	activeModel     *mod.Model
	gisCache        *gc.GisCache
	daiMgr          *tm.DaiMgr
	//refreshTicker        *time.Ticker
	updateAppInfoCB      func(tm.AppInfoList)
	updateScenarioNameCB func(string)
	cleanUpCB            func()
	mutex                sync.Mutex
}

var sbi *DaiSbi = nil

// Init - DAI Service SBI initialization
func Init(cfg SbiCfg) (err error) {
	// Create new SBI instance
	if sbi != nil {
		sbi = nil
	}
	sbi = new(DaiSbi)
	sbi.moduleName = cfg.ModuleName
	sbi.sandboxName = cfg.SandboxName
	sbi.hostUrl = cfg.HostUrl
	sbi.mepName = cfg.MepName
	sbi.scenarioName = ""
	sbi.updateAppInfoCB = cfg.AppInfoList
	sbi.updateScenarioNameCB = cfg.ScenarioNameCb
	sbi.cleanUpCB = cfg.CleanUpCb
	redisAddr = cfg.RedisAddr
	influxAddr = cfg.InfluxAddr
	notifyAppContextDeletion = cfg.NotifyAppContextDeletion

	// Fill locality map
	if len(cfg.Locality) > 0 {
		sbi.locality = make(map[string]bool)
		for _, locality := range cfg.Locality {
			sbi.locality[locality] = true
		}
		sbi.localityEnabled = true
	} else {
		sbi.localityEnabled = false
	}

	// Create message queue
	sbi.mqLocal, err = mq.NewMsgQueue(mq.GetLocalName(sbi.sandboxName), sbi.moduleName, sbi.sandboxName, cfg.RedisAddr)
	if err != nil {
		log.Error("Failed to create Message Queue with error: ", err)
		return err
	}
	log.Info("Message Queue created")

	// Create Swagger API Manager
	log.Info("am.NewSwaggerApiMgr: ", sbi.moduleName, " - ", sbi.sandboxName, " - ", sbi.mepName)
	sbi.apiMgr, err = sam.NewSwaggerApiMgr(sbi.moduleName, sbi.sandboxName, sbi.mepName, sbi.mqLocal)
	if err != nil {
		log.Error("Failed to create Swagger API Manager. Error: ", err)
		return err
	}
	log.Info("Swagger API Manager created")

	// Create new active scenario model
	modelCfg := mod.ModelCfg{
		Name:      "activeScenario",
		Namespace: sbi.sandboxName,
		Module:    moduleName,
		UpdateCb:  nil,
		DbAddr:    cfg.RedisAddr,
	}
	sbi.activeModel, err = mod.NewModel(modelCfg)
	if err != nil {
		log.Error("Failed to create model: ", err.Error())
		return err
	}

	// Connect to GIS cache
	sbi.gisCache, err = gc.NewGisCache(sbi.sandboxName, cfg.RedisAddr)
	if err != nil {
		log.Error("Failed to GIS Cache: ", err.Error())
		return err
	}
	log.Info("Connected to GIS Cache")

	// Connect to DAI Manager
	cfgDai := tm.DaiCfg{
		Name:                     sbi.moduleName,
		Namespace:                sbi.sandboxName,
		User:                     postgresUser,
		Pwd:                      postgresPwd,
		Host:                     cfg.PostgisHost,
		Port:                     cfg.PostgisPort,
		NotifyAppContextDeletion: notifyAppContextDeletionSbi,
	}
	sbi.daiMgr, err = tm.NewDaiMgr(cfgDai)
	if err != nil {
		log.Error("Failed connection to VIS Traffic Manager: ", err)
		return err
	}
	log.Info("Connected to DAI Manager")

	// Delete any old tables
	_ = sbi.daiMgr.DeleteTables()

	// Create new tables
	err = sbi.daiMgr.CreateTables()
	if err != nil {
		log.Error("Failed to create tables: ", err)
		return err
	}
	log.Info("Created new DAI DB tables")

	err = sbi.daiMgr.LoadOnboardedMecApplications(cfg.OnboardedMecApplicationsFolder)
	if err != nil {
		log.Error("Failed to load simulating data: ", err)
		return err
	}
	log.Info("Created existing application")

	// Initialize service
	processActiveScenarioUpdate()

	return nil
}

// Run - MEEP DAI execution
func Run() (err error) {

	// Start Swagger API Manager (provider)
	err = sbi.apiMgr.Start(true, false)
	if err != nil {
		log.Error("Failed to start Swagger API Manager with error: ", err.Error())
		return err
	}
	log.Info("Swagger API Manager started")

	// Add module Swagger APIs
	err = sbi.apiMgr.AddApis()
	if err != nil {
		log.Error("Failed to add Swagger APIs with error: ", err.Error())
		return err
	}
	log.Info("Swagger APIs successfully added")

	// Register Message Queue handler
	handler := mq.MsgHandler{Handler: msgHandler, UserData: nil}
	sbi.handlerId, err = sbi.mqLocal.RegisterHandler(handler)
	if err != nil {
		log.Error("Failed to register message queue handler: ", err.Error())
		return err
	}

	// Start refresh loop
	//startRefreshTicker()

	return nil
}

func Stop() (err error) {
	if sbi == nil {
		return
	}

	// Stop refresh loop
	//stopRefreshTicker()

	if sbi.mqLocal != nil {
		sbi.mqLocal.UnregisterHandler(sbi.handlerId)
	}

	if sbi.apiMgr != nil {
		// Remove APIs
		err = sbi.apiMgr.RemoveApis()
		if err != nil {
			log.Error("Failed to remove APIs with err: ", err.Error())
			return err
		}
	}

	// Flush all DAI tables
	_ = sbi.daiMgr.DeleteTables()

	// Delete DAI instance
	if sbi.daiMgr != nil {
		err = sbi.daiMgr.DeleteDaiMgr()
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}

	return nil
}

// Message Queue handler
func msgHandler(msg *mq.Msg, userData interface{}) {
	switch msg.Message {
	case mq.MsgScenarioActivate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		processActiveScenarioUpdate()
	case mq.MsgScenarioUpdate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		processActiveScenarioUpdate()
	case mq.MsgScenarioTerminate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		processActiveScenarioTerminate()
	default:
		log.Trace("Ignoring unsupported message: ", mq.PrintMsg(msg))
	}
}

func processActiveScenarioTerminate() {
	log.Debug("processActiveScenarioTerminate")

	// Sync with active scenario store
	sbi.activeModel.UpdateScenario()

	sbi.cleanUpCB()
}

func processActiveScenarioUpdate() {

	sbi.mutex.Lock()
	defer sbi.mutex.Unlock()

	log.Debug("processActiveScenarioUpdate")

	// Sync with active scenario store
	sbi.activeModel.UpdateScenario()

	scenarioName := sbi.activeModel.GetScenarioName()

	// Connect to Metric Store
	sbi.updateScenarioNameCB(scenarioName)

	if scenarioName != sbi.scenarioName {
		sbi.scenarioName = scenarioName
		var err error

		metricStore, err = met.NewMetricStore(scenarioName, sbi.sandboxName, influxAddr, redisAddr)
		if err != nil {
			log.Error("Failed connection to metric-store: ", err)
		}
	}
}

// ETSI GS MEC 016 V2.2.1 (2020-04) Clause 7.3.3.1 GET
func GetApplicationListAppList(appNames []string, appProviders []string, appSoftVersions []string, vendorIds []string, serviceConts []string) (appListSbi *map[string]*tm.AppInfoList, err error) {
	log.Debug("GetApplicationListAppList: ", appNames)

	appListSbi, err = GetAllListAppList()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// Apply criteria
	if len(appNames) != 0 {
		filteredAppListSbi := filterAppNames(appNames, appListSbi)
		log.Debug("GetApplicationListAppList: After appName, filteredAppListSbi: ", filteredAppListSbi)
		if len(appProviders) != 0 && len(filteredAppListSbi) != 0 {
			filterExcludeAppProviders(appProviders, &filteredAppListSbi)
			log.Debug("GetApplicationListAppList: After appProvider, filteredAppListSbi: ", filteredAppListSbi)
		}
		if len(appSoftVersions) != 0 && len(filteredAppListSbi) != 0 {
			filterExcludeAppSoftVersion(appSoftVersions, &filteredAppListSbi)
			log.Debug("GetApplicationListAppList: After appSoftVersion, filteredAppListSbi: ", filteredAppListSbi)
		}
		// FIXME if len(vendorIds) != 0 && len(filteredAppListSbi) != 0 {
		if len(serviceConts) != 0 && len(filteredAppListSbi) != 0 {
			filterExcludeServiceConts(serviceConts, &filteredAppListSbi)
			log.Debug("GetApplicationListAppList: After appSoftVersion, filteredAppListSbi: ", filteredAppListSbi)
		}

		return &filteredAppListSbi, nil
	} else if len(appProviders) != 0 {
		filteredAppListSbi := filterAppProviders(appProviders, appListSbi)
		log.Debug("GetApplicationListAppList: After appProvider, filteredAppListSbi: ", filteredAppListSbi)
		if len(appSoftVersions) != 0 && len(filteredAppListSbi) != 0 {
			filterExcludeAppSoftVersion(appSoftVersions, &filteredAppListSbi)
			log.Debug("GetApplicationListAppList: After appSoftVersion, filteredAppListSbi: ", filteredAppListSbi)
		}
		// FIXME if len(vendorIds) != 0 && len(filteredAppListSbi) != 0 {
		if len(serviceConts) != 0 && len(filteredAppListSbi) != 0 {
			filterExcludeServiceConts(serviceConts, &filteredAppListSbi)
			log.Debug("GetApplicationListAppList: After appSoftVersion, filteredAppListSbi: ", filteredAppListSbi)
		}

		return &filteredAppListSbi, nil
	} else if len(appSoftVersions) != 0 {
		filteredAppListSbi := filterAppSoftVersions(appSoftVersions, appListSbi)
		// FIXME if len(vendorIds) != 0 && len(filteredAppListSbi) != 0 {
		if len(serviceConts) != 0 && len(filteredAppListSbi) != 0 {
			filterExcludeServiceConts(serviceConts, &filteredAppListSbi)
			log.Debug("GetApplicationListAppList: After appSoftVersion, filteredAppListSbi: ", filteredAppListSbi)
		}

		return &filteredAppListSbi, nil
		// FIXME } else if len(vendorIds) != 0 {
	} else if len(serviceConts) != 0 {
		filteredAppListSbi := filterServiceConts(serviceConts, appListSbi)

		return &filteredAppListSbi, nil
	}

	return appListSbi, nil
}

func filterAppNames(appNames []string, appListSbi *map[string]*tm.AppInfoList) map[string]*tm.AppInfoList {
	filteredAppListSbi := make(map[string]*tm.AppInfoList)
	for _, appName := range appNames {
		log.Debug("filterAppNames: Processing appName: ", appName)
		// Remove quotes
		if appName != "" {
			appName = appName[1 : len(appName)-1] // Remove quotes
		}
		log.Debug("filterAppNames: After removing quotes: appName: ", appName)
		// Search for the entry
		for _, item := range *appListSbi {
			if item.AppName == appName {
				filteredAppListSbi[appName] = item
			}
		} // End of 'for' statement
	} // End of 'for' statement

	return filteredAppListSbi
}

func filterAppProviders(appProviders []string, appListSbi *map[string]*tm.AppInfoList) map[string]*tm.AppInfoList {
	filteredAppListSbi := make(map[string]*tm.AppInfoList)
	for _, appProvider := range appProviders {
		log.Debug("filterAppProviders: Processing appProvider: ", appProvider)
		// Remove quotes
		if appProvider != "" {
			appProvider = appProvider[1 : len(appProvider)-1] // Remove quotes
		}
		log.Debug("filterAppProviders: After removing quotes: appSoftVersion: ", appProvider)
		// Search for the entry
		for _, item := range *appListSbi {
			if item.AppProvider == appProvider {
				filteredAppListSbi[item.AppName] = item
			}
		} // End of 'for' statement
	} // End of 'for' statement

	return filteredAppListSbi
}

func filterAppSoftVersions(appSoftVersions []string, appListSbi *map[string]*tm.AppInfoList) map[string]*tm.AppInfoList {
	filteredAppListSbi := make(map[string]*tm.AppInfoList)
	for _, appSoftVersion := range appSoftVersions {
		log.Debug("filterAppSoftVersions: Processing appSoftVersion: ", appSoftVersion)
		// Remove quotes
		if appSoftVersion != "" {
			appSoftVersion = appSoftVersion[1 : len(appSoftVersion)-1] // Remove quotes
		}
		log.Debug("filterAppSoftVersions: After removing quotes: appSoftVersion: ", appSoftVersion)
		// Search for the entry
		for _, item := range *appListSbi {
			if item.AppSoftVersion == appSoftVersion {
				filteredAppListSbi[item.AppName] = item
			}
		} // End of 'for' statement
	} // End of 'for' statement

	return filteredAppListSbi
}

func filterServiceConts(serviceConts []string, appListSbi *map[string]*tm.AppInfoList) map[string]*tm.AppInfoList {
	filteredAppListSbi := make(map[string]*tm.AppInfoList)
	for _, serviceCount := range serviceConts {
		log.Debug("filterServiceConts: Processing serviceCount: ", serviceCount)
		// Remove quotes
		if serviceCount != "" {
			serviceCount = serviceCount[1 : len(serviceCount)-1] // Remove quotes
		}
		log.Debug("filterServiceConts: After removing quotes: appSoftVersion: ", serviceCount)
		svcCount, _ := strconv.ParseUint(serviceCount, 10, 32)
		svcCount32 := uint32(svcCount)
		// Search for the entry
		for _, item := range *appListSbi {
			if len(item.AppCharcs) != 0 {
				for _, sc := range item.AppCharcs {
					if sc.ServiceCont != nil && *sc.ServiceCont == svcCount32 {
						filteredAppListSbi[item.AppName] = item
						break
					}
				} // End of 'for' statement
			}
		} // End of 'for' statement
	} // End of 'for' statement

	return filteredAppListSbi
}

func filterExcludeAppProviders(appProviders []string, filteredAppListSbi *map[string]*tm.AppInfoList) {
	for _, appProvider := range appProviders {
		log.Debug("filterExcludeAppProviders: Processing appProvider: ", appProvider)
		// Remove quotes
		if appProvider != "" {
			appProvider = appProvider[1 : len(appProvider)-1] // Remove quotes
		}
		log.Debug("filterExcludeAppProviders: After removing quotes: appProvider: ", appProvider)
		// Search for the entry
		for _, item := range *filteredAppListSbi {
			if item.AppProvider != appProvider {
				log.Debug("filterExcludeAppProviders: Removing entry: ", item.AppName)
				delete(*filteredAppListSbi, item.AppName)
			}
		} // End of 'for' statement
	} // End of 'for' statement
}

func filterExcludeAppSoftVersion(appSoftVersions []string, filteredAppListSbi *map[string]*tm.AppInfoList) {
	for _, appSoftVersion := range appSoftVersions {
		log.Debug("GetApplicationListAppList: Processing appProvider: ", appSoftVersion)
		// Remove quotes
		if appSoftVersion != "" {
			appSoftVersion = appSoftVersion[1 : len(appSoftVersion)-1] // Remove quotes
		}
		log.Debug("GetApplicationListAppList: After removing quotes: appSoftVersion: ", appSoftVersion)
		// Search for the entry
		for _, item := range *filteredAppListSbi {
			if item.AppSoftVersion != appSoftVersion {
				log.Debug("GetApplicationListAppList: Removing entry: ", appSoftVersion)
				delete(*filteredAppListSbi, item.AppName)
			}
		} // End of 'for' statement
	} // End of 'for' statement
}

func filterExcludeServiceConts(serviceConts []string, filteredAppListSbi *map[string]*tm.AppInfoList) {
	for _, serviceCount := range serviceConts {
		log.Debug("filterExcludeServiceConts: Processing serviceCount: ", serviceCount)
		// Remove quotes
		if serviceCount != "" {
			serviceCount = serviceCount[1 : len(serviceCount)-1] // Remove quotes
		}
		log.Debug("filterExcludeServiceConts: After removing quotes: serviceCount: ", serviceCount)
		svcCount, _ := strconv.ParseUint(serviceCount, 10, 32)
		svcCount32 := uint32(svcCount)
		// Search for the entry
		for _, item := range *filteredAppListSbi {
			if len(item.AppCharcs) != 0 {
				for _, sc := range item.AppCharcs {
					if sc.ServiceCont != nil && *sc.ServiceCont != svcCount32 {
						log.Debug("filterExcludeServiceConts: Removing entry: ", item.AppName)
						delete(*filteredAppListSbi, item.AppName)
						break
					}
				} // End of 'for' statement
			}
		} // End of 'for' statement
	} // End of 'for' statement
}

func GetAllListAppList() (appListSbi *map[string]*tm.AppInfoList, err error) {

	// Get list of application
	appInfoList, err := sbi.daiMgr.GetAllAppInfoListEntry()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	appListSbi = &appInfoList

	//sbi.updateAppInfoCB(appListSbi)
	return appListSbi, nil
}

func CreateAppContext(appContextSbi *tm.AppContext) (appContextSbi_ *tm.AppContext, err error) {

	appContextSbi_, err = sbi.daiMgr.CreateAppContext(appContextSbi, sbi.hostUrl, sbi.sandboxName)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return appContextSbi_, nil
}

func DeleteAppContext(contextId string) (err error) {

	err = sbi.daiMgr.DeleteAppContext(contextId)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func PutAppContext(appContextSbi tm.AppContext) (err error) {

	err = sbi.daiMgr.PutAppContext(appContextSbi)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func PosApplicationLocationAvailability(applicationLocationAvailabilitySbi *tm.ApplicationLocationAvailability) (applicationLocationAvailability_ *tm.ApplicationLocationAvailability, err error) {

	// Retrieve the AppInfo data for the specified application
	var appNames = []string{strconv.Quote(applicationLocationAvailabilitySbi.AppInfo.AppName)}
	var appProviders = []string{strconv.Quote(applicationLocationAvailabilitySbi.AppInfo.AppProvider)}
	var appSoftVersions []string
	if applicationLocationAvailabilitySbi.AppInfo.AppSoftVersion != nil {
		appSoftVersions = []string{strconv.Quote(*applicationLocationAvailabilitySbi.AppInfo.AppSoftVersion)}
	}
	var vendorIds []string    // FIXME To be done
	var serviceConts []string // FIXME To be done
	appListSbi, err := GetApplicationListAppList(appNames, appProviders, appSoftVersions, vendorIds, serviceConts)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	if len(*appListSbi) == 0 {
		err = errors.New("No application found for " + applicationLocationAvailabilitySbi.AppInfo.AppName)
		log.Error(err.Error())
		return nil, err
	}

	applicationLocationAvailability_ = applicationLocationAvailabilitySbi
	for _, item := range *appListSbi {
		if item.AppName != applicationLocationAvailabilitySbi.AppInfo.AppName {
			err = errors.New("Wrong application found for " + applicationLocationAvailabilitySbi.AppInfo.AppName + "- " + item.AppName)
			log.Error(err.Error())
			return nil, err
		}
		applicationLocationAvailability_.AppInfo.AvailableLocations = make(tm.AvailableLocations, 1)
		applicationLocationAvailability_.AppInfo.AvailableLocations[0].AppLocation = new(tm.LocationConstraints)
		*applicationLocationAvailability_.AppInfo.AvailableLocations[0].AppLocation = item.AppLocation
		break
	} // End of 'for' statement

	return applicationLocationAvailability_, nil
}

func notifyAppContextDeletionSbi(notifyUrl string, appContextId string) {
	log.Debug(">>> notifyAppContextDeletionSbi: ", appContextId)

	notifyAppContextDeletion(notifyUrl, appContextId)
}

/*func isAppConnected(app string) bool {
	parentNode := sbi.activeModel.GetNodeParent(app)
	if parentNode != nil {
		pl := parentNode.(*dataModel.PhysicalLocation)
		return pl.Connected
	}
	return false
}

func isInLocality(name string) bool {
	if sbi.localityEnabled {
		ctx := sbi.activeModel.GetNodeContext(name)
		if ctx == nil {
			log.Error("Error getting context for: " + name)
			return false
		}
		if _, found := sbi.locality[ctx.Parents[mod.Zone]]; !found {
			return false
		}
	}
	return true
}*/
