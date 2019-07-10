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
import { Grid, GridCell, GridInner } from '@rmwc/grid';
import { Elevation } from '@rmwc/elevation';
import IDCVis from '../idc-vis';
import IDCGraph from '../idc-graph';
import IDCAreaChart from '../idc-area-chart';
import DashboardContainer from '../dashboard-container';
import ExecPageScenarioButtons from './exec-page-scenario-buttons';

import HeadlineBar from '../../components/headline-bar';
import EventCreationPane from './event-creation-pane';
import ExecTable from './exec-table';

import IDDeployScenarioDialog from '../../components/dialogs/id-deploy-scenario-dialog';
import IDTerminateScenarioDialog from '../../components/dialogs/id-terminate-scenario-dialog';
import IDSaveScenarioDialog from '../../components/dialogs/id-save-scenario-dialog';

import {
  execChangeScenarioList,
  execVisFilteredData
} from '../../state/exec';

import {
  uiChangeCurrentDialog,
  uiExecChangeEventCreationMode,
  uiExecChangeCurrentEvent,
  uiExecChangeShowApps,
  IDC_DIALOG_DEPLOY_SCENARIO,
  IDC_DIALOG_TERMINATE_SCENARIO,
  IDC_DIALOG_SAVE_SCENARIO,

  // Event types
  UE_MOBILITY_EVENT,
  NETWORK_CHARACTERISTICS_EVENT,
  PAGE_EXECUTE,
  PAGE_EXPERIMENTAL_EXECUTE
} from '../../state/ui';

import {
  execChangeScenario,
  execChangeScenarioName,
  execChangeScenarioState,
  execChangeOkToTerminate
} from '../../state/exec';

import {
  TYPE_EXEC,

  // States
  EXEC_STATE_IDLE
} from '../../meep-constants';


class ExecPageContainer extends Component {
  constructor(props) {
    super(props);
  }

