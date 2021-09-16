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

import _ from 'lodash';
import { connect } from 'react-redux';
import React, { Component } from 'react';
import autoBind from 'react-autobind';
import axios from 'axios';
import { updateObject, deepCopy } from '../util/object-util';
import { ToolbarFixedAdjust } from '@rmwc/toolbar';

// Import JS dependencies
import * as meepPlatformCtrlRestApiClient from '../../../../../js-packages/meep-platform-ctrl-client/src/index.js';
import * as meepSandboxCtrlRestApiClient from '../../../../../js-packages/meep-sandbox-ctrl-client/src/index.js';
import * as meepMonEngineRestApiClient from '../../../../../js-packages/meep-mon-engine-client/src/index.js';
import * as meepGisEngineRestApiClient from '../../../../../js-packages/meep-gis-engine-client/src/index.js';
import * as meepAuthSvcRestApiClient from '../../../../../js-packages/meep-auth-svc-client/src/index.js';

import MeepTopBar from '@/js/components/meep-top-bar';
import Footer from '@/js/components/footer';
import IDSignInOAuthDialog from '@/js/components/dialogs/id-sign-in-oauth-dialog';
import IDSessionTerminatedDialog from '@/js/components/dialogs/id-session-terminated-dialog';
import CfgPageContainer from '@/js/containers/cfg/cfg-page-container';
import ExecPageContainer from '@/js/containers/exec/exec-page-container';
import SettingsPageContainer from '@/js/containers/settings/settings-page-container';
import MonitorPageContainer from '@/js/containers/monitor/monitor-page-container';
import HomePageContainer from '@/js/containers/home/home-page-container';

import {
  HOST_PATH,
  TYPE_CFG,
  TYPE_EXEC,
  EXEC_STATE_IDLE,
  EXEC_STATE_DEPLOYED,
  NO_SCENARIO_NAME,
  PAGE_CONFIGURE,
  PAGE_EXECUTE,
  PAGE_MONITOR,
  PAGE_SETTINGS,
  PAGE_HOME,
  STATUS_SIGNED_IN,
  STATUS_SIGNING_IN,
  STATUS_SIGNED_OUT,
  STATUS_SIGNIN_NOT_SUPPORTED,
  PAGE_HOME_INDEX,
  PAGE_CONFIGURE_INDEX,
  IDC_DIALOG_SIGN_IN,
  IDC_DIALOG_SESSION_TERMINATED
} from '../meep-constants';

import {
  parseScenario,
  createNewScenario,
  addElementToScenario,
  updateElementInScenario,
  cloneElementInScenario,
  removeElementFromScenario
} from '../util/scenario-utils';

import {
  uiChangeCurrentDialog,
  uiChangeCurrentPage,
  uiExecChangeSandbox,
  uiExecChangeSandboxList,
  uiExecChangeSandboxCfg,
  uiExecChangeEventCreationMode,
  uiExecChangeEventReplayMode,
  uiChangeSignInStatus,
  uiChangeSignInUsername,
  uiChangeCurrentTab
} from '../state/ui';

import {
  execChangeScenario,
  execChangeScenarioState,
  execChangeScenarioPodsPhases,
  execChangeServiceMaps,
  execChangeMapUeList,
  execChangeMapPoaList,
  execChangeMapComputeList,
  execChangeVisData,
  execChangeTable,
  execChangeCorePodsPhases,
  execChangeOkToTerminate,
  corePodsRunning,
  corePodsErrors,
  execVisFilteredData,
  execChangeReplayStatus,
  execChangeAppInstanceTable
} from '../state/exec';

import {
  cfgChangeScenario,
  cfgChangeVisData,
  cfgChangeTable,
  cfgChangeMap
} from '../state/cfg';

import {
  FIELD_CONNECTIVITY_MODEL,
  getElemByName,
  getElemFieldVal
} from '../util/elem-utils';

// REST API Clients
var basepathPlatformCtrl = HOST_PATH + '/platform-ctrl/v1';
meepPlatformCtrlRestApiClient.ApiClient.instance.basePath = basepathPlatformCtrl.replace(/\/+$/,'');
var basepathSandboxCtrl = HOST_PATH + '/sandbox-ctrl/v1';
meepSandboxCtrlRestApiClient.ApiClient.instance.basePath = basepathSandboxCtrl.replace(/\/+$/,'');
var basepathMonEngine = HOST_PATH + '/mon-engine/v1';
meepMonEngineRestApiClient.ApiClient.instance.basePath = basepathMonEngine.replace(/\/+$/,'');
var basepathGisEngine = HOST_PATH + '/gis/v1';
meepGisEngineRestApiClient.ApiClient.instance.basePath = basepathGisEngine.replace(/\/+$/,'');
var basepathAuthSvc = HOST_PATH + '/auth/v1';
meepAuthSvcRestApiClient.ApiClient.instance.basePath = basepathAuthSvc.replace(/\/+$/,'');

