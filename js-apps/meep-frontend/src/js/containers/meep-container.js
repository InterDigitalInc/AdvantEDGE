/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
import _ from 'lodash';
import { connect } from 'react-redux';
import React, { Component }  from 'react';
import axios from 'axios';
import { updateObject } from '../util/object-util';

// Import JS dependencies
import * as meepCtrlRestApiClient from '../../../../../js-packages/meep-ctrl-engine-client/src/index.js';

import MeepDrawer from './meep-drawer';
import MeepTopBar from '../components/meep-top-bar';
import CfgPageContainer from './cfg/cfg-page-container';
import ExecPageContainer from './exec/exec-page-container';
import SettingsPageContainer from './settings/settings-page-container';
import MonitorPageContainer from './monitor/monitor-page-container';

import {
  TYPE_CFG,
  TYPE_EXEC,
  EXEC_STATE_DEPLOYED,
  NO_SCENARIO_NAME
} from '../meep-constants';

import {
  parseScenario,
  createNewScenario,
  addElementToScenario,
  updateElementInScenario,
  removeElementFromScenario
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
  execVisFilteredData
} from '../state/exec';

import {
  cfgChangeScenario,
  cfgChangeVisData,
  cfgChangeTable
} from '../state/cfg';

import {
  PAGE_CONFIGURE,
  PAGE_EXECUTE,
  PAGE_EXPERIMENTAL_EXECUTE,
  PAGE_MONITOR,
  PAGE_SETTINGS
} from '../state/ui';

// MEEP Controller REST API JS client
// var basepath = 'http://' + location.host + location.pathname + 'v1';
var basepath = 'http://10.3.16.73:30000/v1';

meepCtrlRestApiClient.ApiClient.instance.basePath = basepath.replace(/\/+$/, '');

class MeepContainer extends Component {
  constructor(props) {
    super(props);
    this.state = {};
    this.refreshIntervalTimer = null;
    this.meepSettingsApi = new meepCtrlRestApiClient.MEEPSettingsApi();
    this.meepCfgApi = new meepCtrlRestApiClient.ScenarioConfigurationApi();
    this.meepExecApi = new meepCtrlRestApiClient.ScenarioExecutionApi();
  }

  componentDidMount() {
    document.title = 'AdvantEDGE';
    this.props.changeEventCreationMode(false);
    this.refreshScenario();
    // -- Migration -- //
    this.props.changeCurrentPage(PAGE_CONFIGURE);
    this.startRefreshCycle();
  }

  startRefreshCycle() {
    this.startPodsPhasesPeriodicCheck();
    this.monitorTabFocus();
  }

  startPodsPhasesPeriodicCheck() {
    this.podsPhasesIntervalTimer = setInterval(() => this.checkPodsPhases(), 1000);
  }

  stopCorePodsPhasesPeriodicCheck() {
    clearInterval(this.podsPhasesIntervalTimer);
  }

