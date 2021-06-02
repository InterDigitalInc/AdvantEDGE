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

import { connect } from 'react-redux';
import React, { Component } from 'react';
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
  uiExecChangeSourceNodeSelectedView1,
  uiExecChangeDestNodeSelectedView1,
  uiExecChangeSourceNodeSelectedView2,
  uiExecChangeDestNodeSelectedView2,
  uiExecChangeDashboardView1,
  uiExecChangeDashboardView2,
  uiExecChangeSandboxCfg
} from '@/js/state/ui';

import {
  EXEC_VIEW_SELECT,
  VIEW_1,
  VIEW_2,
  VIEW_NAME_NONE,
  MAP_VIEW,
  NET_TOPOLOGY_VIEW,
  DEFAULT_DASHBOARD_OPTIONS,
  PAGE_EXECUTE,
  EXEC_BTN_DASHBOARD_BTN_CLOSE
} from '@/js/meep-constants';

const ViewSelect = props => {
  return (
    <Grid style={styles.field}>
      <GridCell span={12}>
        <Select
          style={styles.select}
          label="View"
          fullwidth="true"
          outlined
          options={props.viewOptions}
          onChange={props.onChange}
          data-cy={EXEC_VIEW_SELECT}
          value={props.value}
        />
      </GridCell>
    </Grid>
  );
};

const ViewSelectFields = props => {
  switch (props.currentView) {
  case VIEW_1:
    return (
      <DashCfgDetailPane
        currentView={props.currentView}
        onSuccess={props.onSuccess}
        onClose={props.onClose}
        index={1}
        appIds={props.appIds}
        ueIds={props.ueIds}
        poaIds={props.poaIds}
        sandbox={props.sandbox}
        viewName={props.view1Name}
        sourceNodeSelected={props.sourceNodeSelectedView1}
        destNodeSelected={props.destNodeSelectedView1}
        changeViewName={props.changeView1}
        changeSourceNodeSelected={props.changeSourceNodeSelectedView1}
        changeDestNodeSelected={props.changeDestNodeSelectedView1}
        dashboardViewsList={props.dashboardViewsList}
        showApps={props.showApps}
      />
    );
  case VIEW_2:
    return (
      <DashCfgDetailPane
        currentView={props.currentView}
        onSuccess={props.onSuccess}
        onClose={props.onClose}
        index={2}
        appIds={props.appIds}
        ueIds={props.ueIds}
        poaIds={props.poaIds}
        sandbox={props.sandbox}
        viewName={props.view2Name}
        sourceNodeSelected={props.sourceNodeSelectedView2}
        destNodeSelected={props.destNodeSelectedView2}
        changeViewName={props.changeView2}
        changeSourceNodeSelected={props.changeSourceNodeSelectedView2}
        changeDestNodeSelected={props.changeDestNodeSelectedView2}
        dashboardViewsList={props.dashboardViewsList}
        showApps={props.showApps}
      />
    );
  default:
    return <div></div>;
  }
};

const showInExecStr = '<exec>';

