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
import autoBind from 'react-autobind';
import { Grid, GridCell } from '@rmwc/grid';
import { Elevation } from '@rmwc/elevation';
// import ReactDOM from 'react-dom';

import IDCMap from '../idc-map';
import IDCVis from '../idc-vis';
import Iframe from 'react-iframe';

import {
  uiExecChangeDashboardView1,
  uiExecChangeDashboardView2,
  uiExecChangeSandboxCfg
} from '../../state/ui';

import {
  TYPE_EXEC,
  VIEW_NAME_NONE,
  MAP_VIEW,
  NET_TOPOLOGY_VIEW,
  DEFAULT_DASHBOARD_OPTIONS
} from '../../meep-constants';

import { updateObject } from '../../util/object-util';

const styles = {
  button: {
    marginRight: 10
  },
  dashboard: {
    height: '70vh',
    minHeight: 600
  }
};

const showInExecStr = '<exec>';
const passVarsStr = '<vars>';

const getUrl = (dashboardName, dashboardOptions) => {
  var url = '';
  if (dashboardOptions) {
    for (var i = 0; i < dashboardOptions.length; i++) {
      var dashboard = dashboardOptions[i];
      if (dashboard.label === dashboardName) {
        url = dashboard.value;
        url = url.replace(showInExecStr, '');
        break;
      }
    }
  }
  return url;
};

const ViewForName = ({
  sandboxName,
  scenarioName,
  selectedSource,
  selectedDest,
  viewName,
  dashboardOptions
}) => {

  // Handle Map view
  if (viewName === MAP_VIEW) {
    return (
      <div style={styles.dashboard}>
        <IDCMap
          type={TYPE_EXEC}
          sandboxName={sandboxName}
        />
      </div>
    );
  }

  // Handle Network Topology view
  if (viewName === NET_TOPOLOGY_VIEW) {
    return (
      <div style={styles.dashboard}>
        <IDCVis
          type={TYPE_EXEC}
          width='100%'
          height='100%'
          onEditElement={() => { }}
        />
      </div>
    );
  }

  // Get URL from Monitoring page dashboard options
  var selectedUrl = getUrl(viewName, DEFAULT_DASHBOARD_OPTIONS);
  if (selectedUrl === '') {
    selectedUrl = getUrl(viewName, dashboardOptions);
  }

  // Add variables if requested
  if (selectedUrl !== '') {
    if (selectedUrl.indexOf(passVarsStr) !== -1) {
      selectedUrl = selectedUrl.replace(passVarsStr, '');
      
      // Prepend sandbox name to scenario name and replace '-' with '_'
      var sandboxScenario = sandboxName + '_' + scenarioName;
      var scenario = sandboxScenario.replace(/-/g, '_');

      var url = new URL(selectedUrl);
      url.searchParams.append('var-database', scenario);
      url.searchParams.append('var-src', selectedSource);
      url.searchParams.append('var-dest', selectedDest);
      selectedUrl = url.href + '&kiosk';
    }

    return (
      <div style={styles.dashboard}>
        <Iframe
          url={selectedUrl}
          id='myId'
          display='initial'
          position='relative'
          allowFullScreen
          width='100%'
          height='100%'
        />
      </div>
    );
  }

  return null;
};

class DashboardContainer extends Component {
  constructor(props) {
    super(props);
    autoBind(this);

    this.state = {
      sourceNodeId: ''
    };
  }

  componentDidMount() { }

  componentWillUnmount() {
    clearInterval(this.dataTimer);
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

  changeShowApps(checked) {
    this.props.onShowAppsChanged(checked);
  }

  populateDashboardList(dashboardViewsList, dashboardOptions) {
    for (var i = 0; i < dashboardOptions.length; i++) {
      var dashboard = dashboardOptions[i];
      if ((dashboard.label !== '') && (dashboard.value !== '') && (dashboard.value.indexOf(showInExecStr) !== -1)) {
        dashboardViewsList.push(dashboard.label);
      }
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
    this.keyForSvg++;

    const view1Name = this.getView1();
    const view2Name = this.getView2();
    const view1Present = view1Name !== VIEW_NAME_NONE;
    const view2Present = view2Name !== VIEW_NAME_NONE;

    let span1 = 12;
    let span2 = 12;
    if (view1Present && view2Present) {
      span1 = 6;
      span2 = 6;
    } else if (!view1Present && !view2Present) {
      span1 = 0;
      span2 = 0;
    }

    const view1 = (
      <ViewForName
        sandboxName={this.props.sandbox}
        scenarioName={this.props.scenarioName}
        selectedSource={this.props.sourceNodeSelectedView1}
        selectedDest={this.props.destNodeSelectedView1}
        viewName={view1Name}
        dashboardOptions={this.props.dashboardOptions}
      />
    );

    const view2 = (
      <ViewForName
        sandboxName={this.props.sandbox}
        scenarioName={this.props.scenarioName}
        selectedSource={this.props.sourceNodeSelectedView2}
        selectedDest={this.props.destNodeSelectedView2}
        viewName={view2Name}
        dashboardOptions={this.props.dashboardOptions}
      />
    );

    // Populate Dashboard view list using links from monitoring tab
    var dashboardViewsList = [
      VIEW_NAME_NONE,
      MAP_VIEW,
      NET_TOPOLOGY_VIEW
    ];
    this.populateDashboardList(dashboardViewsList, DEFAULT_DASHBOARD_OPTIONS);
    this.populateDashboardList(dashboardViewsList, this.props.dashboardOptions);

    return (
      <>
        <Grid>
          {!view1Present ? null : (
            <GridCell span={span1} className='chartContainer'>
              <Elevation
                z={2}
                className='idcc-elevation'
                style={{ padding: 10 }}
              >
                {view1}
              </Elevation>
            </GridCell>
          )}

          {!view2Present ? null : (
            <GridCell
              span={span2}
              style={{ marginLeft: -10, paddingLeft: 10 }}
              className='chartContainer'
            >
              <Elevation
                z={2}
                className='idcc-elevation'
                style={{ padding: 10 }}
              >
                {view2}
              </Elevation>
            </GridCell>
          )}
        </Grid>
      </>
    );
  }
}

const mapStateToProps = state => {
  return {
    displayedScenario: state.exec.displayedScenario,
    sourceNodeSelected: state.ui.sourceNodeSelected,
    destNodeSelected: state.ui.destNodeSelected,
    sourceNodeSelectedView1: state.ui.sourceNodeSelectedView1,
    destNodeSelectedView1: state.ui.destNodeSelectedView1,
    sourceNodeSelectedView2: state.ui.sourceNodeSelectedView2,
    destNodeSelectedView2: state.ui.destNodeSelectedView2,
    eventCreationMode: state.exec.eventCreationMode,
    scenarioState: state.exec.state.scenario,
    view1Name: state.ui.dashboardView1,
    view2Name: state.ui.dashboardView2,
    sandboxCfg: state.ui.sandboxCfg,
    dashboardOptions: state.monitor.dashboardOptions
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeView1: name => dispatch(uiExecChangeDashboardView1(name)),
    changeView2: name => dispatch(uiExecChangeDashboardView2(name)),
    changeSandboxCfg: cfg => dispatch(uiExecChangeSandboxCfg(cfg))
  };
};

const ConnectedDashboardContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(DashboardContainer);

export default ConnectedDashboardContainer;
