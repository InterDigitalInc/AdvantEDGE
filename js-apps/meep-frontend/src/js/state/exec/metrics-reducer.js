/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
import { updateObject } from '../../util/object-util';

const initialState = {
  epochs: [],
  // dataTypeSelected: 'ingressPacketStats',
  dataTypeSelected: 'latency',
  sourceNodeSelected: ''
};

const EXEC_ADD_METRICS_EPOCH = 'EXEC_ADD_METRICS_EPOCH';
function execAddMetricsEpoch(epoch) {
  return {
    type: EXEC_ADD_METRICS_EPOCH,
    payload: epoch
  };
}

const EXEC_CHANGE_SOURCE_NODE_SELECTED = 'EXEC_CHANGE_SOURCE_NODE_SELECTED';
function execChangeSourceNodeSelected(node) {
  return {
    type: EXEC_CHANGE_SOURCE_NODE_SELECTED,
    payload: node
  };
}

const EXEC_CHANGE_DATA_TYPE_SELECTED = 'EXEC_CHANGE_DATA_TYPE_SELECTED';
function execChangeDataTypeSelected(node) {
  return {
    type: EXEC_CHANGE_DATA_TYPE_SELECTED,
    payload: node
  };
}

export { execAddMetricsEpoch, execChangeSourceNodeSelected, execChangeDataTypeSelected };

const NB_EPOCHS_TO_KEEP = 25;
export function metricsReducer(state = initialState, action) {
  const currentId = state.sourceNodeSelected ? state.sourceNodeSelected.data.id : null;
  switch (action.type) {
  case EXEC_ADD_METRICS_EPOCH:
    return updateObject(state, {epochs: state.epochs.splice(-NB_EPOCHS_TO_KEEP).concat([action.payload])});
  case EXEC_CHANGE_SOURCE_NODE_SELECTED:
    if (action.payload.data.id === currentId) {
      return updateObject(state, {sourceNodeSelected: null});
    } else {
      return updateObject(state, {sourceNodeSelected: action.payload});
    }
  case EXEC_CHANGE_DATA_TYPE_SELECTED:
    return updateObject(state, {dataTypeSelected: action.payload});
  default:
    return state;
  }
}
