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

import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Select } from '@rmwc/select';
import { Grid, GridCell } from '@rmwc/grid';
import { Elevation } from '@rmwc/elevation';

import CfgNetworkElementContainer from '@/js/containers/cfg/cfg-network-element-container';
import CancelApplyPair from '@/js/components/helper-components/cancel-apply-pair';

import { validateNetworkElement } from '@/js/containers/cfg/cfg-page-container';

import {
  TYPE_EXEC,
  UE_TYPE_STR,
  UE_APP_TYPE_STR,
  EDGE_APP_TYPE_STR,
  CLOUD_APP_TYPE_STR,
  EXEC_EVT_SU_ACTION,
  EXEC_EVT_SU_REMOVE_ELEM_NAME,
  SCENARIO_UPDATE_ACTION_NONE,
  SCENARIO_UPDATE_ACTION_ADD,
  SCENARIO_UPDATE_ACTION_MODIFY,
  SCENARIO_UPDATE_ACTION_REMOVE,
  ELEMENT_TYPE_UE_APP,
  ELEMENT_TYPE_EDGE_APP,
  ELEMENT_TYPE_CLOUD_APP
} from '@/js/meep-constants';

import {
  uiExecChangeScenarioUpdateAction,
  uiExecScenarioUpdateRemoveEleName
} from '@/js/state/ui';

import {
  execElemNew,
  execElemClear,
  execElemSetErrMsg,
  EXEC_ELEM_MODE_NEW
} from '@/js/state/exec';

import {
  getElemFieldVal,
  FIELD_NAME,
  FIELD_PARENT
} from '@/js/util/elem-utils';

import { updateObject } from '@/js/util/object-util';

import { addElementToScenario, updateElementInScenario } from '@/js/util/scenario-utils';

var elementNames = [];

class ScenarioUpdateEventPane extends Component {
  constructor(props) {
    super(props);

    this.state = {
      actionTypes: [
        SCENARIO_UPDATE_ACTION_NONE,
        SCENARIO_UPDATE_ACTION_ADD,
        // SCENARIO_UPDATE_ACTION_MODIFY,
        SCENARIO_UPDATE_ACTION_REMOVE
      ]
    };
  }

  componentDidMount() {
    this.props.changeActionType(SCENARIO_UPDATE_ACTION_NONE);
    this.props.execElemNew();
    this.props.changeRemoveActionEleType('');
    this.props.changeRemoveActionEleName('');
  }

  changeAction(action) {
    this.props.changeActionType(action);
    if (action === SCENARIO_UPDATE_ACTION_ADD) {
      this.props.execElemNew();
    } else if (action === SCENARIO_UPDATE_ACTION_REMOVE) {
      this.props.changeRemoveActionEleType('');
      this.props.changeRemoveActionEleName('');
    }
  }

  changeElementType(elementType) {
    this.props.changeRemoveActionEleType(elementType);
    this.getElementNames(elementType, this.props.scenario);
    this.props.changeRemoveActionEleName('');
  }

  getElementNames(elementName, scenario) {
    elementNames.length = 0;
    var neType = '';
    switch(elementName) {
    case ELEMENT_TYPE_UE_APP:
      neType = UE_APP_TYPE_STR;
      break;
    case ELEMENT_TYPE_EDGE_APP:
      neType = EDGE_APP_TYPE_STR;
      break;
    case ELEMENT_TYPE_CLOUD_APP:
      neType = CLOUD_APP_TYPE_STR;
      break;
    default:
      return;
    }
    for (var dInd in scenario.deployment.domains) {
      var domain = scenario.deployment.domains[dInd];
      for (var zInd in domain.zones) {
        var zone = domain.zones[zInd];
        for (var nInd in zone.networkLocations) {
          var nl = zone.networkLocations[nInd];
          for (var plInd in nl.physicalLocations) {
            var pl = nl.physicalLocations[plInd];
            for (var prInd in pl.processes) {
              var pr = pl.processes[prInd];
              if (pr.type === neType) {
                elementNames.push(pr.name);
              }
            }
          }
        }
      }
    }
  }

  onSaveElement(element) {
    if (!validateNetworkElement(element, this.props.table.entries, this.props.execElemSetErrMsg)) {
      return;
    }

    var action = '';
    var updatedScenario = updateObject({}, this.props.scenario);
    if (this.props.execConfigMode === EXEC_ELEM_MODE_NEW) {
      addElementToScenario(updatedScenario, element);
      action = SCENARIO_UPDATE_ACTION_ADD;
    } else {
      updateElementInScenario(updatedScenario, element);
      action = SCENARIO_UPDATE_ACTION_MODIFY;
    }

    var pl = this.getPLFromScenario(getElemFieldVal(element, FIELD_NAME), updatedScenario);
    this.sendEvent(getElemFieldVal(element, FIELD_PARENT), pl, action);
    this.props.execElemClear();
    this.props.execElemNew();
  }

  onDeleteElement(e) {
    e.preventDefault();
    var pl = { name: this.props.scenarioUpdateRemoveEleName };
    this.sendEvent('', pl, SCENARIO_UPDATE_ACTION_REMOVE);
    this.props.execElemClear();
    this.props.changeRemoveActionEleName('');
    this.props.changeRemoveActionEleType('');
  }

