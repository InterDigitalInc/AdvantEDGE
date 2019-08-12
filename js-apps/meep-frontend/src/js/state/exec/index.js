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

const execMobTypes = createSelector([execTableElements], elems => {
  return _.filter(elems, elem => (getElemFieldVal(elem, FIELD_TYPE) === 'UE' || getElemFieldVal(elem, FIELD_TYPE) === 'FOG' || getElemFieldVal(elem, FIELD_TYPE) === 'EDGE'));
});

const execEdges = createSelector([execTableElements], elems => {
  return _.filter(elems, elem => getElemFieldVal(elem, FIELD_TYPE) === 'EDGE');
});

const execFogs = createSelector([execTableElements], elems => {
  return _.filter(elems, elem => getElemFieldVal(elem, FIELD_TYPE) === 'FOG');
});

const execZones = createSelector([execTableElements], elems => {
  return _.filter(elems, elem => getElemFieldVal(elem, FIELD_TYPE) === 'ZONE');
});

const execPOAs = createSelector([execTableElements], elems => {
  return _.filter(elems, elem => getElemFieldVal(elem, FIELD_TYPE) === 'POA');
});

export { execUEs, execPOAs, execMobTypes, execEdges, execFogs, execZones };

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
