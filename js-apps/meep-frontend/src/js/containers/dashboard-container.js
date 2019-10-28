import { connect } from 'react-redux';
import React, { Component }  from 'react';

import { Grid, GridCell } from '@rmwc/grid';
import { Elevation } from '@rmwc/elevation';
// import ReactDOM from 'react-dom';
import { Button } from '@rmwc/button';
import { Checkbox } from '@rmwc/checkbox';
import { Slider } from '@rmwc/slider';
import moment from 'moment';
import * as d3 from 'd3';

import { blue } from './graph-utils';
import IDCLineChart from './idc-line-chart';
import IDCGraph from './idc-graph';
import IDCAppsView from './idc-apps-view';
import IDSelect from '../components/helper-components/id-select';
import IDCVis from './idc-vis';
import ResizeableContainer from './resizeable-container';

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
  execChangeMetricsTimeIntervalDuration,
  execClearMetricsEpochs
} from '../state/exec';

import {
  uiExecChangeDashboardView1,
  uiExecChangeDashboardView2,
  uiExecExpandDashboardConfig
} from '../state/ui';


import {
  ME_LATENCY_METRICS,
  ME_THROUGHPUT_METRICS,
  ME_MOBILITY_EVENT,
  TYPE_EXEC,
  DASHBOARD_VIEWS_LIST,
  VIEW_NAME_NONE,
  HIERARCHY_VIEW,
  APPS_VIEW,
  LATENCY_VIEW,
  THROUGHPUT_VIEW,
  VIS_VIEW
} from '../meep-constants';

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

const TimeIntervalConfig = (props) => {
  let  PauseResumeButton = null;
  if (props.slidingWindowStopped) {
    PauseResumeButton = () => (
      <Button outlined
        onClick={() => props.startSlidingWindow()}
      >
        RESUME
      </Button>
    );
  } else {
    PauseResumeButton = () => (
      <Button outlined
        onClick={() => props.stopSlidingWindow()}
      >
        PAUSE
      </Button>
    );
  }
  return (
    <div style={{marginTop: 10}}>
      <Grid>
        <GridCell span={3}>
          <div style={{margin:10}}>
            <div>
              <span className="mdc-typography--headline8">Timeframe in secs </span>
            </div>
            <Slider
              value={props.value}
              onChange={e => props.timeIntervalDurationChanged(e.detail.value)}
              discrete
              min={5}
              max={60}
              step={1}
            />
          </div>
          
        </GridCell>
        <GridCell span={1}>

        </GridCell>
        <GridCell span={8}>
          <div style={{margin:10}}>
            <PauseResumeButton />
          </div>
        </GridCell>
      </Grid>
    </div>
    
    
  );
};

const ConfigurationView = (props) => {
  return (
    <>
    <Grid style={{marginBottom: 10}}>
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
      <GridCell span={3}>
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
      <GridCell span={4}>
        <Checkbox
          checked={props.displayEdgeLabels}
          onChange={() => props.changeDisplayEdgeLabels(!props.displayEdgeLabels)}
        >
                    Show data on edges
        </Checkbox>
      </GridCell>
      <GridCell span={1}>
      </GridCell>
    </Grid>
    <TimeIntervalConfig 
      timeIntervalDurationChanged={(value) => {props.timeIntervalDurationChanged(value);}}
      stopSlidingWindow={props.stopSlidingWindow}
      startSlidingWindow={props.startSlidingWindow}
      slidingWindowStopped={props.slidingWindowStopped}
    />
    </>
  );
};

const buttonStyles = {
  marginRight: 0
};

