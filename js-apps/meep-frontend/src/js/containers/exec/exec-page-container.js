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
import { Grid, GridCell, GridInner } from '@rmwc/grid';
import { Elevation } from '@rmwc/elevation';
import DashboardContainer from './dashboard-container';
import EventContainer from './event-container';
import ExecPageScenarioButtons from './exec-page-scenario-buttons';

import HeadlineBar from '../../components/headline-bar';
import EventCreationPane from './event-creation-pane';
import EventReplayPane from './event-replay-pane';

import ExecTable from './exec-table';

import IDDeployScenarioDialog from '../../components/dialogs/id-deploy-scenario-dialog';
import IDTerminateScenarioDialog from '../../components/dialogs/id-terminate-scenario-dialog';
import IDSaveScenarioDialog from '../../components/dialogs/id-save-scenario-dialog';
import IDSaveReplayDialog from '../../components/dialogs/id-save-replay-dialog';


import { execChangeScenarioList, execVisFilteredData } from '../../state/exec';

import {
  uiChangeCurrentDialog,
  uiExecChangeEventCreationMode,
  uiExecChangeEventReplayMode,
  uiExecChangeDashCfgMode,
  uiExecChangeEventCfgMode,
  uiExecChangeCurrentEvent,
  uiExecChangeShowApps,
  uiExecChangeReplayFilesList
} from '../../state/ui';

import {
  execChangeScenario,
  execChangeScenarioName,
  execChangeScenarioState,
  execChangeOkToTerminate
} from '../../state/exec';

import {
  // States
  EXEC_STATE_IDLE,
  PAGE_EXECUTE,
  IDC_DIALOG_DEPLOY_SCENARIO,
  IDC_DIALOG_TERMINATE_SCENARIO,
  IDC_DIALOG_SAVE_SCENARIO,
  IDC_DIALOG_SAVE_REPLAY,
  MOBILITY_EVENT,
  NETWORK_CHARACTERISTICS_EVENT
} from '../../meep-constants';

class ExecPageContainer extends Component {
  constructor(props) {
    super(props);
  }

  componentDidMount() {
    this.props.changeCurrentEvent(MOBILITY_EVENT);
  }

  /**
   * Callback function to receive the result of the getScenarioList operation.
   * @callback module:api/ScenarioConfigurationApi~getScenarioListCallback
   * @param {String} error Error message, if any.
   * @param {module:model/ScenarioList} data The data returned by the service call.
   */
  getScenarioListDeployCb(error, data) {
    if (error !== null) {
      // TODO: consider showing an alert/toast
      return;
    }

    this.props.changeDeployScenarioList(_.map(data.scenarios, 'name'));
  }

  /**
   * Callback function to receive the result of the activateScenario operation.
   * @callback module:api/ScenarioExecutionApi~activateScenarioCallback
   * @param {String} error Error message, if any.
   */
  activateScenarioCb(error) {
    if (error) {
      // TODO: consider showing an alert/toast
      return;
    }

    this.props.refreshScenario();
  }

  /**
   * Callback function to receive the result of the terminateScenario operation.
   * @callback module:api/ScenarioExecutionApi~terminateScenarioCallback
   * @param {String} error Error message, if any.
   */
  terminateScenarioCb(error) {
    if (error !== null) {
      // TODO consider showing an alert  (i.e. toast)
      return;
    }

    this.props.deleteScenario();
    this.props.changeState(EXEC_STATE_IDLE);
    this.props.execChangeOkToTerminate(false);
  }

  /**
   * Callback function to receive the result of the createScenario operation.
   * @callback module:api/ScenarioConfigurationApi~createScenarioCallback
   * @param {String} error Error message, if any.
   * @param data This operation does not return a value.
   * @param {String} response The complete HTTP response.
   */
  createScenarioCb(/*error, data, response*/) {
    // if (error == null) {
    //   console.log('Scenario successfully created');
    // } else {
    //   console.log('Failed to create scenario');
    // }
    // TODO: consider showing an alert/toast
  }

  /**
   * Callback function to receive the result of the getReplayList operation.
   * @callback module:api/EventReplayApi~getReplayFileListCallback
   * @param {String} error Error message, if any.
   * @param {module:model/ReplayFileList} data The data returned by the service call.
   */
  getReplayFileListCb(error, data) {
    if (error !== null) {
      // TODO: consider showing an alert/toast
      return;
    }
    let replayFiles = data.replayFiles;
    replayFiles.unshift('None');
    this.props.changeReplayFilesList(replayFiles);
  }

  saveScenario(scenarioName) {
    const scenario = this.props.scenario;

    const scenarioCopy = JSON.parse(JSON.stringify(scenario));
    scenarioCopy.name = scenarioName;

    this.props.cfgApi.createScenario(
      scenarioName,
      scenarioCopy,
      (error, data, response) => this.createScenarioCb(error, data, response)
    );
  }

  updateReplayFileList() {
    this.props.replayApi.getReplayFileList((error, data, response) => {
      this.getReplayFileListCb(error, data, response);
    });
  }
  
