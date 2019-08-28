import { LATENCY_METRICS, THROUGHPUT_METRICS, MOBILITY_EVENT } from '../meep-constants';

export const dataAccessorForType = dataType => {
  switch (dataType) {
  case LATENCY_METRICS:
    return p => p.data.latency;
  case THROUGHPUT_METRICS:
    return p => p.data.throughput;
  case MOBILITY_EVENT:
    return p => p;
  default:
    return dataAccessorForType(LATENCY_METRICS);
  }
};

export const dataSetterForType = dataType => {
  switch (dataType) {
  case LATENCY_METRICS:
    return val => p => p.data.latency = val;
  case THROUGHPUT_METRICS:
    return val => p => p.data.throughput = val;
  case MOBILITY_EVENT:
    return () => p => p;
  default:
    return dataSetterForType(LATENCY_METRICS);
  }
};

export const isDataPointOfType = type => p => p.dataType === type;