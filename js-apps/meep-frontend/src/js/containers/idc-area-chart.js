/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

import React, { useRef, useEffect }  from 'react';
import * as d3 from 'd3';

const IDCAreaChart = props => {
  const d3Container = useRef(null);
  const onKeySelected = props.onKeySelected;

  /* The useEffect Hook is for running side effects outside of React,
       for instance inserting elements into the DOM using D3 */
  useEffect(
    () => {
      
      const margin = {top: 20, right: 40, bottom: 30, left: 30};
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

        const colorRange = props.colorRange;
        const strokecolor = colorRange[0];

        const yRange = [0, 400];
        const timeRange = d3.extent(data, d => new Date(d.date));
        const x = d3.scaleTime().domain(timeRange).range([0, width]);
        const y = d3.scaleLinear().domain(yRange).range([height - 50, 0]);
        const z = d3.scaleOrdinal().range(colorRange);
      
        // Axes
        const xAxis = d3.axisBottom(x); //.ticks(d3.timeSeconds);
        const yAxis = d3.axisLeft(y).scale(y)
          .tickSize(0.01);
        // const yAxisr = d3.axisLeft(y);

        const keys = props.sources;
        const stack = d3.stack().keys(keys);

        const area = d3.area()
          .x( (d, i) => x(data[i].date))
          .y0(d => y(d[0]))
          .y1(d => y(d[1]))
          .curve(d3.curveCardinal);

        const layers = stack(data);

        mainGroup.selectAll('.layer')
          .data(layers, d => d.key)
          .join('path').attr('class', 'layer')
        // .transition()
        //     .duration(250)
          .attr('d', d => area(d))
          .style('fill', (d, i) => z(i));
              
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

        const yAxisGroup0 = mainGroup.selectAll('.yaxis0');
        if (yAxisGroup0.size() === 0) {
          mainGroup.append('g')
            .attr('class', 'yaxis0')
            .attr('transform', 'translate(0 , 0)')
            .style('z-index', '18')
            .call(yAxis);
        }
              

        // svg.append('g').append('path')
        //   .attr('class', 'vertical')
        //   .attr('width', '1px')
        //   .attr('d', 'M0 0 L0 380')
        //   .attr('transform', 'translate(300, 0)')
        //   .attr('stroke', '000')
        //   .attr('stroke-width', '10px')
        //   .attr('visibility', 'visible');
              
        // svg.append('g')
        //   .attr('class', 'y axis')
        //   .call(yAxis);
      
        mainGroup.selectAll('.layer')
          // .attr('opacity', 1)
          .on('click', function(d, i, nodes) {
            const node = nodes[i];
            const selected = d3.select(node).classed('selected');
            mainGroup.selectAll('.layer').transition().duration(250)
              .attr('opacity', (d, j) => {
                if (selected) {
                  return 1.0;
                } else {
                  return j !== i ? 0.6 : 1;
                }
              })
            // .classed('hover', (d, i) => {   
            //     return j !== i ? false : true;
            // })
              .attr('stroke', (d, j) => {
                return j !== i ? colorRange[j] : strokecolor;
              })
              .attr('stroke-width', '0.5px');

            d3.select(node).classed('selected', !selected);
            const newSelection = keys[i];
            onKeySelected(newSelection);
          });
        // .on('mousemove', function(d) {
        //   let mousex = d3.mouse(this);
        //   mousex = mousex[0];
        //   const time = x.invert(mousex);
        //   const millisecs = time.getTime();
        //   //   invertedx = invertedx.getMonth() + invertedx.getDate();
        //   //   var selected = (d.values);
        //   let index = 0;
        //   for (; index < data.length; index++) {
        //     if (data[index].date.getTime() >= millisecs) {
        //       break;
        //     }
        //   }
          
        // const value = d[index][1] - d[index][0];
          
        // mainGroup.select(this)
        //     .classed('hover', true)
        //     .attr('stroke', strokecolor)
        //     .attr('stroke-width', '0.5px');

        // tooltip.html( '<p>' + d.key + '<br>' + value + '<br>' + time + '</p>' ).style('visibility', 'visible');
        // })
        // .on('mouseout', function() {
        //   mainGroup.selectAll('.layer')
        //     .transition()
        //     .duration(250)
        //     .attr('opacity', '1');
        //   mainGroup.select(this)
        //     .classed('hover', false)
        //     .attr('stroke-width', '0px');              
        //   // tooltip.html( '<p>' + '</p>' ).style('visibility', 'hidden');
        // });
      
        // var vertical = d3.select('.chart')
        //   .append('div')
        //   .attr('class', 'remove')
        //   .style('position', 'absolute')
        //   .style('z-index', '19')
        //   .style('width', '1px')
        //   .style('height', '380px')
        //   .style('top', '10px')
        //   .style('bottom', '30px')
        //   .style('left', '0px')
        //   .style('background', '#fff');

        mainGroup.select('.chart')
          .on('mousemove', function(){  
            let mousex = d3.mouse(this);
            mousex = mousex[0] + 5;
            const vertical = d3.select('.vertical');

            vertical
              .attr('transform', `translate(${mousex + 5}, 0)`)
              .attr('visibility', 'visible');
          })
          .on('mouseover', function(){  
            let mousex = d3.mouse(this);
            mousex = mousex[0] + 5;
            d3.select('.vertical')
              .attr('transform', `translate(${mousex + 5}, 0)`)
              .attr('visibility', 'visible');
          });
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

export default IDCAreaChart;