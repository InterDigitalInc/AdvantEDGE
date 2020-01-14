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
import { Slider } from '@rmwc/slider';
import moment from 'moment';
import * as d3 from 'd3';

import { blue } from '../graph-utils';
import IDCLineChart from './idc-line-chart';
import IDCAppsView from '../idc-apps-view';
import IDSelect from '../../components/helper-components/id-select';
import IDCVis from '../idc-vis';
import Iframe from 'react-iframe';
import ResizeableContainer from '../resizeable-container';

import { getScenarioNodeChildren, isApp } from '../../util/scenario-utils';

import { isDataPointOfType } from '../../util/metrics';

import {
  execFakeChangeSelectedDestination,
  execChangeSourceNodeSelected,
  execChangeDestNodeSelected,
  execChangeMetricsTimeIntervalDuration,
  execClearMetricsEpochs
} from '../../state/exec';

import {
  uiExecChangeDashboardView1,
  uiExecChangeDashboardView2
} from '../../state/ui';

import {
  ME_LATENCY_METRICS,
  ME_THROUGHPUT_METRICS,
  ME_MOBILITY_EVENT,
  TYPE_EXEC,
  DASHBOARD_VIEWS_LIST,
  VIEW_NAME_NONE,
  APPS_VIEW,
  METRICS_VIEW,
  LATENCY_VIEW,
  THROUGHPUT_VIEW,
  VIS_VIEW
} from '../../meep-constants';

const TIME_FORMAT = moment.HTML5_FMT.DATETIME_LOCAL_MS;
const MIN_TIME_RANGE_VALUE = 15;
const MAX_TIME_RANGE_VALUE = 60;

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
    const colorPoint = colorStart + i * intervalSize;
    colorArray.push(colorScale(colorPoint));
  }

  return colorArray;
}

const buildSeriesFromEpoch = (series, epoch) => {
  epoch.data.forEach(p => {
    if (!series[p.dest]) {
      series[p.dest] = [];
    }
    series[p.dest].push(p);
  });

  return series;
};

const epochsToSeries = epochs => {
  let series = epochs.reduce((s, current) => {
    return buildSeriesFromEpoch(s, current);
  }, {});
  return series;
};

const TimeIntervalConfig = props => {
  let PauseResumeButton = null;
  if (props.slidingWindowStopped) {
    PauseResumeButton = () => (
      <Button outlined onClick={() => props.startSlidingWindow()}>
       RESUME
      </Button>
    );
  } else {
    PauseResumeButton = () => (
      <Button outlined onClick={() => props.stopSlidingWindow()}>
        PAUSE
      </Button>
    );
  }

  return (
    <div style={{ marginTop: 10 }}>
      <Grid>
        <GridCell span={3}>
          <div style={styles.slider.container}>
            <div style={styles.slider.title}>
              <span className="mdc-typography--headline8">
                Timeframe in secs{' '}
              </span>
            </div>
            <Grid>
              <GridCell span={1} style={styles.slider.boundaryValues}>
                <span>{MIN_TIME_RANGE_VALUE}</span>
              </GridCell>
              <GridCell span={10}>
                <Slider
                  value={props.value}
                  onChange={e =>
                    props.changeTimeIntervalDuration(e.detail.value)
                  }
                  discrete
                  min={MIN_TIME_RANGE_VALUE}
                  max={MAX_TIME_RANGE_VALUE}
                  step={1}
                />
              </GridCell>
              <GridCell span={1} style={styles.slider.boundaryValues}>
                <span>{MAX_TIME_RANGE_VALUE}</span>
              </GridCell>
            </Grid>
          </div>
        </GridCell>
        <GridCell span={1}></GridCell>
        <GridCell span={8}>
          <div style={{ margin: 10 }}>
            <PauseResumeButton />
          </div>
        </GridCell>
      </Grid>
    </div>
  );
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
            checked={props.displayEdgeLabels}
            onChange={() =>
              props.changeDisplayEdgeLabels(!props.displayEdgeLabels)
            }
          >
            Show Link Data
          </Checkbox>
        </GridCell>
        <GridCell span={2}>
          <Checkbox
            checked={props.showApps}
            onChange={e => props.changeShowApps(e.target.checked)}
          >
            Show Apps
          </Checkbox>
        </GridCell>
        <GridCell span={12}>
          <TimeIntervalConfig
            changeTimeIntervalDuration={value => {
              props.changeTimeIntervalDuration(value);
            }}
            stopSlidingWindow={props.stopSlidingWindow}
            startSlidingWindow={props.startSlidingWindow}
            slidingWindowStopped={props.slidingWindowStopped}
          />
        </GridCell>
      </Grid>
    </>
  );
};

