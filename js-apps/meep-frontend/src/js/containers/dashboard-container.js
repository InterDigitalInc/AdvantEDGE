import _ from 'lodash';
import { connect } from 'react-redux';
import React, { Component, useState }  from 'react';

import { Grid, GridCell, GridInner } from '@rmwc/grid';
import { Elevation } from '@rmwc/elevation';
import { Graph } from 'react-d3-graph';
import ReactDOM from 'react-dom';
import { Button } from '@rmwc/button';
import { Checkbox } from '@rmwc/checkbox';
import { TextField, TextFieldHelperText } from '@rmwc/textfield';
import * as d3 from 'd3';
import axios from 'axios';

// import IDCAreaChart from './idc-area-chart';
import IDCLineChart from './idc-line-chart';
import IDCGraph from './idc-graph';
import IDCAppsView from './idc-apps-view';
import IDSelect from '../components/helper-components/id-select';

import {
  getScenarioNodeChildren,
  isApp
} from '../util/scenario-utils';

import {
  dataAccessorForType,
  dataSetterForType,
  isDataPointOfType
} from '../util/metrics';

import {
  execFakeChangeSelectedDestination,
  execChangeSourceNodeSelected,
  execAddMetricsEpoch
} from '../state/exec';

import {
  LATENCY_METRICS,
  THROUGHPUT_METRICS,
  MOBILITY_EVENT
} from '../meep-constants';

const VIEW_NAME_NONE = 'none';

function colorArray(dataLength) {
  const colorScale = d3.interpolateInferno;
  // const colorScale = d3.interpolateMagma;
  // const colorScale = d3.interpolateCool;
  // const colorScale = d3.interpolateWarm;
  // const colorScale = d3.interpolateCubehelixDefault;
  // interpolateViridis
  // const colorScale = d3.interpolateCubehelixDefault;
  
  let colorArray = [];

  const colorStart = 0.2;
  const colorEnd = 0.8;
  const colorRange = colorEnd - colorStart;
  var intervalSize = colorRange / dataLength;
  for (let i = 0; i < dataLength; i++) {
    const colorPoint = colorStart + i*intervalSize;
    colorArray.push(colorScale(colorPoint));
  }

  return colorArray;
}

const metricsBasePath = 'http://10.3.16.73:30008/v1';

const dataPointFromEpochDataPoints = destinations => sourceNodeId => dataAccessor => epochDataPoints => {
  if (!epochDataPoints.length) {
    return null;
  }
  let dp = {
    date: epochDataPoints[0].timestamp
  };

  const avgForDest = dataPoints => acc => dest => {
    const hasSource = src => p => p.src === src;
    const hasDestination = dest => p => p.dest === dest;
    
    const dataPointsForDestSource = dataPoints
      .filter(hasSource(sourceNodeId))
      .filter(hasDestination(dest));
    const avg = d3.mean(dataPointsForDestSource, acc);
    return avg;
  };
  
  destinations.forEach(dest => {
    dp[dest] = avgForDest(epochDataPoints)(dataAccessor)(dest) || 0;
  });

  return dp;
};

const notNull = x => x;
const epochsToDataPoints = epochs => nb => destinations => dataAccessor => sourceNodeId => {
  const selectedEpochs = epochs.length ? epochs.slice(-nb) : [];

  if (selectedEpochs.length === 0) {
    console.log('epoch length is 0');
  }
  const dataPoints = selectedEpochs.map(dataPointFromEpochDataPoints(destinations)(sourceNodeId)(dataAccessor)).filter(notNull);
  return dataPoints;
};

