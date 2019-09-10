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
import moment from 'moment';
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
  isDataPointOfType,
  valueOfPoint
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

// const dataPointFromEpochDataPoints = destinations => sourceNodeId => dataAccessor => epochDataPoints => {
//   if (!epochDataPoints.length) {
//     return null;
//   }
//   let dp = {
//     date: epochDataPoints[0].timestamp
//   };

//   const avgForDest = dataPoints => acc => dest => {
//     const hasSource = src => p => p.src === src;
//     const hasDestination = dest => p => p.dest === dest;
    
//     const dataPointsForDestSource = dataPoints
//       .filter(hasSource(sourceNodeId))
//       .filter(hasDestination(dest));
//     const avg = d3.mean(dataPointsForDestSource, acc);
//     return avg;
//   };
  
//   destinations.forEach(dest => {
//     dp[dest] = avgForDest(epochDataPoints)(dataAccessor)(dest) || 0;
//   });

//   return dp;
// };

const notNull = x => x;

const buildSeriesFromEpoch = (series, epoch) => {
  epoch.data.forEach(p => {
    if (! series[p.dest]) {
      series[p.dest] = [];
    }
    series[p.dest].push(p);
  });

  return series;
};

const epochsToSeries = (epochs) => {
  let series = epochs.reduce((s, current) => {
    return buildSeriesFromEpoch(s, current);
  }, {});
  return series;
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
    series,
    startTime,
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
        series={series}
        startTime={startTime}
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
        series={series}
        startTime={startTime}
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
        series={series}
        startTime={startTime}
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

const filterSeries = keys => filter => series => {
  let newSeries = {};
  keys.forEach(key => {
    if (series[key]) {
      newSeries[key] = removeDuplicatePoints(series[key].filter(filter));
    }
  });

  return newSeries;
};

const removeDuplicatePoints = sequence => {
  let timestampsMap = {};
  let newSequence = [];
  sequence.forEach(p => {
    if (!timestampsMap[p.timestamp]) {
      timestampsMap[p.timestamp] = true;
      newSequence.push(p);
    }
  });

  return newSequence;
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
    const startTime = moment().add(-7, 'seconds').format(moment.HTML5_FMT.DATETIME_LOCAL_MS);
    const stopTime = moment().add(-6, 'seconds').format(moment.HTML5_FMT.DATETIME_LOCAL_MS);
    // const now = moment().format(moment.HTML5_FMT.DATETIME_LOCAL_MS);
    return axios.get(`${metricsBasePath}/metrics?startTime=${startTime}&stopTime=${stopTime}`)
      .then(res => {

        let epoch = {
          data: res.data.logResponse || [], //.sort((a, b) => new Date(a).getTime() - new Date(b).getTime() || []),
          startTime: startTime
        };
  
        this.props.addMetricsEpoch(epoch);
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

    // Determine first and last epochs
    const firstEpoch = this.props.epochs.length ? this.props.epochs[0] : {
      data: [],
      startTime: null
    };
    let lastEpoch = this.props.epochs.length ? this.props.epochs.slice(-1)[0] : {
      data: [],
      startTime: null
    };
 
    // Determine startTime of first epoch and endTime of last epoch
    const startTime = firstEpoch.data.length ? firstEpoch.startTime : null;
    const endTime = lastEpoch.data.length ? new Date(new Date(lastEpoch.startTime).getTime() + 1000).toString() : null;
    const series = epochsToSeries(this.props.epochs, selectedSource);

    const withTypeAndSource = type => source => point => {
      return point.dataType === type && point.src === source;
    };

    // For view 1
    const view1DataType = dataTypeForView(this.state.view1Name);
    const series1 =  filterSeries(appIds)(withTypeAndSource(view1DataType)(selectedSource))(series);
    const lastEpochData1 = lastEpoch.data.filter(isDataOfType(view1DataType));
    // const max1 = d3.max(data1, p => p.value);
    // const min1 = d3.min(data1, p => p.value);

    // For view2
    const view2DataType = dataTypeForView(this.state.view2Name);
    const series2 =  filterSeries(appIds)(withTypeAndSource(view2DataType)(selectedSource))(series);
    const lastEpochData2 = lastEpoch.data.filter(isDataOfType(view2DataType));

    // Mobility events
    const extractPointsOfType = type => epoch => epoch.data.filter(isDataPointOfType(type));
    const extractMobilityEvents = extractPointsOfType(MOBILITY_EVENT);
    const mobilityEvents = this.props.epochs.flatMap(extractMobilityEvents);

    if (mobilityEvents.length) {
      console.log('Some mobility events ...');
    }
    
    // const max2 = d3.max(data2, view2Accessor);
    // const min2 = d3.min(data2, view2Accessor);
    
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
        data={lastEpochData1}
        series={series1}
        startTime={startTime}
        endTime={endTime}
        mobilityEvents={mobilityEvents}
        // min={min1}
        // max={max1}
        // dataAccessor={view1Accessor}
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
        data={lastEpochData2}
        series={series2}
        startTime={startTime}
        endTime={endTime}
        mobilityEvents={mobilityEvents}
        // min={min2}
        // max={max2}
        // dataAccessor={view2Accessor}
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