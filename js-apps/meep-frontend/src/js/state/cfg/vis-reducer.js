/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
 import { updateObject } from '../../util/update';

// CHANGE_VIS
const CFG_CHANGE_VIS = 'CFG_CHANGE_VIS';
function cfgChangeVis(vis) {
  return {
    type: CFG_CHANGE_VIS,
    payload: vis
  };
}

// CHANGE_VIS
const CFG_CHANGE_VIS_DATA = 'CFG_CHANGE_VIS_DATA';
function cfgChangeVisData(data) {
  return {
    type: CFG_CHANGE_VIS_DATA,
    payload: data
  };
}

export {
  // Action creators
  cfgChangeVis,
  cfgChangeVisData,

  CFG_CHANGE_VIS
};

const initialState = {
  network: {},
  options: {},
  data: {
    nodes: [],
    edges: [],
  },
  showConfig: false
};

export function cfgVisReducer(state = initialState, action) {
  switch (action.type) {
  case CFG_CHANGE_VIS:
    return action.payload;
  case CFG_CHANGE_VIS_DATA:
    return updateObject(state, {data: action.payload});
  default:
    return state;
  }
}
