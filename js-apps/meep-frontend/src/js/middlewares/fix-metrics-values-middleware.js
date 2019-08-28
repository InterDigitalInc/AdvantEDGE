import {
  dataAccessorForType,
  dataSetterForType
} from '../util/metrics';

import {
  EXEC_ADD_METRICS_EPOCH
} from '../state/exec';

import { LATENCY_METRICS, THROUGHPUT_METRICS, MOBILITY_EVENT } from '../meep-constants';

const fixMissingValue = p => {
  const accessor = dataAccessorForType(p.dataType);
  const setter = dataSetterForType(p.dataType);
  const value = accessor(p);
  if (!value) {
    setter(0.5)(p);
    // console.log(`Setting ${p.dataType} value to 0.5 for src ${p.src} and dest ${p.dest}`);
  }

  if (p.dataType !==LATENCY_METRICS && p.dataType !==THROUGHPUT_METRICS) {
    console.log('Other data type: ', p.dataType);
  }

  return p;
};

const introduceMobilityEvents = p => {
  const accessor = dataAccessorForType(p.dataType);
  const setter = dataSetterForType(p.dataType);
  const value = accessor(p);

  const randVal = Math.random()*4000;

  if (randVal < 1) {
    p.dataType = MOBILITY_EVENT
    // console.log(`Setting ${p.dataType} value to 0.5 for src ${p.src} and dest ${p.dest}`);
  }

  return p;
};

const fixMetricsValuesMiddleware = store => next => action => {
  if (action.type === EXEC_ADD_METRICS_EPOCH) {
    let newEpoch = action.payload.map(fixMissingValue);
    newEpoch = action.payload.map(introduceMobilityEvents);
    action.payload = newEpoch;
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