const SESSION_KEEPALIVE_INTERVAL = 600000; // 10 min

class MeepContainer extends Component {
  constructor(props) {
    super(props);
    autoBind(this);

    this.pduSessions = null;

    this.sessionKeepaliveTimer = null;
    this.platformRefreshIntervalTimer = null;
    this.execPageRefreshIntervalTimer = null;
    this.replayStatusRefreshIntervalTimer = null;
    this.meepScenarioConfigurationApi = new meepPlatformCtrlRestApiClient.ScenarioConfigurationApi();
    this.meepSandboxControlApi = new meepPlatformCtrlRestApiClient.SandboxControlApi();  
    this.meepAppInfoApi = new meepSandboxCtrlRestApiClient.ApplicationsApi();
    this.meepActiveScenarioApi = new meepSandboxCtrlRestApiClient.ActiveScenarioApi();
    this.meepEventsApi = new meepSandboxCtrlRestApiClient.EventsApi();
    this.meepEventReplayApi = new meepSandboxCtrlRestApiClient.EventReplayApi();
    this.meepEventAutomationApi = new meepGisEngineRestApiClient.AutomationApi();
    this.meepGeoDataApi = new meepGisEngineRestApiClient.GeospatialDataApi();
    this.meepAuthApi = new meepAuthSvcRestApiClient.AuthApi();
    this.meepConnectivityApi = new meepSandboxCtrlRestApiClient.ConnectivityApi();
  }

  componentDidMount() {
    document.title = 'AdvantEDGE';
    this.setBasepath(this.props.sandbox);
    this.refreshScenario();
    this.monitorTabFocus();

    this.meepAuthApi.loginSupported((_, __, response) => {
      if (response.status === 404) {
        this.props.changeSignInStatus(STATUS_SIGNIN_NOT_SUPPORTED);
      } else if (response.status === 200) {
        this.props.changeSignInStatus(STATUS_SIGNED_IN);
      } else {
        this.props.changeSignInStatus(STATUS_SIGNED_OUT);
        this.logout();
      }
      this.startTimers();
    });
  }

  componentWillMount() {
    // Handle OAuth login in progress
    if (this.props.signInStatus === STATUS_SIGNING_IN) {
      let params = (new URL(document.location)).searchParams;
      let userName = params.get('user');
      if (userName) {
        this.props.changeSignInUsername(userName);
        window.history.replaceState({}, document.title, '/');
        this.props.changeSignInStatus(STATUS_SIGNED_IN);
        this.props.changeCurrentPage(PAGE_CONFIGURE);
        this.props.changeTabIndex(PAGE_CONFIGURE_INDEX);
      } else {
        // Sign in failed
        this.logout();
        this.props.changeSignInStatus(STATUS_SIGNED_OUT);
      }
    }
  }
  
  // Timers
  startTimers() {
    if (this.props.signInStatus === STATUS_SIGNED_IN || this.props.signInStatus === STATUS_SIGNIN_NOT_SUPPORTED) {
      this.startPlatformRefresh();
      this.startExecPageRefresh();
      this.startReplayStatusRefresh();
    }
  }
  stopTimers() {
    this.stopReplayStatusRefresh();
    this.stopExecPageRefresh();
    this.stopPlatformRefresh();
  }

  // Platform refresh
  startPlatformRefresh() {
    this.platformRefreshIntervalTimer = setInterval(
      () => {
        // Make sure watchdog timer is running if we are signed in
        if (this.props.signInStatus === STATUS_SIGNED_IN) {
          if (!this.sessionKeepaliveTimer) {
            this.startSessionKeepaliveTimer();
          }
        }
        this.checkPlatformStatus();
      },
      1000
    );
  }
  stopPlatformRefresh() {
    if (this.platformRefreshIntervalTimer) {
      clearInterval(this.platformRefreshIntervalTimer);
      this.platformRefreshIntervalTimer = null;
    }
  }

