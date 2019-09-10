/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

import _ from 'lodash';
import { connect } from 'react-redux';
import React, { useRef, useEffect }  from 'react';
import ReactDOM from 'react-dom';
import moment from 'moment';
import * as d3 from 'd3';
import {Axis, axisPropsFromTickScale, LEFT} from 'react-d3-axis';
import uuid from 'uuid';
import { uiChangeCurrentDialog } from '../state/ui';
import { execFakeChangeSelectedDestination } from '../state/exec';
import { LATENCY_METRICS, THROUGHPUT_METRICS } from '../meep-constants';

// const Axis = props => {
//   const axisRef = axis => {
//     axis && props.axisCreator(select(axis));
//   };

//   return <g className={props.className} ref={axisRef} />;
// };

const notNull = x => x;
const IDCLineChart = props => {

  const margin = {top: 20, right: 40, bottom: 30, left: 60};
  const width = props.width - margin.left - margin.right;
  const height = props.height - margin.top - margin.bottom;

  const min = props.min;
  const max = props.max;

  const maxOfYScale = Math.ceil(max/100.0) * 100.0;

  const destinations = props.selectedSource ? props.destinations.slice(-props.destinations.length) : [];
  const colorRange = destinations.map(s => props.colorForApp[s]);

  const yRange = [0, 200];

  const flattenSeries = series => {
    return _.flatMap(Object.values(series));
  };
  const timeRange = d3.extent(flattenSeries(props.series), d => new Date(d.timestamp));
  const x = d3.scaleTime().domain(timeRange).range([0, width]);
  const y = d3.scaleLinear().domain(yRange).range([height - 50, 0]);
  const z = d3.scaleOrdinal().range(colorRange);

  // Axes
  const xAxis = d3.axisBottom(x); //.ticks(d3.timeSeconds);
  const yAxis = d3.axisLeft(y).scale(y).tickSize(0.01);

  const dataLineFromSeries = series => key => {
    if (!series[key]) {
      console.log('Check');
    }
    let line;
    
    if (series[key]) {
      line = series[key].filter(notNull).filter(p => p.value)
        .sort((a, b) => {
          return x(new Date(a.timestamp)) - x(new Date(b.timestamp));
        });
    } else {
      line = [];
    }
    
    //TODO: add point at props.startTime and props.endTime
    
    line.key = key;
    return line;
  };

  let dataLines = destinations.map(dataLineFromSeries(props.series));

  // dataLines = dataLines.length ? [dataLines[0]] : [];

  if (dataLines.length < 5) {
    console.log('Too few dataLines: ', dataLines.length);
  }

  const valueLine = d3.line()
    .x(function(d) {
      return margin.left + x(new Date(d.timestamp));
    })
    .y(function(d) {
      return y(d.value) + margin.top;
    })
    .curve(d3.curveMonotoneX);

  const lines = dataLines.map((dl, i) => {
    return (
      <path
        className='line'
        key={`linechart${i}`}
        d={valueLine(dl)}
        style={{fill: 'none', 'strokeWidth': 3, 'stroke': z(i)}}
      />
    );
  });

  return (
    <svg
      height={width}
      width={height}
    >
      <>
      <g
        transform={`translate(${margin.left}, ${margin.top})`}
      >
        <g className='axisBottom'>

        </g>
        {lines}
      </g>
        
      </>
    </svg>
  );
  
};

export default IDCLineChart;