const ConfigurationView = (props) => {
  return (
    <Grid>
      <GridCell span={2}>
        <IDSelect
          label={'Select View 1'}
          outlined
          options={props.dashboardViewsList}
          onChange={(e) => {
            props.changeView1(e.target.value);
          }}
          // disabled={props.disabled}
          value={props.view1}
        />
      </GridCell>
      <GridCell span={2}>
        <IDSelect
          label={'Select View 2'}
          outlined
          options={props.dashboardViewsList}
          onChange={(e) => {
            props.changeView2(e.target.value);
          }}
          // disabled={props.disabled}
          value={props.view1}
        />
      </GridCell>
      <GridCell span={2}>
        <IDSelect
          label={'Select Source Node'}
          outlined
          options={props.nodeIds}
          onChange={(e) => {
            props.changeSourceNodeSelected(e.target.value);
          }}
          // disabled={props.disabled}
          value={props.sourceNodeSelected ? props.sourceNodeSelected.data.id : ''}
        />
      </GridCell>
      <GridCell span={1}>
        <Checkbox
          checked={props.displayEdgeLabels}
          onChange={() => props.changeDisplayEdgeLabels(!props.displayEdgeLabels)}
        >
                    Show data on edges
        </Checkbox>
      </GridCell>
      <GridCell span={5}>
      </GridCell>
    </Grid>
  );
};

const DATA_CONFIGURATION = 'DATA_CONFIGURATION';
const MAIN_CONFIGURATION = 'MAIN_CONFIGURATION';

const buttonStyles = {
  marginRight: 0
};

const HIERARCHY_VIEW = 'HIERARCHY_VIEW';
const APPS_VIEW = 'APPS_VIEW';
const LATENCY_VIEW = 'LATENCY_VIEW';
const THROUGHPUT_VIEW = 'THROUGHPUT_VIEW';

const DASHBOARD_VIEWS_LIST = [VIEW_NAME_NONE, HIERARCHY_VIEW, APPS_VIEW, LATENCY_VIEW, THROUGHPUT_VIEW];

const ViewForName = (
  {
    apps,
    colorRange,
    width,
    height,
    min,
    max,
    data,
    mobilityEvents,
    dataPoints,
    dataAccessor,
    dataType,
    selectedSource,
    colorForApp,
    changeSourceNodeSelected,
    viewName,
    displayEdgeLabels
  }
) => {
  
  const appIds = apps.map(app => app.data.id);
  switch(viewName) {
  case HIERARCHY_VIEW:
    return (
      <IDCGraph 
        width={width}
        height={600}
      />
    );
  case APPS_VIEW:
    return (
      <IDCAppsView
        apps={apps}
        colorRange={colorRange}
        width={width}
        height={600}
        data={data}
        dataAccessor={dataAccessor}
        dataType={dataType}
        selectedSource={selectedSource}
        colorForApp={colorForApp}
        onNodeClicked={(e) => {
          console.log('Node clicked is: ', e.node);
          changeSourceNodeSelected(e.node);
        }}
        displayEdgeLabels={displayEdgeLabels}
      />
    );
  case LATENCY_VIEW:
    return (
      <IDCLineChart
        data={dataPoints}
        mobilityEvents={mobilityEvents}
        width={width} height={600}
        destinations={appIds}
        colorRange={colorRange}
        selectedSource={selectedSource}
        dataType={dataType}
        // Specify units
        // Specify label
        min={min}
        max={max}
        colorForApp={colorForApp}
      />
    );
  case THROUGHPUT_VIEW:
    return (
      <IDCLineChart
        data={dataPoints}
        mobilityEvents={mobilityEvents}
        width={width} height={600}
        destinations={appIds}
        colorRange={colorRange}
        selectedSource={selectedSource}
        dataType={dataType}
        // Specify units
        // Specify label
        min={min}
        max={max}
        colorForApp={colorForApp}
      />
    );
  default:
    return null;
  }
};

const DashboardConfiguration = (props) => {
  let configurationView = null;
  
  if(props.configurationType) {
    configurationView = (
      <ConfigurationView
        dashboardViewsList={props.dashboardViewsList}
        view1Type={props.view1Type}
        view2Type={props.view2Type}
        changeView1={props.changeView1}
        changeView2={props.changeView2}

        nodeIds={props.nodeIds}
        sourceNodeSelected={props.sourceNodeSelected}
        changeSourceNodeSelected={props.changeSourceNodeSelected}
        changeDisplayEdgeLabels={props.changeDisplayEdgeLabels}
        displayEdgeLabels={props.displayEdgeLabels}
      />
    );
  }

  const buttonConfig = !props.configurationType
    ? (
      <Button outlined style={buttonStyles} onClick={props.displayConfiguration}>
          Configuration
      </Button>
    )
    : null;

  const buttonClose = props.configurationType
    ? (
      <Button outlined style={buttonStyles} onClick={props.hideConfiguration}>
          Close
      </Button>
    )
    : null;

  const backgroundColor = 'ffffff'; // props.configurationType ? '#e4e4e4' : 'ffffff';
  return (
    <div style={{border: '1px solid #e4e4e4', padding: 10, marginBottom: 10, backgroundColor: backgroundColor}}>
      <Grid>
        <GridCell span={10}>
        </GridCell>
        <GridCell span={2}>
          {buttonConfig}
          {buttonClose}
        </GridCell>
      </Grid>

      {configurationView}
    </div>
  );
};

