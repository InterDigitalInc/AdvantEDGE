import _ from 'lodash';
import { connect } from 'react-redux';
import React, { Component }  from 'react';
import { Grid, GridCell, GridInner } from '@rmwc/grid';
import { Elevation } from '@rmwc/elevation';
import { Graph } from 'react-d3-graph';
import ReactDOM from 'react-dom';
import { Button } from '@rmwc/button';
import * as d3 from 'd3';
import axios from 'axios';

// import IDCAreaChart from './idc-area-chart';
import IDCLineChart from './idc-line-chart';
import IDCGraph from './idc-graph';
import IDCAppsView from './idc-apps-view';

import {
  getScenarioNodeChildren,
  isApp
} from '../util/scenario-utils';

import {
  execFakeChangeSelectedDestination,
  execChangeSourceNodeSelected,
  execAddMetricsEpoch
} from '../state/exec';

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
  const dataPoints = selectedEpochs.map(dataPointFromEpochDataPoints(destinations)(sourceNodeId)(dataAccessor)).filter(notNull);
  return dataPoints;
};

const dataAccessorForType = dataType => {
  switch (dataType) {
  case 'latency':
    return p => p.data.latency;
  case 'ingressPacketStats':
    return p => p.data.throughput;
  default:
    return dataAccessorForType('latency');
  }
};

class DashboardContainer extends Component {
  constructor(props) {
    super(props);

    this.state = {
      
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

  render() {
    const root = this.getRoot();
    const nodes = root.descendants();
   
    const apps = nodes.filter(isApp);
    const destinations = apps.map(a => a.data.id);
    const colorRange = colorArray(destinations.length);
    const nbEpochs = 25;

    const selectedNodeId = this.props.sourceNodeSelected ? this.props.sourceNodeSelected.data.id : null;
    const dataAccessor = dataAccessorForType(this.props.dataTypeSelected);
    const dataPoints = epochsToDataPoints(this.props.epochs)(nbEpochs)(destinations)(dataAccessor)(selectedNodeId);

    const showApps = this.props.showAppsView;
    const span = showApps ? 6 : 12;

    const colorForApp = apps.reduce((res, val, i) => {
      // res[val.data.id] = colorRange[i];
      return {...res, [val.data.id]: colorRange[i]};
    }, {});

    const lastEpoch = this.props.epochs.length ? this.props.epochs.slice(-1)[0] : [];
    const isDataOfType = type => dataPoint => dataPoint.dataType === type;
    const data = lastEpoch.filter(isDataOfType(this.props.dataTypeSelected));

    let graph = null;

    if (showApps) {
      graph = (
        <IDCAppsView
          apps={apps}
          colorRange={colorRange}
          width={700}
          height={600}
          data={data}
          dataAccessor={dataAccessor}
          dataType={this.props.dataTypeSelected}
          selectedSource={selectedNodeId}
          colorForApp={colorForApp}
          onNodeClicked={(e) => {
            console.log('Node clicked is: ', e.node);
            this.props.changeSourceNodeSelected(e.node);
          }}
        />
      );
    } else {
      graph = (<IDCGraph 
        width={1000}
        height={600}
      />);
    }

    return (
      <Grid>
        <GridCell span={span} style={{marginLeft: -10}}>
          <Elevation z={4}>
            {graph}
          </Elevation>
        </GridCell>

        {showApps ? (<GridCell span={6}  style={{marginRight: -10}}>
          <Elevation z={4}>
            <IDCLineChart
              data={dataPoints}
              width={700} height={600}
              destinations={destinations}
              colorRange={colorRange}
              sourceSelected={this.props.sourceNodeSelected}
              // min={min}
              // max={max}
              colorForApp={colorForApp}
            />
          </Elevation>
          
        </GridCell>) 
          : null}
      </Grid>
      
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