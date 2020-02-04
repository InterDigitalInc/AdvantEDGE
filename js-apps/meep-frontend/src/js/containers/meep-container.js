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
import axios from 'axios';
import moment from 'moment';
import { updateObject, deepCopy } from '../util/object-util';

// Import JS dependencies
import * as meepCtrlRestApiClient from '../../../../../js-packages/meep-ctrl-engine-client/src/index.js';

import MeepDrawer from './meep-drawer';
import MeepTopBar from '../components/meep-top-bar';
import CfgPageContainer from './cfg/cfg-page-container';
import ExecPageContainer from './exec/exec-page-container';
import SettingsPageContainer from './settings/settings-page-container';
import MonitorPageContainer from './monitor/monitor-page-container';

import {
  FIELD_PARENT,
  FIELD_NAME,
  getElemFieldVal,
  setElemFieldVal,
  createUniqueName
} from '../util/elem-utils';

import {
  TYPE_CFG,
  TYPE_EXEC,
  EXEC_STATE_DEPLOYED,
  NO_SCENARIO_NAME,
  VIS_VIEW,
  VIEW_NAME_NONE,
  PAGE_CONFIGURE,
  PAGE_EXECUTE,
  PAGE_MONITOR,
  PAGE_SETTINGS,
  COMMON_ZONE_TYPE_STR,
  DEFAULT_NL_TYPE_STR
} from '../meep-constants';

import {
  parseScenario,
  createNewScenario,
  addElementToScenario,
  updateElementInScenario,
  removeElementFromScenario,
  getElementFromScenario
} from '../util/scenario-utils';

import {
  uiChangeCurrentPage,
  uiExecChangeEventCreationMode,
  uiToggleMainDrawer
} from '../state/ui';

import {
  execChangeScenario,
  execChangeScenarioState,
  execChangeScenarioPodsPhases,
  execChangeServiceMaps,
  execChangeVisData,
  execChangeTable,
  execChangeCorePodsPhases,
  execChangeOkToTerminate,
  corePodsRunning,
  corePodsErrors,
  execVisFilteredData,
  execAddMetricsEpoch
} from '../state/exec';

import {
  cfgChangeScenario,
  cfgChangeVisData,
  cfgChangeTable
} from '../state/cfg';

import { idlog } from '../util/functional';

// MEEP Controller REST API JS client
var basepath = 'http://' + location.host + location.pathname + 'v1';
// const basepath = 'http://10.3.16.137:30000/v1';

const metricsBasePath = 'http://' + location.hostname + ':30008/v1';
// const metricsBasePath = 'http://10.3.16.137:30008/v1';

const TIME_FORMAT = moment.HTML5_FMT.DATETIME_LOCAL_MS;

meepCtrlRestApiClient.ApiClient.instance.basePath = basepath.replace(
  /\/+$/,
  ''
);

class MeepContainer extends Component {
  constructor(props) {
    super(props);
    this.state = {};
    this.refreshIntervalTimer = null;
    this.meepCfgApi = new meepCtrlRestApiClient.ScenarioConfigurationApi();
    this.meepExecApi = new meepCtrlRestApiClient.ScenarioExecutionApi();
    this.metricsPollingEnabled = false;
  }

  componentDidMount() {
    document.title = 'AdvantEDGE';
    this.props.changeEventCreationMode(false);
    this.refreshScenario();
    if (this.props.automaticRefresh) {
      this.startAutomaticRefresh();
    }
    this.startRefreshCycle();
    this.setMetricsPolling();
  }

  startRefreshCycle() {
    this.startPodsPhasesPeriodicCheck();
    this.monitorTabFocus();
  }

  startPodsPhasesPeriodicCheck() {
    this.podsPhasesIntervalTimer = setInterval(
      () => this.checkPodsPhases(),
      1000
    );
  }

