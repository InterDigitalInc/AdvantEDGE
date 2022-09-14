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
import IDCSeq from '../idc-seq';
import IDCDataflow from '../idc-dataflow';
import Iframe from 'react-iframe';

import {
  TYPE_EXEC,
  VIEW_NAME_NONE,
  VIEW_1,
  VIEW_2,
  MAP_VIEW,
  NET_TOPOLOGY_VIEW,
  SEQ_DIAGRAM_VIEW,
  DATAFLOW_DIAGRAM_VIEW,
  DEFAULT_DASHBOARD_OPTIONS
} from '../../meep-constants';

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
  viewCfg,
  dashboardOptions
}) => {

  // Handle Map view
  if (viewCfg.viewType === MAP_VIEW) {
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
  if (viewCfg.viewType === NET_TOPOLOGY_VIEW) {
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

  // Handle Sequence Diagram view
  if (viewCfg.viewType === SEQ_DIAGRAM_VIEW) {
    return (
      <div style={styles.dashboard}>
        <IDCSeq
          participants={viewCfg.participants}
        />
      </div>
    );
  }

  // Handle Sequence Diagram view
  if (viewCfg.viewType === DATAFLOW_DIAGRAM_VIEW) {
    return (
      <div style={styles.dashboard}>
        <IDCDataflow/>
      </div>
    );
  }

  // Get URL from Monitoring page dashboard options
  var selectedUrl = getUrl(viewCfg.viewType, DEFAULT_DASHBOARD_OPTIONS);
  if (selectedUrl === '') {
    selectedUrl = getUrl(viewCfg.viewType, dashboardOptions);
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
      url.searchParams.append('var-src', viewCfg.sourceNodeSelected);
      url.searchParams.append('var-dest', viewCfg.destNodeSelected);
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

  componentWillUnmount() {
    clearInterval(this.dataTimer);
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

  populateDashboardList(dashboardViewsList, dashboardOptions) {
    for (var i = 0; i < dashboardOptions.length; i++) {
      var dashboard = dashboardOptions[i];
      if ((dashboard.label !== '') && (dashboard.value !== '') && (dashboard.value.indexOf(showInExecStr) !== -1)) {
        dashboardViewsList.push(dashboard.label);
      }
    }
  }

  render() {
    this.keyForSvg++;

    const view1Cfg = this.getViewCfg(VIEW_1);
    const view2Cfg = this.getViewCfg(VIEW_2);
    const view1Present = (view1Cfg && view1Cfg.viewType !== VIEW_NAME_NONE);
    const view2Present = (view2Cfg && view2Cfg.viewType !== VIEW_NAME_NONE);

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
        viewCfg={view1Cfg}
        dashboardOptions={this.props.dashboardOptions}
      />
    );

    const view2 = (
      <ViewForName
        sandboxName={this.props.sandbox}
        scenarioName={this.props.scenarioName}
        viewCfg={view2Cfg}
        dashboardOptions={this.props.dashboardOptions}
      />
    );

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
    eventCreationMode: state.exec.eventCreationMode,
    scenarioState: state.exec.state.scenario,
    sandboxCfg: state.ui.sandboxCfg,
    dashboardOptions: state.monitor.dashboardOptions
  };
};

const mapDispatchToProps = () => {
  return {
    
  };
};

const ConnectedDashboardContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(DashboardContainer);

export default ConnectedDashboardContainer;
