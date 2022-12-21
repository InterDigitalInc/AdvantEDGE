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

import { connect } from 'react-redux';
import React, { Component } from 'react';
import autoBind from 'react-autobind';
import { Select } from '@rmwc/select';
import { Button } from '@rmwc/button';
import { Grid, GridCell, GridInner } from '@rmwc/grid';
import { Typography } from '@rmwc/typography';
import { updateObject } from '@/js/util/object-util';
import * as d3 from 'd3';
import DashCfgDetailPane from './dashboard-cfg-detail-pane';

import {
  getScenarioNodeChildren,
  isApp,
  isUe,
  isPoa
} from '../../util/scenario-utils';

import {
  uiExecChangeCurrentView,
  uiExecChangeSandboxCfg
} from '@/js/state/ui';

import {
  EXEC_VIEW_SELECT,
  VIEW_1,
  VIEW_2,
  VIEW_NAME_NONE,
  MAP_VIEW,
  NET_TOPOLOGY_VIEW,
  SEQ_DIAGRAM_VIEW,
  DATAFLOW_DIAGRAM_VIEW,
  DEFAULT_DASHBOARD_OPTIONS,
  PAGE_EXECUTE,
  EXEC_BTN_DASHBOARD_BTN_CLOSE
} from '@/js/meep-constants';

const showInExecStr = '<exec>';

class DashCfgPane extends Component {
  constructor(props) {
    super(props);
    autoBind(this);

    if (!this.props.currentView) {
      this.props.changeView('');
    }
  }

  getViewCfg(view) {
    const sandboxName = this.props.sandbox;
    const sandboxCfg = this.props.sandboxCfg;
    if (sandboxCfg && sandboxCfg[sandboxName]) {
      if (view === VIEW_1) {
        return sandboxCfg[sandboxName].dashView1;
      } else if (view === VIEW_2) {
        return sandboxCfg[sandboxName].dashView2;
      }
    }
    return null;
  }

  changeViewCfg(view, viewCfg) {
    const sandboxName = this.props.sandbox;
    var newSandboxCfg = updateObject({}, this.props.sandboxCfg);
    if (newSandboxCfg && newSandboxCfg[sandboxName]) {
      if (view === VIEW_1) {
        newSandboxCfg[sandboxName].dashView1 = viewCfg;
      } else if (view === VIEW_2) {
        newSandboxCfg[sandboxName].dashView2 = viewCfg;
      }
      this.props.changeSandboxCfg(newSandboxCfg);
    }
  }

  changeView(event) {
    this.props.changeView(event.target.value);
  }

  populateDashboardList(dashboardViewsList, dashboardOptions) {
    for (var i = 0; i < dashboardOptions.length; i++) {
      var dashboard = dashboardOptions[i];
      if ((dashboard.label !== '') && (dashboard.value !== '') && (dashboard.value.indexOf(showInExecStr) !== -1)) {
        dashboardViewsList.push(dashboard.label);
      }
    }
  }

  onDashCfgClose(e) {
    e.preventDefault();
    this.props.onClose(e);
  }

  getRoot() {
    return d3.hierarchy(this.props.displayedScenario, getScenarioNodeChildren);
  }

  render() {
    if (this.props.page !== PAGE_EXECUTE || this.props.hide) {
      return null;
    }

    var currentViewCfg = this.getViewCfg(this.props.currentView);

    const root = this.getRoot();
    const nodes = root.descendants();
    const apps = nodes.filter(isApp);
    const appIds = apps.map(a => a.data.name).sort();
    const ues = nodes.filter(isUe);
    const ueIds = ues.map(a => a.data.name).sort();
    const poas = nodes.filter(isPoa);
    const poaIds = poas.map(a => a.data.name).sort();
    appIds.unshift('None');
    ueIds.unshift('None');
    poaIds.unshift('None');

    // Populate Dashboard view list using links from monitoring tab
    var dashboardViewsList = [
      VIEW_NAME_NONE,
      MAP_VIEW,
      NET_TOPOLOGY_VIEW,
      SEQ_DIAGRAM_VIEW,
      DATAFLOW_DIAGRAM_VIEW
    ];
    this.populateDashboardList(dashboardViewsList, DEFAULT_DASHBOARD_OPTIONS);
    this.populateDashboardList(dashboardViewsList, this.props.dashboardOptions);

    return (
      <div style={{ padding: 10 }}>
        <div style={styles.block}>
          <Typography use="headline6">Dashboard View</Typography>
        </div>
        <Grid style={styles.field}>
          <GridCell span={12}>
            <Select
              style={styles.select}
              label="View"
              fullwidth="true"
              outlined
              options={this.props.viewOptions}
              onChange={this.changeView}
              data-cy={EXEC_VIEW_SELECT}
              value={this.props.currentView}
            />
          </GridCell>
        </Grid>
        {currentViewCfg &&
          <DashCfgDetailPane
            currentView={this.props.currentView}
            viewCfg={currentViewCfg}
            dashboardViewsList={dashboardViewsList}
            changeViewCfg={this.changeViewCfg}
            onClose={this.onDashCfgClose}
            appIds={appIds}
            ueIds={ueIds}
            poaIds={poaIds}
            metricsApi={this.props.metricsApi}
          />
        }
        <div>
          <Grid style={{ marginTop: 20 }}>
            <GridInner>
              <GridCell span={12}>
                <Button
                  outlined
                  style={styles.button}
                  onClick={this.onDashCfgClose}
                  data-cy={EXEC_BTN_DASHBOARD_BTN_CLOSE}
                >
                  Close
                </Button>
              </GridCell>
            </GridInner>
          </Grid>
        </div>
      </div>
    );
  }
}

const styles = {
  block: {
    marginBottom: 20
  },
  field: {
    marginBottom: 20
  },
  select: {
    width: '100%'
  }
};

const mapStateToProps = state => {
  return {
    currentView: state.ui.execCurrentView,
    page: state.ui.page,
    displayedScenario: state.exec.displayedScenario,
    scenarioState: state.exec.state.scenario,
    sandboxCfg: state.ui.sandboxCfg,
    dashboardOptions: state.monitor.dashboardOptions
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeView: view => dispatch(uiExecChangeCurrentView(view)),
    changeSandboxCfg: cfg => dispatch(uiExecChangeSandboxCfg(cfg))
  };
};

const ConnectedDashCfgPane = connect(
  mapStateToProps,
  mapDispatchToProps
)(DashCfgPane);

export default ConnectedDashCfgPane;
