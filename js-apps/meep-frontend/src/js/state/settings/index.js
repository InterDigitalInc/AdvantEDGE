/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
 const initialState = {
  debug: false
};

const CHANGE_SETTINGS = 'CHANGE_SETTINGS';
export function changeSettings(settings) {
  return {
    type: CHANGE_SETTINGS,
    payload: settings
  };
}

export default function settingsReducer(state = initialState, action) {
  switch (action.type) {
  case CHANGE_SETTINGS:
    return action.payload;
  default:
    return state;
  }
}