const ViewForName = ({
  scenarioName,
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
  selectedDest,
  colorForApp,
  changeSourceNodeSelected,
  viewName,
  displayEdgeLabels
}) => {
  const appIds = apps.map(app => app.data.id);

  const dashboard = 'http://' + location.hostname + ':30009/d/100/metrics-dashboard?orgId=1';
  const datasource = '&var-datasource=meep-influxdb';
  const database = '&var-database=' + scenarioName;
  const refreshInterval = '&refresh=1s';
  const srcApp = '&var-src=' + selectedSource;
  const destApp = '&var-dest=' + selectedDest;
  const viewMode = '&kiosk';
  const theme = '&theme=light';
  const dashboardUrl = dashboard + datasource + database + refreshInterval + srcApp + destApp + viewMode + theme;
  
  switch (viewName) {
  case APPS_VIEW:
    return (
      <ResizeableContainer key={keyForSvg}>
        {(width, height) => (
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
            onNodeClicked={e => {
              changeSourceNodeSelected(e.node);
            }}
            displayEdgeLabels={displayEdgeLabels}
          />
        )}
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
            width={width}
            height={height}
            destinations={appIds}
            colorRange={colorRange}
            selectedSource={selectedSource}
            dataType={dataType}
            min={min}
            max={max}
            colorForApp={colorForApp}
          />
        )}
      </ResizeableContainer>
    );
  case THROUGHPUT_VIEW:
    return (
      <ResizeableContainer key={keyForSvg}>
        {(width, height) => (
          <IDCLineChart
            keyForSvg={keyForSvg}
            data={dataPoints}
            series={series}
            startTime={startTime}
            mobilityEvents={mobilityEvents}
            width={width}
            height={height}
            destinations={appIds}
            colorRange={colorRange}
            selectedSource={selectedSource}
            dataType={dataType}
            min={min}
            max={max}
            colorForApp={colorForApp}
          />
        )}
      </ResizeableContainer>
    );
  case METRICS_VIEW:
    return (
      <div style={{ height: '70vh' }}>
        <Iframe
          url={dashboardUrl}
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
      <IDCVis
        type={TYPE_EXEC}
        width='100%'
        height='100%'
        onEditElement={() => { }}
      />
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
      changeDisplayEdgeLabels={props.changeDisplayEdgeLabels}
      displayEdgeLabels={props.displayEdgeLabels}
      changeShowApps={props.changeShowApps}
      showApps={props.showApps}
      changeTimeIntervalDuration={props.changeTimeIntervalDuration}
      stopSlidingWindow={props.stopSlidingWindow}
      startSlidingWindow={props.startSlidingWindow}
      slidingWindowStopped={props.slidingWindowStopped}
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
  count: { color: blue },
  eventName: { color: '#6e6e6e' },
  arrow: { color: '#6e6e6e' },
  element: { color: blue }
};

