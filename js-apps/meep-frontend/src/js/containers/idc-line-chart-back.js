/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

import _ from 'lodash';
import React, { useRef, useEffect }  from 'react';
import * as d3 from 'd3';
import { LATENCY_METRICS, THROUGHPUT_METRICS } from '../meep-constants';

const notNull = x => x;
const IDCLineChartBack = props => {
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

      const flattenSeries = series => {
        return _.flatMap(Object.values(series));
      };
      
      const chart = (series) => {
        const destinations = props.selectedSource ? props.destinations.slice(-props.destinations.length) : [];
        const colorRange = destinations.map(s => props.colorForApp[s]);

        const yRange = [0, 200];
        const timeRange = d3.extent(flattenSeries(series), d => new Date(d.timestamp));
        const x = d3.scaleTime().domain(timeRange).range([0, width]);
        const y = d3.scaleLinear().domain(yRange).range([height - 50, 0]);
        const z = d3.scaleOrdinal().range(colorRange);
      
        // Axes
        const xAxis = d3.axisBottom(x); //.ticks(d3.timeSeconds);
        const yAxis = d3.axisLeft(y).scale(y)
          .tickSize(0.01);
        // const yAxisr = d3.axisLeft(y);

        // const dataLinePointFromDataPoint = key => point => {
        //   if (point[key] === undefined) {
        //     console.log('point[key] is undefined, for key ' + key + ' and point ', point);
        //     return null;
        //   }
        //   return {
        //     date: point.date,
        //     value: point[key]
        //   };
        // };
        const dataLineFromSeries = series => key => {
          
          const line = series[key].filter(notNull).filter(p => p.value);
          // .sort((a, b) => {
          //   return x(new Date(a.timestamp)) - x(new Date(b.timestamp));
          // });

          //TODO: add point at props.startTime and props.endTime
          
          line.key = key;
          return line;
        };

        let dataLines = destinations.map(dataLineFromSeries(props.series));

        // TODO: remove
        dataLines = dataLines.length ? [dataLines[0]] : [];

        const valueLine = d3.line()
          .x(function(d) {
            return margin.left + x(new Date(d.timestamp));
          })
          .y(function(d) {
            return y(d.value) + margin.top;
          })
          .curve(d3.curveMonotoneX);

        mainGroup.selectAll('.line')
          .data(dataLines)
          .join('path').attr('class', 'line')
          .attr('d', valueLine)
          .style('stroke', (d, i) => z(i))
          .style('fill', 'none')
          .style('stroke-width', 3);

        // if (mainGroup.selectAll('.line').size() > 0) {
        //   const linesMarginLeft = mainGroup.selectAll('.line').attr('margin-left');
        //   console.log('linesMarginLeft: ', linesMarginLeft);
        // }

        // Mobility events
        // const mobilityEventLine = d => `M${x(new Date(d.timestamp)) + margin.left},${y(yRange[0]) + margin.top} L${x(new Date(d.timestamp)) + margin.left},${y(yRange[1]) + margin.top}`;
        const mobilityEventLine = d => `M${x(new Date(d.timestamp)) + margin.left},${y(yRange[1]) + margin.top} L${x(new Date(d.timestamp)) + margin.left},${y(yRange[0]) + margin.top}`;
        mainGroup.selectAll('.mobilityEventLine')
          .data(props.mobilityEvents)
          .join('path')
          .attr('class', 'mobilityEventLine')
          .attr('d', mobilityEventLine)
          .attr('id', d => d.timestamp)
          .style('stroke', 'gray')
          .style('stroke-width', 1)
          .style('fill', 'none');
          
        mainGroup.selectAll('.mobilityEventLineText')
          .data(props.mobilityEvents)
          .join('text')
          .attr('class', 'mobilityEventLineText')
          .style('stroke','gray')
          .style('stroke-width', 1)
          .style('fill','gray');
        // .attr('x', d => x(new Date(d.timestamp)) + margin.left)
        // .attr('dy',50 + margin.top)
            
      
        mainGroup.selectAll('.mobilityEventLineTextPath').remove();
        mainGroup.selectAll('.mobilityEventLineText')
          .data(props.mobilityEvents)
          .append('textPath')
          .attr('class', 'mobilityEventLineTextPath')
          .attr('xlink:href', d => `#${d.timestamp}`)
          .attr('stroke','gray')
          .attr('fill','gray')
          .text(d => `Mobility Event:  ${d.src} to ${d.dest}`)
          .attr('transform', 'rotate(-180)');
          

        
        const xAxisGroup = mainGroup.selectAll('.xaxis');
        if (xAxisGroup.size() === 0) {
          mainGroup.append('g')
            .attr('class', 'xaxis')
            .attr('transform', 'translate(0,' + height + ')').call(xAxis);
        } else {
          xAxisGroup.attr('transform', 'translate(0,' + height + ')').call(xAxis);
        }

        mainGroup.selectAll('.xaxis').call(xAxis);
         
        const yAxisGroup = mainGroup.selectAll('.yaxis');
        if (yAxisGroup.size() === 0) {
          mainGroup.append('g')
            .attr('class', 'yaxis')
            .attr('transform', 'translate(' + width + ', 0)')
            .style('z-index', '18')
            .call(yAxis);
        } else {
          yAxisGroup.attr('transform', 'translate(' + width + ', 0)');
        }

        // text label for the y axis
        const labelForType = type => {
          switch (type) {
          case LATENCY_METRICS:
            return 'Latency (ms)';
          case THROUGHPUT_METRICS:
            return 'Throughput (kbs)';
          default:
            return '';
          }
        };

        const yAxisLabel = labelForType(props.dataType);
        if (!mainGroup.selectAll('.yLabel').size()) {
          mainGroup.append('text')
            .attr('class', 'yLabel')
            .attr('transform', 'rotate(-90)')
            .attr('y', 0 - margin.left + 10)
            .attr('x', 0 - (height / 2))
            .attr('dy', '1em')
            .style('text-anchor', 'middle')
            .text(yAxisLabel);
        } else {
          mainGroup.selectAll('.yLabel')
            .text(yAxisLabel);
        }


       
        const yAxisGroup0 = mainGroup.selectAll('.yaxis0');
        if (yAxisGroup0.size() === 0) {
          mainGroup.append('g')
            .attr('class', 'yaxis0')
            .attr('transform', 'translate(0, 0)')
            .style('z-index', '18')
            .call(yAxis);
        } else {
          yAxisGroup0.attr('transform', 'translate(0, 0)').style('z-index', '18');
        }
      };
        
      chart(props.series);
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

export default IDCLineChartBack;