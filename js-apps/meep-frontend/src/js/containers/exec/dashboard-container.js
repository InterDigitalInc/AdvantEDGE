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
import { Grid, GridCell } from '@rmwc/grid';
import { Elevation } from '@rmwc/elevation';
// import ReactDOM from 'react-dom';
import { Button } from '@rmwc/button';
import { Checkbox } from '@rmwc/checkbox';
import * as d3 from 'd3';

import IDSelect from '../../components/helper-components/id-select';
import IDCVis from '../idc-vis';
import Iframe from 'react-iframe';

import { getScenarioNodeChildren, isApp } from '../../util/scenario-utils';

import {
  execChangeSourceNodeSelected,
  execChangeDestNodeSelected
} from '../../state/exec';

import {
  uiExecChangeDashboardView1,
  uiExecChangeDashboardView2
} from '../../state/ui';

import {
  TYPE_EXEC,
  DASHBOARD_VIEWS_LIST,
  VIEW_NAME_NONE,
  METRICS_VIEW,
  PLAYER_METRICS_1_VIEW,
  PLAYER_METRICS_1_NOFPS_VIEW,
  PLAYER_METRICS_2_VIEW,
  PLAYER_METRICS_2_NOFPS_VIEW,
  VIS_VIEW
} from '../../meep-constants';

const greyColor = 'grey';

const styles = {
  button: {
    marginRight: 0
  },
  slider: {
    container: {
      marginTop: 10,
      marginBottom: 10,
      color: greyColor
    },
    boundaryValues: {
      marginTop: 15
    },
    title: {
      marginBottom: 0
    }
  }
};

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
            options={props.nodeIds}
            onChange={e => {
              props.changeSourceNodeSelected(e.target.value);
            }}
            value={
              props.sourceNodeSelected ? props.sourceNodeSelected.data.id : ''
            }
          />
        </GridCell>
        <GridCell span={2}>
          <IDSelect
            label={'Destination Node'}
            outlined
            options={props.nodeIds}
            onChange={e => {
              props.changeDestNodeSelected(e.target.value);
            }}
            value={
              props.destNodeSelected ? props.destNodeSelected.data.id : ''
            }
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

const ViewForName = ({
  scenarioName,
  selectedSource,
  selectedDest,
  viewName
}) => {

  // Remove '-' from scenario name
  var scenario = scenarioName.replace(/-/g, '');

  const datasource = '&var-datasource=meep-influxdb';
  const database = '&var-database=' + scenario;
  const refreshInterval = '&refresh=1s';
  const srcApp = '&var-src=' + selectedSource;
  const destApp = '&var-dest=' + selectedDest;
  const viewMode = '&kiosk';
  const theme = '&theme=light';

  const metricsDashboard = 'http://' + location.hostname + ':30009/d/100/metrics-dashboard?orgId=1';
  const metricsDashboardUrl = metricsDashboard + datasource + database + refreshInterval + srcApp + destApp + viewMode + theme;
  
  const playerMetrics1Dashboard = 'http://' + location.hostname + ':30009/d/MWC2020-P12M-1/player-metrics-1?orgId=1';
  const playerMetrics1DashboardUrl = playerMetrics1Dashboard + datasource + database + refreshInterval + viewMode + theme;

  const playerMetrics1NoFpsDashboard = 'http://' + location.hostname + ':30009/d/MWC2020-P12M-1-nofps/player-metrics-1-nofps?orgId=1';
  const playerMetrics1NoFpsDashboardUrl = playerMetrics1NoFpsDashboard + datasource + database + refreshInterval + viewMode + theme;
  
  const playerMetrics2Dashboard = 'http://' + location.hostname + ':30009/d/MWC2020-P12M-2/player-metrics-2?orgId=1';
  const playerMetrics2DashboardUrl = playerMetrics2Dashboard + datasource + database + refreshInterval + viewMode + theme;

  const playerMetrics2NoFpsDashboard = 'http://' + location.hostname + ':30009/d/MWC2020-P12M-2nofps/player-metrics-2-nofps?orgId=1';
  const playerMetrics2NoFpsDashboardUrl = playerMetrics2NoFpsDashboard + datasource + database + refreshInterval + viewMode + theme;

  switch (viewName) {
  case METRICS_VIEW:
    return (
      <div style={{ height: '80vh' }}>
        <Iframe
          url={metricsDashboardUrl}
          id="myId"
          display="initial"
          position="relative"
          allowFullScreen
          width='100%'
          height='100%'
        />
      </div>
    );
  case PLAYER_METRICS_1_VIEW:
    return (
      <div style={{ height: '80vh' }}>
        <Iframe
          url={playerMetrics1DashboardUrl}
          id="myId"
          display="initial"
          position="relative"
          allowFullScreen
          width='100%'
          height='100%'
        />
      </div>
    );
  case PLAYER_METRICS_1_NOFPS_VIEW:
    return (
      <div style={{ height: '80vh' }}>
        <Iframe
          url={playerMetrics1NoFpsDashboardUrl}
          id="myId"
          display="initial"
          position="relative"
          allowFullScreen
          width='100%'
          height='100%'
        />
      </div>
    );
  case PLAYER_METRICS_2_VIEW:
    return (
      <div style={{ height: '80vh' }}>
        <Iframe
          url={playerMetrics2DashboardUrl}
          id="myId"
          display="initial"
          position="relative"
          allowFullScreen
          width='100%'
          height='100%'
        />
      </div>
    );
  case PLAYER_METRICS_2_NOFPS_VIEW:
    return (
      <div style={{ height: '80vh' }}>
        <Iframe
          url={playerMetrics2NoFpsDashboardUrl}
          id="myId"
          display="initial"
          position="relative"
          allowFullScreen
          width='100%'
          height='100%'
        />
      </div>
    );
  case VIS_VIEW:
    return (
      <div style={{ height: '80vh' }}>
        <IDCVis
          type={TYPE_EXEC}
          width='100%'
          height='100%'
          onEditElement={() => { }}
        />
      </div>
    );
  default:
    return null;
  }
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
      nodeIds={props.nodeIds}
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
      className="component-style"
      style={{ padding: 10, marginBottom: 10 }}
    >
      <Grid>
        <GridCell span={11}>
          <div style={{ marginBottom: 10 }}>
            <span className="mdc-typography--headline6">
              Dashboard Configuration
            </span>
          </div>
        </GridCell>
        <GridCell span={1}>
          <Button
            outlined
            style={styles.button}
            onClick={() => props.onCloseDashCfg()}
          >
            Close
          </Button>
        </GridCell>
      </Grid>
      {configurationView}
    </Elevation>
  );
};

