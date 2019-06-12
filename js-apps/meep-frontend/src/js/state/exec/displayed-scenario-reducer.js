/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
import { updateObject } from '../../util/object-util';

// EXEC_CHANGE_SCENARIO
export const EXEC_CHANGE_DISPLAYED_SCENARIO = 'EXEC_CHANGE_DISPLAYED_SCENARIO';
function execChangeDisplayedScenario(scenario) {
  return {
    type: EXEC_CHANGE_DISPLAYED_SCENARIO,
    payload: scenario
  };
}

export {
  // Action creators
  execChangeDisplayedScenario
};

const initialState = {
  
};

export function displayedScenarioReducer(state = initialState, action) {
  switch (action.type) {
  case EXEC_CHANGE_DISPLAYED_SCENARIO:
    return action.payload;
  default:
    return state;
  }
}
