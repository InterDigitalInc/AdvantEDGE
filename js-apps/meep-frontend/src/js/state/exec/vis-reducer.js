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

import _ from 'lodash';
import { updateObject } from '../../util/object-util';
import { createSelector } from 'reselect';

import {
  ELEMENT_TYPE_UE_APP,
  ELEMENT_TYPE_EDGE_APP,
  ELEMENT_TYPE_CLOUD_APP,
  TYPE_EXEC
} from '../../meep-constants';

import {
  getElemFieldVal,
  FIELD_TYPE
} from '../../util/elem-utils';

const initialState = {
  type: TYPE_EXEC,
  network: {},
  options: {},
  data: {
    nodes: [],
    edges: []
  },
  showConfig: false,
  showApps: false
};

// CHANGE_VIS
const EXEC_CHANGE_VIS = 'EXEC_CHANGE_VIS';
export function execChangeVis(vis) {
  return {
    type: EXEC_CHANGE_VIS,
    payload: vis
  };
}

// CHANGE_VIS
const EXEC_CHANGE_VIS_DATA = 'EXEC_CHANGE_VIS_DATA';
export function execChangeVisData(data) {
  return {
    type: EXEC_CHANGE_VIS_DATA,
    payload: data
  };
}

const dataSelector = state => state.exec.vis.data;
const tableSelector = state => state.exec.table;
const showAppsSelector = state => state.ui.execShowApps;
export const execVisFilteredData = createSelector(
  [dataSelector, tableSelector, showAppsSelector],
  (data, table, showApps) => {
    var appTypes = [
      ELEMENT_TYPE_UE_APP,
      ELEMENT_TYPE_EDGE_APP,
      ELEMENT_TYPE_CLOUD_APP
    ];

    var types = {};
    _.each(table.entries, entry => {
      types[entry.id] = getElemFieldVal(entry,FIELD_TYPE);
    });

    var newNodes = data.nodes;
    var newEdges = data.edges;

    if (!showApps) {
      if (data.nodes.length) {
        _.forOwn(data.nodes.get(), (elem) => {
          if (_.includes(appTypes, types[elem.id])) {
            newNodes.remove(elem.id);
          }
        });
      }
      
      if (data.nodes.length) {
        _.forOwn(data.edges.get(), (edge) => {
          if (
            _.includes(appTypes, types[edge.from]) ||
            _.includes(appTypes, types[edge.to])
          ) {
            newEdges.remove(edge.id);
          }
        });
      }
    }

    return { nodes: newNodes, edges: newEdges };
  }
);

export function execVisReducer(state = initialState, action) {
  switch (action.type) {
  case EXEC_CHANGE_VIS:
    return action.payload;
  case EXEC_CHANGE_VIS_DATA:
    return updateObject(state, { data: action.payload });
  default:
    return state;
  }
}
