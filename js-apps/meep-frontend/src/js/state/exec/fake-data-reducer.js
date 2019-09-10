/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
import { updateObject } from '../../util/object-util';

const initialState = {
  pingBuckets: [],
  selectedDestination: ''
};

const EXEC_FAKE_CHANGE_PING_BUCKETS = 'EXEC_FAKE_CHANGE_PING_BUCKETS';
function execFakeChangePingBuckets(pings) {
  return {
    type: EXEC_FAKE_CHANGE_PING_BUCKETS,
    payload: pings
  };
}

const EXEC_FAKE_ADD_PING_BUCKET = 'EXEC_FAKE_ADD_PING_BUCKET';
function execFakeAddPingBucket(bucket) {
  return {
    type: EXEC_FAKE_ADD_PING_BUCKET,
    payload: bucket
  };
}

const EXEC_FAKE_CHANGE_SELECTED_DESTINATION = 'EXEC_FAKE_CHANGE_SELECTED_DESTINATION';
function execFakeChangeSelectedDestination(dest) {
  return {
    type: EXEC_FAKE_CHANGE_SELECTED_DESTINATION,
    payload: dest
  };
}

export { execFakeChangePingBuckets, execFakeAddPingBucket, execFakeChangeSelectedDestination, EXEC_FAKE_CHANGE_PING_BUCKETS };

export function fakeDataReducer(state = initialState, action) {
  switch (action.type) {
  case EXEC_FAKE_CHANGE_PING_BUCKETS:
    return updateObject(state, {pingBuckets: action.payload});
  case EXEC_FAKE_ADD_PING_BUCKET:
    return updateObject(state, {pingBuckets: state.pingBuckets.concat([action.payload])});
  case EXEC_FAKE_CHANGE_SELECTED_DESTINATION:
    return updateObject(state, {selectedDestination: action.payload});
  default:
    return state;
  }
}
