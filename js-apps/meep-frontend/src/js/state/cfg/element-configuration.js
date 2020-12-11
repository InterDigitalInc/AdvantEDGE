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
import { createElem } from '../../util/elem-utils';

export const CFG_ELEM_MODE_NEW = 'CFG_ELEM_MODE_NEW';
export const CFG_ELEM_MODE_EDIT = 'CFG_ELEM_MODE_EDIT';
export const CFG_ELEM_MODE_CLONE = 'CFG_ELEM_MODE_CLONE';

const initialState = {
  configuredElement: null,
  configurationMode: null,
  isModified: false,
  errorMessage: ''
};

const CFG_ELEM_NEW = 'CFG_ELEM_NEW';
export function cfgElemNew() {
  return {
    type: CFG_ELEM_NEW
  };
}

const CFG_ELEM_EDIT = 'CFG_ELEM_EDIT';
export function cfgElemEdit(elem) {
  return {
    type: CFG_ELEM_EDIT,
    payload: elem
  };
}

const CFG_ELEM_CLONE = 'CFG_ELEM_CLONE';
export function cfgElemClone(elem) {
  return {
    type: CFG_ELEM_CLONE,
    payload: elem
  };
}

const CFG_ELEM_CLEAR = 'CFG_ELEM_CLEAR';
export function cfgElemClear() {
  return {
    type: CFG_ELEM_CLEAR
  };
}

const CFG_ELEM_UPDATE = 'CFG_ELEM_UPDATE';
export function cfgElemUpdate(elem) {
  return {
    type: CFG_ELEM_UPDATE,
    payload: elem
  };
}

const CFG_ELEM_SET_ERR_MSG = 'CFG_ELEM_SET_ERR_MSG';
export function cfgElemSetErrMsg(msg) {
  return {
    type: CFG_ELEM_SET_ERR_MSG,
    payload: msg
  };
}

export function cfgElementConfigurationReducer(state = initialState, action) {
  switch (action.type) {
  case CFG_ELEM_NEW:
    return updateObject(state, {
      configuredElement: createElem(''),
      configurationMode: CFG_ELEM_MODE_NEW,
      errorMessage: ''
    });
  case CFG_ELEM_EDIT:
    return updateObject(state, {
      configuredElement: action.payload,
      configurationMode: CFG_ELEM_MODE_EDIT,
      errorMessage: '',
      isModified: false
    });
  case CFG_ELEM_CLONE:
    return updateObject(state, {
      configuredElement: action.payload,
      configurationMode: CFG_ELEM_MODE_CLONE,
      errorMessage: '',
      isModified: true
    });
  case CFG_ELEM_CLEAR:
    return updateObject(state, initialState);
  case CFG_ELEM_UPDATE:
    return updateObject(state, { 
      configuredElement: action.payload,
      isModified: true
    });
  case CFG_ELEM_SET_ERR_MSG:
    return updateObject(state, { errorMessage: action.payload });
  default:
    return state;
  }
}
