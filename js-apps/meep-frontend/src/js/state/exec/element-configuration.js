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

import { updateObject } from '../../util/object-util';
import { createElem } from '../../util/elem-utils';

export const EXEC_ELEM_MODE_NEW = 'CFG_ELEM_MODE_NEW';
export const EXEC_ELEM_MODE_EDIT = 'CFG_ELEM_MODE_EDIT';
export const EXEC_ELEM_MODE_CLONE = 'CFG_ELEM_MODE_CLONE';

const initialState = {
  configuredElement: null,
  configurationMode: null,
  isModified: false,
  errorMessage: ''
};

const EXEC_ELEM_NEW = 'EXEC_ELEM_NEW';
export function execElemNew() {
  return {
    type: EXEC_ELEM_NEW
  };
}

const EXEC_ELEM_EDIT = 'EXEC_ELEM_EDIT';
export function execElemEdit(elem) {
  return {
    type: EXEC_ELEM_EDIT,
    payload: elem
  };
}

const EXEC_ELEM_CLONE = 'EXEC_ELEM_CLONE';
export function execElemClone(elem) {
  return {
    type: EXEC_ELEM_CLONE,
    payload: elem
  };
}

const EXEC_ELEM_CLEAR = 'EXEC_ELEM_CLEAR';
export function execElemClear() {
  return {
    type: EXEC_ELEM_CLEAR
  };
}

const EXEC_ELEM_UPDATE = 'EXEC_ELEM_UPDATE';
export function execElemUpdate(elem) {
  return {
    type: EXEC_ELEM_UPDATE,
    payload: elem
  };
}

const EXEC_ELEM_SET_ERR_MSG = 'EXEC_ELEM_SET_ERR_MSG';
export function execElemSetErrMsg(msg) {
  return {
    type: EXEC_ELEM_SET_ERR_MSG,
    payload: msg
  };
}

export function execElementConfigurationReducer(state = initialState, action) {
  switch (action.type) {
  case EXEC_ELEM_NEW:
    return updateObject(state, {
      configuredElement: createElem(''),
      configurationMode: EXEC_ELEM_MODE_NEW,
      errorMessage: ''
    });
  case EXEC_ELEM_EDIT:
    return updateObject(state, {
      configuredElement: action.payload,
      configurationMode: EXEC_ELEM_MODE_EDIT,
      errorMessage: '',
      isModified: false
    });
  case EXEC_ELEM_CLONE:
    return updateObject(state, {
      configuredElement: action.payload,
      configurationMode: EXEC_ELEM_MODE_CLONE,
      errorMessage: '',
      isModified: true
    });
  case EXEC_ELEM_CLEAR:
    return updateObject(state, initialState);
  case EXEC_ELEM_UPDATE:
    return updateObject(state, { 
      configuredElement: action.payload,
      isModified: true
    });
  case EXEC_ELEM_SET_ERR_MSG:
    return updateObject(state, { errorMessage: action.payload });
  default:
    return state;
  }
}
