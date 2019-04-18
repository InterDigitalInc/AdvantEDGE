import _ from 'lodash';
import { combineReducers } from 'redux';
import { createSelector } from 'reselect';

import { typeReducer } from './type-reducer';
import { stateReducer } from './state-reducer';
import { scenarioReducer } from './scenario-reducer';
import { execVisReducer } from './vis-reducer';
import { execTableReducer } from './table-reducer';
import { execSelectedScenarioElement } from './selected-scenario-element';
import { execApiResultsReducer } from './api-results';
import { getElemFieldVal, FIELD_TYPE } from '../../util/elem-utils';

export * from './type-reducer';
export * from './state-reducer';
export * from './scenario-reducer';
export * from './vis-reducer';
export * from './table-reducer';
export * from './selected-scenario-element';
export * from './api-results';

const execTableElements = state => state.exec.table.entries;
const execUEs = createSelector([execTableElements], elems => {
  return _.filter(elems, elem => getElemFieldVal(elem, FIELD_TYPE) === 'UE');
});

const execPOAs = createSelector([execTableElements], elems => {
  return _.filter(elems, elem => getElemFieldVal(elem, FIELD_TYPE) === 'POA');
});

export { execUEs, execPOAs };

const execReducer = combineReducers({
  type: typeReducer,
  state: stateReducer,
  scenario: scenarioReducer,
  vis: execVisReducer,
  table: execTableReducer,
  selectedScenarioElement: execSelectedScenarioElement,
  apiResults: execApiResultsReducer
});

export default execReducer;
