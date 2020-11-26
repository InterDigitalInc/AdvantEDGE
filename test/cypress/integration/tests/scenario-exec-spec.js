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

// Import MEEP Contstants
import * as meep from '../../../../js-apps/meep-frontend/src/js/meep-constants';

// Import element utils
import {
  // Field Names
  FIELD_PARENT,
  FIELD_NAME,
  FIELD_INT_DOM_LATENCY,
  FIELD_INT_DOM_LATENCY_VAR,
  FIELD_INT_DOM_LATENCY_DIST,
  FIELD_INT_DOM_THROUGHPUT_DL,
  FIELD_INT_DOM_THROUGHPUT_UL,
  FIELD_INT_DOM_PKT_LOSS,
  FIELD_INT_ZONE_LATENCY,
  FIELD_INT_ZONE_LATENCY_VAR,
  FIELD_INT_ZONE_THROUGHPUT_DL,
  FIELD_INT_ZONE_THROUGHPUT_UL,
  FIELD_INT_ZONE_PKT_LOSS,
  FIELD_INTRA_ZONE_LATENCY,
  FIELD_INTRA_ZONE_LATENCY_VAR,
  FIELD_INTRA_ZONE_THROUGHPUT_DL,
  FIELD_INTRA_ZONE_THROUGHPUT_UL,
  FIELD_INTRA_ZONE_PKT_LOSS,
  FIELD_TERM_LINK_LATENCY,
  FIELD_TERM_LINK_LATENCY_VAR,
  FIELD_TERM_LINK_THROUGHPUT_DL,
  FIELD_TERM_LINK_THROUGHPUT_UL,
  FIELD_TERM_LINK_PKT_LOSS,

  getElemFieldVal,
} from '../../../../js-apps/meep-frontend/src/js/util/elem-utils';

// Import Test utility functions
import { selector, click, check, type, select, verify, verifyEnabled, verifyForm, verifyChecked } from '../util/util';

