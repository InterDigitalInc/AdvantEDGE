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

import IDSelect from '../../components/helper-components/id-select';

import {
  VIEW_NAME_NONE,
  NET_TOPOLOGY_VIEW,
  NET_METRICS_PTP_VIEW,
  NET_METRICS_AGG_VIEW,
  WIRELESS_METRICS_PTP_VIEW,
  WIRELESS_METRICS_AGG_VIEW
} from '@/js/meep-constants';

import {
  uiExecChangeShowApps,
  uiExecChangeSandboxCfg
} from '@/js/state/ui';

import {
  updateObject
} from '@/js/util/object-util';

class DashCfgDetailPane extends Component {
  constructor(props) {
    super(props);
    autoBind(this);
  }

  componentDidUpdate() {
    // Create sandbox config if it does not exist
    const sandboxName = this.props.sandbox;
    const sandboxCfg = this.props.sandboxCfg;
    if (!sandboxName || (sandboxCfg && sandboxCfg[sandboxName])) {
      return;
    } else {
      var newSandboxCfg = updateObject({}, sandboxCfg);
      newSandboxCfg[sandboxName] = {
        dashboardView1: NET_TOPOLOGY_VIEW,
        dashboardView2: VIEW_NAME_NONE
      };
      this.props.changeSandboxCfg(newSandboxCfg);
    }
  }

  onClose(e) {
    e.preventDefault();
    this.props.onClose(e);
  }

  changeView1(viewName) {
    var sandboxCfg = updateObject({}, this.props.sandboxCfg);
    sandboxCfg[this.props.sandbox].dashboardView1 = viewName;
    this.props.changeSandboxCfg(sandboxCfg);
  }

  changeView2(viewName) {
    var sandboxCfg = updateObject({}, this.props.sandboxCfg);
    sandboxCfg[this.props.sandbox].dashboardView2 = viewName;
    this.props.changeSandboxCfg(sandboxCfg);
  }

  render() {
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
                this.props.index === 1 ? this.changeView1(e.target.value) : this.changeView2(e.target.value);  
              }}
              value={this.props.viewName}
            />
          </GridCell>

          { this.props.viewName === NET_TOPOLOGY_VIEW ?
            <>
              <GridCell span={2}>
                <Checkbox
                  checked={this.props.showApps}
                  onChange={e => this.props.changeShowApps(e.target.checked)}
                >
                  Show Apps
                </Checkbox>
              </GridCell>
            </> : null
          }

          { this.props.viewName === NET_METRICS_PTP_VIEW ?
            <>
              <GridCell span={12}>
                <IDSelect
                  label={'Source Node'}
                  outlined
                  options={netSrcNodeIds}
                  onChange={e => {
                    this.props.changeSourceNodeSelected(e.target.value);
                  }}
                  value={netSrcNodeIds.includes(this.props.sourceNodeSelected) ? this.props.sourceNodeSelected : 'None'}
                />
              </GridCell>
              <GridCell span={12}>
                <IDSelect
                  label={'Destination Node'}
                  outlined
                  options={netDstNodeIds}
                  onChange={e => {
                    this.props.changeDestNodeSelected(e.target.value);
                  }}
                  value={netDstNodeIds.includes(this.props.destNodeSelected) ? this.props.destNodeSelected : 'None'}
                />
              </GridCell>
            </> : null
          }

          { this.props.viewName === NET_METRICS_AGG_VIEW ?
            <>
              <GridCell span={12}>
                <IDSelect
                  label={'Source Node'}
                  outlined
                  options={netSrcNodeIds}
                  onChange={e => {
                    this.props.changeSourceNodeSelected(e.target.value);
                  }}
                  value={netSrcNodeIds.includes(this.props.sourceNodeSelected) ? this.props.sourceNodeSelected : 'None'}
                />
              </GridCell>
            </> : null
          }

          { this.props.viewName === WIRELESS_METRICS_PTP_VIEW ?
            <>
              <GridCell span={12}>
                <IDSelect
                  label={'Source Node'}
                  outlined
                  options={wirelessSrcNodeIds}
                  onChange={e => {
                    this.props.changeSourceNodeSelected(e.target.value);
                  }}
                  value={wirelessSrcNodeIds.includes(this.props.sourceNodeSelected) ? this.props.sourceNodeSelected : 'None'}
                />
              </GridCell>
              <GridCell span={12}>
                <IDSelect
                  label={'Destination Node'}
                  outlined
                  options={wirelessDstNodeIds}
                  onChange={e => {
                    this.props.changeDestNodeSelected(e.target.value);
                  }}
                  value={wirelessDstNodeIds.includes(this.props.destNodeSelected) ? this.props.destNodeSelected : 'None'}
                />
              </GridCell>
            </> : null
          }

          { this.props.viewName === WIRELESS_METRICS_AGG_VIEW ?
            <>
              <GridCell span={12}>
                <IDSelect
                  label={'Source Node'}
                  outlined
                  options={wirelessSrcNodeIds}
                  onChange={e => {
                    this.props.changeSourceNodeSelected(e.target.value);
                  }}
                  value={wirelessSrcNodeIds.includes(this.props.sourceNodeSelected) ? this.props.sourceNodeSelected : 'None'}
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
    sandboxCfg: state.ui.sandboxCfg
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeSandboxCfg: cfg => dispatch(uiExecChangeSandboxCfg(cfg)),
    changeShowApps: show => dispatch(uiExecChangeShowApps(show))
  };
};

const ConnectedDashCfgDetailPane = connect(
  mapStateToProps,
  mapDispatchToProps
)(DashCfgDetailPane);

export default ConnectedDashCfgDetailPane;
