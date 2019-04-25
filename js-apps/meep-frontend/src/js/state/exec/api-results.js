/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
import { updateObject } from '../../util/update';

const EXEC_CHANGE_SCENARIO_LIST = 'EXEC_CHANGE_SCENARIO_LIST';
function execChangeScenarioList(scenarios) {
  return {
    type: EXEC_CHANGE_SCENARIO_LIST,
    payload: scenarios
  };
}

export {
  EXEC_CHANGE_SCENARIO_LIST,
  execChangeScenarioList
};

const initialState = {
  scenarios: []
};

export function execApiResultsReducer(state = initialState, action) {
  switch (action.type) {
  case EXEC_CHANGE_SCENARIO_LIST:
    return updateObject(state, {scenarios: action.payload});
  default:
    return state;
  }
}