  // Exec page refresh 
  startExecPageRefresh() {
    this.execPageRefreshIntervalTimer = setInterval(
      () => {
        if (this.props.page === PAGE_EXECUTE) {
          this.refreshSandboxList();
          if (this.props.sandbox) {
            this.checkScenarioStatus();
            this.refreshPduSessions();
            this.refreshScenario();
            this.refreshMap();
            // Only update while scenario is running
            if (this.props.execScenarioState === 'DEPLOYED') {
              this.refreshAppInstancesTable();
            }
          }
        }
      },
      1000
    );
  }

  stopExecPageRefresh() {
    if (this.execPageRefreshIntervalTimer) {
      clearInterval(this.execPageRefreshIntervalTimer);
      this.execPageRefreshIntervalTimer = null;
    }
  }

  // Replay status refresh
  startReplayStatusRefresh() {
    this.replayStatusRefreshIntervalTimer = setInterval(
      () => this.checkReplayStatus(),
      1000
    );
  }
  stopReplayStatusRefresh() {
    if (this.replayStatusRefreshIntervalTimer) {
      clearInterval(this.replayStatusRefreshIntervalTimer);
      this.replayStatusRefreshIntervalTimer = null;
    }
  }

  monitorTabFocus() {
    var hidden, visibilityChange;
    if (typeof document.hidden !== 'undefined') {
      // Opera 12.10 and Firefox 18 and later support
      hidden = 'hidden';
      visibilityChange = 'visibilitychange';
    } else if (typeof document.msHidden !== 'undefined') {
      hidden = 'msHidden';
      visibilityChange = 'msvisibilitychange';
    } else if (typeof document.webkitHidden !== 'undefined') {
      hidden = 'webkitHidden';
      visibilityChange = 'webkitvisibilitychange';
    }

    const handleVisibilityChange = () => {
      if (document[hidden]) {
        this.stopTimers();
      } else {
        this.startTimers();
      }
    };

    // Warn if the browser doesn't support addEventListener or the Page Visibility API
    if (
      typeof document.addEventListener === 'undefined' ||
      hidden === undefined
    ) {
      // TODO: consider showing an alert
      // console.log('This demo requires a browser, such as Google Chrome or Firefox, that supports the Page Visibility API.');
    } else {
      // Handle page visibility change
      document.addEventListener(
        visibilityChange,
        handleVisibilityChange,
        false
      );
    }
  }

  checkPlatformStatus() {
    // Core pods
    axios
      .get(`${basepathMonEngine}/states?long=true&type=core&sandbox=all`)
      .then(res => {
        this.props.changeCorePodsPhases(res.data.podStatus);
      })
      .catch(error => {
        this.props.changeCorePodsPhases([]);

        // Log out if session was terminated
        if (this.props.signInStatus === STATUS_SIGNED_IN) {
          if (error.response.status === 401) {
            this.logout();
            this.props.changeCurrentDialog(IDC_DIALOG_SESSION_TERMINATED);
          }
        }
      });
  }

  checkScenarioStatus() {
    // Scenario pods
    axios
      .get(`${basepathMonEngine}/states?long=true&type=scenario&sandbox=${this.props.sandbox}`)
      .then(res => {
        var scenarioPodsPhases = res.data.podStatus;
        this.props.changeScenarioPodsPhases(scenarioPodsPhases);
      })
      .catch(() => {
        this.props.changeScenarioPodsPhases([]);
      });

    // Service maps
    axios
      .get(`${basepathSandboxCtrl}/active/serviceMaps`)
      .then(res => {
        var serviceMaps = res.data;
        this.props.changeServiceMaps(serviceMaps);
      })
      .catch(() => {
        this.props.changeServiceMaps([]);
      });
  }

  /**
   * Callback function to receive the result of the getSandboxList operation.
   * @callback module:api/SandboxControlApi~getSandboxListCallback
   * @param {String} error Error message, if any.
   * @param {module:model/SandboxList} data The data returned by the service call.
   * @param {String} response The complete HTTP response.
   */
  getSandboxListCb(error, data) {
    if (error !== null) {
      // TODO: consider showing an alert
      return;
    }

    // Update list of sandboxes, if any
    var orderedSandboxList = _.map(data.sandboxes, 'name');
    if ((orderedSandboxList.length !== this.props.sandboxes.length) ||
      orderedSandboxList.every((value, index) => value !== this.props.sandboxes[index])) {
      this.props.changeSandboxList(orderedSandboxList);
    }
  }

  refreshSandboxList() {
    this.meepSandboxControlApi.getSandboxList((error, data, response) => {
      this.getSandboxListCb(error, data, response);
    });
  }

