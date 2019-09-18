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

import * as cfgTable from '../../../src/js/state/cfg/table-reducer';
import * as execTable from '../../../src/js/state/exec/table-reducer';
import * as ui from '../../../src/js/state/ui';
import * as exec from  '../../../src/js/state/exec';
import uiReducer from '../../../src/js/state/ui';

import {
  scenarioPodsPending,
  scenarioPodsTerminating,
  scenarioPodsTerminated
} from '../../../src/js/state/exec';

describe('Reducers', () => {
  const testTable = [1, 2, 3, 4];

  it('should change exec table properly', () => {
    
    const action = {
      type: execTable.EXEC_CHANGE_TABLE,
      payload: {data: [1, 2, 3]}
    };

    var state = testTable;
    var expectedState = {
      data: [1, 2, 3]
    };

    var newState = execTable.execTableReducer(state, action);

    expect(newState).toEqual(expectedState);
  });

  it('should change cfg table properly', () => {
    
    const action = {
      type: cfgTable.CFG_CHANGE_TABLE,
      payload: {data: [1, 2, 3]}
    };

    var state = testTable;
    var expectedState = {
      data: [1, 2, 3]
    };

    var newState = cfgTable.cfgTableReducer(state, action);

    expect(newState).toEqual(expectedState);
  });

  it('should change ui current event correctly', () => {
    const action = {
      type: ui.EXEC_CHANGE_CURRENT_EVENT,
      payload: ui.UE_MOBILITY_EVENT
    };

    const state = {
      execCurrentEvent: 'Unknown',
      otherElement: 'other'
    };

    var newState = uiReducer(state, action);

    const expectedState = {
      execCurrentEvent: ui.UE_MOBILITY_EVENT,
      otherElement: 'other'
    };

    expect(newState).toEqual(expectedState);
  });

  it('should change exec core pods phases correctly', () => {
    const action = {
      type: exec.EXEC_CHANGE_CORE_PODS_PHASES,
      payload: [1, 2, 3]
    };

    const state = {
      scenario: 'Unknown',
      corePodsPhases: [1, 2],
      scenarioPods: []
    };

    var newState = exec.stateReducer(state, action);

    const expectedState = {
      scenario: 'Unknown',
      corePodsPhases: [1, 2, 3],
      scenarioPods: []
    };

    expect(newState).toEqual(expectedState);
  });

  it('should change exec state correctly', () => {
    const action = {
      type: exec.EXEC_CHANGE_SCENARIO_PODS_PHASES,
      payload: [1, 2, 3, 4]
    };

    const state = {
      scenario: 'Unknown',
      scenarioPodsPhases: [1, 2, 3]
    };

    var newState = exec.stateReducer(state, action);

    const expectedState = {
      scenario: 'Unknown',
      scenarioPodsPhases: [1, 2, 3, 4]
    };

    expect(newState).toEqual(expectedState);
  });

  it('should change vis data correctly', () => {
    
    const action = {
      type: exec.EXEC_CHANGE_VIS_DATA,
      payload: {edges: [1, 2, 3], nodes: [4, 5, 6]}
    };

    const state = {
      other: 'Unknown',
      data: {edges: [1, 2], nodes: [1, 2]}
    };

    var newState = exec.execVisReducer(state, action);

    const expectedState = {
      other: 'Unknown',
      data: {edges: [1, 2, 3], nodes: [4, 5, 6]}
    };

    expect(newState).toEqual(expectedState);
  });

  it('should calculate corePodsRunning selector value properly', () => {
    
    const stateFalse = {
      exec: {
        state: {
          corePodsPhases: [{phase: exec.CORE_PODS_PHASE_RUNNING}, {phase: exec.CORE_PODS_PHASE_PENDING}]
        }
      }
    };
    const stateTrue = {
      exec: {
        state: {
          corePodsPhases: [{phase: exec.CORE_PODS_PHASE_RUNNING}, {phase: exec.CORE_PODS_PHASE_RUNNING}]
        }
      }
    };

    expect(exec.corePodsRunning(stateTrue)).toEqual(true);
    expect(exec.corePodsRunning(stateFalse)).toEqual(false);
  });

  it('should calculate podsWithServiceMaps selector value properly', () => {
    
    const state = {
      exec: {
        state: {
          scenarioPodsPhases: [{name: 'ext-app', phase: exec.CORE_PODS_PHASE_RUNNING}, {name: 'name2', phase: exec.CORE_PODS_PHASE_PENDING}],
          serviceMaps : [{
            client: 'ext-app',
            serviceMap: [
              {
                externalPort: 31101,
                name: 'svc1-edge1',
                port: 8080,
                protocol: 'TCP'
              },
              {
                externalPort: 31102,
                name: 'svc2-edge2',
                port: 8080,
                protocol: 'TCP'
              }
            ]
          }]
        }
      }
    };

    var expectedPodsWithServiceMaps = [
      {
        name: 'ext-app',
        phase: exec.CORE_PODS_PHASE_RUNNING,
        serviceMaps: [{
          externalPort: 31101,
          name: 'svc1-edge1',
          port: 8080,
          protocol: 'TCP'
        },
        {
          externalPort: 31102,
          name: 'svc2-edge2',
          port: 8080,
          protocol: 'TCP'
        }]
      },
      {name: 'name2', phase: exec.CORE_PODS_PHASE_PENDING}
    ];

    expect(exec.podsWithServiceMaps(state)).toEqual(expectedPodsWithServiceMaps);
  });

  it('should calculate scenario pods phases selectors', () => {
    
    const state1 = {
      exec: {
        state: {
          scenarioPodsPhases: [
            {
              logicalState: exec.SCENARIO_PODS_PHASE_RUNNING
            },
            {
              logicalState: exec.SCENARIO_PODS_PHASE_PENDING
            },
            {
              logicalState: exec.SCENARIO_PODS_PHASE_TERMINATING
            }
          ]
        }
      }
    };
    const state2 = {
      exec: {
        state: {
          scenarioPodsPhases: [
            {
              logicalState: exec.SCENARIO_PODS_PHASE_RUNNING
            },
            {
              logicalState: exec.SCENARIO_PODS_PHASE_RUNNING
            },
            {
              logicalState: exec.SCENARIO_PODS_PHASE_RUNNING
            }
          ]
        }
      }
    };

    const state3 = {
      exec: {
        state: {
          scenarioPodsPhases: []
        }
      }
    };


    // state1
    expect(scenarioPodsTerminating(state1)).toEqual(true);
    expect(exec.scenarioPodsPending(state1)).toEqual(true);
    expect(exec.scenarioPodsTerminated(state1)).toEqual(false);

    // state2
    expect(exec.scenarioPodsTerminating(state2)).toEqual(false);
    expect(exec.scenarioPodsPending(state2)).toEqual(false);
    expect(exec.scenarioPodsTerminated(state2)).toEqual(false);

    // state3
    expect(exec.scenarioPodsTerminating(state3)).toEqual(false);
    expect(exec.scenarioPodsPending(state3)).toEqual(false);
    expect(exec.scenarioPodsTerminated(state3)).toEqual(true);
  });

});