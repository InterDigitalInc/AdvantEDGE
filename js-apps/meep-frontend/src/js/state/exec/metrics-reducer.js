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
  sourceNodeSelected: '',
  destNodeSelected: ''
};


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

export {
  execChangeSourceNodeSelected,
  execChangeDestNodeSelected
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
  default:
    return state;
  }
}
