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

import * as d3 from 'd3';

import { valueOfPoint } from '../util/metrics';

import { EXEC_ADD_METRICS_EPOCH } from '../state/exec';

// Compute avg for each triplet src, dest, dataType
// Create point for each triplet and fill it with avg time and avg value

let mobilityEventIndex = 1;
const mergeEpochPoints = epoch => {
  let pointsMap = epoch.data.reduce((acc, point) => {
    const key = `${point.src},${point.dest},${point.dataType}`;
    if (!acc[key]) {
      acc[key] = [];
    }
    acc[key].push(point);
    return acc;
  }, {});

  const consolidatedEpochData = Object.keys(pointsMap).map(key => {
    const points = pointsMap[key];
    const avgTimestamp = new Date(
      d3.mean(points, p => new Date(p.timestamp).getTime())
    );
    let p = {
      src: points[0].src,
      dest: points[0].dest,
      timestamp: avgTimestamp,
      value: d3.mean(points, p => valueOfPoint(p)),
      dataType: points[0].dataType
    };

    if (p.dataType === 'mobilityEvent') {
      p.dest = points[0].data.newPoa;
      p.mobilityEventIndex = mobilityEventIndex++;
    }

    return p;
  });

  epoch.data = consolidatedEpochData;
  return epoch;
};

const fixMetricsValuesMiddleware = () => next => action => {
  if (action.type === EXEC_ADD_METRICS_EPOCH) {
    mergeEpochPoints(action.payload); // Will also fix missing latency values through the latency accessor
  }

  next(action);
};

export { fixMetricsValuesMiddleware };
