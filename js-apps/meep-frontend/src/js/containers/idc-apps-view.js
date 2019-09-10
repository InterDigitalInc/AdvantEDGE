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
import React, { useState }  from 'react';
import ReactDOM from 'react-dom';
import * as d3 from 'd3';

import IDCNode from './idc-node.js';

import {
  lineGeneratorNodes
} from './graph-utils';

const edgesFromData = (data, dataAccessor, colorForApp, selectedSource) => {
  const pings = data;
  let m = {};
  _.each(pings, p => {
    if (!m[p.src]) {
      m[p.src] = {};
    }
 
    if (!m[p.src][p.dest]) {
      m[p.src][p.dest] = {
        pings: []
      };
    }
 
    const o = m[p.src][p.dest];
    o.pings.push(p);
  });

  const apps = Object.keys(m);
 
  const edgesFromSource = dataAccessor => src => {
    const rowObject = m[src];
    if (!rowObject) {
      return [];
    }
    const destinations = Object.keys(m[src]);

    const edgesFromDestinations = (dest) => {
      // To debug
      const dataFromPing = p => {
        if (dataAccessor(p)) {
          console.log('Bad value!');  
        }
        return dataAccessor(p);
      };

      if (!d3.mean(rowObject[dest].pings, dataAccessor)) {
        console.log('Bad value!');
      }
      return  {
        src: src,
        dest: dest,
        count: rowObject[dest].pings.length,
        color: colorForApp[dest],
        avgData: d3.mean(rowObject[dest].pings, dataAccessor)
      };
    };
    return _.map(destinations, edgesFromDestinations);
  };

  const outwardEdgesIfSourceSelected = e => {
    if (selectedSource) {
      return e.src === selectedSource;
    } else {
      return true;
    }
  };
  const edges = _.flatMap(apps.map(edgesFromSource(dataAccessor))).filter(outwardEdgesIfSourceSelected);

  return edges; 
};

const positionAppsCircle = ({apps, width, height}) => {
  const cx = width/2.0;
  const cy = height/2.0;
  const PI = 3.141592653598793846264;
  const r = 0.5*height*0.8;
  
  _.each(apps, (app, i) => {
    const theta = (i/apps.length)*(2*PI);
    app.X = cx + r*Math.cos(theta);
    app.Y = cy + r*Math.sin(theta);
  });
};

const edgeLabelForDataType = type => {
  switch(type) {
  case 'latency':
    return 'Latency: ';
  case 'ingressPacketStats':
    return 'Throughput: ';
  default:
    return '';
  }
};

const unitsForDataType = type => {
  switch(type) {
  case 'latency':
    return 'ms';
  case 'ingressPacketStats':
    return 'Kbps';
  default:
    return '';
  }
};

const IDCAppsView = (
  {
    apps,
    colorRange,
    selectedSource,
    data,
    dataAccessor,
    dataType,
    width,
    height,
    onNodeClicked,
    colorForApp,
    displayEdgeLabels
  }
) => {
  
  const [positioningNeeded, setPositioningNeeded] = useState(true);

  //if (positioningNeeded) {
  // copyAttributesRecursive(data)(this.root);
  positionAppsCircle({apps: apps, height: height, width: width});
  //setPositioningNeeded(false);
  //}

  const appsMap = {};
  _.each(apps, a => appsMap[a.data.id] = a);

  const edges = edgesFromData(data.filter(dataAccessor), dataAccessor, colorForApp, selectedSource);

  const edgeLabel = edgeLabelForDataType(dataType);
  const edgeUnits = unitsForDataType(dataType);

  const lineDefs = 
    <defs>
      {
        _.map(edges, (e, i) => {
          return <path
            key={'path' + i}
            id={'textPathDef' + i}
            d={lineGeneratorNodes(appsMap[e.src])(appsMap[e.dest])}
            style={{fill: 'none', 'strokeWidth': e.count*0.1}}
            className='line'
          />;
        })
      }
    </defs>;

  const lines = _.map(edges, (e, i) => {
    return <path
      key={'path' + i}
      id={'path' + i}
      d={lineGeneratorNodes(appsMap[e.src])(appsMap[e.dest])}
      style={{fill: 'none', 'strokeWidth': 0.5, 'stroke': e.color}}
      className='line'
    />;
  });

  const textPaths = _.map(edges, (e,i) =>
    <text key={'textPath' + i} style={{stroke: e.color}}>
      <textPath
        xlinkHref={`#textPathDef${i}`}
        startOffset={'45%'}
      >
        {displayEdgeLabels ? `${edgeLabel} ${e.avgData.toFixed(2)} ${edgeUnits}` : null}
      </textPath>
    </text>
  );

  const nodes = apps
    .map((d, i) =>
      <IDCNode
        collapsible={false}
        key={`node${i}`}
        d={d}
        stroke={colorRange[i]}
        updateParent={() => {}}
        onClick={onNodeClicked}
      />
    );
         
  return (
    <svg
      height={height}
      width={width}
    >
      <>
        {lines}
        {lineDefs}
        {textPaths}
        {nodes}
      </>
    </svg>
  );
};

export default IDCAppsView;