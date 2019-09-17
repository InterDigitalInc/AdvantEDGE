import { connect } from 'react-redux';
import React, { Component }  from 'react';

import { Grid, GridCell } from '@rmwc/grid';
import { Elevation } from '@rmwc/elevation';
// import ReactDOM from 'react-dom';
import { Button } from '@rmwc/button';
import { Checkbox } from '@rmwc/checkbox';
import moment from 'moment';
import * as d3 from 'd3';
import axios from 'axios';

import IDCLineChart from './idc-line-chart';
import IDCGraph from './idc-graph';
import IDCAppsView from './idc-apps-view';
import IDSelect from '../components/helper-components/id-select';

import {
  idlog
} from '../util/functional';

import {
  getScenarioNodeChildren,
  isApp
} from '../util/scenario-utils';

import {
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
const TIME_FORMAT = moment.HTML5_FMT.DATETIME_LOCAL_MS;

function colorArray(dataLength) {
  const colorScale = d3.interpolateInferno;
  // Other possible color scales:
  // const colorScale = d3.interpolateMagma;
  // const colorScale = d3.interpolateCool;
  // const colorScale = d3.interpolateWarm;
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

  return (
    <Elevation z={2}
      style={{padding: 10, marginBottom: 10}}
    >
      <Grid>
        <GridCell span={10}>
        </GridCell>
        <GridCell span={2}>
          {buttonConfig}
          {buttonClose}
        </GridCell>
      </Grid>
      {configurationView}
    </Elevation>
     
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
    const startTime = moment().utc().add(-7, 'seconds').format(TIME_FORMAT);
    const stopTime = moment().utc().add(-6, 'seconds').format(TIME_FORMAT);
    return axios.get(`${metricsBasePath}/metrics?startTime=${startTime}&stopTime=${stopTime}`)
      .then(res => {

        let epoch = {
          data: res.data.logResponse || [],
          startTime: startTime
        };
  
        this.props.addMetricsEpoch(epoch);
      }).catch((e) => {
        idlog('Error while fetching metrics')(e);
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

    const selectedSource = this.props.sourceNodeSelected ? this.props.sourceNodeSelected.data.id : null;

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
    const endTime = lastEpoch.data.length ? moment(lastEpoch.startTime).add(1, 'seconds').format(TIME_FORMAT) : null;
    const series = epochsToSeries(this.props.epochs, selectedSource);

    const withTypeAndSource = type => source => point => {
      return point.dataType === type && point.src === source;
    };

    // For view 1
    const view1DataType = dataTypeForView(this.state.view1Name);
    const series1 =  filterSeries(appIds)(withTypeAndSource(view1DataType)(selectedSource))(series);
    const lastEpochData1 = lastEpoch.data.filter(isDataOfType(view1DataType));

    // For view2
    const view2DataType = dataTypeForView(this.state.view2Name);
    const series2 =  filterSeries(appIds)(withTypeAndSource(view2DataType)(selectedSource))(series);
    const lastEpochData2 = lastEpoch.data.filter(isDataOfType(view2DataType));

    // Mobility events
    const extractPointsOfType = type => epoch => epoch.data.filter(isDataPointOfType(type));
    const extractMobilityEvents = extractPointsOfType(MOBILITY_EVENT);
    const mobilityEvents = this.props.epochs.flatMap(extractMobilityEvents);

    if (mobilityEvents.length) {
      // console.log('Some mobility events ...');
    }
  
    
    const height = 600;

    let span1 = 6;
    let width1 = 700;
    let width2 = 700;

    if (this.state.view1Name === VIEW_NAME_NONE) {
      span1 = 0;
      width1 = 0;
      width2 = 1400;
    }

    if (this.state.view2Name === VIEW_NAME_NONE) {
      span1 = 12;
      width1 = 1400;
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
          
          <GridCell span={span1} style={{paddingRight: 10}}>
            <Elevation z={2}
              style={{padding: 10}}
            >
              {view1}
            </Elevation>
          </GridCell>
          
          <GridCell span={span1} style={{marginLeft: -10, paddingLeft: 10}}>
            <Elevation z={2}
              style={{padding: 10}}
            >
              {view2}
            </Elevation>
          </GridCell>
        </Grid>
      
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