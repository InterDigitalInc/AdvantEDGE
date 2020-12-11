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

import { updateObject } from '../../util/object-util';

const initialState = {
  data: [],
  entries: {},
  selected: [],
  order: 'asc',
  orderBy: 'name',
  rowsPerPage: 10,
  page: 0,
  refresh: () => {}
};

const EXEC_CHANGE_TABLE = 'EXEC_CHANGE_TABLE';
export function execChangeTable(table) {
  return {
    type: EXEC_CHANGE_TABLE,
    payload: table
  };
}

export function execTableReducer(state = initialState, action) {
  switch (action.type) {
  case EXEC_CHANGE_TABLE:
    return updateObject({}, action.payload);
  default:
    return state;
  }
}
