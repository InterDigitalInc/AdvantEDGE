/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
const initialState = {
  dashboardUrl: ''
};
  
const CHANGE_DASHBOARD_URL = 'CHANGE_DASHBOARD_URL';
export function changeDashboardUrl(url) {
  return {
    type: CHANGE_DASHBOARD_URL,
    payload: url
  };
}
  
export default function settingsReducer(state = initialState, action) {
  switch (action.type) {
  case CHANGE_DASHBOARD_URL:
    return {...state, dashboardUrl: action.payload};
  default:
    return state;
  }
}
  