  /**
   * Callback function to receive the result of the getReplayStatus operation.
   * @callback module:api/EventReplayApi~getReplayStatusCallback
   * @param {String} error Error message, if any.
   * @param {module:model/Replay} data The data returned by the service call.
   */
  getReplayStatusCb(error, data) {
    this.props.changeReplayStatus((error === null) ? data : null);
  }

  checkReplayStatus() {
    if (this.props.execScenarioState === EXEC_STATE_IDLE) {
      return;
    }

    if (this.props.eventCfgMode || this.props.eventReplayMode) {
      this.meepEventReplayApi.getReplayStatus((error, data, response) => {
        this.getReplayStatusCb(error, data, response);
      });
    }
  }

  setMainContent(targetId) {
    this.props.changeCurrentPage(targetId);
  }

  /**
   * Callback function to receive the result of the getActiveScenario operation.
   * @callback module:api/ScenarioExecutionApi~getActiveScenarioCallback
   * @param {String} error Error message, if any.
   * @param {module:model/Scenario} data The data returned by the service call.
   */
  getActiveScenarioCb(error, data) {
    if ((error !== null) || (!data.deployment)) {
      this.props.execChangeScenarioState(EXEC_STATE_IDLE);
      this.props.execChangeOkToTerminate(false);
      return;
    }

    // Store & Process deployed scenario
    this.execSetScenario(data);

    // TODO set a timer of 2 seconds
    this.props.execChangeScenarioState(EXEC_STATE_DEPLOYED);
    setTimeout(() => {
      if (this.props.execScenarioState === EXEC_STATE_DEPLOYED) {
        this.props.execChangeOkToTerminate(true);
      }
    }, 2000);
  }

  changeScenario(pageType, scenario) {
    this.updateScenario(pageType, scenario, false);
  }

  // Change & process scenario
  updateScenario(pageType, scenario, reInitVisView) {
    // Change scenario state
    if (pageType === TYPE_CFG) {
      this.props.cfgChangeScenario(scenario);
    } else {
      this.props.execChangeScenario(scenario);
    }

    // Parse Scenario object to retrieve visualization data and scenario table
    var page = (pageType === TYPE_CFG) ? this.props.cfg : this.props.exec;
    var pduSessions = (pageType === TYPE_CFG) ? null : this.pduSessions;
    var scenarioName = (page.scenario) ? page.scenario.name : '';
    var parsedScenario = parseScenario(page.scenario, pduSessions);
    var updatedMapData = updateObject({}, parsedScenario.mapData);
    var updatedVisData = updateObject(page.vis.data, parsedScenario.visData);
    // updatedVisData.nodes._data.sort();
    // updatedVisData.edges._data.sort();
    var updatedTable = updateObject(page.table, parsedScenario.table);

    // Update connectivity mode
    this.updateConnectivityMode(pageType, updatedTable, scenarioName);

    // Dispatch state updates
    if (pageType === TYPE_CFG) {
      this.props.cfgChangeVisData(updatedVisData);
      this.props.cfgChangeTable(updatedTable);
      // Update map after table to make sure latest entries are available
      this.props.cfgChangeMap(updatedMapData);

      const vis = this.props.cfgVis;
      if (vis && vis.network && vis.network.setData) {
        //save the canvas position and scale level in vis
        var view;
        if (!reInitVisView) {
          view = deepCopy(vis.network.canvas.body.view);
        }
        vis.network.setData(updatedVisData);
        if (view) {
          //restore the canvas position and scale in vis
          vis.network.canvas.body.view = view;
        }
      }
    } else {
      this.props.execChangeVisData(updatedVisData);
      this.props.execChangeTable(updatedTable);

      const vis = this.props.execVis;
      if (vis && vis.network && vis.network.setData) {
        _.defer(() => {
          //save the canvas position and scale level in vis
          const view = deepCopy(vis.network.canvas.body.view);
          vis.network.setData(this.props.execVisData);
          //restore the canvas position and scale in vis
          vis.network.canvas.body.view = view;
        });
      } 
    }
  }