  monitorTabFocus() {
    var hidden, visibilityChange;
    if (typeof document.hidden !== 'undefined') { // Opera 12.10 and Firefox 18 and later support
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
    if (typeof document.addEventListener === 'undefined' || hidden === undefined) {
      // TODO: consider showing an alert
      // console.log('This demo requires a browser, such as Google Chrome or Firefox, that supports the Page Visibility API.');
    } else {
      // Handle page visibility change
      document.addEventListener(visibilityChange, handleVisibilityChange, false);
    }
  }

  // -- Migration --
  checkPodsPhases() {
    // Core pods
    axios.get(`${basepath}/states?long=true&type=core`)
      .then(res => {
        this.props.changeCorePodsPhases(res.data.podStatus);
      }).catch(() => {
        this.props.changeCorePodsPhases([]);
      });

    // Scenario pods
    axios.get(`${basepath}/states?long=true&type=scenario`)
      .then(res => {
        var scenarioPodsPhases = res.data.podStatus;
        this.props.changeScenarioPodsPhases(scenarioPodsPhases);
      }).catch(() => {
        this.props.changeScenarioPodsPhases([]);
      });

    // Service maps
    axios.get(`${basepath}/active/serviceMaps`)
      .then(res => {
        var serviceMaps = res.data;
        this.props.changeServiceMaps(serviceMaps);
      }).catch(() => {
        this.props.changeServiceMaps([]);
      });
  }

  setMainContent(targetId) {
    this.props.changeCurrentPage(targetId);
  }

  // Periodic visualization update handler
  refreshMeepController() {
    if ((this.props.page === PAGE_EXECUTE) && this.props.automaticRefresh) {
      this.refreshScenario();
    }
  }

  startAutomaticRefresh() {
    _.defer(() => {
      var value = this.props.refreshInterval;
      if (isNaN(value) || value < 500 || value > 60000) {
        clearInterval(this.refreshIntervalTimer);
      } else {
        clearInterval(this.refreshIntervalTimer);
        this.refreshIntervalTimer = setInterval(() => this.refreshMeepController(), value);
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

  // Change & process scenario
  changeScenario(pageType, scenario) {
    // Change scenario state
    if (pageType === TYPE_CFG) {
      this.props.cfgChangeScenario(scenario);
    } else {
      this.props.execChangeScenario(scenario);
    }

    // Parse Scenario object to retrieve visualization data and scenario table 
    var page = (pageType === TYPE_CFG) ? this.props.cfg : this.props.exec;
    var parsedScenario = parseScenario(page.scenario);
    var updatedVisData = updateObject(page.vis.data, parsedScenario.visData);
    var updatedTable = updateObject(page.table, parsedScenario.table);

    // Dispatch state updates
    if (pageType === TYPE_CFG) {
      this.props.cfgChangeVisData(updatedVisData);
      this.props.cfgChangeTable(updatedTable);

      const vis = this.props.cfgVis;
      if (vis && vis.network && vis.network.setData) {
        vis.network.setData(updatedVisData);
      }
    } else {
      this.props.execChangeVisData(updatedVisData);
      this.props.execChangeTable(updatedTable);

      const vis = this.props.execVis;
      if (vis && vis.network && vis.network.setData) {
        _.defer(() => {
          vis.network.setData(this.props.execVisData);
        });
      }
    }
  }

  // Create, store & process new scenario
  createScenario(pageType, name) {
    var scenario = createNewScenario(name);
    this.changeScenario(pageType, scenario);
  }

  // Set & process scenario
  setScenario(pageType, scenario) {
    this.changeScenario(pageType, scenario);
  }

  // Delete & process scenario
  deleteScenario(pageType) {
    var scenario = createNewScenario(NO_SCENARIO_NAME);
    this.changeScenario(pageType, scenario);
  }

  // Refresh Active scenario
  refreshScenario() {
    this.meepExecApi.getActiveScenario((error, data) => this.getActiveScenarioCb(error, data));
  }

  // Add new element to scenario
  newScenarioElem(pageType, element) {
    var scenario = (pageType === TYPE_CFG) ? this.props.cfg.scenario : this.props.exec.scenario;
    var updatedScenario = updateObject({}, scenario);
    addElementToScenario(updatedScenario, element);
    this.changeScenario(pageType, updatedScenario);
  }

  // Update element in scenario
  updateScenarioElem(pageType, element) {
    var scenario = (pageType === TYPE_CFG) ? this.props.cfg.scenario : this.props.exec.scenario;
    var updatedScenario = updateObject({}, scenario);
    updateElementInScenario(updatedScenario, element);
    this.changeScenario(pageType, updatedScenario);
  }

  // Delete element in scenario (also deletes child elements)
  deleteScenarioElem(pageType, element) {
    var scenario = (pageType === TYPE_CFG) ? this.props.cfg.scenario : this.props.exec.scenario;
    var updatedScenario = updateObject({}, scenario);
    removeElementFromScenario(updatedScenario, element);
    this.changeScenario(pageType, updatedScenario);
  }

  renderPage() {
    switch(this.props.page) {
    case PAGE_CONFIGURE:
      return (
        <CfgPageContainer style={{width: '100%'}}
          api={this.meepCfgApi}
          createScenario={(name) => {this.createScenario(TYPE_CFG, name);}}
          setScenario={(scenario) => {this.setScenario(TYPE_CFG, scenario);}}
          deleteScenario={() => {this.deleteScenario(TYPE_CFG);}}
          newScenarioElem={(elem) => {this.newScenarioElem(TYPE_CFG, elem);}}
          updateScenarioElem={(elem) => {this.updateScenarioElem(TYPE_CFG, elem);}}
          deleteScenarioElem={(elem) => {this.deleteScenarioElem(TYPE_CFG, elem);}}
        />
      );

    case PAGE_EXECUTE:
      return (
        <>
        <ExecPageContainer style={{width: '100%'}}
          api={this.meepExecApi}
          cfgApi={this.meepCfgApi}
          refreshScenario={() => {this.refreshScenario();}}
          deleteScenario={() => {this.deleteScenario(TYPE_EXEC);}}
        />
        </>
      );

    case PAGE_EXPERIMENTAL_EXECUTE:
      return (
        <>
        <ExecPageContainer style={{width: '100%'}}
          experimental={true}
          api={this.meepExecApi}
          cfgApi={this.meepCfgApi}
          refreshScenario={() => {this.refreshScenario();}}
          deleteScenario={() => {this.deleteScenario(TYPE_EXEC);}}
        />
        </>
      );

    case PAGE_SETTINGS:
      return (
        <SettingsPageContainer style={{width: '100%'}}
          startRefresh={() => this.startAutomaticRefresh()}
          stopRefresh={() => this.stopAutomaticRefresh()}
        />
      );

    case PAGE_MONITOR:
      return (
        <MonitorPageContainer style={{paddingRight: '100%'}} />
      );

    default:
      return null;
    }
  }

  render() {
    const flexString = this.props.mainDrawerOpen ? '0 0 250px' : '0 0 0px';
    return (
      <div style={{width: '100%'}}>
        <MeepTopBar
          title=""
          toggleMainDrawer={() => this.props.toggleMainDrawer()}
          corePodsRunning={this.props.corePodsRunning}
          corePodsErrors={this.props.corePodsErrors}
        />

        <div style={{display: 'flex'}}>
          <div className="component-style" style={{overflow: 'hidden', flex: flexString, marginTop: '5px'}}>
            <MeepDrawer open={this.props.mainDrawerOpen} style={{flex:1}}/>
          </div>
          <div style={{flex: '1', padding: 10}}>
            {this.renderPage()}
          </div>
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
    corePodsRunning: corePodsRunning(state),
    corePodsErrors: corePodsErrors(state),
    execVisData: execVisFilteredData(state)
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeCurrentPage: (page) => dispatch(uiChangeCurrentPage(page)),
    changeEventCreationMode: (mode) => dispatch(uiExecChangeEventCreationMode(mode)),
    cfgChangeScenario: (scenario) => dispatch(cfgChangeScenario(scenario)),
    execChangeScenario: (scenario) => dispatch(execChangeScenario(scenario)),
    execChangeScenarioState: (s) => dispatch(execChangeScenarioState(s)),
    changeScenarioPodsPhases: (phases) => dispatch(execChangeScenarioPodsPhases(phases)),
    changeCorePodsPhases: (phases) => dispatch(execChangeCorePodsPhases(phases)),
    changeServiceMaps: (maps) => dispatch(execChangeServiceMaps(maps)),
    execChangeVisData: (data) => dispatch(execChangeVisData(data)),
    execChangeTable: (table) => dispatch(execChangeTable(table)),
    cfgChangeVisData: (data) => dispatch(cfgChangeVisData(data)),
    cfgChangeTable: (data) => dispatch(cfgChangeTable(data)),
    execChangeOkToTerminate: (ok) => dispatch(execChangeOkToTerminate(ok)),
    toggleMainDrawer: () => dispatch(uiToggleMainDrawer())
  };
};

const ConnectedMeepContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(MeepContainer);

export default ConnectedMeepContainer;