const ViewForName = (
  {
    keyForSvg,
    apps,
    colorRange,
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
      <ResizeableContainer key={keyForSvg}>
        {(width, height) => (
          <IDCGraph 
            keyForSvg={keyForSvg}
            width={width}
            height={height}
          />)}
      </ResizeableContainer>
    );
  case APPS_VIEW:
    return (
      <ResizeableContainer key={keyForSvg}>
        {
          (width, height) => (
            <IDCAppsView
              keyForSvg={keyForSvg}
              apps={apps}
              colorRange={colorRange}
              width={width}
              height={height}
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
          )
        }
      </ResizeableContainer>
      
    );
  case LATENCY_VIEW:
    return (
      <ResizeableContainer key={keyForSvg}>
        {(width, height) => (
          <IDCLineChart
            keyForSvg={keyForSvg}
            data={dataPoints}
            series={series}
            startTime={startTime}
            mobilityEvents={mobilityEvents}
            width={width} height={height}
            destinations={appIds}
            colorRange={colorRange}
            selectedSource={selectedSource}
            dataType={dataType}
            min={min}
            max={max}
            colorForApp={colorForApp}
          />
        )
        }
      </ResizeableContainer>
      
    );
  case THROUGHPUT_VIEW:
    return (
      <ResizeableContainer key={keyForSvg}>
        {
          (width, height) => (
            <IDCLineChart
              keyForSvg={keyForSvg}
              data={dataPoints}
              series={series}
              startTime={startTime}
              mobilityEvents={mobilityEvents}
              width={width} height={height}
              destinations={appIds}
              colorRange={colorRange}
              selectedSource={selectedSource}
              dataType={dataType}
              min={min}
              max={max}
              colorForApp={colorForApp}
            />
          )
        }
      </ResizeableContainer>
    );
  case VIS_VIEW:
    return (
      <ResizeableContainer>
        {
          (width, height) => (
            <IDCVis 
              type={TYPE_EXEC}
              width={width}
              height={height}
              onEditElement={() => {}}
            />
          )
        }
        
      </ResizeableContainer>
    );
  default:
    return null;
  }
};