  // Create, store & process new scenario
  cfgCreateScenario(name) {
    var scenario = createNewScenario(name);
    this.updateScenario(TYPE_CFG, scenario, true);
  }
  // Set & process scenario
  cfgSetScenario(scenario) {
    this.updateScenario(TYPE_CFG, scenario, true);
  }
  // Delete & process scenario
  cfgDeleteScenario() {
    var scenario = createNewScenario(NO_SCENARIO_NAME);
    this.updateScenario(TYPE_CFG, scenario, true);
  }
  // Add new element to scenario
  cfgNewScenarioElem(element, scenarioUpdate) {
    var scenario = this.props.cfg.scenario;
    var updatedScenario = updateObject({}, scenario);
    addElementToScenario(updatedScenario, element);
    if (scenarioUpdate) {
      this.changeScenario(TYPE_CFG, updatedScenario);
    }
  }
  // Update element in scenario
  cfgUpdateScenarioElem(element) {
    var scenario = this.props.cfg.scenario;
    var updatedScenario = updateObject({}, scenario);
    updateElementInScenario(updatedScenario, element);
    this.changeScenario(TYPE_CFG, updatedScenario);
  }
  // Delete element in scenario (also deletes child elements)
  cfgDeleteScenarioElem(element) {
    var scenario = this.props.cfg.scenario;
    var updatedScenario = updateObject({}, scenario);
    removeElementFromScenario(updatedScenario, element);
    this.changeScenario(TYPE_CFG, updatedScenario);
  }
  // Clone element in scenario
  cfgCloneScenarioElem(element) {
    var updatedScenario = updateObject({}, this.props.cfg.scenario);
    cloneElementInScenario(updatedScenario, element, this.props.cfg.table);
    this.changeScenario(TYPE_CFG, updatedScenario);
  }

  // Set & process scenario
  execSetScenario(scenario) {
    this.updateScenario(TYPE_EXEC, scenario, true);
  }
  // Delete & process scenario
  execDeleteScenario() {
    var scenario = createNewScenario(NO_SCENARIO_NAME);
    this.updateScenario(TYPE_EXEC, scenario, true);
  }

  // Refresh Active scenario
  refreshScenario() {
    this.meepActiveScenarioApi.getActiveScenario(null, (error, data) =>
      this.getActiveScenarioCb(error, data)
    );
  }

  /**
   * Callback function to receive the result of the postHttpQuery operation.
   * @callback module:api/AppsApi~applicationsGET
   * @param {String} error Error message, if any.
   * @param {module:model/ApplicationInfo} data The data returned by the service call.
   * @param {String} response The complete HTTP response.
   */
  getAppInstancesCb(error, data) {
    if (error !== null) {
      this.props.changeAppInstanceTable([]);
      return;
    }

    // Update App Instance table only if data is different 
    var appInstances = data ? data : [];
    const isArrayEqual = (x, y) => _.isEmpty(_.xorWith(x, y, _.isEqual));
    if (!isArrayEqual(this.props.appInstanceTable,appInstances)) {
      this.props.changeAppInstanceTable(appInstances);
    }
  }

  // Refresh App Instances
  refreshAppInstancesTable() {
    this.meepAppInfoApi.applicationsGET(null, (error, data, response) => {
      this.getAppInstancesCb(error, data, response);
    });
  }

  /**
   * Callback function to receive the result of the getAssetData operation.
   * @callback module:api/GeospatialDataApi~getAssetDataCallback
   * @param {String} error Error message, if any.
   * @param {module:model/GeoDataAssetList} data The data returned by the service call.
   * @param {String} response The complete HTTP response.
   */
  getAssetDataCb(error, data) {

    if (error !== null) {
      return;
    }

    var ueList = [];
    var poaList = [];
    var computeList = [];

    // Extract assets by type

    if (data.geoDataAssets) {
      _.forEach(data.geoDataAssets, asset => {
        switch (asset.assetType) {
        case 'UE':
          ueList.push(asset);
          break;
        case 'POA':
          poaList.push(asset);
          break;
        case 'COMPUTE':
          computeList.push(asset);
          break;
        default:
          break;
        }
      });
    }
    
    // Update asset map
    var assetMap = {
      ueList: _.sortBy(ueList, ['assetName']),
      poaList: _.sortBy(poaList, ['assetName']),
      computeList: _.sortBy(computeList, ['assetName'])
    };

    //Update UE LIST, COMPUTE LIST, POA list
    this.props.execChangeMapUeList(assetMap.ueList? _.sortBy(assetMap.ueList, ['assetName']) : []);
    this.props.execChangeMapPoaList(assetMap.poaList? _.sortBy(assetMap.poaList, ['assetName']) : []);
    this.props.execChangeMapComputeList(assetMap.computeList? _.sortBy(assetMap.computeList, ['assetName']) : []);
  }

