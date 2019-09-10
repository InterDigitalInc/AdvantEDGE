import * as d3 from 'd3';

import {
  dataAccessorForType,
  dataSetterForType,
  valueOfPoint
} from '../util/metrics';

import {
  EXEC_ADD_METRICS_EPOCH
} from '../state/exec';

import { LATENCY_METRICS, THROUGHPUT_METRICS, MOBILITY_EVENT } from '../meep-constants';

const extractValue = p => {
  p.value = valueOfPoint(p);
  delete p.data; return p;
};

const introduceMobilityEvents = p => {
  const randVal = Math.random()*4000;

  if (randVal < 1) {
    p.dataType = MOBILITY_EVENT;
    // console.log(`Setting ${p.dataType} value to 0.5 for src ${p.src} and dest ${p.dest}`);
  }

  return p;
};

// Compute avg for each triplet src, dest, dataType
// Create point for each triplet and fill it with avg time and avg value

const mergeEpochPoints = epoch => {
  let pointsMap = epoch.data.reduce((acc, point) => {
    const key = `${point.src},${point.dest},${point.dataType}`;
    if (! acc[key]) {
      acc[key] = [];
    }
    acc[key].push(point);
    return acc;
  }, {});

  const consolidatedEpochData = Object.keys(pointsMap).map(key => {
    const points = pointsMap[key];
    const avgTimestamp = new Date(d3.mean(points, p => new Date(p.timestamp).getTime()));
    let p = {
      src: points[0].src,
      dest: points[0].dest,
      timestamp: avgTimestamp,
      value: d3.mean(points, p => valueOfPoint(p)),
      dataType: points[0].dataType
    };

    return p;
  });

  epoch.data = consolidatedEpochData;
  return epoch;
};

const fixMetricsValuesMiddleware = () => next => action => {
  if (action.type === EXEC_ADD_METRICS_EPOCH) {
    mergeEpochPoints(action.payload); // Will also fix missing latency values through the latency accessor
    let newEpochData = action.payload.data.map(introduceMobilityEvents);
    action.payload.data = newEpochData;
      
    action.changed = true;
    // store.dispatch({
    //   type: EXEC_ADD_METRICS_EPOCH,
    //   payload: newEpoch
    // });
  }

  next(action);
};

export {
  fixMetricsValuesMiddleware
};