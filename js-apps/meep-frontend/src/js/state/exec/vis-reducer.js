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
import { updateObject } from '../../util/object-util';
import {createSelector } from 'reselect';
import * as vis from 'vis';

import {
  ELEMENT_TYPE_UE_APP,
  ELEMENT_TYPE_EDGE_APP,
  ELEMENT_TYPE_CLOUD_APP
} from '../../meep-constants';

import {
  getElemFieldVal,
  FIELD_NAME,
  FIELD_TYPE
} from '../../util/elem-utils';

// CHANGE_VIS
const EXEC_CHANGE_VIS = 'EXEC_CHANGE_VIS';
function execChangeVis(vis) {
  return {
    type: EXEC_CHANGE_VIS,
    payload: vis
  };
}

// CHANGE_VIS
const EXEC_CHANGE_VIS_DATA = 'EXEC_CHANGE_VIS_DATA';
function execChangeVisData(data) {
  return {
    type: EXEC_CHANGE_VIS_DATA,
    payload: data
  };
}

export {
  // Action creators
  execChangeVis,
  execChangeVisData,
  EXEC_CHANGE_VIS,
  EXEC_CHANGE_VIS_DATA
};

const dataSelector = state => state.exec.vis.data;
const tableSelector = state => state.exec.table;
const showAppsSelector = state => state.ui.execShowApps;
const execVisFilteredData = createSelector([dataSelector, tableSelector, showAppsSelector], (data, table, showApps) => {
  var appTypes = [
    ELEMENT_TYPE_UE_APP,
    ELEMENT_TYPE_EDGE_APP,
    ELEMENT_TYPE_CLOUD_APP
  ];

  var types = {};
  _.each(table.entries, entry => {
    types[getElemFieldVal(entry, FIELD_NAME)] = getElemFieldVal(entry, FIELD_TYPE);
  });

  if (showApps) {
    newNodes = data.nodes;
    newEdges = data.edges;
  } else {

    var newNodes = new vis.DataSet(_.values(data.nodes._data));
    _.forOwn(data.nodes._data, (elem, key) => {
      if (_.includes(appTypes, types[elem.id])) {
        newNodes.remove(key);
      }
    });

    var newEdges = new vis.DataSet(_.values(data.edges._data));
    _.forOwn(data.edges._data, (edge, key) => {
      if (_.includes(appTypes, types[edge.from]) || _.includes(appTypes, types[edge.to])) {
        newEdges.remove(key);
      }
    });
  }
  return {nodes: newNodes, edges: newEdges};
});

export {
  execVisFilteredData
};

const initialState = {
  network: {},
  options: {},
  data: {
    nodes: [],
    edges: []
  },
  showConfig: false,
  showApps: false
};

export function execVisReducer(state = initialState, action) {
  switch (action.type) {
  case EXEC_CHANGE_VIS:
    return action.payload;
  case EXEC_CHANGE_VIS_DATA:
    return updateObject(state, {data: action.payload});
  default:
    return state;
  }
}