  /**
   * Callback function to receive the result of the getAssetData operation.
   * @callback module:api/GeospatialDataApi~getAssetDataCallback
   * @param {String} error Error message, if any.
   * @param {module:model/GeoDataAssetList} data The data returned by the service call.
   * @param {String} response The complete HTTP response.
   */
  getUeAssetDataCb(error, data) {
    if (error !== null) {
      return;
    }
    
    // Update UE list
    this.props.execChangeMapUeList(data.geoDataAssets ? _.sortBy(data.geoDataAssets, ['assetName']) : []);
  }

  /**
   * Callback function to receive the result of the getAssetData operation.
   * @callback module:api/GeospatialDataApi~getAssetDataCallback
   * @param {String} error Error message, if any.
   * @param {module:model/GeoDataAssetList} data The data returned by the service call.
   * @param {String} response The complete HTTP response.
   */
  getPoaAssetDataCb(error, data) {
    if (error !== null) {
      return;
    }

    // Update POA list
    this.props.execChangeMapPoaList(data.geoDataAssets ? _.sortBy(data.geoDataAssets, ['assetName']) : []);
  }

  /**
   * Callback function to receive the result of the getAssetData operation.
   * @callback module:api/GeospatialDataApi~getAssetDataCallback
   * @param {String} error Error message, if any.
   * @param {module:model/GeoDataAssetList} data The data returned by the service call.
   * @param {String} response The complete HTTP response.
   */
  getComputeAssetDataCb(error, data) {
    if (error !== null) {
      return;
    }

    // Update Compute list
    this.props.execChangeMapComputeList(data.geoDataAssets ? _.sortBy(data.geoDataAssets, ['assetName']) : []);
  }

  // Refresh Map
  refreshMap() { 
    this.meepGeoDataApi.getAssetData({}, (error, data) =>
      this.getAssetDataCb(error, data)
    );
  }

  /**
   * Callback function to receive the result of the getPduSessionList operation.
   * @callback module:api/ConnectivityApi~getPduSessionListCallback
   * @param {String} error Error message, if any.
   * @param {Array.<module:model/PDUSessionList>} data The data returned by the service call.
   * @param {String} response The complete HTTP response.
   */
  getPduSessionListCb(error, data) {
    if (error !== null) {
      return;
    }

    // Update PDU session map
    this.pduSessions = data;
  }

  // Refresh PDU Sessions
  refreshPduSessions() {
    if (this.connectivityMode === 'PDU') {
      this.meepConnectivityApi.getPduSessionList(null, (error, data) =>
        this.getPduSessionListCb(error, data)
      );
    }
  }

  // Update connectivity mode
  updateConnectivityMode(pageType, table, scenarioName) {
    if (pageType === TYPE_CFG) {
      this.connectivityMode = '';
    } else {
      this.connectivityMode = getElemFieldVal(getElemByName(table.entries, scenarioName), FIELD_CONNECTIVITY_MODEL);
    }
  }

  // Set sandox-specific API basepath
  setBasepath(sandboxName) {
    var sandboxPath = (sandboxName) ? '/' + sandboxName : '';
    basepathSandboxCtrl = HOST_PATH + sandboxPath + '/sandbox-ctrl/v1';
    meepSandboxCtrlRestApiClient.ApiClient.instance.basePath = basepathSandboxCtrl.replace(/\/+$/,'');
    basepathGisEngine = HOST_PATH + sandboxPath + '/gis/v1';
    meepGisEngineRestApiClient.ApiClient.instance.basePath = basepathGisEngine.replace(/\/+$/,'');
  }

  /**
   * Callback function to receive the result of the createSandboxWithName operation.
   * @callback module:api/SandboxControlApi~createSandboxWithNameCallback
   * @param {String} error Error message, if any.
   * @param data This operation does not return a value.
   * @param {String} response The complete HTTP response.
   */
  createSandboxWithNameCb(error) {
    if (error) {
      this.props.changeSandbox('');
      return;
    }

    // Set active sandbox
    this.setBasepath(this.props.sandbox);
    this.refreshScenario();
  }

  // Create a new sandbox
  createSandbox(name) {
    this.props.changeSandbox(name);
    this.meepSandboxControlApi.createSandboxWithName(name, {}, (error, data, response) => {
      this.createSandboxWithNameCb(error, data, response);
    });
  }

  // Set active sandbox
  setSandbox(name) {
    this.setBasepath(name);
    this.refreshScenario();
    this.refreshMap();
    this.props.changeSandbox(name);
  }