const DashboardConfiguration = (props) => {
  if (!props.showConfig) {
    return null;
  }

  let configurationView = null;
  
  if(props.dashboardConfigExpanded) {
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
        timeIntervalDurationChanged={props.timeIntervalDurationChanged}
        stopSlidingWindow={props.stopSlidingWindow}
        startSlidingWindow={props.startSlidingWindow}
        slidingWindowStopped={props.slidingWindowStopped}
      />
    );
  }

  const buttonConfig = !props.dashboardConfigExpanded
    ? (
      <Button outlined style={buttonStyles} onClick={() => props.expandDashboardConfig(true)}>
          Open
      </Button>
    )
    : null;

  const buttonClose = props.dashboardConfigExpanded
    ? (
      <Button outlined style={buttonStyles} onClick={() => props.expandDashboardConfig(false)}>
          Close
      </Button>
    )
    : null;

  return (
    <Elevation z={2}
      style={{padding: 10, marginBottom: 10}}
    >
    
      <Grid>
        <GridCell span={11}>
          <div style={{marginBottom:10}}>
            <span className="mdc-typography--headline6">Dashboard Configuration</span>
          </div>
        </GridCell>
        <GridCell span={1}>
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

const eventLogStyle = {
  padding: 10,
  marginTop: 10,
  marginLeft: 10,
  marginRight: 10,
  marginBottom: 10,
  border: '1px solid #e4e4e4',
  count: {color: blue},
  eventName: {color: '#6e6e6e'},
  arrow: {color: '#6e6e6e'},
  element: {color: blue}
};

// let eventCount=0;
const EventLog = (props) => {
  // TODO: generalize function for other types of events.
  // Now it creates a description for Mobility Events
  const descriptionFromEvent = (event) => {
    // eventCount++;
    return (
      <div key={event.mobilityEventIndex}>
        <span style={eventLogStyle.count}>{event.mobilityEventIndex}.</span>
        <span style={eventLogStyle.eventName}>{' Mobility: '}</span>
        <span style={eventLogStyle.element}>{` ${event.src} `}</span>
        <span style={eventLogStyle.arrow}>{' -> '}</span>
        <span style={eventLogStyle.element}>{` ${event.dest}`}</span>
      </div>
    );
  };
  return (
    <>
    <span className="mdc-typography--headline8" style={{marginLeft: 10}}>Events
    </span>
    <div style={eventLogStyle}>
      {props.events.map(descriptionFromEvent)}
    </div>
    </>
  );
};

class DashboardContainer extends Component {
  constructor(props) {
    super(props);

    this.keyForSvg = 0;

    this.state = {
      sourceNodeId: '',
      nbSecondsToDisplay: 25,
      displayEdgeLabels: false,
      slidingWindowStopped: false
    };

    this.epochs = [];
  }

  componentDidMount() {
    
  }

  componentWillUnmount() {
    clearInterval(this.dataTimer);
  }

  getRoot() {
    return d3.hierarchy(this.props.displayedScenario, getScenarioNodeChildren);
  }

  changeDisplayEdgeLabels(val) {
    this.setState({displayEdgeLabels: val});
  }

  changeMetricsTimeIntervalDuration(duration) {
    this.props.changeMetricsTimeIntervalDuration(duration);
  }


  stopSlidingWindow() {
    this.setState({slidingWindowStopped: true});
  }

  startSlidingWindow() {
    this.setState({slidingWindowStopped: false});
  }

  render() {

    let epochs = null;
    if (!this.state.slidingWindowStopped) {
      this.epochs = this.props.epochs.slice();
      epochs = this.epochs;
    } else {
      epochs = this.epochs;
    }

    this.keyForSvg++;
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
        return ME_LATENCY_METRICS;
      case THROUGHPUT_VIEW:
        return ME_THROUGHPUT_METRICS;
      default:
        return ME_LATENCY_METRICS;
      }
    };

    // Determine first and last epochs
    const firstEpoch = epochs.length ? epochs[0] : {
      data: [],
      startTime: null
    };
    let lastEpoch = epochs.length ? epochs.slice(-1)[0] : {
      data: [],
      startTime: null
    };
 
    // Determine startTime of first epoch and endTime of last epoch
    const startTime = firstEpoch.data.length ? firstEpoch.startTime : null;
    const endTime = lastEpoch.data.length ? moment(lastEpoch.startTime).add(1, 'seconds').format(TIME_FORMAT) : null;
    const series = epochsToSeries(epochs, selectedSource);

    const withTypeAndSource = type => source => point => {
      return point.dataType === type && point.src === source;
    };

    // For view 1
    const view1DataType = dataTypeForView(this.props.view1Name);
    const series1 =  filterSeries(appIds)(withTypeAndSource(view1DataType)(selectedSource))(series);
    const lastEpochData1 = lastEpoch.data.filter(isDataOfType(view1DataType));

    // For view2
    const view2DataType = dataTypeForView(this.props.view2Name);
    const series2 =  filterSeries(appIds)(withTypeAndSource(view2DataType)(selectedSource))(series);
    const lastEpochData2 = lastEpoch.data.filter(isDataOfType(view2DataType));

    // Mobility events
    const extractPointsOfType = type => epoch => epoch.data.filter(isDataPointOfType(type));
    const extractMobilityEvents = extractPointsOfType(ME_MOBILITY_EVENT);
    const mobilityEvents = epochs.flatMap(extractMobilityEvents);

    if (mobilityEvents.length) {
      console.log('Some mobility events ...');
    }
  
    
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
        keyForSvg={this.keyForSvg}
        apps={apps}
        colorRange={colorRange}
        // width={width1}
        // height={height}
        data={lastEpochData1}
        series={series1}
        startTime={startTime}
        endTime={endTime}
        mobilityEvents={mobilityEvents}
        dataType={view1DataType}
        selectedSource={selectedSource}
        colorForApp={colorForApp}
        changeSourceNodeSelected={(node) => this.props.changeSourceNodeSelected(node)}
        viewName={this.props.view1Name}
        displayEdgeLabels={this.state.displayEdgeLabels}
      />
    );

    const view2 = (
      <ViewForName
        keyForSvg={this.keyForSvg}
        apps={apps}
        colorRange={colorRange}
        // width={width2}
        // height={height}
        data={lastEpochData2}
        series={series2}
        startTime={startTime}
        endTime={endTime}
        mobilityEvents={mobilityEvents}
        dataType={view2DataType}
        selectedSource={selectedSource}
        colorForApp={colorForApp}
        changeSourceNodeSelected={(node) => this.props.changeSourceNodeSelected(node)}
        viewName={this.props.view2Name}
        displayEdgeLabels={this.state.displayEdgeLabels}
      >
      </ViewForName>
    );

    const EventLogComponent = () => (
      <EventLog 
        events={mobilityEvents}
      />
    );

    return (
      <>
      
        <DashboardConfiguration
          showConfig={this.props.showConfig}
          dashboardConfigExpanded={this.props.dashboardConfigExpanded}
          expandDashboardConfig={(show) => this.props.expandDashboardConfig(show)}
          nodeIds={appIds}
          sourceNodeSelected={this.props.sourceNodeSelected}
          changeSourceNodeSelected={(nodeId) => this.props.changeSourceNodeSelected(appMap[nodeId])}
          timeIntervalDurationChanged={(duration) => {this.changeMetricsTimeIntervalDuration(duration);}}
          stopSlidingWindow={() => this.stopSlidingWindow()}
          startSlidingWindow={() => this.startSlidingWindow()}
          slidingWindowStopped={this.state.slidingWindowStopped}
          dashboardViewsList={DASHBOARD_VIEWS_LIST}
          changeView1={(viewName) => this.props.changeView1(viewName)}
          changeView2={(viewName) => this.props.changeView2(viewName)}
          displayEdgeLabels={this.state.displayEdgeLabels}
          changeDisplayEdgeLabels={(display) => this.changeDisplayEdgeLabels(display)}
        />
        
        <Grid>

          {!view1Present ? null : (
            <GridCell span={span1} style={{paddingRight: 10}} className='chartContainer'>
              <Elevation z={2}
                style={{padding: 10}}
              >
                {view1}
              </Elevation>
            </GridCell>
          )}
          
          {!view2Present ? null : (
            <GridCell span={span2} style={{marginLeft: -10, paddingLeft: 10}} className='chartContainer'>
              <Elevation z={2}
                style={{padding: 10}}
              >
                {view2}
                <EventLogComponent />
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
    epochs: state.exec.metrics.epochs,
    sourceNodeSelected: state.exec.metrics.sourceNodeSelected,
    dataTypeSelected: state.exec.metrics.dataTypeSelected,
    eventCreationMode: state.exec.eventCreationMode,
    metricsTimeIntervalDuration: state.exec.metrics.timeIntervalDuration,
    scenarioState: state.exec.state.scenario,
    showConfig: state.ui.showDashboardConfig,
    dashboardConfigExpanded: state.ui.dashboardConfigExpanded,
    view1Name: state.ui.dashboardView1,
    view2Name: state.ui.dashboardView2
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeSelectedDestination: (dest) => dispatch(execFakeChangeSelectedDestination(dest)),
    changeSourceNodeSelected: (src) => dispatch(execChangeSourceNodeSelected(src)),
    changeMetricsTimeIntervalDuration: (duration) => dispatch(execChangeMetricsTimeIntervalDuration(duration)),
    clearMetricsEpochs: () => dispatch(execClearMetricsEpochs()),
    changeView1: (name) => dispatch(uiExecChangeDashboardView1(name)),
    changeView2: (name) => dispatch(uiExecChangeDashboardView2(name)) ,
    expandDashboardConfig: (expand) => dispatch(uiExecExpandDashboardConfig(expand))
  };
};

const ConnectedDashboardContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(DashboardContainer);

export default ConnectedDashboardContainer;