  stopCorePodsPhasesPeriodicCheck() {
    clearInterval(this.podsPhasesIntervalTimer);
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
        this.stopCorePodsPhasesPeriodicCheck();
        this.stopAutomaticRefresh();
      } else {
        this.startPodsPhasesPeriodicCheck();

        if (this.props.automaticRefresh) {
          this.startAutomaticRefresh();
        }
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

  fetchMetrics() {
    const delta = -7;
    const startTime = moment()
      .utc()
      .add(delta, 'seconds')
      .format(TIME_FORMAT);
    const stopTime = moment()
      .utc()
      .add(delta + 1, 'seconds')
      .format(TIME_FORMAT);
    return axios
      .get(
        `${metricsBasePath}/metrics?startTime=${startTime}&stopTime=${stopTime}`
      )
      .then(res => {
        let epoch = {
          data: res.data.logResponse || [],
          startTime: startTime
        };

        this.props.addMetricsEpoch(epoch);
      })
      .catch(e => {
        idlog('Error while fetching metrics')(e);
      });
  }

  setMetricsPolling() {
    if (
      this.props.dashboardView1 === VIS_VIEW &&
      this.props.dashboardView2 === VIEW_NAME_NONE
    ) {
      this.stopMetricsPolling();
    } else {
      this.startMetricsPolling();
    }
  }
  stopMetricsPolling() {
    if (this.metricsPollingEnabled) {
      clearInterval(this.dataTimer);
      this.metricsPollingEnabled = false;
    }
  }
  startMetricsPolling() {
    if (!this.metricsPollingEnabled) {
      this.epochCount = 0;
      const nextData = () => {
        this.epochCount += 1;
        this.fetchMetrics();
      };
      this.dataTimer = setInterval(nextData, 1000);
      this.metricsPollingEnabled = true;
    }
  }

  checkPodsPhases() {
    // Core pods
    axios
      .get(`${basepath}/states?long=true&type=core`)
      .then(res => {
        this.props.changeCorePodsPhases(res.data.podStatus);
      })
      .catch(() => {
        this.props.changeCorePodsPhases([]);
      });

    // Scenario pods
    axios
      .get(`${basepath}/states?long=true&type=scenario`)
      .then(res => {
        var scenarioPodsPhases = res.data.podStatus;
        this.props.changeScenarioPodsPhases(scenarioPodsPhases);
      })
      .catch(() => {
        this.props.changeScenarioPodsPhases([]);
      });

    // Service maps
    axios
      .get(`${basepath}/active/serviceMaps`)
      .then(res => {
        var serviceMaps = res.data;
        this.props.changeServiceMaps(serviceMaps);
      })
      .catch(() => {
        this.props.changeServiceMaps([]);
      });
  }

  setMainContent(targetId) {
    this.props.changeCurrentPage(targetId);
  }

  // Periodic visualization update handler
  refreshMeepController() {
    if (this.props.page === PAGE_EXECUTE && this.props.automaticRefresh) {
      this.refreshScenario();
    }
  }

  startAutomaticRefresh() {
    _.defer(() => {
      var value = this.props.refreshInterval;
      clearInterval(this.refreshIntervalTimer);
      if (!isNaN(value) && value >= 500 && value <= 60000) {
        this.refreshIntervalTimer = setInterval(
          () => this.refreshMeepController(),
          value
        );
      }
    });
  }

  stopAutomaticRefresh() {
    clearInterval(this.refreshIntervalTimer);
  }

  /**
   * Callback function to receive the result of the getActiveScenario operation.
   * @callback module:api/ScenarioExecutionApi~getActiveScenarioCallback
   * @param {String} error Error message, if any.
   * @param {module:model/Scenario} data The data returned by the service call.
   */
  getActiveScenarioCb(error, data) {
    if (error !== null) {
      // console.log(error);
      // TODO consider showing an alert
      return;
    }

    if (!data.deployment) {
      return;
    }

    // Store & Process deployed scenario
    this.setScenario(TYPE_EXEC, data);

    // TODO set a timer of 2 seconds
    this.props.execChangeScenarioState(EXEC_STATE_DEPLOYED);
    setTimeout(() => {
      this.props.execChangeOkToTerminate(true);
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
    var page = pageType === TYPE_CFG ? this.props.cfg : this.props.exec;
    var parsedScenario = parseScenario(page.scenario);
    var updatedVisData = updateObject(page.vis.data, parsedScenario.visData);
    var updatedTable = updateObject(page.table, parsedScenario.table);

    // Dispatch state updates
    if (pageType === TYPE_CFG) {
      this.props.cfgChangeVisData(updatedVisData);
      this.props.cfgChangeTable(updatedTable);

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
  createScenario(pageType, name) {
    var scenario = createNewScenario(name);
    this.updateScenario(pageType, scenario, true);
  }

  // Set & process scenario
  setScenario(pageType, scenario) {
    this.updateScenario(pageType, scenario, true);
  }

  // Delete & process scenario
  deleteScenario(pageType) {
    var scenario = createNewScenario(NO_SCENARIO_NAME);
    this.updateScenario(pageType, scenario, true);
  }

  // Refresh Active scenario
  refreshScenario() {
    this.meepExecApi.getActiveScenario((error, data) =>
      this.getActiveScenarioCb(error, data)
    );
  }

  // Add new element to scenario
  newScenarioElem(pageType, element, scenarioUpdate) {
    var scenario =
      pageType === TYPE_CFG
        ? this.props.cfg.scenario
        : this.props.exec.scenario;
    var updatedScenario = updateObject({}, scenario);
    addElementToScenario(updatedScenario, element);
    if (scenarioUpdate) {
      this.changeScenario(pageType, updatedScenario);
    }
  }

  // Update element in scenario
  updateScenarioElem(pageType, element) {
    var scenario =
      pageType === TYPE_CFG
        ? this.props.cfg.scenario
        : this.props.exec.scenario;
    var updatedScenario = updateObject({}, scenario);
    updateElementInScenario(updatedScenario, element);
    this.changeScenario(pageType, updatedScenario);
  }

  // Delete element in scenario (also deletes child elements)
  deleteScenarioElem(pageType, element) {
    var scenario =
      pageType === TYPE_CFG
        ? this.props.cfg.scenario
        : this.props.exec.scenario;
    var updatedScenario = updateObject({}, scenario);
    removeElementFromScenario(updatedScenario, element);
    this.changeScenario(pageType, updatedScenario);
  }

  // Clone element in scenario
  cloneScenarioElem(element) {

    var inDomainCloneBranch = false, inZoneCloneBranch = false, inNlCloneBranch = false, inPlCloneBranch = false;
    var newZoneRootParentName = '';
    var newNlRootParentName = '';
    var newPlRootParentName = '';
    var newProcessRootParentName = '';
    var elementFromScenario;

    var scenario = this.props.cfg.scenario;
    // Domains
    for (var i in scenario.deployment.domains) {
      var domain = scenario.deployment.domains[i];

      // Add domain to graph and table (ignore public domain)
      if (domain.id === element.id) {
        newZoneRootParentName = this.cloneElement(element, getElemFieldVal(element, FIELD_PARENT), true);
        inDomainCloneBranch = true;
      } else {
        inDomainCloneBranch = false;
      }

      // Zones
      for (var j in domain.zones) {
        var zone = domain.zones[j];

        if (inDomainCloneBranch) {
          if (zone.name.indexOf(COMMON_ZONE_TYPE_STR) !== -1) {
            newNlRootParentName = newZoneRootParentName + COMMON_ZONE_TYPE_STR;
          } else {
            elementFromScenario = getElementFromScenario(scenario, zone.id);
            newNlRootParentName = this.cloneElement(elementFromScenario, newZoneRootParentName, false);
          }
        } else {
          if (zone.id === element.id) {
            newNlRootParentName = this.cloneElement(element, getElemFieldVal(element, FIELD_PARENT), true);
            inZoneCloneBranch = true;
          } else {
            inZoneCloneBranch = false;
          }
        }

        // Network Locations
        for (var k in zone.networkLocations) {
          var nl = zone.networkLocations[k];

          if (inDomainCloneBranch || inZoneCloneBranch) {
            if (nl.name.indexOf(DEFAULT_NL_TYPE_STR) !== -1) {
              newPlRootParentName = newNlRootParentName;
            } else {
              elementFromScenario = getElementFromScenario(scenario, nl.id);
              newPlRootParentName = this.cloneElement(elementFromScenario, newNlRootParentName, false);
            }
          } else {
            if (nl.id === element.id) {
              newPlRootParentName = this.cloneElement(element, getElemFieldVal(element, FIELD_PARENT, true));
              inNlCloneBranch = true;
            } else {
              inNlCloneBranch = false;
            }
          }

          // Physical Locations
          for (var l in nl.physicalLocations) {
            var pl = nl.physicalLocations[l];

            if (inDomainCloneBranch || inZoneCloneBranch || inNlCloneBranch) {
              elementFromScenario = getElementFromScenario(scenario, pl.id);
              newProcessRootParentName = this.cloneElement(elementFromScenario, newPlRootParentName, false);
            } else {
              if (pl.id === element.id) {
                newProcessRootParentName = this.cloneElement(element, getElemFieldVal(element, FIELD_PARENT, true));
                inPlCloneBranch = true;
              } else {
                inPlCloneBranch = false;
              }
            }

            // Processes
            for (var m in pl.processes) {
              var proc = pl.processes[m];

              if (inDomainCloneBranch || inZoneCloneBranch || inNlCloneBranch || inPlCloneBranch) {
                elementFromScenario = getElementFromScenario(scenario, proc.id);
                this.cloneElement(elementFromScenario, newProcessRootParentName, false);
              } else {
                if (proc.id === element.id) {
                  this.cloneElement(element, getElemFieldVal(element, FIELD_PARENT, true));
                }
              }
            }
          }
        }
      }

      if(inDomainCloneBranch || inZoneCloneBranch || inNlCloneBranch || inPlCloneBranch) {
        break;
      }
    }
    this.changeScenario(TYPE_CFG, scenario);
  }

  // CLONE ELEMENT, return new element name
  cloneElement(element, newParentName, isRoot) {
    let newElement = deepCopy(element);

    var name = getElemFieldVal(element, FIELD_NAME);
    if (isRoot === false) {
      name = createUniqueName(this.props.cfg.table.entries, name + '-copy');
      setElemFieldVal(newElement, FIELD_NAME, name);
    }
    setElemFieldVal(newElement, FIELD_PARENT, newParentName);

    // add new element to scenario
    // new id and label will be created as part of the addNewElementToScenario called by newScenarioElem
    this.newScenarioElem(TYPE_CFG, newElement, false);
    return name;
  }

  renderPage() {
    switch (this.props.page) {
    case PAGE_CONFIGURE:
      return (
        <CfgPageContainer
          style={{ width: '100%' }}
          api={this.meepCfgApi}
          createScenario={name => {
            this.createScenario(TYPE_CFG, name);
          }}
          setScenario={scenario => {
            this.setScenario(TYPE_CFG, scenario);
          }}
          deleteScenario={() => {
            this.deleteScenario(TYPE_CFG);
          }}
          newScenarioElem={(elem, update) => {
            this.newScenarioElem(TYPE_CFG, elem, update);
          }}
          cloneScenarioElem={elem => {
            this.cloneScenarioElem(elem);
          }}
          updateScenarioElem={elem => {
            this.updateScenarioElem(TYPE_CFG, elem);
          }}
          deleteScenarioElem={elem => {
            this.deleteScenarioElem(TYPE_CFG, elem);
          }}
        />
      );

    case PAGE_EXECUTE:
      return (
          <>
            <ExecPageContainer
              style={{ width: '100%' }}
              api={this.meepExecApi}
              cfgApi={this.meepCfgApi}
              refreshScenario={() => {
                this.refreshScenario();
              }}
              deleteScenario={() => {
                this.deleteScenario(TYPE_EXEC);
              }}
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
      return <MonitorPageContainer style={{ paddingRight: '100%' }} />;

    default:
      return null;
    }
  }

  render() {
    const flexString = this.props.mainDrawerOpen ? '0 0 250px' : '0 0 0px';
    this.setMetricsPolling();

    return (
      <div style={{ width: '100%' }}>
        <MeepTopBar
          title=""
          toggleMainDrawer={() => this.props.toggleMainDrawer()}
          corePodsRunning={this.props.corePodsRunning}
          corePodsErrors={this.props.corePodsErrors}
        />

        <div style={{ display: 'flex' }}>
          <div
            className="component-style"
            style={{ overflow: 'hidden', flex: flexString, marginTop: '5px' }}
          >
            <MeepDrawer open={this.props.mainDrawerOpen} style={{ flex: 1 }} />
          </div>
          <div style={{ flex: '1', padding: 10 }}>{this.renderPage()}</div>
        </div>
      </div>
    );
  }
}

const mapStateToProps = state => {
  return {
    cfg: state.cfg,
    cfgVis: state.cfg.vis,
    exec: state.exec,
    execVis: state.exec.vis,
    page: state.ui.page,
    automaticRefresh: state.ui.automaticRefresh,
    refreshInterval: state.ui.refreshInterval,
    devMode: state.ui.devMode,
    mainDrawerOpen: state.ui.mainDrawerOpen,
    dashboardView1: state.ui.dashboardView1,
    dashboardView2: state.ui.dashboardView2,
    corePodsRunning: corePodsRunning(state),
    corePodsErrors: corePodsErrors(state),
    execVisData: execVisFilteredData(state)
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeCurrentPage: page => dispatch(uiChangeCurrentPage(page)),
    changeEventCreationMode: mode =>
      dispatch(uiExecChangeEventCreationMode(mode)),
    cfgChangeScenario: scenario => dispatch(cfgChangeScenario(scenario)),
    execChangeScenario: scenario => dispatch(execChangeScenario(scenario)),
    execChangeScenarioState: s => dispatch(execChangeScenarioState(s)),
    changeScenarioPodsPhases: phases =>
      dispatch(execChangeScenarioPodsPhases(phases)),
    changeCorePodsPhases: phases => dispatch(execChangeCorePodsPhases(phases)),
    changeServiceMaps: maps => dispatch(execChangeServiceMaps(maps)),
    execChangeVisData: data => dispatch(execChangeVisData(data)),
    execChangeTable: table => dispatch(execChangeTable(table)),
    cfgChangeVisData: data => dispatch(cfgChangeVisData(data)),
    cfgChangeTable: data => dispatch(cfgChangeTable(data)),
    execChangeOkToTerminate: ok => dispatch(execChangeOkToTerminate(ok)),
    toggleMainDrawer: () => dispatch(uiToggleMainDrawer()),
    addMetricsEpoch: epoch => dispatch(execAddMetricsEpoch(epoch))
  };
};

const ConnectedMeepContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(MeepContainer);

export default ConnectedMeepContainer;
