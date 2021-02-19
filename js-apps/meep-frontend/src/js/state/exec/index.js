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
import { displayedScenarioReducer } from './displayed-scenario-reducer';
import { execMapReducer } from './map-reducer';
import { execVisReducer } from './vis-reducer';
import { execTableReducer } from './table-reducer';
import { execSelectedScenarioElement } from './selected-scenario-element';
import { execApiResultsReducer } from './api-results';
import { execElementConfigurationReducer } from './element-configuration';

import {
  getElemFieldVal,
  FIELD_GROUP,
  FIELD_TYPE
} from '../../util/elem-utils';

import {
  ELEMENT_TYPE_ZONE,
  ELEMENT_TYPE_POA,
  ELEMENT_TYPE_POA_4G,
  ELEMENT_TYPE_POA_5G,
  ELEMENT_TYPE_POA_WIFI,
  ELEMENT_TYPE_EDGE,
  ELEMENT_TYPE_FOG,
  ELEMENT_TYPE_UE,
  ELEMENT_TYPE_EDGE_APP
} from '../../meep-constants';

export * from './type-reducer';
export * from './state-reducer';
export * from './scenario-reducer';
export * from './displayed-scenario-reducer';
export * from './map-reducer';
export * from './vis-reducer';
export * from './table-reducer';
export * from './selected-scenario-element';
export * from './api-results';
export * from './element-configuration';

const execTableElements = state => state.exec.table.entries;
const execUEs = createSelector(
  [execTableElements],
  elems => {
    return _.filter(elems, elem => getElemFieldVal(elem, FIELD_TYPE) === ELEMENT_TYPE_UE);
  }
);

const execMobTypes = createSelector(
  [execTableElements],
  elems => {
    return _.filter(
      elems,
      elem =>
        getElemFieldVal(elem, FIELD_TYPE) === ELEMENT_TYPE_UE ||
        getElemFieldVal(elem, FIELD_TYPE) === ELEMENT_TYPE_FOG ||
        getElemFieldVal(elem, FIELD_TYPE) === ELEMENT_TYPE_EDGE ||
        (getElemFieldVal(elem, FIELD_TYPE) === ELEMENT_TYPE_EDGE_APP &&
          getElemFieldVal(elem, FIELD_GROUP) === '')
    );
  }
);

const execFogEdges = createSelector(
  [execTableElements],
  elems => {
    return _.filter(
      elems,
      elem =>
        getElemFieldVal(elem, FIELD_TYPE) === ELEMENT_TYPE_FOG ||
        getElemFieldVal(elem, FIELD_TYPE) === ELEMENT_TYPE_EDGE
    );
  }
);

const execEdgeApps = createSelector(
  [execTableElements],
  elems => {
    return _.filter(
      elems,
      elem => getElemFieldVal(elem, FIELD_TYPE) === ELEMENT_TYPE_EDGE_APP
    );
  }
);

const execEdges = createSelector(
  [execTableElements],
  elems => {
    return _.filter(
      elems,
      elem => getElemFieldVal(elem, FIELD_TYPE) === ELEMENT_TYPE_EDGE
    );
  }
);

const execFogs = createSelector(
  [execTableElements],
  elems => {
    return _.filter(elems, elem => getElemFieldVal(elem, FIELD_TYPE) === ELEMENT_TYPE_FOG);
  }
);

const execZones = createSelector(
  [execTableElements],
  elems => {
    return _.filter(
      elems,
      elem => getElemFieldVal(elem, FIELD_TYPE) === ELEMENT_TYPE_ZONE
    );
  }
);

const execPOAs = createSelector(
  [execTableElements],
  elems => {
    return _.filter(
      elems,
      elem =>
        getElemFieldVal(elem, FIELD_TYPE) === ELEMENT_TYPE_POA ||
        getElemFieldVal(elem, FIELD_TYPE) === ELEMENT_TYPE_POA_4G ||
        getElemFieldVal(elem, FIELD_TYPE) === ELEMENT_TYPE_POA_5G ||
        getElemFieldVal(elem, FIELD_TYPE) === ELEMENT_TYPE_POA_WIFI
    );
  }
);

export {
  execUEs,
  execPOAs,
  execMobTypes,
  execEdges,
  execFogs,
  execZones,
  execEdgeApps,
  execFogEdges
};

const execReducer = combineReducers({
  type: typeReducer,
  state: stateReducer,
  scenario: scenarioReducer,
  displayedScenario: displayedScenarioReducer,
  map: execMapReducer,
  vis: execVisReducer,
  table: execTableReducer,
  selectedScenarioElement: execSelectedScenarioElement,
  apiResults: execApiResultsReducer,
  elementConfiguration: execElementConfigurationReducer
});

export default execReducer;