class DashboardContainer extends Component {
  constructor(props) {
    super(props);

    this.state = {
      configurationType: null,
      view1Name: APPS_VIEW,
      view2Name: LATENCY_VIEW,
      sourceNodeId: '',
      nbSecondsToDisplay: 25,
      displayEdgeLabels: false
    };
  }

  componentDidMount() {
    this.epochCount = 0;
    const nextData = () => {
      this.epochCount += 1;
      this.fetchMetrics();
    };
    this.dataTimer = setInterval(nextData, 1000);
  }

  componentWillUnmount() {
    clearInterval(this.dataTimer);
  }

  fetchMetrics() {
    return axios.get(`${metricsBasePath}/metrics?startTime=now-6s&stopTime=now`)
      .then(res => {
        this.props.addMetricsEpoch(res.data.dataResponse || []);
      }).catch((e) => {
        console.log('Error while fetching metrics', e);
      });
  }

  getRoot() {
    return d3.hierarchy(this.props.displayedScenario, getScenarioNodeChildren);
  }

  changeView1(name) {
    this.setState({
      view1Name: name
    });
  }

  changeView2(name) {
    this.setState({
      view2Name: name
    });
  }

  changeDisplayEdgeLabels(val) {
    this.setState({displayEdgeLabels: val});
  }

