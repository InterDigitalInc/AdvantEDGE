/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
 import { updateObject } from '../../util/update';
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
    type: CFG_ELEM_NEW,
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
    type: CFG_ELEM_CLEAR,
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