  /**
   * Callback function to receive the result of the deleteSandbox operation.
   * @callback module:api/SandboxControlApi~deleteSandboxCallback
   * @param {String} error Error message, if any.
   * @param data This operation does not return a value.
   * @param {String} response The complete HTTP response.
   */
  deleteSandboxCb(error) {
    if (error !== null) {
      // TODO consider showing an alert  (i.e. toast)
      return;
    }

    // Reset sandbox
    this.props.changeSandbox(null);
    this.setBasepath(null);

    // Delete the active scenario
    this.execDeleteScenario(TYPE_EXEC);
    this.props.execChangeScenarioState(EXEC_STATE_IDLE);
    this.props.execChangeOkToTerminate(false);
  }

  // Delete the active sandbox
  deleteSandbox() {
    this.meepSandboxControlApi.deleteSandbox(this.props.sandbox, (error, data, response) => {
      this.deleteSandboxCb(error, data, response);
    });
  }

  // Session Keep-alive
  startSessionKeepaliveTimer() {
    if (!this.sessionKeepaliveTimer) {
      this.meepAuthApi.triggerWatchdog();

      // Start keepalive timer
      this.sessionKeepaliveTimer = setInterval(() => {
        this.meepAuthApi.triggerWatchdog();
      },
      SESSION_KEEPALIVE_INTERVAL
      );
    }
  }

  stopSessionKeepaliveTimer() {
    if (this.sessionKeepaliveTimer) {
      clearInterval(this.sessionKeepaliveTimer);
      this.sessionKeepaliveTimer = null;
    }
  }

  /**
   * Callback function to receive the result of the logout operation.
   * @callback module:api/AuthenticationApi~logout
   * @param {String} error Error message, if any.
   * @param none
   * @param {String} response The complete HTTP response.
   */
  logoutCb() {
    this.props.changeSignInStatus(STATUS_SIGNED_OUT);
    if (this.props.currentPage !== PAGE_HOME) {
      this.props.changeCurrentPage(PAGE_HOME);
      this.props.changeTabIndex(PAGE_HOME_INDEX);
    }
  }

  logout() {
    this.stopSessionKeepaliveTimer();
    this.meepAuthApi.logout((error, data, response) => {
      this.logoutCb(error, data, response);
    });
  }

  signInProcedure() {
    if (this.props.signInStatus === STATUS_SIGNED_IN) {
      this.logout();
    } else {
      this.props.changeCurrentDialog(IDC_DIALOG_SIGN_IN);
      // this.props.changeCurrentPage(PAGE_HOME);
      // this.props.changeTabIndex(PAGE_HOME_INDEX);
    }
  }

  signInOAuth(provider) {
    // Set state to signing in
    this.props.changeSignInStatus(STATUS_SIGNING_IN);
    this.closeDialog();
    window.location.href = HOST_PATH + '/auth/v1/login?provider=' + provider;
  }

  // CLOSE DIALOG
  closeDialog() {
    this.props.changeCurrentDialog(Math.random());
  }
  
  renderDialogs() {
    return (
      <>
        <IDSignInOAuthDialog
          title='Sign in with'
          open={this.props.currentDialog === IDC_DIALOG_SIGN_IN}
          onSignIn={provider => this.signInOAuth(provider)}
          onClose={() => this.closeDialog()}
        />
        <IDSessionTerminatedDialog
          title='Session ended'
          open={this.props.currentDialog === IDC_DIALOG_SESSION_TERMINATED}
          onClose={() => this.closeDialog()}
        />
      </>
    );
  }

  renderPage() {
    switch (this.props.page) {
    case PAGE_CONFIGURE:
      return (
        <CfgPageContainer
          api={this.meepScenarioConfigurationApi}
          createScenario={this.cfgCreateScenario}
          setScenario={this.cfgSetScenario}
          deleteScenario={this.cfgDeleteScenario}
          newScenarioElem={this.cfgNewScenarioElem}
          cloneScenarioElem={this.cfgCloneScenarioElem}
          updateScenarioElem={this.cfgUpdateScenarioElem}
          deleteScenarioElem={this.cfgDeleteScenarioElem}
        />
      );

    case PAGE_EXECUTE:
      return (
          <>
            <ExecPageContainer
              api={this.meepActiveScenarioApi}
              eventsApi={this.meepEventsApi}
              automationApi={this.meepEventAutomationApi}
              replayApi={this.meepEventReplayApi}
              cfgApi={this.meepScenarioConfigurationApi}
              sandboxApi={this.meepSandboxControlApi}
              sandbox={this.props.sandbox}
              sandboxes={this.props.sandboxes}
              createSandbox={this.createSandbox}
              setSandbox={this.setSandbox}
              deleteSandbox={this.deleteSandbox}
              refreshScenario={this.refreshScenario}
              deleteScenario={this.execDeleteScenario}
            />
          </>
      );

    case PAGE_SETTINGS:
      return (
        <SettingsPageContainer
          style={{ width: '100%' }}
          startRefresh={() => this.startAutomaticRefresh()}
          stopRefresh={() => this.stopAutomaticRefresh()}
        />
      );

    case PAGE_MONITOR:
      return <MonitorPageContainer />;
    
    case PAGE_HOME:
      return <HomePageContainer />;

    default:
      return null;
    }
  }

