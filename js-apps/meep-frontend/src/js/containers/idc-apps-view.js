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

import {
  getScenarioNodeChildren
} from '../util/scenario-utils';


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

const IDCAppsView = (
  {
    apps,
    colorRange,
    selectedSource,
    pingBuckets,
    width,
    height,
    onNodeClicked
  }
) => {
  
  const [positioningNeeded, setPositioningNeeded] = useState(true);

  

  const colorForApp = apps.reduce((res, val, i) => {
    // res[val.data.id] = colorRange[i];
    return {...res, [val.data.id]: colorRange[i]};
  }, {});

  //if (positioningNeeded) {
    // copyAttributesRecursive(data)(this.root);
    positionAppsCircle({apps: apps, height: height, width: width});
    //setPositioningNeeded(false);
  //}

  const pingBucket = _.last(pingBuckets);

  if (!pingBucket) {
    return null;
  }

  const appsMap = {};
  _.each(apps, a => appsMap[a.data.id] = a);

  const pings = pingBucket.pings;
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

  const edges = _.flatMap(apps
    .map((d) => {
      const rowObject = m[d.data.id];
      if (!rowObject) {
        return [];
      }
      const destinations = Object.keys(m[d.data.id]);
      return _.map(destinations, (dest) => {
        return  {
          src: d.data.id,
          dest: dest,
          count: rowObject[dest].pings.length,
          color: colorForApp[dest],
          avgLatency: d3.mean(rowObject[dest].pings, d => d.delay)
        };
      });
    }
    )
  ).filter(e => {
    // return nbSelected ? appsMap[e.src].selected : true;
    // console.log(`${appsMap[e.src].data.id}:`, appsMap[e.src]);
    if (selectedSource) {
      return e.src === selectedSource;
    } else {
      return false;
    }
  });

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
        {`Avg lat: ${e.avgLatency.toFixed(2)} ms`}
          
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