  saveReplay(state) {
    const scenarioName = this.props.scenario.name;
    var replayInfo = {
      scenarioName: '',
      description: ''
    };

    replayInfo.scenarioName = scenarioName;
    replayInfo.description = state.description;

    this.props.replayApi.createReplayFileFromScenarioExec(state.replayName, replayInfo, (error) => {
      if (error) {
        // TODO consider showing an alert
        // console.log(error);
      }
      // Refresh file list
      this.updateReplayFileList();
    });
  }

  // CLOSE DIALOG
  closeDialog() {
    this.props.changeCurrentDialog(Math.random());
  }

  // DEPLOY DIALOG
  onDeployScenario() {
    // Retrieve list of available scenarios
    this.props.cfgApi.getScenarioList((error, data, response) => {
      this.getScenarioListDeployCb(error, data, response);
    });
    this.props.changeCurrentDialog(IDC_DIALOG_DEPLOY_SCENARIO);
  }

  // SAVE SCENARIO
  onSaveScenario() {
    this.props.changeCurrentDialog(IDC_DIALOG_SAVE_SCENARIO);
  }

  // TERMINATE DIALOG
  onTerminateScenario() {
    this.props.changeCurrentDialog(IDC_DIALOG_TERMINATE_SCENARIO);
  }

  // SAVE REPLAY FILE
  onSaveReplay() {
    this.props.changeCurrentDialog(IDC_DIALOG_SAVE_REPLAY);
  }

  // SHOW REPLAY PANE
  onShowReplay() {
    this.updateReplayFileList();
  }

  // CLOSE CREATE EVENT PANE
  onQuitEventCreationMode() {
    this.props.changeEventCreationMode(false);
  }

  // CLOSE REPLAY EVENT PANE
  onQuitEventReplayMode() {
    this.props.changeEventReplayMode(false);
  }

  // CONFIGURE DASHBOARD
  onOpenDashCfg() {
    this.props.changeDashCfgMode(true);
  }

  // STOP CONFIGURE DASHBOARD
  onCloseDashCfg() {
    this.props.changeDashCfgMode(false);
  }

  // CONFIGURE EVENTS
  onOpenEventCfg() {
    this.props.changeEventCfgMode(true);
  }

  // STOP CONFIGURE EVENTS
  onCloseEventCfg() {
    this.props.changeEventCfgMode(false);
  }

  // Terminate Active scenario
  terminateScenario() {
    this.props.api.terminateScenario((error, data, response) =>
      this.terminateScenarioCb(error, data, response)
    );
  }

  showApps(show) {
    this.props.changeShowApps(show);
    _.defer(() => {
      this.props.execVis.network.setData(this.props.execVisData);
    });
  }

  renderDialogs() {
    return (
      <>
        <IDDeployScenarioDialog
          title="Open Scenario"
          open={this.props.currentDialog === IDC_DIALOG_DEPLOY_SCENARIO}
          options={this.props.scenarios}
          onClose={() => {
            this.closeDialog();
          }}
          api={this.props.api}
          activateScenarioCb={(error, data, response) =>
            this.activateScenarioCb(error, data, response)
          }
        />

        <IDSaveScenarioDialog
          title="Save Scenario as ..."
          open={this.props.currentDialog === IDC_DIALOG_SAVE_SCENARIO}
          onClose={() => {
            this.closeDialog();
          }}
          api={this.props.api}
          saveScenario={name => this.saveScenario(name)}
          scenarioNameRequired={true}
        />

        <IDTerminateScenarioDialog
          title="Terminate Scenario"
          open={this.props.currentDialog === IDC_DIALOG_TERMINATE_SCENARIO}
          scenario={this.props.scenario}
          onClose={() => {
            this.closeDialog();
          }}
          onSubmit={() => {
            this.terminateScenario();
          }}
        />

        <IDSaveReplayDialog
          title="Save Events as Replay file"
          open={this.props.currentDialog === IDC_DIALOG_SAVE_REPLAY}
          onClose={() => {
            this.closeDialog();
          }}
          api={this.props.replayApi}
          saveReplay={replayInfo => this.saveReplay(replayInfo)}
          replayNameRequired={true}
        />

      </>
    );
  }

