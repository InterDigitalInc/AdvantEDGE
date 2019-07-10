import _ from 'lodash';
import { connect } from 'react-redux';
import React, { Component }  from 'react';
import { Grid, GridCell, GridInner } from '@rmwc/grid';
import { Graph } from 'react-d3-graph';
import ReactDOM from 'react-dom';
import { Button } from '@rmwc/button';
import * as d3 from 'd3';

import IDCAreaChart from './idc-area-chart';
import IDCGraph from './idc-graph';

import {
  getScenarioNodeChildren,
  isApp
} from '../util/scenario-utils';

import {
  execFakeChangeSelectedDestination
} from '../state/exec';

const newDataPoint = (date) => {
  const newDate = date || new Date();
  const secs = newDate.getSeconds();
  const newDateString = newDate.toString();
  return {
    'date':newDate,
    'AR':Math.abs(Math.random()*Math.sin(0.05*secs)),
    'DJ':Math.abs(Math.random()*Math.cos(0.05*secs)),
    'MS':Math.abs(Math.random()*Math.sin(0.5*secs)),
    'RC':Math.abs(Math.random()*Math.cos(0.1*secs)),
    'CG':Math.abs(Math.random()*Math.sin(0.2*secs)),
    'RI':Math.abs(Math.random()*Math.sin(3.0*secs))
  };
};

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

// const dataStr = '[{"date":"01/08/13","AR":0.1,"DJ":0.35,"MS":0.21,"RC":0.1,"CG":0.1,"RI":0.1},{"date":"01/09/13","AR":0.15,"DJ":0.36,"MS":0.25,"RC":0.15,"CG":0.15,"RI":0.15},{"date":"01/10/13","AR":0.35,"DJ":0.37,"MS":0.27,"RC":0.35,"CG":0.35,"RI":0.35},{"date":"01/11/13","AR":0.38,"DJ":0.22,"MS":0.23,"RC":0.38,"CG":0.38,"RI":0.38},{"date":"01/12/13","AR":0.22,"DJ":0.24,"MS":0.24,"RC":0.22,"CG":0.22,"RI":0.22},{"date":"01/13/13","AR":0.16,"DJ":0.26,"MS":0.21,"RC":0.16,"CG":0.16,"RI":0.16},{"date":"01/14/13","AR":0.07,"DJ":0.34,"MS":0.35,"RC":0.07,"CG":0.07,"RI":0.07},{"date":"01/15/13","AR":0.02,"DJ":0.21,"MS":0.39,"RC":0.02,"CG":0.02,"RI":0.02},{"date":"01/16/13","AR":0.17,"DJ":0.18,"MS":0.4,"RC":0.17,"CG":0.17,"RI":0.17},{"date":"01/17/13","AR":0.33,"DJ":0.45,"MS":0.36,"RC":0.33,"CG":0.33,"RI":0.33},{"date":"01/18/13","AR":0.4,"DJ":0.32,"MS":0.33,"RC":0.4,"CG":0.4,"RI":0.4},{"date":"01/19/13","AR":0.32,"DJ":0.35,"MS":0.43,"RC":0.32,"CG":0.32,"RI":0.32},{"date":"01/20/13","AR":0.26,"DJ":0.3,"MS":0.4,"RC":0.26,"CG":0.26,"RI":0.26},{"date":"01/21/13","AR":0.35,"DJ":0.28,"MS":0.34,"RC":0.35,"CG":0.35,"RI":0.35},{"date":"01/22/13","AR":0.4,"DJ":0.27,"MS":0.28,"RC":0.4,"CG":0.4,"RI":0.4},{"date":"01/23/13","AR":0.32,"DJ":0.26,"MS":0.26,"RC":0.32,"CG":0.32,"RI":0.32},{"date":"01/24/13","AR":0.26,"DJ":0.15,"MS":0.37,"RC":0.26,"CG":0.26,"RI":0.26},{"date":"01/25/13","AR":0.22,"DJ":0.3,"MS":0.41,"RC":0.22,"CG":0.22,"RI":0.22},{"date":"01/26/13","AR":0.16,"DJ":0.35,"MS":0.46,"RC":0.16,"CG":0.16,"RI":0.16},{"date":"01/27/13","AR":0.22,"DJ":0.42,"MS":0.47,"RC":0.22,"CG":0.22,"RI":0.22},{"date":"01/28/13","AR":0.1,"DJ":0.42,"MS":0.41,"RC":0.1,"CG":0.1,"RI":0.1}]';
// const theData = JSON.parse(dataStr);


