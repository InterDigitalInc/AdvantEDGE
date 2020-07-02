/*
 * Copyright (c) 2020  InterDigital Communications, Inc
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

// CFG_CHANGE_MAP
const CFG_CHANGE_MAP = 'CFG_CHANGE_MAP';
export function cfgChangeMap(map) {
  return {
    type: CFG_CHANGE_MAP,
    payload: map
  };
}

// CFG_CHANGE_UE_LIST
const CFG_CHANGE_UE_LIST = 'CFG_CHANGE_UE_LIST';
export function cfgChangeMapUeList(ueList) {
  return {
    type: CFG_CHANGE_UE_LIST,
    payload: ueList
  };
}

// CFG_CHANGE_POA_LIST
const CFG_CHANGE_POA_LIST = 'CFG_CHANGE_POA_LIST';
export function cfgChangeMapPoaList(poaList) {
  return {
    type: CFG_CHANGE_POA_LIST,
    payload: poaList
  };
}

// CHANGE_COMPUTE_LIST
const CFG_CHANGE_COMPUTE_LIST = 'CFG_CHANGE_COMPUTE_LIST';
export function cfgChangeMapComputeList(computeList) {
  return {
    type: CFG_CHANGE_COMPUTE_LIST,
    payload: computeList
  };
}

export function cfgMapReducer(state = initialState, action) {
  switch (action.type) {
  case CFG_CHANGE_MAP:
    return action.payload;
  case CFG_CHANGE_UE_LIST:
    return updateObject(state, { ueList: action.payload });
  case CFG_CHANGE_POA_LIST:
    return updateObject(state, { poaList: action.payload });
  case CFG_CHANGE_COMPUTE_LIST:
    return updateObject(state, { computeList: action.payload });
  default:
    return state;
  }
}