  render() {
    return (
      <div className='ui-background'>
        {this.renderDialogs()}

        <MeepTopBar
          title=""
          corePodsRunning={this.props.corePodsRunning}
          corePodsErrors={this.props.corePodsErrors}
          onClickSignIn={() => this.signInProcedure()}
        />

        <ToolbarFixedAdjust />

        <div className='idcc-page'>
          {this.renderPage()}
        </div>

        <Footer />
      </div>
    );
  }
}

const mapStateToProps = state => {
  return {
    cfg: state.cfg,
    cfgVis: state.cfg.vis,
    exec: state.exec,
    execScenarioState: state.exec.state.scenario,
    execVis: state.exec.vis,
    currentDialog: state.ui.currentDialog,
    page: state.ui.page,
    sandbox: state.ui.sandbox,
    sandboxes: state.ui.sandboxes,
    sandboxCfg: state.ui.sandboxCfg,
    automaticRefresh: state.ui.automaticRefresh,
    refreshInterval: state.ui.refreshInterval,
    devMode: state.ui.devMode,
    eventReplayMode: state.ui.eventReplayMode,
    eventCfgMode: state.ui.eventCfgMode,
    corePodsRunning: corePodsRunning(state),
    corePodsErrors: corePodsErrors(state),
    execVisData: execVisFilteredData(state),
    signInStatus: state.ui.signInStatus
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeCurrentDialog: type => dispatch(uiChangeCurrentDialog(type)),
    changeCurrentPage: page => dispatch(uiChangeCurrentPage(page)),
    changeSandbox: name => dispatch(uiExecChangeSandbox(name)),
    changeSandboxList: list => dispatch(uiExecChangeSandboxList(list)),
    changeSandboxCfg: cfg => dispatch(uiExecChangeSandboxCfg(cfg)),
    changeEventCreationMode: mode => dispatch(uiExecChangeEventCreationMode(mode)),
    changeEventReplayMode: mode => dispatch(uiExecChangeEventReplayMode(mode)),
    changeReplayStatus: status => dispatch(execChangeReplayStatus(status)),
    cfgChangeScenario: scenario => dispatch(cfgChangeScenario(scenario)),
    execChangeScenario: scenario => dispatch(execChangeScenario(scenario)),
    execChangeScenarioState: s => dispatch(execChangeScenarioState(s)),
    changeScenarioPodsPhases: phases => dispatch(execChangeScenarioPodsPhases(phases)),
    changeCorePodsPhases: phases => dispatch(execChangeCorePodsPhases(phases)),
    changeServiceMaps: maps => dispatch(execChangeServiceMaps(maps)),
    execChangeVisData: data => dispatch(execChangeVisData(data)),
    execChangeTable: table => dispatch(execChangeTable(table)),
    execChangeMapUeList: list => dispatch(execChangeMapUeList(list)),
    execChangeMapPoaList: list => dispatch(execChangeMapPoaList(list)),
    execChangeMapComputeList: list => dispatch(execChangeMapComputeList(list)),
    cfgChangeMap: map => dispatch(cfgChangeMap(map)),
    cfgChangeVisData: data => dispatch(cfgChangeVisData(data)),
    cfgChangeTable: data => dispatch(cfgChangeTable(data)),
    execChangeOkToTerminate: ok => dispatch(execChangeOkToTerminate(ok)),
    changeSignInStatus: status => dispatch(uiChangeSignInStatus(status)),
    changeSignInUsername: name => dispatch(uiChangeSignInUsername(name)),
    changeTabIndex: index => dispatch(uiChangeCurrentTab(index)),
    changeAppInstanceTable: value => dispatch(execChangeAppInstanceTable(value))
  };
};

const ConnectedMeepContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(MeepContainer);

export default ConnectedMeepContainer;
