import {
  ME_LATENCY_METRICS,
  ME_THROUGHPUT_METRICS,
  ME_MOBILITY_EVENT
} from '../meep-constants';

export const dataAccessorForType = dataType => {
  switch (dataType) {
  case ME_LATENCY_METRICS:
    return p => p.data.latency || 0.5;
  case ME_THROUGHPUT_METRICS:
    return p => p.data.throughput;
  case ME_MOBILITY_EVENT:
    return p => p;
  default:
    return dataAccessorForType(ME_LATENCY_METRICS);
  }
};

export const dataSetterForType = dataType => {
  switch (dataType) {
  case ME_LATENCY_METRICS:
    return val => p => p.data.latency = val;
  case ME_THROUGHPUT_METRICS:
    return val => p => p.data.throughput = val;
  case ME_MOBILITY_EVENT:
    return () => p => p;
  default:
    return dataSetterForType(ME_LATENCY_METRICS);
  }
};

export const isDataPointOfType = type => p => p.dataType === type;

export const valueOfPoint = p => dataAccessorForType(p.dataType)(p);