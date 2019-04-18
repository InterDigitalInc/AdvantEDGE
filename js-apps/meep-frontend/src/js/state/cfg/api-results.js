import { updateObject } from '../../util/update';

const CFG_CHANGE_SCENARIO_LIST = 'CFG_CHANGE_SCENARIO_LIST';
function cfgChangeScenarioList(scenarios) {
  return {
    type: CFG_CHANGE_SCENARIO_LIST,
    payload: scenarios
  };
}

export {
  CFG_CHANGE_SCENARIO_LIST,
  cfgChangeScenarioList
};

const initialState = {
  scenarios: []
};

export function cfgApiResultsReducer(state = initialState, action) {
  switch (action.type) {
  case CFG_CHANGE_SCENARIO_LIST:
    return updateObject(state, {scenarios: action.payload});
  default:
    return state;
  }
}