  componentDidMount() {
    this.props.changeCurrentEvent(UE_MOBILITY_EVENT);
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

  saveScenario(scenarioName) {
    const scenario = this.props.scenario;

    const scenarioCopy = JSON.parse(JSON.stringify(scenario));
    scenarioCopy.name = scenarioName;

    this.props.cfgApi.createScenario(scenarioName, scenarioCopy, (error, data, response) => this.createScenarioCb(error, data, response));
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

  // CLOSE DIALOG
  closeDialog() {
    this.props.changeCurrentDialog(Math.random());
  }

  // DEPLOY DIALOG
  onDeployScenario() {
    // Retrieve list of available scenarios
    this.props.cfgApi.getScenarioList((error, data, response) => {this.getScenarioListDeployCb(error, data, response);});
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

  // CREATE EVENT
  onCreateEvent() {
    this.props.changeEventCreationMode(true);
  }

  // STOP CREATING EVENT
  onQuitEventCreationMode() {
    this.props.changeEventCreationMode(false);
  }

  // Terminate Active scenario
  terminateScenario() {
    this.props.api.terminateScenario((error, data, response) => this.terminateScenarioCb(error, data, response));
  }

  showApps(show) {
    this.props.changeShowApps(show);
    // _.defer(() => {
    //   this.props.execVis.network.setData(this.props.execVisData);
    // });
  }

  renderDialogs() {
    return (
      <>
          <IDDeployScenarioDialog
            title='Open Scenario'
            open={this.props.currentDialog===IDC_DIALOG_DEPLOY_SCENARIO}
            options={this.props.scenarios}
            onClose={() => {this.closeDialog();}}
            api={this.props.api}
            activateScenarioCb={(error, data, response) => this.activateScenarioCb(error, data, response)}
          />

          <IDSaveScenarioDialog
            title='Save Scenario as ...'
            open={this.props.currentDialog===IDC_DIALOG_SAVE_SCENARIO}
            onClose={() => {this.closeDialog();}}
            api={this.props.api}
            saveScenario={(name) => this.saveScenario(name)}
            scenarioNameRequired={true}
          />

          <IDTerminateScenarioDialog
            title='Terminate Scenario'
            open={this.props.currentDialog===IDC_DIALOG_TERMINATE_SCENARIO}
            scenario={this.props.scenario}
            onClose={() => {this.closeDialog();}}
            onSubmit={() => {this.terminateScenario();}}
          />
      </>
    );
  }

  render() {
    if (this.props.page !== PAGE_EXECUTE && this.props.page !== PAGE_EXPERIMENTAL_EXECUTE) {
      return null;
    }

    const scenarioName = (this.props.page === PAGE_EXECUTE)
      ? this.props.execScenarioName
      : this.props.cfgScenarioName;

    const spanLeft = this.props.eventCreationMode ? 8 : 12;
    const spanRight = this.props.eventCreationMode ? 4 : 0;
    return (
      <div style={{width: '100%'}}>
        {this.renderDialogs()}

        <div style={{width: '100%'}}>
          <Grid style={styles.headlineGrid}>
            <GridCell span={12}>
              <Elevation className="component-style" z={2} style={styles.headline}>
                <GridInner>
                  <GridCell align={'middle'} span={4}>
                    <HeadlineBar
                      titleLabel="Deployed Scenario"
                      scenarioName={scenarioName}
                    />
                  </GridCell>
                  <GridCell span={8}>
                    <GridInner align={'right'}>
                      <GridCell align={'middle'} span={12}>
                        <ExecPageScenarioButtons
                          onDeploy={() => this.onDeployScenario()}
                          onSaveScenario={() => this.onSaveScenario()}
                          onTerminate={() => this.onTerminateScenario()}
                          onRefresh={this.props.refreshScenario}
                          onCreateEvent={() => this.onCreateEvent()}
                          onShowAppsChanged={(show) => this.showApps(show)}
                          showApps={this.props.showApps}
                        />
                      </GridCell>
                    </GridInner>
                  </GridCell>
                </GridInner>
              </Elevation>
            </GridCell>
          </Grid>
        </div>

        {this.props.exec.state.scenario !== EXEC_STATE_IDLE &&
          <>
              <Grid style={{width: '100%'}}>
                <GridInner>
                  <GridCell span={spanLeft}>
                    <Elevation className="component-style" z={2}>
                      <div style={{padding: 10}}>
                        {this.props.experimental ? (<DashboardContainer />) : (<IDCVis type={TYPE_EXEC} />)}
                        
                      </div>
                    </Elevation>
                  </GridCell>
                  <GridCell span={spanRight} hidden={!this.props.eventCreationMode} style={styles.inner}>
                    <Elevation className="component-style" z={2}>
                      <EventCreationPane
                        eventTypes={[UE_MOBILITY_EVENT, NETWORK_CHARACTERISTICS_EVENT]}
                        api={this.props.api}
                        onSuccess={() => {this.props.refreshScenario();}}
                        onClose={() => this.onQuitEventCreationMode()}
                      />
                    </Elevation>
                  </GridCell>
                </GridInner>
              </Grid>
          </>
        }
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
    page: state.ui.page,
    execScenarioName: state.exec.scenario.name,
    cfgScenarioName: state.cfg.scenario.name,
    execVisData: execVisFilteredData(state)
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeCurrentDialog: (type) => dispatch(uiChangeCurrentDialog(type)),
    changeScenario: (scenario) => dispatch(execChangeScenario(scenario)),
    changeDeployScenarioList: (scenarios) => dispatch(execChangeScenarioList(scenarios)),
    changeScenarioName: (name) => dispatch(execChangeScenarioName(name)),
    changeState: (s) => dispatch(execChangeScenarioState(s)),
    changeEventCreationMode: (val) => dispatch(uiExecChangeEventCreationMode(val)), // (true or false)
    changeCurrentEvent: (e) => dispatch(uiExecChangeCurrentEvent(e)),
    execChangeOkToTerminate: (ok) => dispatch(execChangeOkToTerminate(ok)),
    changeShowApps: (show) => dispatch(uiExecChangeShowApps(show))
  };
};

const ConnectedExecPageContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(ExecPageContainer);

export default ConnectedExecPageContainer;
