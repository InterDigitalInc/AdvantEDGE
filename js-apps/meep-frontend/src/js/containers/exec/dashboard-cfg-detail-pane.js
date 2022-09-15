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
import autoBind from 'react-autobind';
import { Grid, GridCell } from '@rmwc/grid';
import { Button } from '@rmwc/button';
import { Checkbox } from '@rmwc/checkbox';
import { TextField, TextFieldIcon, TextFieldHelperText } from '@rmwc/textfield';
import { updateObject } from '@/js/util/object-util';

import IDSelect from '../../components/helper-components/id-select';


import {
  NET_TOPOLOGY_VIEW,
  SEQ_DIAGRAM_VIEW,
  DATAFLOW_DIAGRAM_VIEW,
  NET_METRICS_PTP_VIEW,
  NET_METRICS_AGG_VIEW,
  WIRELESS_METRICS_PTP_VIEW,
  WIRELESS_METRICS_AGG_VIEW
} from '@/js/meep-constants';

import {
  uiExecChangeShowApps,
  uiExecChangePauseSeq,
  uiExecChangePauseDataflow
} from '@/js/state/ui';

import {
  DASH_CFG_VIEW_TYPE,
  DASH_CFG_SOURCE_NODE_SELECTED,
  DASH_CFG_DEST_NODE_SELECTED,
  DASH_CFG_PARTICIPANTS,
  DASH_CFG_MAX_MSG_COUNT,
  getDashCfgFieldVal,
  setDashCfgField
} from '@/js/util/dashboard-utils';

// COMPONENTS
const DashCfgTextField = props => {
  var err = props.viewCfg[props.fieldName] ? props.viewCfg[props.fieldName].err: null;

  return (
    <>
      <TextField
        outlined
        style={{ width: '100%', marginBottom: 0 }}
        label={props.label}
        withLeadingIcon={!props.icon ? null : 
          <TextFieldIcon
            tabIndex="0"
            icon={props.icon}
            onClick={props.onIconClick}
          />
        }
        type={props.type}
        onChange={event => {
          var err = props.validate ? props.validate(event.target.value) : null;
          var val =
            event.target.value && props.isNumber && !err
              ? Number(event.target.value)
              : event.target.value;
          props.onUpdate(props.fieldName, val, err);
        }}
        invalid={err}
        value={
          props.viewCfg[props.fieldName]
            ? props.viewCfg[props.fieldName].val
            : ''
        }
        disabled={props.disabled}
        data-cy={props.cydata}
      />
      <TextFieldHelperText validationMsg={true}>
        <span>
          {props.viewCfg[props.fieldName] ? props.viewCfg[props.fieldName].err : ''}
        </span>
      </TextFieldHelperText>
    </>
  );
};

class DashCfgDetailPane extends Component {
  constructor(props) {
    super(props);
    autoBind(this);
  }

  onClose(e) {
    e.preventDefault();
    this.props.onClose(e);
  }

  onUpdateViewType(event) {
    this.onUpdateDashCfg(DASH_CFG_VIEW_TYPE, event.target.value, null);
  }
  onUpdateSourceNodeSelected(event) {
    this.onUpdateDashCfg(DASH_CFG_SOURCE_NODE_SELECTED, event.target.value, null);
  }
  onUpdateDestNodeSelected(event) {
    this.onUpdateDashCfg(DASH_CFG_DEST_NODE_SELECTED, event.target.value, null);
  }
  onUpdateDashCfg(name, val, err) {
    var newViewCfg = updateObject({}, this.props.viewCfg);
    setDashCfgField(newViewCfg, name, val, err);
    this.props.changeViewCfg(this.props.currentView, newViewCfg);
  }

  onChangeShowApps(event) {
    this.props.changeShowApps(event.target.checked);
  }
  onChangePauseSeq(event) {
    this.props.changePauseSeq(event.target.checked);
  }
  onChangePauseDataflow(event) {
    this.props.changePauseDataflow(event.target.checked);
  }