// Scenario Execution Tests
describe('Scenario Execution', function () {

  // Test Variables
  let defaultScenario = 'None';
  let sandbox = 'sbox-test';
  let scenario = 'demo1';
  let scenario2 = 'demo2';
  let replayEventsName = 'replaydemo1'
  let replayEventsDescription = 'replay demo1 description'

  // ------------------------------
  //            TESTS
  // ------------------------------

  // Demo1 scenario testing (Virt-Engine)
  it('Deploy DEMO1 scenario', function () {
    openDefaultMeepUrl();

    // Create Sandbox
    createSandbox(sandbox);

    // Deploy demo scenario
    deployScenario(scenario);
  });

  // Test manual events
  it('Test Event - Manual button - Cancel', function () {
    openDefaultMeepUrl();
    click(meep.MEEP_TAB_EXEC);
    cy.wait(1000);
    select(meep.EXEC_SELECT_SANDBOX, sandbox);
    cy.wait(3000);

    testCancelEvent();
  });

  it('Test Event - Manual button - Mobility event', function () {
    openDefaultMeepUrl();
    click(meep.MEEP_TAB_EXEC);
    cy.wait(1000);
    select(meep.EXEC_SELECT_SANDBOX, sandbox);
    cy.wait(3000);

    testMobilityEvent();
  });

  it('Test Event - Manual button - NetChar event', function () {
    openDefaultMeepUrl();
    click(meep.MEEP_TAB_EXEC);
    cy.wait(1000);
    select(meep.EXEC_SELECT_SANDBOX, sandbox);
    cy.wait(3000);

    testNetCharEvent(scenario);
  });

  // Test Save events
  it('Test Event - Save button', function () {
    openDefaultMeepUrl();
    click(meep.MEEP_TAB_EXEC);
    cy.wait(1000);
    select(meep.EXEC_SELECT_SANDBOX, sandbox);
    cy.wait(3000);

    testSaveReplayEvents(replayEventsName, replayEventsDescription);
  });

  // Test Auto-Replay
  it('Test Event - Auto-Replay button', function () {
    openDefaultMeepUrl();
    click(meep.MEEP_TAB_EXEC);
    cy.wait(1000);
    select(meep.EXEC_SELECT_SANDBOX, sandbox);
    cy.wait(3000);

    testAutoReplay(replayEventsName);
  });

  // Test Automation
  it('Test Event - Automation button', function () {
    openDefaultMeepUrl();
    click(meep.MEEP_TAB_EXEC);
    cy.wait(1000);
    select(meep.EXEC_SELECT_SANDBOX, sandbox);
    cy.wait(3000);

    testEventAutomation();
  });

  // Open/Close panes
  it('Test Event and Dashboard panes visibility', function () {
    openDefaultMeepUrl();
    click(meep.MEEP_TAB_EXEC);
    cy.wait(1000);
    select(meep.EXEC_SELECT_SANDBOX, sandbox);
    cy.wait(3000);

    testEventPane();
    testDashboardPane();
    testPanes();
  });

  it('DEMO1 scenario - terminate', function () {
    openDefaultMeepUrl();
    click(meep.MEEP_TAB_EXEC);
    cy.wait(1000);
    select(meep.EXEC_SELECT_SANDBOX, sandbox);
    cy.wait(3000);

    // Terminate demo scenario
    terminateScenario(scenario);

    // Destroy Sandbox
    destroySandbox(sandbox);
  });

  // Demo2 scenario testing (User Charts)
  it('Deploy & Test DEMO2 scenario', function () {
    openDefaultMeepUrl();
    cy.wait(3000);

    // Create Sandbox
    createSandbox(sandbox);

    // Deploy demo scenario
    deployScenario(scenario2);

    // Test events
    testCancelEvent();
    testMobilityEvent();
    testNetCharEvent(scenario2);

  });

  //separating deployment and termination to force termination to be called if the test above fails during execution
  it('DEMO2 scenario - terminate', function () {
    openDefaultMeepUrl();
    click(meep.MEEP_TAB_EXEC);
    cy.wait(1000);
    select(meep.EXEC_SELECT_SANDBOX, sandbox);
    cy.wait(3000);

    // Terminate demo scenario
    terminateScenario(scenario2);

    // Destroy Sandbox
    destroySandbox(sandbox);
  });

  // ------------------------------
  //          FUNCTIONS
  // ------------------------------

  // Opens the default meepUrl page
  function openDefaultMeepUrl() {
    var meepUrl = Cypress.env('meep_url');
    if (meepUrl == null) {
      meepUrl = 'http://127.0.0.1';
    }

    cy.viewport(1920, 1080);
    cy.visit(meepUrl);
    cy.wait(1000);
  }

  // Create sandbox with provided name
  function createSandbox(name) {
    // Go to execution page
    cy.log('Go to execution page');
    click(meep.MEEP_TAB_EXEC);
    cy.wait(1000);
    verify(meep.MEEP_LBL_SCENARIO_NAME, 'contain', defaultScenario);
    verifyEnabled(meep.EXEC_BTN_NEW_SANDBOX, true);
    verifyEnabled(meep.EXEC_BTN_DELETE_SANDBOX, false);
    verifyEnabled(meep.EXEC_BTN_DEPLOY, false);
    verifyEnabled(meep.EXEC_BTN_TERMINATE, false);
    verifyEnabled(meep.EXEC_BTN_EVENT, false);

    // Create sandbox
    cy.log('Create sandbox: ' + name);
    click(meep.EXEC_BTN_NEW_SANDBOX);
    type(meep.MEEP_DLG_NEW_SANDBOX_NAME, name);
    click(meep.MEEP_DLG_NEW_SANDBOX, 'Ok');
    cy.wait(15000);
    verifyEnabled(meep.EXEC_BTN_NEW_SANDBOX, true);
    verifyEnabled(meep.EXEC_BTN_DELETE_SANDBOX, true);
    verifyEnabled(meep.EXEC_BTN_DEPLOY, true);
    verifyEnabled(meep.EXEC_BTN_TERMINATE, false);
    verifyEnabled(meep.EXEC_BTN_EVENT, false);
  }

  // Destroy sandbox with provided name
  function destroySandbox(name) {
    cy.log('Destroy Sandbox: ' + name);
    select(meep.EXEC_SELECT_SANDBOX, name);
    cy.wait(1000);
    click(meep.EXEC_BTN_DELETE_SANDBOX);
    click(meep.MEEP_DLG_DELETE_SANDBOX, 'Ok');
    cy.wait(20000);
    verifyEnabled(meep.EXEC_BTN_NEW_SANDBOX, true);
    verifyEnabled(meep.EXEC_BTN_DELETE_SANDBOX, false);
    verifyEnabled(meep.EXEC_BTN_DEPLOY, false);
    verifyEnabled(meep.EXEC_BTN_TERMINATE, false);
    verifyEnabled(meep.EXEC_BTN_EVENT, false);
    verify(meep.MEEP_LBL_SCENARIO_NAME, 'contain', defaultScenario);
  }

  // Deploy scenario with provided name
  function deployScenario(name) {
    // Go to execution page
    cy.log('Go to execution page');
    click(meep.MEEP_TAB_EXEC);
    cy.wait(5000);
    verify(meep.MEEP_LBL_SCENARIO_NAME, 'contain', defaultScenario);
    verifyEnabled(meep.EXEC_BTN_DEPLOY, true);
    verifyEnabled(meep.EXEC_BTN_TERMINATE, false);
    verifyEnabled(meep.EXEC_BTN_EVENT, false);

    // Deploy scenario
    cy.log('Deploy scenario: ' + name);
    click(meep.EXEC_BTN_DEPLOY);
    cy.wait(1000);
    select(meep.MEEP_DLG_DEPLOY_SCENARIO_SELECT, name);
    click(meep.MEEP_DLG_DEPLOY_SCENARIO, 'Ok');
    cy.wait(15000);
    verifyEnabled(meep.EXEC_BTN_EVENT, true, 30000);
    verifyEnabled(meep.EXEC_BTN_DEPLOY, false);
    verifyEnabled(meep.EXEC_BTN_TERMINATE, true);
    verify(meep.MEEP_LBL_SCENARIO_NAME, 'contain', name);
  }

  // Terminate deployed scenario
  function terminateScenario(name) {
    cy.log('Terminate Scenario: ' + name);
    click(meep.EXEC_BTN_TERMINATE);
    click(meep.MEEP_DLG_TERMINATE_SCENARIO, 'Ok');
    cy.wait(10000);
    verifyEnabled(meep.EXEC_BTN_DEPLOY, true, 120000);
    verifyEnabled(meep.EXEC_BTN_TERMINATE, false);
    verifyEnabled(meep.EXEC_BTN_EVENT, false);
    verify(meep.MEEP_LBL_SCENARIO_NAME, 'contain', defaultScenario);
  }

  // Cancel Event creation
  function testCancelEvent() {
    cy.log('Cancel event creation');
    click(meep.EXEC_BTN_EVENT);
    click(meep.EXEC_BTN_EVENT_BTN_MANUAL_REPLAY);
    verifyForm(meep.EXEC_EVT_TYPE, true);
    verifyEnabled(meep.MEEP_BTN_CANCEL, true);
    // verifyEnabled(meep.MEEP_BTN_APPLY, false)
    click(meep.MEEP_BTN_CANCEL);
    cy.wait(1000);
  }

  // Create & Validate Mobility events
  function testMobilityEvent() {
    cy.log('Create Mobility events');
    createMobilityEvent('ue1', 'zone1-poa2');
    createMobilityEvent('ue1', 'zone2-poa1');
    createMobilityEvent('ue1', 'zone1-poa1');
    createMobilityEvent('ue2-ext', 'zone1-poa2');
    createMobilityEvent('ue2-ext', 'zone2-poa1');
    createMobilityEvent('ue2-ext', 'zone1-poa1');
  }

  // Create Network Characteristic events
  function testNetCharEvent(scenario) {
    cy.log('Create & Validate Network Characteristic event');
    createNetCharEvent('SCENARIO', scenario, 60, 5, 1, 200000, 'Pareto');
    createNetCharEvent('OPERATOR', 'operator1', 10, 3, 2, 90000, '');
    createNetCharEvent('ZONE', 'zone1', 6, 2, 1, 70000, '');
    createNetCharEvent('ZONE', 'zone2', 6, 2, 1, 70000, '');
    createNetCharEvent('POA', 'zone1-poa1', 2, 3, 4, 10000, '');
    createNetCharEvent('POA', 'zone1-poa2', 40, 5, 2, 20000, '');
    createNetCharEvent('POA', 'zone2-poa1', 0, 0, 1, 15000, '');
  }

  // Create a Mobility event
  function createMobilityEvent(elem, dest) {
    cy.log('Moving ' + elem + ' --> ' + dest);
    click(meep.EXEC_BTN_EVENT);
    click(meep.EXEC_BTN_EVENT_BTN_MANUAL_REPLAY);
    select(meep.EXEC_EVT_TYPE, meep.MOBILITY_EVENT);
    select(meep.EXEC_EVT_MOB_TARGET, elem);
    select(meep.EXEC_EVT_MOB_DEST, dest);
    click(meep.MEEP_BTN_APPLY);

    // Validate event
    cy.wait(1000);
    validateMobilityEvent(elem, dest);
  }

  // Create a Network Characteristic event
  function createNetCharEvent(elemType, name, l, lv, pl, tp, dist) {
    cy.log('Setting Net Char for type[' + elemType + '] name[' + name + '] latency[' + l +
      '] variation[' + lv + '] packetLoss[' + pl + '] throughput[' + tp + ']');
    click(meep.EXEC_BTN_EVENT);
    click(meep.EXEC_BTN_EVENT_BTN_MANUAL_REPLAY);
    select(meep.EXEC_EVT_TYPE, meep.NETWORK_CHARACTERISTICS_EVENT);
    select(meep.EXEC_EVT_NC_TYPE, elemType);
    select(meep.EXEC_EVT_NC_NAME, name);
    cy.wait(1000);
    type(meep.CFG_ELEM_LATENCY, l);
    type(meep.CFG_ELEM_LATENCY_VAR, lv);
    type(meep.CFG_ELEM_PKT_LOSS, pl);
    type(meep.CFG_ELEM_THROUGHPUT_DL, tp);
    type(meep.CFG_ELEM_THROUGHPUT_UL, tp-1);
    if (dist != '') {
      select(meep.CFG_ELEM_LATENCY_DIST, dist);
    }
    click(meep.MEEP_BTN_APPLY);

    // Validate event
    cy.wait(1000);
    validateNetCharEvent(elemType, name, l, lv, pl, tp, tp-1, dist);
  }

  // Save Replay Events
  function testSaveReplayEvents(name, desc) {
    cy.log('Save events');
    click(meep.EXEC_BTN_EVENT);
    click(meep.EXEC_BTN_EVENT_BTN_SAVE_REPLAY);
    type(meep.MEEP_DLG_SAVE_REPLAY_NAME, name);
    type(meep.MEEP_DLG_SAVE_REPLAY_DESCRIPTION, desc);
    click(meep.MEEP_DLG_SAVE_REPLAY, 'Ok');
    cy.wait(1000);
    //no way to check except by executing auto-replay test that gets the file
    click(meep.EXEC_BTN_EVENT_BTN_AUTO_REPLAY);
    select(meep.EXEC_BTN_EVENT_BTN_AUTO_REPLAY_EVT_REPLAY_FILES, name);
    cy.wait(500);
    //check that by loading the saved replay event file that the description gets populated to confirm it was stored properly
    //and check play button is enabled
    cy.contains(desc)
    verifyEnabled(meep.EXEC_BTN_EVENT_BTN_AUTO_REPLAY_BTN_REPLAY_START, true);
  }

  // Auto Replay
  function testAutoReplay(name) {
    cy.log('Auto Replay');

    testAutoReplayExecution(name, false)
    cy.wait(1000);
    //confirm execution started
    verifyEnabled(meep.EXEC_BTN_EVENT_BTN_AUTO_REPLAY_BTN_REPLAY_START, false);
    cy.wait(75000);
    //confirm it executed completely
    verifyEnabled(meep.EXEC_BTN_EVENT_BTN_AUTO_REPLAY_BTN_REPLAY_START, true);

    testAutoReplayExecution(name, true)
    cy.wait(1000);
    //confirm execution started
    verifyEnabled(meep.EXEC_BTN_EVENT_BTN_AUTO_REPLAY_BTN_REPLAY_START, false);
    cy.wait(75000);
    //confirm it executed completely but restarted
    verifyEnabled(meep.EXEC_BTN_EVENT_BTN_AUTO_REPLAY_BTN_REPLAY_START, false);

    click(meep.EXEC_BTN_EVENT_BTN_AUTO_REPLAY_BTN_REPLAY_STOP)
    cy.wait(500);
    verifyEnabled(meep.EXEC_BTN_EVENT_BTN_AUTO_REPLAY_BTN_REPLAY_START, true);
  }

  function testAutoReplayExecution(name, loop) {
    click(meep.EXEC_BTN_EVENT);
    click(meep.EXEC_BTN_EVENT_BTN_AUTO_REPLAY);
    if (loop) {
      check(meep.EXEC_BTN_EVENT_BTN_AUTO_REPLAY_CHKBOX_LOOP, true);
    } else {
      check(meep.EXEC_BTN_EVENT_BTN_AUTO_REPLAY_CHKBOX_LOOP, false);
    }
    select(meep.EXEC_BTN_EVENT_BTN_AUTO_REPLAY_EVT_REPLAY_FILES, name);
    cy.wait(500);
    verifyEnabled(meep.EXEC_BTN_EVENT_BTN_AUTO_REPLAY_BTN_REPLAY_START, true);
    click(meep.EXEC_BTN_EVENT_BTN_AUTO_REPLAY_BTN_REPLAY_START)
    //check right number of events will be executed
    cy.contains("EVENT COUNT");
    cy.contains("/ 13");
  }

  // Event pane display
  function testEventPane() {
    cy.log('Event Pane show/hide');

    click(meep.EXEC_BTN_EVENT);
    cy.wait(500);
    cy.contains("Event");
    click(meep.EXEC_BTN_EVENT_BTN_CLOSE);
    cy.wait(500);
    cy.contains("Event").should('not.visible');
  }

  // Dashboard pane display
  function testDashboardPane() {
    cy.log('Dashboard Pane show/hide');

    click(meep.EXEC_BTN_DASHBOARD);
    cy.wait(500);
    cy.contains("Dashboard");
    click(meep.EXEC_BTN_DASHBOARD_BTN_CLOSE);
    cy.wait(500);
    cy.contains("Dashboard").should('not.visible');
  }

  // Panes mix display
  function testPanes() {
    cy.log('Mix Panes show/hide');

    click(meep.EXEC_BTN_EVENT);
    click(meep.EXEC_BTN_DASHBOARD);
    cy.wait(500);
    cy.contains("Event");
    cy.contains("Dashboard");
    click(meep.EXEC_BTN_EVENT_BTN_CLOSE);
    click(meep.EXEC_BTN_DASHBOARD_BTN_CLOSE);
    cy.wait(500);
    cy.contains("Event").should('not.visible');
    cy.contains("Dashboard").should('not.visible');
  }

  // Event automation
  function testEventAutomation() {
    cy.log('Event Automation');

    click(meep.EXEC_BTN_EVENT);
    click(meep.EXEC_BTN_EVENT_BTN_AUTOMATION);
    check(meep.EXEC_BTN_EVENT_BTN_AUTOMATION_CHKBOX_MOVEMENT, true);
    check(meep.EXEC_BTN_EVENT_BTN_AUTOMATION_CHKBOX_MOBILITY, true);
    check(meep.EXEC_BTN_EVENT_BTN_AUTOMATION_CHKBOX_POAS_IN_RANGE, true);
    check(meep.EXEC_BTN_EVENT_BTN_AUTOMATION_CHKBOX_NETCHAR, true);

    cy.wait(2000);

    verifyChecked(meep.EXEC_BTN_EVENT_BTN_AUTOMATION_CHKBOX_MOVEMENT, true);
    verifyChecked(meep.EXEC_BTN_EVENT_BTN_AUTOMATION_CHKBOX_MOBILITY, true);
    verifyChecked(meep.EXEC_BTN_EVENT_BTN_AUTOMATION_CHKBOX_POAS_IN_RANGE, true);
    verifyChecked(meep.EXEC_BTN_EVENT_BTN_AUTOMATION_CHKBOX_NETCHAR, true);

    check(meep.EXEC_BTN_EVENT_BTN_AUTOMATION_CHKBOX_MOVEMENT, false);
    check(meep.EXEC_BTN_EVENT_BTN_AUTOMATION_CHKBOX_MOBILITY, false);
    check(meep.EXEC_BTN_EVENT_BTN_AUTOMATION_CHKBOX_POAS_IN_RANGE, false);
    check(meep.EXEC_BTN_EVENT_BTN_AUTOMATION_CHKBOX_NETCHAR, false);
    cy.wait(2000);

    verifyChecked(meep.EXEC_BTN_EVENT_BTN_AUTOMATION_CHKBOX_MOVEMENT, false);
    verifyChecked(meep.EXEC_BTN_EVENT_BTN_AUTOMATION_CHKBOX_MOBILITY, false);
    verifyChecked(meep.EXEC_BTN_EVENT_BTN_AUTOMATION_CHKBOX_POAS_IN_RANGE, false);
    verifyChecked(meep.EXEC_BTN_EVENT_BTN_AUTOMATION_CHKBOX_NETCHAR, false);
  }

  // Retrieve Element entry from Application table
  function getEntry(entries, name) {
    if (entries) {
      return entries[name] ? entries[name] : null;
    }
    return null;
  }

  // Validate that new UE parent matches destination
  function validateMobilityEvent(elem, dest) {
    cy.window().then((win) => {
      var entry = getEntry(win.meepStore.getState().exec.table.entries, elem);
      assert.isNotNull(entry);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), dest);
    });
  }

  // Validate that network characteristics were correctly applied
  function validateNetCharEvent(elemType, name, l, lv, pl, tpDl, tpUl, dist) {
    cy.window().then((win) => {
      var entry = getEntry(win.meepStore.getState().exec.table.entries, name);
      assert.isNotNull(entry);

      switch (elemType) {
        case 'SCENARIO':
          assert.equal(getElemFieldVal(entry, FIELD_INT_DOM_LATENCY), l);
          assert.equal(getElemFieldVal(entry, FIELD_INT_DOM_LATENCY_VAR), lv);
          assert.equal(getElemFieldVal(entry, FIELD_INT_DOM_PKT_LOSS), pl);
          assert.equal(getElemFieldVal(entry, FIELD_INT_DOM_LATENCY_DIST), dist);
          assert.equal(getElemFieldVal(entry, FIELD_INT_DOM_THROUGHPUT_DL), tpDl);
          assert.equal(getElemFieldVal(entry, FIELD_INT_DOM_THROUGHPUT_UL), tpUl);

          break;
        case 'OPERATOR':
          assert.equal(getElemFieldVal(entry, FIELD_INT_ZONE_LATENCY), l);
          assert.equal(getElemFieldVal(entry, FIELD_INT_ZONE_LATENCY_VAR), lv);
          assert.equal(getElemFieldVal(entry, FIELD_INT_ZONE_PKT_LOSS), pl);
          assert.equal(getElemFieldVal(entry, FIELD_INT_ZONE_THROUGHPUT_DL), tpDl);
          assert.equal(getElemFieldVal(entry, FIELD_INT_ZONE_THROUGHPUT_UL), tpUl);
          break;
        case 'ZONE':
          assert.equal(getElemFieldVal(entry, FIELD_INTRA_ZONE_LATENCY), l);
          assert.equal(getElemFieldVal(entry, FIELD_INTRA_ZONE_LATENCY_VAR), lv);
          assert.equal(getElemFieldVal(entry, FIELD_INTRA_ZONE_PKT_LOSS), pl);
          assert.equal(getElemFieldVal(entry, FIELD_INTRA_ZONE_THROUGHPUT_DL), tpDl);
          assert.equal(getElemFieldVal(entry, FIELD_INTRA_ZONE_THROUGHPUT_UL), tpUl);
          break;
        case 'POA':
          assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_LATENCY), l);
          assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_LATENCY_VAR), lv);
          assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_PKT_LOSS), pl);
          assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_THROUGHPUT_DL), tpDl);
          assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_THROUGHPUT_UL), tpUl);
          break;
        default:
          assert.isOk(false, 'Unsupported element type');
      }
    });
  }

});


