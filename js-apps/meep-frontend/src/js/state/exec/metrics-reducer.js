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

import { updateObject } from '../../util/object-util';

const initialState = {
  epochs: [],
  // dataTypeSelected: 'ingressPacketStats',
  dataTypeSelected: 'latency',
  sourceNodeSelected: '',
  timeIntervalDuration: 25
};

export const EXEC_ADD_METRICS_EPOCH = 'EXEC_ADD_METRICS_EPOCH';
function execAddMetricsEpoch(epoch) {
  return {
    type: EXEC_ADD_METRICS_EPOCH,
    payload: epoch
  };
}

export const EXEC_CHANGE_SOURCE_NODE_SELECTED =
  'EXEC_CHANGE_SOURCE_NODE_SELECTED';
function execChangeSourceNodeSelected(node) {
  return {
    type: EXEC_CHANGE_SOURCE_NODE_SELECTED,
    payload: node
  };
}

export const EXEC_CHANGE_DEST_NODE_SELECTED =
  'EXEC_CHANGE_DEST_NODE_SELECTED';
function execChangeDestNodeSelected(node) {
  return {
    type: EXEC_CHANGE_DEST_NODE_SELECTED,
    payload: node
  };
}

export const EXEC_CHANGE_DATA_TYPE_SELECTED = 'EXEC_CHANGE_DATA_TYPE_SELECTED';
function execChangeDataTypeSelected(node) {
  return {
    type: EXEC_CHANGE_DATA_TYPE_SELECTED,
    payload: node
  };
}

export const EXEC_CHANGE_METRICS_TIME_INTERVAL_DURATION =
  'EXEC_CHANGE_METRICS_TIME_INTERVAL_DURATION';
function execChangeMetricsTimeIntervalDuration(duration) {
  return {
    type: EXEC_CHANGE_METRICS_TIME_INTERVAL_DURATION,
    payload: duration
  };
}

export const EXEC_CLEAR_METRICS_EPOCHS = 'EXEC_CLEAR_METRICS_EPOCHS';
function execClearMetricsEpochs() {
  return {
    type: EXEC_CLEAR_METRICS_EPOCHS
  };
}

export {
  execAddMetricsEpoch,
  execChangeSourceNodeSelected,
  execChangeDestNodeSelected,
  execChangeDataTypeSelected,
  execChangeMetricsTimeIntervalDuration,
  execClearMetricsEpochs
};

// const NB_EPOCHS_TO_KEEP = 25;
export function metricsReducer(state = initialState, action) {
  const currentSourceNodeId = state.sourceNodeSelected
    ? state.sourceNodeSelected.data.id
    : null;
  const currentDestNodeId = state.destNodeSelected
    ? state.destNodeSelected.data.id
    : null;
  switch (action.type) {
  case EXEC_ADD_METRICS_EPOCH:
    return updateObject(state, {
      epochs: state.epochs
        .splice(-state.timeIntervalDuration)
        .concat([action.payload])
    });
  case EXEC_CHANGE_SOURCE_NODE_SELECTED:
    if (action.payload.data.id === currentSourceNodeId) {
      return updateObject(state, { sourceNodeSelected: null });
    } else {
      return updateObject(state, { sourceNodeSelected: action.payload });
    }
  case EXEC_CHANGE_DEST_NODE_SELECTED:
    if (action.payload.data.id === currentDestNodeId) {
      return updateObject(state, { destNodeSelected: null });
    } else {
      return updateObject(state, { destNodeSelected: action.payload });
    }
  case EXEC_CHANGE_DATA_TYPE_SELECTED:
    return updateObject(state, { dataTypeSelected: action.payload });
  case EXEC_CHANGE_METRICS_TIME_INTERVAL_DURATION:
    return updateObject(state, { timeIntervalDuration: action.payload });
  case EXEC_CLEAR_METRICS_EPOCHS:
    return updateObject(state, { epochs: [] });
  default:
    return state;
  }
}
