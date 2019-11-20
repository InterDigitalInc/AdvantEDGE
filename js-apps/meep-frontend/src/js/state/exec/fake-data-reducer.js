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

const EXEC_FAKE_CHANGE_SELECTED_DESTINATION =
  'EXEC_FAKE_CHANGE_SELECTED_DESTINATION';
function execFakeChangeSelectedDestination(dest) {
  return {
    type: EXEC_FAKE_CHANGE_SELECTED_DESTINATION,
    payload: dest
  };
}

export {
  execFakeChangePingBuckets,
  execFakeAddPingBucket,
  execFakeChangeSelectedDestination,
  EXEC_FAKE_CHANGE_PING_BUCKETS
};

export function fakeDataReducer(state = initialState, action) {
  switch (action.type) {
  case EXEC_FAKE_CHANGE_PING_BUCKETS:
    return updateObject(state, { pingBuckets: action.payload });
  case EXEC_FAKE_ADD_PING_BUCKET:
    return updateObject(state, {
      pingBuckets: state.pingBuckets.concat([action.payload])
    });
  case EXEC_FAKE_CHANGE_SELECTED_DESTINATION:
    return updateObject(state, { selectedDestination: action.payload });
  default:
    return state;
  }
}