  render() {
    const root = this.getRoot();
    const nodes = root.descendants();
   
    const apps = nodes.filter(isApp);
    const appIds = apps.map(a => a.data.id);
    const appMap = apps.reduce((acc, app) => {acc[app.data.id] = app; return acc;}, {});
    const colorRange = colorArray(appIds.length);
    const nbEpochs = 25;

    const selectedSource = this.props.sourceNodeSelected ? this.props.sourceNodeSelected.data.id : null;
 
    const showApps = this.props.showAppsView;
    const span = showApps ? 6 : 12;

    const colorForApp = apps.reduce((res, val, i) => {
      return {...res, [val.data.id]: colorRange[i]};
    }, {});

    let lastEpoch = this.props.epochs.length ? this.props.epochs.slice(-1)[0] : [];
    const hasValue = p => {
      const accessor = dataAccessorForType(p.dataType);
      if (! accessor(p)) {
        console.log(`No value for src ${p.src} and dest ${p.dest}`);
      }
      return accessor(p);
    };
    lastEpoch = lastEpoch.filter(hasValue);

    const isDataOfType = type => dataPoint => dataPoint.dataType === type;
    
    const dataTypeForView = view => {
      switch (view) {
      case LATENCY_VIEW:
        return LATENCY_METRICS;
      case THROUGHPUT_VIEW:
        return THROUGHPUT_METRICS;
      default:
        return LATENCY_METRICS;
      }
    };

    const view1DataType = dataTypeForView(this.state.view1Name);
    const view1Accessor = dataAccessorForType(view1DataType);
    const view1DataPoints = epochsToDataPoints(this.props.epochs)(nbEpochs)(appIds)(view1Accessor)(selectedSource);
    const data1 = lastEpoch.filter(isDataOfType(view1DataType));
    const max1 = d3.max(data1, view1Accessor);
    const min1 = d3.min(data1, view1Accessor);

    const view2DataType = dataTypeForView(this.state.view2Name);
    const view2Accessor = dataAccessorForType(view2DataType);
    const view2DataPoints = epochsToDataPoints(this.props.epochs)(nbEpochs)(appIds)(view2Accessor)(selectedSource);
    const data2 = lastEpoch.filter(isDataOfType(view2DataType));

    const extractPointsOfType = type => epoch => epoch.filter(isDataPointOfType(type));
    const extractMobilityEvents = extractPointsOfType(MOBILITY_EVENT);
    const mobilityEvents = this.props.epochs.flatMap(extractMobilityEvents);

    if (mobilityEvents.length) {
      console.log('Some mobility events ...');
    }
    data2.forEach((d) => {
      const dd = view1Accessor(d);
      if (!dd) {
        console.log(`Null data: ${dd}. `);
      }
    });

    const max2 = d3.max(data2, view2Accessor);
    const min2 = d3.min(data2, view2Accessor);
    
    const width = 700;
    const height = 600;

    let span1 = 6;
    let width1 = 700;
    let span2 = 6;
    let width2 = 700;

    if (this.state.view1Name === VIEW_NAME_NONE) {
      span1 = 0;
      width1 = 0;
      span2 = 12;
      width2 = 1200;
    }

    if (this.state.view2Name === VIEW_NAME_NONE) {
      span1 = 12;
      width1 = 1200;
      span2 = 0;
      width2 = 0;
    }

    const view1 = (
      <ViewForName
        apps={apps}
        colorRange={colorRange}
        width={width1}
        height={height}
        data={data1}
        mobilityEvents={mobilityEvents}
        min={min1}
        max={max1}
        dataPoints={view1DataPoints}
        dataAccessor={view1Accessor}
        dataType={view1DataType}
        selectedSource={selectedSource}
        colorForApp={colorForApp}
        changeSourceNodeSelected={(node) => this.props.changeSourceNodeSelected(node)}
        viewName={this.state.view1Name}
        displayEdgeLabels={this.state.displayEdgeLabels}
      />
    );

    const view2 = (
      <ViewForName
        apps={apps}
        colorRange={colorRange}
        width={width2}
        height={height}
        data={data2}
        mobilityEvents={mobilityEvents}
        min={min2}
        max={max2}
        dataPoints={view2DataPoints}
        dataAccessor={view2Accessor}
        dataType={view2DataType}
        selectedSource={selectedSource}
        colorForApp={colorForApp}
        changeSourceNodeSelected={(node) => this.props.changeSourceNodeSelected(node)}
        viewName={this.state.view2Name}
        displayEdgeLabels={this.state.displayEdgeLabels}
      >
      </ViewForName>
    );

    return (
      <>
      <Elevation z={4}
        style={{padding: 10}}
      >
        <DashboardConfiguration
          configurationType={this.state.configurationType}
          displayConfiguration={
            () => {
              this.setState({configurationType: MAIN_CONFIGURATION});
            }}
          hideConfiguration={() => {this.setState({configurationType: ''});}}
          nodeIds={appIds}
          sourceNodeSelected={this.props.sourceNodeSelected}
          changeSourceNodeSelected={(nodeId) => this.props.changeSourceNodeSelected(appMap[nodeId])}
          dashboardViewsList={DASHBOARD_VIEWS_LIST}
          changeView1={(viewName) => this.changeView1(viewName)}
          changeView2={(viewName) => this.changeView2(viewName)}
          displayEdgeLabels={this.state.displayEdgeLabels}
          changeDisplayEdgeLabels={(display) => this.changeDisplayEdgeLabels(display)}
        />
        
        <Grid>
          <GridCell span={span1} style={{marginLeft: -10}}>
            {view1}
          </GridCell>

          <GridCell span={span2} style={{marginLeft: -10}}>
            {view2}
          </GridCell>
        </Grid>
      </Elevation>
      
      </>
    );
  }
}

const mapStateToProps = state => {
  return {
    displayedScenario: state.exec.displayedScenario,
    epochs: state.exec.metrics.epochs,
    sourceNodeSelected: state.exec.metrics.sourceNodeSelected,
    dataTypeSelected: state.exec.metrics.dataTypeSelected
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeSelectedDestination: (dest) => dispatch(execFakeChangeSelectedDestination(dest)),
    changeSourceNodeSelected: (src) => dispatch(execChangeSourceNodeSelected(src)),
    addMetricsEpoch: (epoch) => dispatch(execAddMetricsEpoch(epoch))
  };
};

const ConnectedDashboardContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(DashboardContainer);

export default ConnectedDashboardContainer;