class DashboardContainer extends Component {
  constructor(props) {
    super(props);
    this.state = {
      sourceNodeId: ''
    };
  }

  componentDidMount() { }

  componentWillUnmount() {
    clearInterval(this.dataTimer);
  }

  getRoot() {
    return d3.hierarchy(this.props.displayedScenario, getScenarioNodeChildren);
  }

  changeShowApps(checked) {
    this.props.onShowAppsChanged(checked);
  }

  render() {
    this.keyForSvg++;
    const root = this.getRoot();
    const nodes = root.descendants();

    const apps = nodes.filter(isApp);
    const appIds = apps.map(a => a.data.id);
    const appMap = apps.reduce((acc, app) => {
      acc[app.data.id] = app;
      return acc;
    }, {});

    const selectedSource = this.props.sourceNodeSelected
      ? this.props.sourceNodeSelected.data.id
      : null;

    const selectedDest = this.props.destNodeSelected
      ? this.props.destNodeSelected.data.id
      : null;

    // For view 1
    const view1Name = this.props.view1Name;

    // For view2
    const view2Name = this.props.view2Name;

    // const height = 600;

    let span1 = 12;
    let span2 = 12;
    // let width1 = 700;
    // let width2 = 700;

    const view1Present = this.props.view1Name !== VIEW_NAME_NONE;
    const view2Present = this.props.view2Name !== VIEW_NAME_NONE;

    if (view1Present && view2Present) {
      span1 = 6;
      span2 = 6;
    } else if (!view1Present && !view2Present) {
      span1 = 0;
      span2 = 0;
    }

    const view1 = (
      <ViewForName
        scenarioName={this.props.scenarioName}
        selectedSource={selectedSource}
        selectedDest={selectedDest}
        viewName={view1Name}
      />
    );

    const view2 = (
      <ViewForName
        scenarioName={this.props.scenarioName}
        selectedSource={selectedSource}
        selectedDest={selectedDest}
        viewName={view2Name}
      />
    );

    return (
      <>
        <DashboardConfiguration
          dashCfgMode={this.props.dashCfgMode}
          onCloseDashCfg={this.props.onCloseDashCfg}
          nodeIds={appIds}
          view1Name={view1Name}
          view2Name={view2Name}
          sourceNodeSelected={this.props.sourceNodeSelected}
          destNodeSelected={this.props.destNodeSelected}
          changeSourceNodeSelected={nodeId =>
            this.props.changeSourceNodeSelected(appMap[nodeId])
          }
          changeDestNodeSelected={nodeId =>
            this.props.changeDestNodeSelected(appMap[nodeId])
          }
          dashboardViewsList={DASHBOARD_VIEWS_LIST}
          changeView1={viewName => this.props.changeView1(viewName)}
          changeView2={viewName => this.props.changeView2(viewName)}
          changeShowApps={checked => this.changeShowApps(checked)}
          showApps={this.props.showApps}
        />

        <Grid>
          {!view1Present ? null : (
            <GridCell span={span1} className="chartContainer">
              <Elevation
                z={2}
                className="component-style"
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
              className="chartContainer"
            >
              <Elevation
                z={2}
                className="component-style"
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
    sourceNodeSelected: state.exec.metrics.sourceNodeSelected,
    destNodeSelected: state.exec.metrics.destNodeSelected,
    eventCreationMode: state.exec.eventCreationMode,
    scenarioState: state.exec.state.scenario,
    view1Name: state.ui.dashboardView1,
    view2Name: state.ui.dashboardView2
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeSourceNodeSelected: src => dispatch(execChangeSourceNodeSelected(src)),
    changeDestNodeSelected: dest => dispatch(execChangeDestNodeSelected(dest)),
    changeView1: name => dispatch(uiExecChangeDashboardView1(name)),
    changeView2: name => dispatch(uiExecChangeDashboardView2(name))
  };
};

const ConnectedDashboardContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(DashboardContainer);

export default ConnectedDashboardContainer;