  render() {
    if (this.props.page !== PAGE_EXECUTE) {
      return null;
    }

    const scenarioName =
      this.props.page === PAGE_EXECUTE
        ? this.props.execScenarioName
        : this.props.cfgScenarioName;

    const spanLeft = (this.props.eventCreationMode || this.props.eventReplayMode) ? 9 : 12;
    const spanRight = (this.props.eventCreationMode || this.props.eventReplayMode) ? 3 : 0;
    return (
      <div style={{ width: '100%' }}>
        {this.renderDialogs()}

        <div style={{ width: '100%' }}>
          <Grid style={styles.headlineGrid}>
            <GridCell span={12}>
              <Elevation
                className="component-style"
                z={2}
                style={styles.headline}
              >
                <GridInner>
                  <GridCell align={'middle'} span={4}>
                    <HeadlineBar
                      titleLabel="Deployed Scenario"
                      scenarioName={scenarioName}
                    />
                  </GridCell>
                  <GridCell align={'middle'} span={8}>
                    <GridInner align={'right'}>
                      <GridCell align={'middle'} span={12}>
                        <ExecPageScenarioButtons
                          onDeploy={() => this.onDeployScenario()}
                          onSaveScenario={() => this.onSaveScenario()}
                          onTerminate={() => this.onTerminateScenario()}
                          onRefresh={this.props.refreshScenario}
                          onOpenDashCfg={() => this.onOpenDashCfg()}
                          onOpenEventCfg={() => this.onOpenEventCfg()}
                        />
                      </GridCell>
                    </GridInner>
                  </GridCell>
                </GridInner>
              </Elevation>
            </GridCell>
          </Grid>
        </div>

        {this.props.exec.state.scenario !== EXEC_STATE_IDLE && (
          <>
            <Grid style={{ width: '100%' }}>
              <GridCell span={spanLeft}>
                <div>
                  <EventContainer
                    scenarioName={this.props.execScenarioName}
                    eventCfgMode={this.props.eventCfgMode}
                    onCloseEventCfg={() => this.onCloseEventCfg()}
                    onSaveReplay={() => this.onSaveReplay()}
                    onShowReplay={() => this.onShowReplay()}
                    api={this.props.replayApi}
                  />

                  <DashboardContainer
                    scenarioName={this.props.execScenarioName}
                    onShowAppsChanged={show => this.showApps(show)}
                    showApps={this.props.showApps}
                    dashCfgMode={this.props.dashCfgMode}
                    onCloseDashCfg={() => this.onCloseDashCfg()}
                  />
                </div>
              </GridCell>
              <GridCell
                span={spanRight}
                hidden={!this.props.eventCreationMode && !this.props.eventReplayMode}
                style={styles.inner}
              >
                <Elevation className="component-style" z={2}>
                  <EventReplayPane
                    api={this.props.replayApi}
                    hide={!this.props.eventReplayMode}
                    onClose={() => this.onQuitEventReplayMode()}
                  />
                </Elevation>
                <Elevation className="component-style" z={2}>
                  <EventCreationPane
                    eventTypes={[MOBILITY_EVENT, NETWORK_CHARACTERISTICS_EVENT]}
                    api={this.props.api}
                    hide={!this.props.eventCreationMode}
                    onSuccess={() => {
                      this.props.refreshScenario();
                    }}
                    onClose={() => this.onQuitEventCreationMode()}
                  />
                </Elevation>
              </GridCell>
            </Grid>
          </>
        )}
        <ExecTable />
      </div>
    );
  }
}

const styles = {
  headlineGrid: {
    marginBottom: 10
  },
  headline: {
    padding: 10
  },
  page: {
    height: 1500,
    marginBottom: 10
  }
};

const mapStateToProps = state => {
  return {
    exec: state.exec,
    showApps: state.ui.execShowApps,
    execVis: state.exec.vis,
    configuredElement: state.cfg.elementConfiguration.configuredElement,
    table: state.exec.table,
    currentDialog: state.ui.currentDialog,
    scenario: state.exec.scenario,
    scenarios: state.exec.apiResults.scenarios,
    eventCreationMode: state.ui.eventCreationMode,
    eventReplayMode: state.ui.eventReplayMode,
    dashCfgMode: state.ui.dashCfgMode,
    eventCfgMode: state.ui.eventCfgMode,
    page: state.ui.page,
    execScenarioName: state.exec.scenario.name,
    cfgScenarioName: state.cfg.scenario.name,
    execVisData: execVisFilteredData(state)
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeCurrentDialog: type => dispatch(uiChangeCurrentDialog(type)),
    changeScenario: scenario => dispatch(execChangeScenario(scenario)),
    changeDeployScenarioList: scenarios => dispatch(execChangeScenarioList(scenarios)),
    changeScenarioName: name => dispatch(execChangeScenarioName(name)),
    changeState: s => dispatch(execChangeScenarioState(s)),
    changeEventCreationMode: val => dispatch(uiExecChangeEventCreationMode(val)), // (true or false)
    changeEventReplayMode: val => dispatch(uiExecChangeEventReplayMode(val)), // (true or false)
    changeDashCfgMode: val => dispatch(uiExecChangeDashCfgMode(val)), // (true or false)
    changeEventCfgMode: val => dispatch(uiExecChangeEventCfgMode(val)), // (true or false)
    changeCurrentEvent: e => dispatch(uiExecChangeCurrentEvent(e)),
    execChangeOkToTerminate: ok => dispatch(execChangeOkToTerminate(ok)),
    changeShowApps: show => dispatch(uiExecChangeShowApps(show)),
    changeReplayFilesList: list => dispatch(uiExecChangeReplayFilesList(list))
  };
};

const ConnectedExecPageContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(ExecPageContainer);

export default ConnectedExecPageContainer;
