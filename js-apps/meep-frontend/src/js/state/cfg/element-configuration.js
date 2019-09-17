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

const CFG_ELEM_NEW = 'CFG_ELEM_NEW';
const CFG_ELEM_EDIT = 'CFG_ELEM_EDIT';
const CFG_ELEM_CLEAR = 'CFG_ELEM_CLEAR';
const CFG_ELEM_UPDATE = 'CFG_ELEM_UPDATE';
const CFG_ELEM_SET_ERR_MSG = 'CFG_ELEM_SET_ERR_MSG';

const CFG_ELEM_MODE_NEW = 'CFG_ELEM_MODE_NEW';
const CFG_ELEM_MODE_EDIT = 'CFG_ELEM_MODE_EDIT';

function cfgElemNew() {
  return {
    type: CFG_ELEM_NEW
  };
}

function cfgElemEdit(elem) {
  return {
    type: CFG_ELEM_EDIT,
    payload: elem
  };
}

function cfgElemClear() {
  return {
    type: CFG_ELEM_CLEAR
  };
}

function cfgElemUpdate(elem) {
  return {
    type: CFG_ELEM_UPDATE,
    payload: elem
  };
}

function cfgElemSetErrMsg(msg) {
  return {
    type: CFG_ELEM_SET_ERR_MSG,
    payload: msg
  };
}

export {
  cfgElemNew,
  cfgElemEdit,
  cfgElemClear,
  cfgElemUpdate,
  cfgElemSetErrMsg,
  CFG_ELEM_MODE_NEW,
  CFG_ELEM_MODE_EDIT
};

const initialState = {
  configuredElement: null,
  configurationMode: null,
  errorMessage: ''
};

export function cfgElementConfigurationReducer(state = initialState, action) {
  switch (action.type) {
  case CFG_ELEM_NEW:
    return updateObject(state, {configuredElement: createElem(''), configurationMode: CFG_ELEM_MODE_NEW, errorMessage: ''});
  case CFG_ELEM_EDIT:
    return updateObject(state, {configuredElement: action.payload, configurationMode: CFG_ELEM_MODE_EDIT, errorMessage: ''});
  case CFG_ELEM_CLEAR:
    return updateObject(state, initialState);
  case CFG_ELEM_UPDATE:
    return updateObject(state, {configuredElement: action.payload});
  case CFG_ELEM_SET_ERR_MSG:
    return updateObject(state, {errorMessage: action.payload});
  default:
    return state;
  }
}