// let eventCount=0;
const EventLog = props => {
  // TODO: generalize function for other types of events.
  // Now it creates a description for Mobility Events
  const descriptionFromEvent = event => {
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
      <span className="mdc-typography--headline8" style={{ marginLeft: 10 }}>
        Events
      </span>
      <div style={eventLogStyle}>{props.events.map(descriptionFromEvent)}</div>
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

  componentDidMount() { }

  componentWillUnmount() {
    clearInterval(this.dataTimer);
  }

  getRoot() {
    return d3.hierarchy(this.props.displayedScenario, getScenarioNodeChildren);
  }

  changeDisplayEdgeLabels(val) {
    this.setState({ displayEdgeLabels: val });
  }

  changeShowApps(checked) {
    this.props.onShowAppsChanged(checked);
  }

  changeMetricsTimeIntervalDuration(duration) {
    this.props.changeMetricsTimeIntervalDuration(duration);
  }

  stopSlidingWindow() {
    this.setState({ slidingWindowStopped: true });
  }

  startSlidingWindow() {
    this.setState({ slidingWindowStopped: false });
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
    const appMap = apps.reduce((acc, app) => {
      acc[app.data.id] = app;
      return acc;
    }, {});
    const colorRange = colorArray(appIds.length);

    const selectedSource = this.props.sourceNodeSelected
      ? this.props.sourceNodeSelected.data.id
      : null;

    const selectedDest = this.props.destNodeSelected
      ? this.props.destNodeSelected.data.id
      : null;

    const colorForApp = apps.reduce((res, val, i) => {
      return { ...res, [val.data.id]: colorRange[i] };
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
    const firstEpoch = epochs.length
      ? epochs[0]
      : {
        data: [],
        startTime: null
      };
    let lastEpoch = epochs.length
      ? epochs.slice(-1)[0]
      : {
        data: [],
        startTime: null
      };

    // Determine startTime of first epoch and endTime of last epoch
    const startTime = firstEpoch.data.length ? firstEpoch.startTime : null;
    const endTime = lastEpoch.data.length
      ? moment(lastEpoch.startTime)
        .add(1, 'seconds')
        .format(TIME_FORMAT)
      : null;
    const series = epochsToSeries(epochs, selectedSource);

    const withTypeAndSource = type => source => point => {
      return point.dataType === type && point.src === source;
    };

    // For view 1
    const view1Name = this.props.view1Name;
    const view1DataType = dataTypeForView(view1Name);
    const series1 = filterSeries(appIds)(
      withTypeAndSource(view1DataType)(selectedSource)
    )(series);
    const lastEpochData1 = lastEpoch.data.filter(isDataOfType(view1DataType));

    // For view2
    const view2Name = this.props.view2Name;
    const view2DataType = dataTypeForView(view2Name);
    const series2 = filterSeries(appIds)(
      withTypeAndSource(view2DataType)(selectedSource)
    )(series);
    const lastEpochData2 = lastEpoch.data.filter(isDataOfType(view2DataType));

    // Mobility events
    const extractPointsOfType = type => epoch =>
      epoch.data.filter(isDataPointOfType(type));
    const extractMobilityEvents = extractPointsOfType(ME_MOBILITY_EVENT);
    const mobilityEvents = epochs.flatMap(extractMobilityEvents);

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
        selectedDest={selectedDest}
        colorForApp={colorForApp}
        changeSourceNodeSelected={node =>
          this.props.changeSourceNodeSelected(node)
        }
        viewName={view1Name}
        displayEdgeLabels={this.state.displayEdgeLabels}
      />
    );

    const view2 = (
      <ViewForName
        scenarioName={this.props.scenarioName}
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
        selectedDest={selectedDest}
        colorForApp={colorForApp}
        changeSourceNodeSelected={node =>
          this.props.changeSourceNodeSelected(node)
        }
        viewName={view2Name}
        displayEdgeLabels={this.state.displayEdgeLabels}
      />
    );

    const EventLogComponent = () => <EventLog events={mobilityEvents} />;

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
          changeTimeIntervalDuration={duration => {
            this.changeMetricsTimeIntervalDuration(duration);
          }}
          stopSlidingWindow={() => this.stopSlidingWindow()}
          startSlidingWindow={() => this.startSlidingWindow()}
          slidingWindowStopped={this.state.slidingWindowStopped}
          dashboardViewsList={DASHBOARD_VIEWS_LIST}
          changeView1={viewName => this.props.changeView1(viewName)}
          changeView2={viewName => this.props.changeView2(viewName)}
          displayEdgeLabels={this.state.displayEdgeLabels}
          changeDisplayEdgeLabels={display =>
            this.changeDisplayEdgeLabels(display)
          }
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

                {(view2Name === LATENCY_VIEW) || (view2Name === THROUGHPUT_VIEW) ?  <EventLogComponent /> : null}
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
    destNodeSelected: state.exec.metrics.destNodeSelected,
    dataTypeSelected: state.exec.metrics.dataTypeSelected,
    eventCreationMode: state.exec.eventCreationMode,
    metricsTimeIntervalDuration: state.exec.metrics.timeIntervalDuration,
    scenarioState: state.exec.state.scenario,
    view1Name: state.ui.dashboardView1,
    view2Name: state.ui.dashboardView2
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeSelectedDestination: dest =>
      dispatch(execFakeChangeSelectedDestination(dest)),
    changeSourceNodeSelected: src =>
      dispatch(execChangeSourceNodeSelected(src)),
    changeDestNodeSelected: dest =>
      dispatch(execChangeDestNodeSelected(dest)),
    changeMetricsTimeIntervalDuration: duration =>
      dispatch(execChangeMetricsTimeIntervalDuration(duration)),
    clearMetricsEpochs: () => dispatch(execClearMetricsEpochs()),
    changeView1: name => dispatch(uiExecChangeDashboardView1(name)),
    changeView2: name => dispatch(uiExecChangeDashboardView2(name))
  };
};

const ConnectedDashboardContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(DashboardContainer);

export default ConnectedDashboardContainer;