// const timeParse = d3.timeParse('%m/%d/%y');
// const startingTime = new Date();
// theData.forEach(function(d, i) {
//     const interval = 1000;
//     d.date = new Date(startingTime.getTime() + i*interval);
// });

// const recentData = end => start => dataPoint => {
//     const dataPointMilli  = dataPoint.date.getTime();
//     return start.getTime() <= dataPointMilli && dataPointMilli <= end.getTime();
// };

const updateData = (data) => {

  let newData;
  if (!data.length) {
    newData = [newDataPoint()];
  } else {
    newData = data.slice(1).concat([newDataPoint()]);
  }

  return newData;
};

const dataPointFromBucket = keys => b => {
  let dp = {
    date: b.date
  };

  const accessor = p => p.delay;
  const avgForKey = pings => acc => key => {
    const pingsForKeyDestination = pings.filter(p => p.dest === key);
    const avg = d3.mean(pingsForKeyDestination, acc);
    return avg;
  };
  
  keys.forEach(k => {
    dp[k] = avgForKey(b.pings)(accessor)(k) || 0;
  });

  return dp;
};

const pingBucketsToData = pingBuckets => nb => keys => {
  const buckets = pingBuckets.slice(-nb);
  const dataPoints = buckets.map(dataPointFromBucket(keys));
  return dataPoints;
};

const maxValue = pingBuckets => {
  const max = d3.max(pingBuckets, b => {
    return d3.max(b.pings, p => p.delay);
  });
  return max;
};

const minValue = pingBuckets => {
  const min = d3.min(pingBuckets, b => {
    return d3.min(b.pings, p => p.delay);
  });
  return min;
};

class DashboardContainer extends Component {
  constructor(props) {
    super(props);

    this.state = {
      
    };
  }

  componentDidMount() {

    // Initial data
    // let theData = [];
    // const initialNbPoints = 25;
    // const now = new Date();
    // for (let i=0; i< initialNbPoints; i++) {
    //   theData.push(newDataPoint(new Date(now.getTime() + (i - initialNbPoints)*1000)));
    // }

    // this.setState({data: theData});
    // const that = this;
    // this.timer = setInterval(() => {
    //   that.setState({
    //     data: updateData(that.state.data)
    //   });
    // }, 1000);
  }

  componentWillUnmount() {
    clearInterval(this.timer);
  }

  calculateDestinations() {
    this.root = this.props.displayedScenario;
    const data = this.root; // || this.props.displayedScenario;
    this.root = d3.hierarchy(data, getScenarioNodeChildren);
    const apps = this.root.descendants().filter(isApp);
    return apps.map(a => a.data.id);
  }

  render() {
   
    const destinations = this.calculateDestinations();

    const nbBuckets = this.props.nbBuckets || 25;
    const max = maxValue(this.props.pingBuckets);
    const min = minValue(this.props.pingBuckets);

    const dataPoints = pingBucketsToData(this.props.pingBuckets)(nbBuckets)(destinations);

    const colorRange = colorArray(destinations.length);
    return (
      <Grid>
        <GridCell span={6}>
          <IDCGraph 
            width={700}
            height={600}
            renderApps={true}
            colorRange={colorRange}
          />
        </GridCell>
        <GridCell span={6}>
          <IDCAreaChart
            data={dataPoints}
            width={700} height={600}
            destinations={destinations}
            colorRange={colorRange}
            onKeySelected={(dest) => this.props.changeSelectedDestination(dest)}
            min={min}
            max={max}
          />    
        </GridCell>
      </Grid>
      
    );
  }
}

const mapStateToProps = state => {
  return {
    pingBuckets: state.exec.fakeData.pingBuckets,
    displayedScenario: state.exec.displayedScenario
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeSelectedDestination: (dest) => dispatch(execFakeChangeSelectedDestination(dest))
  };
};

const ConnectedDashboardContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(DashboardContainer);

export default ConnectedDashboardContainer;