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

import _ from 'lodash';
import * as d3 from 'd3';
import React from 'react';
import {Axis, axisPropsFromTickScale, LEFT, BOTTOM} from 'react-d3-axis';
import { ME_LATENCY_METRICS, ME_THROUGHPUT_METRICS } from '../meep-constants';
import { blue } from './graph-utils';
// const Axis = props => {
//   const axisRef = axis => {
//     axis && props.axisCreator(select(axis));
//   };

//   return <g className={props.className} ref={axisRef} />;
// };

const notNull = x => x;
const IDCLineChart = (props) => {
  const keyForSvg=props.keyForSvg;
  let width = props.width;
  let yClipping = 45;

  const margin = {top: 20, right: 40, bottom: 30, left: 60};
  // const width = props.width; // - margin.left - margin.right;
  const height = props.height; // - margin.top - margin.bottom;

  const maxForKey = series => key => d3.max(series[key], p => p.value);
  const maxes = Object.keys(props.series).map(maxForKey(props.series));
  const max = d3.max(maxes);
  const maxOfYScale = Math.ceil(max/50.0) * 50.0;
  const yRange = [0, maxOfYScale];

  const destinations = props.selectedSource ? props.destinations.slice(-props.destinations.length) : [];
  const colorRange = destinations.map(s => props.colorForApp[s]);

  const flattenSeries = series => {
    return _.flatMap(Object.values(series));
  };
  const timeRange = d3.extent(flattenSeries(props.series), d => new Date(d.timestamp));
  const x = d3.scaleTime().domain(timeRange).range([0, width]);
  const y = d3.scaleLinear().domain(yRange).range([height - yClipping, 0]);
  const z = d3.scaleOrdinal().range(colorRange);

  // Compute data lines
  const dataLineFromSeries = series => key => {
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

  // Chart title
  const chartTitleForType = type => {
    switch (type) {
    case ME_LATENCY_METRICS:
      return 'Latency Chart';
    case ME_THROUGHPUT_METRICS:
      return 'Throughput Chart';
    default:
      return '';
    }
  };
  
  const chartTitle = chartTitleForType(props.dataType);

  const axisWidthOffset = 12;
  const meX = d => x(new Date(d.timestamp)) + axisWidthOffset;

  const mobilityEventLine = me => `M${meX(me) + margin.left},${y(yRange[1]) + margin.top} L${meX(me) + margin.left},${y(yRange[0]) + margin.top}`;
  const mobilityEventText = me => `ME from ${me.src} to ${me.dest}`;

  const mobilityEventLines = props.mobilityEvents.map(me => {
    return (
      <path
        className='mobilityEventLine'
        d={mobilityEventLine(me)}
        id={me.timestamp}
        key={me.timestamp}
        style={{stroke: blue, strokeWidth: 2, fill: 'none', textAnchor: 'middle'}}
      />
    );
  });

  const mobilityEventTextPathDefs = 
  <defs>
    {
      props.mobilityEvents.map((me, i) => {
        return <path
          key={'mobilityEventLinePathDef' + i}
          id={'mobilityEventLinePathDef' + i}
          d={mobilityEventLine(me)}
          className='mobilityEventLinePathDef'
        />;
      })
    }
  </defs>;

  const mobilityEventTexts = props.mobilityEvents.map((me, i) => {
    return(
      <text
        key={'mobilityEventLinePath' + i}
        style={{stroke: 'gray', strokeWidth: 1, fill: 'none'}}
      >
        <textPath
          xlinkHref={'mobilityEventLinePathDef' + i}
          startOffset={'45%'}
        >
          {mobilityEventText(me)}
        </textPath>
      </text>
    );
  
  });

  // text label for the y axis
  const labelForType = type => {
    switch (type) {
    case ME_LATENCY_METRICS:
      return 'Latency (ms)';
    case ME_THROUGHPUT_METRICS:
      return 'Throughput (kbs)';
    default:
      return '';
    }
  };
  
  const yAxisLabel = labelForType(props.dataType);

  return (
    <svg
      key={keyForSvg}
      height={height}
      width={'100%'}
    >
      <>
      <g
        transform={`translate(${margin.left}, ${margin.top})`}
      >
        <Axis 
          {...axisPropsFromTickScale(y, 10)} 
          style={{
            orient: LEFT, tickSizeInner: -width,
            strokeColor: '#BBBBBB'
          }}/>
      </g>

      <g
        transform={`translate(${margin.left}, ${height - margin.top - 5})`}
      >
        <Axis
          {...axisPropsFromTickScale(x, 10)}
          style={{
            orient: BOTTOM,
            tickSizeInner: -(height - yClipping),
            strokeColor: '#BBBBBB'
          }}/>

        {/* <Axis
          {...axisPropsFromTickScale(x, 10)}
          style={{
            orient: TOP,
            tickSizeInner: -(height - 2*yClipping),
            strokeColor: '#BBBBBB'
          }}/> */}
      </g>

      <text
        className='chartTitle'
        y={0 + margin.top - 8 }
        x={width / 2}
        style={{textAnchor: 'middle'}}
      >
        {chartTitle}
      </text>

      <text
        className='yLabel'
        transform='rotate(-90)'
        y={0}
        x={0 - (height / 2)}
        dy='1em'
        style={{textAnchor: 'middle'}}
      >
        {yAxisLabel}
      </text>
      
      {lines}
      {mobilityEventLines}
      {mobilityEventTexts}
      {mobilityEventTextPathDefs}
      </>
    </svg>
  );
  
};

export default IDCLineChart;