class DashCfgPane extends Component {
  constructor(props) {
    super(props);
    this.state = {};

    if (!this.props.currentView) {
      this.props.changeView('');
    }
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

  getView1() {
    const sandboxName = this.props.sandbox;
    const sandboxCfg = this.props.sandboxCfg;
    if (sandboxCfg && sandboxCfg[sandboxName] && sandboxCfg[sandboxName].dashboardView1) {
      return sandboxCfg[sandboxName].dashboardView1;
    } else {
      return VIEW_NAME_NONE;
    }
  }

  getView2() {
    const sandboxName = this.props.sandbox;
    const sandboxCfg = this.props.sandboxCfg;
    if (sandboxCfg && sandboxCfg[sandboxName] && sandboxCfg[sandboxName].dashboardView2) {
      return sandboxCfg[sandboxName].dashboardView2;
    } else {
      return VIEW_NAME_NONE;
    }
  }

  changeView1(viewName, sdbCfg) {
    var sandboxCfg = updateObject({}, sdbCfg);
    sandboxCfg[this.props.sandbox].dashboardView1 = viewName;
    this.props.changeSandboxCfg(sandboxCfg);
  }

  changeView2(viewName) {
    var sandboxCfg = updateObject({}, this.props.sandboxCfg);
    sandboxCfg[this.props.sandbox].dashboardView2 = viewName;
    this.props.changeSandboxCfg(sandboxCfg);
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
    this.props.changeView('');
    this.props.onClose(e);
  }

  getRoot() {
    return d3.hierarchy(this.props.displayedScenario, getScenarioNodeChildren);
  }

  render() {
    if (this.props.page !== PAGE_EXECUTE || this.props.hide) {
      return null;
    }
    
    const view1Name = this.getView1();
    const view2Name = this.getView2();

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
      NET_TOPOLOGY_VIEW
    ];
    this.populateDashboardList(dashboardViewsList, DEFAULT_DASHBOARD_OPTIONS);
    this.populateDashboardList(dashboardViewsList, this.props.dashboardOptions);

    return (
      <div style={{ padding: 10 }}>
        <div style={styles.block}>
          <Typography use="headline6">Dashboard View</Typography>
        </div>
        <ViewSelect
          viewOptions={this.props.viewOptions}
          onChange={event => {
            this.props.changeView(event.target.value);
          }}
          value={this.props.currentView}
        />
        <ViewSelectFields
          currentView={this.props.currentView}
          onSuccess={this.props.onSuccess}
          onClose={e => this.onDashCfgClose(e)}
          sandbox={this.props.sandbox}
          appIds={appIds}
          ueIds={ueIds}
          poaIds={poaIds}
          view1Name={view1Name}
          view2Name={view2Name}
          sourceNodeSelectedView1={this.props.sourceNodeSelectedView1}
          destNodeSelectedView1={this.props.destNodeSelectedView1}
          sourceNodeSelectedView2={this.props.sourceNodeSelectedView2}
          destNodeSelectedView2={this.props.destNodeSelectedView2}
          changeView1={this.changeView1}
          changeView2={this.changeView2}
          changeSourceNodeSelectedView1={this.props.changeSourceNodeSelectedView1}
          changeDestNodeSelectedView1={this.props.changeDestNodeSelectedView1}
          changeSourceNodeSelectedView2={this.props.changeSourceNodeSelectedView2}
          changeDestNodeSelectedView2={this.props.changeDestNodeSelectedView2}
          dashboardViewsList={dashboardViewsList}
          showApps={this.props.showApps}
        />
        <div>
          <Grid style={{ marginTop: 10 }}>
            <GridInner>
              <GridCell span={12}>
                <Button
                  outlined
                  style={styles.button}
                  onClick={e => this.onDashCfgClose(e)}
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
    sourceNodeSelectedView1: state.ui.sourceNodeSelectedView1,
    destNodeSelectedView1: state.ui.destNodeSelectedView1,
    sourceNodeSelectedView2: state.ui.sourceNodeSelectedView2,
    destNodeSelectedView2: state.ui.destNodeSelectedView2,
    eventCreationMode: state.exec.eventCreationMode,
    scenarioState: state.exec.state.scenario,
    view1Name: state.ui.dashboardView1,
    view2Name: state.ui.dashboardView2,
    sandboxCfg: state.ui.sandboxCfg,
    dashboardOptions: state.monitor.dashboardOptions,
    showApps: state.ui.execShowApps
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeView: view => dispatch(uiExecChangeCurrentView(view)),
    changeSourceNodeSelectedView1: src => dispatch(uiExecChangeSourceNodeSelectedView1(src)),
    changeDestNodeSelectedView1: dest => dispatch(uiExecChangeDestNodeSelectedView1(dest)),
    changeSourceNodeSelectedView2: src => dispatch(uiExecChangeSourceNodeSelectedView2(src)),
    changeDestNodeSelectedView2: dest => dispatch(uiExecChangeDestNodeSelectedView2(dest)),
    changeView1: name => dispatch(uiExecChangeDashboardView1(name)),
    changeView2: name => dispatch(uiExecChangeDashboardView2(name)),
    changeSandboxCfg: cfg => dispatch(uiExecChangeSandboxCfg(cfg))
  };
};

const ConnectedDashCfgPane = connect(
  mapStateToProps,
  mapDispatchToProps
)(DashCfgPane);

export default ConnectedDashCfgPane;
