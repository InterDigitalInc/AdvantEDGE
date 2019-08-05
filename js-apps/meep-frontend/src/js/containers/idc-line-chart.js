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
import uuid from 'uuid';
import { uiChangeCurrentDialog } from '../state/ui';
import { execFakeChangeSelectedDestination } from '../state/exec';

const notNull = x => x;
const IDCLineChart = props => {
  const d3Container = useRef(null);

  /* The useEffect Hook is for running side effects outside of React,
       for instance inserting elements into the DOM using D3 */
  useEffect(
    () => {
      
      const margin = {top: 20, right: 40, bottom: 30, left: 60};
      const width = props.width - margin.left - margin.right;
      const height = props.height - margin.top - margin.bottom;

      let mainGroup = d3.select(d3Container.current);
      if (mainGroup.select('g').size() === 0) {
        mainGroup = mainGroup.append('g')
          .attr('width', props.width + margin.left + margin.right)
          .attr('height', props.height + margin.top + margin.bottom)
          .attr('transform', `translate(${margin.left}, ${margin.top})`);
      }
      
      const chart = (data) => {
        const destinations = props.sourceSelected ? props.destinations.slice(-props.destinations.length) : [];
        const colorRange = destinations.map(s => props.colorForApp[s]);

        const yRange = [0, 200];
        const timeRange = d3.extent(data, d => new Date(d.date));
        const x = d3.scaleTime().domain(timeRange).range([0, width]);
        const y = d3.scaleLinear().domain(yRange).range([height - 50, 0]);
        const z = d3.scaleOrdinal().range(colorRange);
      
        // Axes
        const xAxis = d3.axisBottom(x); //.ticks(d3.timeSeconds);
        const yAxis = d3.axisLeft(y).scale(y)
          .tickSize(0.01);
        // const yAxisr = d3.axisLeft(y);

        const dataLinePointFromDataPoint = key => point => {
          if (!point[key]) {
            return null;
          }
          return {
            date: point.date,
            value: point[key]
          };
        };
        const dataLineFromData = key => {
          const line = data.map(dataLinePointFromDataPoint(key)).filter(notNull);
          line.key = key;
          return line;
        };

        
        const dataLines = destinations.map(dataLineFromData);
        const valueLine = d3.line()
          .x(function(d) { return margin.left + x(new Date(d.date)); })
          .y(function(d) { return y(d.value) + margin.top;  })
          .curve(d3.curveMonotoneX);
        mainGroup.selectAll('.line')
          .data(dataLines)
          .join('path').attr('class', 'line')
          .attr('d', valueLine)
          .style('stroke', (d, i) => z(i))
          .style('fill', 'none')
          .style('stroke-width', 3)
          .style('margin-left', 35)
          .style('margin-top', 20);
              
        const xAxisGroup = mainGroup.selectAll('.xaxis');
        if (xAxisGroup.size() === 0) {
          mainGroup.append('g')
            .attr('class', 'xaxis')
            .attr('transform', 'translate(0,' + height + ')').call(xAxis);
        }

        mainGroup.selectAll('.xaxis').call(xAxis);
         
        const yAxisGroup = mainGroup.selectAll('.yaxis');
        if (yAxisGroup.size() === 0) {
          mainGroup.append('g')
            .attr('class', 'yaxis')
            .attr('transform', 'translate(' + width + ', 0)')
            .style('z-index', '18')
            .call(yAxis);
        }

        // text label for the y axis
        if (!mainGroup.selectAll('.yLabel').size()) {
          mainGroup.append('text')
            .attr('transform', 'rotate(-90)')
            .attr('y', 0 - margin.left + 10)
            .attr('x', 0 - (height / 2))
            .attr('dy', '1em')
            .attr('class', 'yLabel')
            .style('text-anchor', 'middle')
            .text('Latency (ms)');
        }
       

        const yAxisGroup0 = mainGroup.selectAll('.yaxis0');
        if (yAxisGroup0.size() === 0) {
          mainGroup.append('g')
            .attr('class', 'yaxis0')
            .attr('transform', 'translate(0 , 0)')
            .style('z-index', '18')
            .call(yAxis);
        }
      };
        
      chart(props.data);
    },

    /*
            useEffect has a dependency array (below). It's a list of dependency
            variables for this useEffect block. The block will run after mount
            and whenever any of these variables change. We still have to check
            if the variables are valid, but we do not have to compare old props
            to next props to decide whether to rerender.
        */
    [props.data, d3Container.current]);

  return (
    <div className='chart'>
      <svg
        //viewBox='0 -20 200 33'
        ref={d3Container}
        className='d3-component'
        height={props.height}
        width={props.width}
      >
      
      </svg>

    </div>
  );  
};

export default IDCLineChart;