  render() {
    this.viewCfg = this.props.viewCfg;
    this.netSrcNodeIds = this.props.appIds;
    this.netDstNodeIds = this.props.appIds;
    this.wirelessSrcNodeIds = this.props.ueIds;
    this.wirelessDstNodeIds = this.props.poaIds;
    this.viewType = getDashCfgFieldVal(this.viewCfg, DASH_CFG_VIEW_TYPE);
    this.sourceNodeSelected = getDashCfgFieldVal(this.viewCfg, DASH_CFG_SOURCE_NODE_SELECTED);
    this.destNodeSelected = getDashCfgFieldVal(this.viewCfg, DASH_CFG_DEST_NODE_SELECTED);

    return (
      <div>
        <Grid>
          <GridCell span={12} style={{marginBottom: 10}}>
            <IDSelect
              label={'View Type'}
              outlined
              options={this.props.dashboardViewsList}
              onChange={this.onUpdateViewType}
              value={this.viewType}
            />
          </GridCell>

          { this.viewType === NET_TOPOLOGY_VIEW ?
            <>
              <GridCell span={6}>
                <Checkbox
                  checked={this.props.showApps}
                  onChange={this.onChangeShowApps}
                >
                  Show Apps
                </Checkbox>
              </GridCell>
            </> : null
          }

          { this.viewType === SEQ_DIAGRAM_VIEW ?
            <>
              <GridCell span={12}>
                <DashCfgTextField
                  onUpdate={this.onUpdateDashCfg}
                  viewCfg={this.viewCfg}
                  // validate={validateParticipants}
                  isNumber={false}
                  label={'Participants'}
                  fieldName={DASH_CFG_PARTICIPANTS}
                />     
              </GridCell>
              <GridCell span={12}>
                <DashCfgTextField
                  onUpdate={this.onUpdateDashCfg}
                  viewCfg={this.viewCfg}
                  // validate={validateParticipants}
                  isNumber={true}
                  label={'Max message count'}
                  fieldName={DASH_CFG_MAX_MSG_COUNT}
                /> 
              </GridCell>
              <GridCell span={12} style={{marginBottom: 10}}>
                <Button
                  outlined
                  // onClick={}
                >
                  Clear
                </Button>
                <Button
                  outlined
                  style={{marginLeft: 10}}
                  // onClick={}
                >
                  Fetch all
                </Button>
              </GridCell>
              <GridCell span={12}>
                <Checkbox
                  // style={{marginTop: 20}}
                  checked={this.props.pauseSeq}
                  onChange={this.onChangePauseSeq}
                >
                  Pause
                </Checkbox>
              </GridCell>
            </> : null
          }

          { this.viewType === DATAFLOW_DIAGRAM_VIEW ?
            <>
              <GridCell span={12}>
                <Checkbox
                  checked={this.props.pauseDataflow}
                  onChange={this.onChangePauseDataflow}
                >
                  Pause
                </Checkbox>
              </GridCell>
            </> : null
          }

          { this.viewType === NET_METRICS_PTP_VIEW ?
            <>
              <GridCell span={12} style={{marginBottom: 10}}>
                <IDSelect
                  label={'Source Node'}
                  outlined
                  options={this.netSrcNodeIds}
                  onChange={this.onUpdateSourceNodeSelected}
                  value={this.netSrcNodeIds.includes(this.sourceNodeSelected) ? this.sourceNodeSelected : 'None'}
                />
              </GridCell>
              <GridCell span={12} style={{marginBottom: 10}}>
                <IDSelect
                  label={'Destination Node'}
                  outlined
                  options={this.netDstNodeIds}
                  onChange={this.onUpdateDestNodeSelected}
                  value={this.netDstNodeIds.includes(this.destNodeSelected) ? this.destNodeSelected : 'None'}
                />
              </GridCell>
            </> : null
          }

          { this.viewType === NET_METRICS_AGG_VIEW ?
            <>
              <GridCell span={12} style={{marginBottom: 10}}>
                <IDSelect
                  label={'Source Node'}
                  outlined
                  options={this.netSrcNodeIds}
                  onChange={this.onUpdateSourceNodeSelected}
                  value={this.netSrcNodeIds.includes(this.sourceNodeSelected) ? this.sourceNodeSelected : 'None'}
                />
              </GridCell>
            </> : null
          }

          { this.viewType === WIRELESS_METRICS_PTP_VIEW ?
            <>
              <GridCell span={12} style={{marginBottom: 10}}>
                <IDSelect
                  label={'UE'}
                  outlined
                  options={this.wirelessSrcNodeIds}
                  onChange={this.onUpdateSourceNodeSelected}
                  value={this.wirelessSrcNodeIds.includes(this.sourceNodeSelected) ? this.sourceNodeSelected : 'None'}
                />
              </GridCell>
              <GridCell span={12} style={{marginBottom: 10}}>
                <IDSelect
                  label={'POA'}
                  outlined
                  options={this.wirelessDstNodeIds}
                  onChange={this.onUpdateDestNodeSelected}
                  value={this.wirelessDstNodeIds.includes(this.destNodeSelected) ? this.destNodeSelected : 'None'}
                />
              </GridCell>
            </> : null
          }

          { this.viewType === WIRELESS_METRICS_AGG_VIEW ?
            <>
              <GridCell span={12} style={{marginBottom: 10}}>
                <IDSelect
                  label={'UE'}
                  outlined
                  options={this.wirelessSrcNodeIds}
                  onChange={this.onUpdateSourceNodeSelected}
                  value={this.wirelessSrcNodeIds.includes(this.sourceNodeSelected) ? this.sourceNodeSelected : 'None'}
                />
              </GridCell>
            </> : null
          }

        </Grid>
      </div>
    );
  }
}

const mapStateToProps = state => {
  return {
    showApps: state.ui.execShowApps,
    pauseSeq: state.ui.execPauseSeq,
    pauseDataflow: state.ui.execPauseDataflow
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeShowApps: show => dispatch(uiExecChangeShowApps(show)),
    changePauseSeq: pause => dispatch(uiExecChangePauseSeq(pause)),
    changePauseDataflow: pause => dispatch(uiExecChangePauseDataflow(pause))
  };
};

const ConnectedDashCfgDetailPane = connect(
  mapStateToProps,
  mapDispatchToProps
)(DashCfgDetailPane);

export default ConnectedDashCfgDetailPane;
