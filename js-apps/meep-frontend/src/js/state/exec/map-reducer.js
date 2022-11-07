/*
 * Copyright (c) 2022  The AdvantEDGE Authors
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
  ueList: [],
  poaList: [],
  computeList: []
};

// CHANGE_MAP
const EXEC_CHANGE_MAP = 'EXEC_CHANGE_MAP';
export function execChangeMap(map) {
  return {
    type: EXEC_CHANGE_MAP,
    payload: map
  };
}

// EXEC_CHANGE_UE_LIST
const EXEC_CHANGE_UE_LIST = 'EXEC_CHANGE_UE_LIST';
export function execChangeMapUeList(ueList) {
  return {
    type: EXEC_CHANGE_UE_LIST,
    payload: ueList
  };
}

// EXEC_CHANGE_POA_LIST
const EXEC_CHANGE_POA_LIST = 'EXEC_CHANGE_POA_LIST';
export function execChangeMapPoaList(poaList) {
  return {
    type: EXEC_CHANGE_POA_LIST,
    payload: poaList
  };
}

// EXEC_CHANGE_COMPUTE_LIST
const EXEC_CHANGE_COMPUTE_LIST = 'EXEC_CHANGE_COMPUTE_LIST';
export function execChangeMapComputeList(computeList) {
  return {
    type: EXEC_CHANGE_COMPUTE_LIST,
    payload: computeList
  };
}

export function execMapReducer(state = initialState, action) {
  switch (action.type) {
  case EXEC_CHANGE_MAP:
    return action.payload;
  case EXEC_CHANGE_UE_LIST:
    return updateObject(state, { ueList: action.payload });
  case EXEC_CHANGE_POA_LIST:
    return updateObject(state, { poaList: action.payload });
  case EXEC_CHANGE_COMPUTE_LIST:
    return updateObject(state, { computeList: action.payload });
  default:
    return state;
  }
}
