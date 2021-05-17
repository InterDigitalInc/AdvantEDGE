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
import { Button } from '@rmwc/button';
import { Checkbox } from '@rmwc/checkbox';
import * as d3 from 'd3';

import IDSelect from '../../components/helper-components/id-select';
import IDCMap from '../idc-map';
import IDCVis from '../idc-vis';
import Iframe from 'react-iframe';

import { getScenarioNodeChildren, isApp, isUe, isPoa } from '../../util/scenario-utils';

import {
  uiExecChangeDashboardView1,
  uiExecChangeDashboardView2,
  uiExecChangeSourceNodeSelected,
  uiExecChangeDestNodeSelected,
  uiExecChangeSandboxCfg
} from '../../state/ui';

import {
  TYPE_EXEC,
  VIEW_NAME_NONE,
  MAP_VIEW,
  NET_TOPOLOGY_VIEW,
  DEFAULT_DASHBOARD_OPTIONS,
  EXEC_BTN_DASHBOARD_BTN_CLOSE
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

const ConfigurationView = props => {
  return (
    <>
      <Grid style={{ marginBottom: 10 }}>
        <GridCell span={2}>
          <IDSelect
            label={'View 1'}
            outlined
            options={props.dashboardViewsList}
            onChange={e => {
              props.changeView1(e.target.value);
            }}
            value={props.view1Name}
          />
        </GridCell>
        <GridCell span={2}>
          <IDSelect
            label={'View 2'}
            outlined
            options={props.dashboardViewsList}
            onChange={e => {
              props.changeView2(e.target.value);
            }}
            value={props.view2Name}
          />
        </GridCell>
        <GridCell span={2}>
          <IDSelect
            label={'Source Node'}
            outlined
            options={props.srcIds}
            onChange={e => {
              props.changeSourceNodeSelected(e.target.value);
            }}
            value={props.sourceNodeSelected}
          />
        </GridCell>
        <GridCell span={2}>
          <IDSelect
            label={'Destination Node'}
            outlined
            options={props.destIds}
            onChange={e => {
              props.changeDestNodeSelected(e.target.value);
            }}
            value={props.destNodeSelected}
          />
        </GridCell>
        <GridCell span={2}>
          <Checkbox
            checked={props.showApps}
            onChange={e => props.changeShowApps(e.target.checked)}
          >
            Show Apps
          </Checkbox>
        </GridCell>
      </Grid>
    </>
  );
};

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

const DashboardConfiguration = props => {
  if (!props.dashCfgMode) {
    return null;
  }

  let configurationView = null;

  configurationView = (
    <ConfigurationView
      dashboardViewsList={props.dashboardViewsList}
      view1Name={props.view1Name}
      view2Name={props.view2Name}
      changeView1={props.changeView1}
      changeView2={props.changeView2}
      srcIds={props.srcIds}
      destIds={props.destIds}
      sourceNodeSelected={props.sourceNodeSelected}
      destNodeSelected={props.destNodeSelected}
      changeSourceNodeSelected={props.changeSourceNodeSelected}
      changeDestNodeSelected={props.changeDestNodeSelected}
      changeShowApps={props.changeShowApps}
      showApps={props.showApps}
    />
  );
  return (
    <Elevation
      z={2}
      className='idcc-elevation'
      style={{ padding: 10, marginBottom: 10 }}
    >
      <Grid>
        <GridCell span={6}>
          <div style={{ marginBottom: 10 }}>
            <span className='mdc-typography--headline6'>
              Dashboard
            </span>
          </div>
        </GridCell>
        <GridCell span={6}>
          <div align={'right'}>
            <Button
              outlined
              style={styles.button}
              onClick={() => props.onCloseDashCfg()}
              data-cy={EXEC_BTN_DASHBOARD_BTN_CLOSE}
            >
            Close
            </Button>
          </div>
        </GridCell>
      </Grid>
      {configurationView}
    </Elevation>
  );
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

  getRoot() {
    return d3.hierarchy(this.props.displayedScenario, getScenarioNodeChildren);
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
    const root = this.getRoot();
    const nodes = root.descendants();
    const apps = nodes.filter(isApp);
    const appIds = apps.map(a => a.data.name);
    const ues = nodes.filter(isUe);
    const ueIds = ues.map(a => a.data.name);
    const poas = nodes.filter(isPoa);
    const poaIds = poas.map(a => a.data.name);
    const srcIds = appIds.concat(ueIds).sort();
    const destIds = appIds.concat(poaIds).sort();

    appIds.unshift('None');
    srcIds.unshift('None');
    destIds.unshift('None');

    const selectedSource = srcIds.includes(this.props.sourceNodeSelected) ? this.props.sourceNodeSelected : 'None';
    const selectedDest = destIds.includes(this.props.destNodeSelected) ? this.props.destNodeSelected : 'None';
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
        selectedSource={selectedSource}
        selectedDest={selectedDest}
        viewName={view1Name}
        dashboardOptions={this.props.dashboardOptions}
      />
    );

    const view2 = (
      <ViewForName
        sandboxName={this.props.sandbox}
        scenarioName={this.props.scenarioName}
        selectedSource={selectedSource}
        selectedDest={selectedDest}
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
        <DashboardConfiguration
          dashCfgMode={this.props.dashCfgMode}
          onCloseDashCfg={this.props.onCloseDashCfg}
          srcIds={srcIds}
          destIds={destIds}
          view1Name={view1Name}
          view2Name={view2Name}
          sourceNodeSelected={selectedSource}
          destNodeSelected={selectedDest}
          changeSourceNodeSelected={this.props.changeSourceNodeSelected}
          changeDestNodeSelected={this.props.changeDestNodeSelected}
          dashboardViewsList={dashboardViewsList}
          changeView1={this.changeView1}
          changeView2={this.changeView2}
          changeShowApps={this.changeShowApps}
          showApps={this.props.showApps}
        />

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
    changeSourceNodeSelected: src => dispatch(uiExecChangeSourceNodeSelected(src)),
    changeDestNodeSelected: dest => dispatch(uiExecChangeDestNodeSelected(dest)),
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
