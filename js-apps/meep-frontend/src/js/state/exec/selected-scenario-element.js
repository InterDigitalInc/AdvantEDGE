/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
import { updateObject } from '../../util/update';
import { createElem } from '../../util/elem-utils';

const EXEC_CHANGE_SELECTED_SCENARIO_ELEMENT = 'EXEC_CHANGE_SELECTED_SCENARIO_ELEMENT';

// CFG_SET_EDITED_ELEMENT
function execChangeSelectedScenarioElement(element) {
  return {
    type: EXEC_CHANGE_SELECTED_SCENARIO_ELEMENT,
    payload: element
  };
}

export {
  execChangeSelectedScenarioElement
};

const initialState = createElem('dummy');

export function execSelectedScenarioElement(state = initialState, action) {
  switch (action.type) {
  case EXEC_CHANGE_SELECTED_SCENARIO_ELEMENT:
    return updateObject({}, action.payload);
  default:
    return state;
  }
}
