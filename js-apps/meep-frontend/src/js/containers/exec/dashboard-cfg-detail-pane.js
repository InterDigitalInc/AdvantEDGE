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
import { Checkbox } from '@rmwc/checkbox';
import { TextField, TextFieldHelperText } from '@rmwc/textfield';
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

class DashCfgDetailPane extends Component {
  constructor(props) {
    super(props);
    autoBind(this);
  }

  onClose(e) {
    e.preventDefault();
    this.props.onClose(e);
  }

  changeViewType(val) {
    var newViewCfg = updateObject({}, this.props.viewCfg);
    newViewCfg.viewType = val;
    this.props.changeViewCfg(this.props.currentView, newViewCfg);
  }

  changeSourceNodeSelected(val) {
    var newViewCfg = updateObject({}, this.props.viewCfg);
    newViewCfg.sourceNodeSelected = val;
    this.props.changeViewCfg(this.props.currentView, newViewCfg);
  }

  changeDestNodeSelected(val) {
    var newViewCfg = updateObject({}, this.props.viewCfg);
    newViewCfg.destNodeSelected = val;
    this.props.changeViewCfg(this.props.currentView, newViewCfg);
  }

  changeParticipants(val) {
    var newViewCfg = updateObject({}, this.props.viewCfg);
    newViewCfg.participants = val;
    this.props.changeViewCfg(this.props.currentView, newViewCfg);
  }

  render() {
    var viewCfg = this.props.viewCfg;
    var netSrcNodeIds = this.props.appIds;
    var netDstNodeIds = this.props.appIds;
    var wirelessSrcNodeIds = this.props.ueIds;
    var wirelessDstNodeIds = this.props.poaIds;

    return (
      <div>
        <Grid>
          <GridCell span={12}>
            <IDSelect
              label={'View Type'}
              outlined
              options={this.props.dashboardViewsList}
              onChange={e => {
                this.changeViewType(e.target.value);
              }}
              value={viewCfg.viewType}
            />
          </GridCell>

          { viewCfg.viewType === NET_TOPOLOGY_VIEW ?
            <>
              <GridCell span={6}>
                <Checkbox
                  checked={this.props.showApps}
                  onChange={e => this.props.changeShowApps(e.target.checked)}
                >
                  Show Apps
                </Checkbox>
              </GridCell>
            </> : null
          }

          { viewCfg.viewType === SEQ_DIAGRAM_VIEW ?
            <>
              <GridCell span={12}>
                <TextField
                  outlined
                  style={{ width: '100%', marginTop: 10 }}
                  label={'Participants'}
                  onChange={e => this.changeParticipants(e.target.value)}
                  value={viewCfg.participants ? viewCfg.participants : ''}
                />
                <TextFieldHelperText validationMsg={true}>
                  <span>{'comma-separated, ordered list'}</span>
                </TextFieldHelperText>        
              </GridCell>
              <GridCell span={12}>
                <Checkbox
                  checked={this.props.pauseSeq}
                  onChange={e => this.props.changePauseSeq(e.target.checked)}
                >
                  Pause
                </Checkbox>
              </GridCell>
            </> : null
          }

          { viewCfg.viewType === DATAFLOW_DIAGRAM_VIEW ?
            <>
              <GridCell span={12}>
                <Checkbox
                  checked={this.props.pauseDataflow}
                  onChange={e => this.props.changePauseDataflow(e.target.checked)}
                >
                  Pause
                </Checkbox>
              </GridCell>
            </> : null
          }

          { viewCfg.viewType === NET_METRICS_PTP_VIEW ?
            <>
              <GridCell span={12}>
                <IDSelect
                  label={'Source Node'}
                  outlined
                  options={netSrcNodeIds}
                  onChange={e => {
                    this.changeSourceNodeSelected(e.target.value);
                  }}
                  value={netSrcNodeIds.includes(viewCfg.sourceNodeSelected) ? viewCfg.sourceNodeSelected : 'None'}
                />
              </GridCell>
              <GridCell span={12}>
                <IDSelect
                  label={'Destination Node'}
                  outlined
                  options={netDstNodeIds}
                  onChange={e => {
                    this.changeDestNodeSelected(e.target.value);
                  }}
                  value={netDstNodeIds.includes(viewCfg.destNodeSelected) ? viewCfg.destNodeSelected : 'None'}
                />
              </GridCell>
            </> : null
          }

          { viewCfg.viewType === NET_METRICS_AGG_VIEW ?
            <>
              <GridCell span={12}>
                <IDSelect
                  label={'Source Node'}
                  outlined
                  options={netSrcNodeIds}
                  onChange={e => {
                    this.changeSourceNodeSelected(e.target.value);
                  }}
                  value={netSrcNodeIds.includes(viewCfg.sourceNodeSelected) ? viewCfg.sourceNodeSelected : 'None'}
                />
              </GridCell>
            </> : null
          }

          { viewCfg.viewType === WIRELESS_METRICS_PTP_VIEW ?
            <>
              <GridCell span={12}>
                <IDSelect
                  label={'UE'}
                  outlined
                  options={wirelessSrcNodeIds}
                  onChange={e => {
                    this.changeSourceNodeSelected(e.target.value);
                  }}
                  value={wirelessSrcNodeIds.includes(viewCfg.sourceNodeSelected) ? viewCfg.sourceNodeSelected : 'None'}
                />
              </GridCell>
              <GridCell span={12}>
                <IDSelect
                  label={'POA'}
                  outlined
                  options={wirelessDstNodeIds}
                  onChange={e => {
                    this.changeDestNodeSelected(e.target.value);
                  }}
                  value={wirelessDstNodeIds.includes(viewCfg.destNodeSelected) ? viewCfg.destNodeSelected : 'None'}
                />
              </GridCell>
            </> : null
          }

          { viewCfg.viewType === WIRELESS_METRICS_AGG_VIEW ?
            <>
              <GridCell span={12}>
                <IDSelect
                  label={'UE'}
                  outlined
                  options={wirelessSrcNodeIds}
                  onChange={e => {
                    this.changeSourceNodeSelected(e.target.value);
                  }}
                  value={wirelessSrcNodeIds.includes(viewCfg.sourceNodeSelected) ? viewCfg.sourceNodeSelected : 'None'}
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