  onCancelElement(e) {
    e.preventDefault();
    this.props.changeActionType(SCENARIO_UPDATE_ACTION_NONE);
    this.props.execElemClear();
    this.props.changeRemoveActionEleName('');
    this.props.onClose(e);
  }

  getPLFromScenario (elementName, scenario) {
    if (elementName === null) {
      return null;
    }

    for (var dInd in scenario.deployment.domains) {
      var domain = scenario.deployment.domains[dInd];
      for (var zInd in domain.zones) {
        var zone = domain.zones[zInd];
        for (var nInd in zone.networkLocations) {
          var nl = zone.networkLocations[nInd];
          for (var plInd in nl.physicalLocations) {
            var pl = nl.physicalLocations[plInd];
            if (pl.name === elementName) {
              return pl;
            }
            for (var prInd in pl.processes) {
              var pr = pl.processes[prInd];
              if (pr.name === elementName) {
                return pl;
              }
            }
          }
        }
      }
    }
    return null;
  }

  sendEvent(parentVal, pl, action) {
    if (pl === null || parentVal === null) {
      return;
    }

    var meepEvent = {
      name: 'name',
      type: this.props.currentEvent,
      eventScenarioUpdate: {
        action: action,
        node: {
          type: UE_TYPE_STR, // Change it to actual Procress Type
          parent: parentVal,
          nodeDataUnion: {
            physicalLocation: pl
          }
        }
      }
    };

    this.props.api.sendEvent(this.props.currentEvent, meepEvent, error => {
      if (!error) {
        this.props.onSuccess();
      }
    });
  }

  render() {
    return (
      <div style={styles.page}>
        <Grid style={styles.field}>
          <GridCell span="8">
            <Select
              style={styles.select}
              label="Action Type"
              outlined
              data-cy={EXEC_EVT_SU_ACTION}
              options={this.state.actionTypes}
              onChange={e => { this.changeAction(e.target.value); }}
              value={this.props.scenarioUpdateAction}
            />
          </GridCell>
          <GridCell span="4"></GridCell>
        </Grid>
        { this.props.scenarioUpdateAction === SCENARIO_UPDATE_ACTION_ADD ||
          this.props.scenarioUpdateAction === SCENARIO_UPDATE_ACTION_MODIFY ?
          <Grid>
            <GridCell span={12} style={styles.inner}>
              <Elevation className="component-style" z={2}>
                <CfgNetworkElementContainer
                  style={{ height: '100%' }}
                  onNewElement={() => {}}
                  onSaveElement={elem => this.onSaveElement(elem)}
                  onDeleteElement={() => {}}
                  onApplyCloneElement={() => {}}
                  onCancelElement={e => this.onCancelElement(e)}
                  onEditLocation={() => {}}
                  onEditPath={() => {}}
                  type={TYPE_EXEC}
                />
              </Elevation>
            </GridCell>
          </Grid> : null
        }
        { this.props.scenarioUpdateAction === 'REMOVE' ?          
          <div>
            <Grid style={styles.block}>
              <GridCell span="8">
                <Select
                  style={styles.select}
                  label="Element Name"
                  outlined
                  options={elementNames}
                  onChange={e => { this.props.changeRemoveActionEleName(e.target.value); }}
                  data-cy={EXEC_EVT_SU_REMOVE_ELEM_NAME}
                  value={this.props.scenarioUpdateRemoveEleName}
                />
              </GridCell>
              <GridCell span="4"></GridCell>
            </Grid>
          </div> : null
        }
        { this.props.scenarioUpdateAction === SCENARIO_UPDATE_ACTION_NONE ||
          this.props.scenarioUpdateAction === SCENARIO_UPDATE_ACTION_REMOVE ?
          <CancelApplyPair
            cancelText="Cancel"
            applyText="Apply"
            onCancel={e => this.onCancelElement(e)}
            onApply={e => this.onDeleteElement(e)}
            saveDisabled={
              this.props.scenarioUpdateRemoveEleName === '' ||
              this.props.scenarioUpdateAction === SCENARIO_UPDATE_ACTION_NONE
            }
            removeCyCancel={true}
          /> : null
        }
      </div>
    );
  }
}

const styles = {
  field: {
    marginBottom: 10,
    width: '100%'
  },
  select: {
    width: '100%'
  },
  inner: {
    height: '100%'
  },
  page: {
    height: '100%',
    marginBottom: 10,
    width: '100%',
    marginRight: 100
  },
  block: {
    marginBottom: 20
  }
};

const mapStateToProps = state => {
  return {
    scenarioUpdateAction: state.ui.scenarioUpdateAction,
    scenarioUpdateRemoveEleName: state.ui.scenarioUpdateRemoveEleName,
    execConfigMode: state.exec.elementConfiguration.configurationMode,
    table: state.exec.table,
    scenario: state.exec.scenario
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeActionType: event => dispatch(uiExecChangeScenarioUpdateAction(event)),
    changeRemoveActionEleName: event => dispatch(uiExecScenarioUpdateRemoveEleName(event)),
    execElemNew: elem => dispatch(execElemNew(elem)),
    execElemClear: elem => dispatch(execElemClear(elem)),
    execElemSetErrMsg: msg => dispatch(execElemSetErrMsg(msg))
  };
};

const ConnectedScenarioUpdateEventPane = connect(
  mapStateToProps,
  mapDispatchToProps
)(ScenarioUpdateEventPane);

export default ConnectedScenarioUpdateEventPane;
