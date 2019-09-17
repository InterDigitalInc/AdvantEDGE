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
import * as states from '../../../../js-apps/meep-frontend/src/js/state/ui/index';

// Import element utils
import {
  // Field Names
  FIELD_TYPE,
  FIELD_PARENT,
  FIELD_NAME,
  FIELD_IMAGE,
  FIELD_PORT,
  FIELD_PROTOCOL,
  FIELD_GROUP,
  FIELD_SVC_MAP,
  FIELD_ENV_VAR,
  FIELD_CMD,
  FIELD_CMD_ARGS,
  FIELD_EXT_PORT,
  FIELD_IS_EXTERNAL,
  FIELD_CHART_ENABLED,
  FIELD_CHART_LOC,
  FIELD_CHART_VAL,
  FIELD_CHART_GROUP,
  FIELD_INT_DOM_LATENCY,
  FIELD_INT_DOM_LATENCY_VAR,
  FIELD_INT_DOM_THROUGPUT,
  FIELD_INT_DOM_PKT_LOSS,
  FIELD_INT_ZONE_LATENCY,
  FIELD_INT_ZONE_LATENCY_VAR,
  FIELD_INT_ZONE_THROUGPUT,
  FIELD_INT_ZONE_PKT_LOSS,
  FIELD_INT_EDGE_LATENCY,
  FIELD_INT_EDGE_LATENCY_VAR,
  FIELD_INT_EDGE_THROUGPUT,
  FIELD_INT_EDGE_PKT_LOSS,
  FIELD_INT_FOG_LATENCY,
  FIELD_INT_FOG_LATENCY_VAR,
  FIELD_INT_FOG_THROUGPUT,
  FIELD_INT_FOG_PKT_LOSS,
  FIELD_EDGE_FOG_LATENCY,
  FIELD_EDGE_FOG_LATENCY_VAR,
  FIELD_EDGE_FOG_THROUGPUT,
  FIELD_EDGE_FOG_PKT_LOSS,
  FIELD_LINK_LATENCY,
  FIELD_LINK_LATENCY_VAR,
  FIELD_LINK_THROUGPUT,
  FIELD_LINK_PKT_LOSS,

  getElemFieldVal,
} from '../../../../js-apps/meep-frontend/src/js/util/elem-utils';

// Import Test utility functions
import { selector, click, type, select, verify, verifyEnabled, verifyForm } from '../util/util';

// Scenario Execution Tests
describe('Scenario Execution', function() {

  // Test Variables
  let defaultScenario = 'None';

  // ------------------------------
  //            TESTS
  // ------------------------------

  // It is not possible to connect to multiple domains in one test case, but it is possible in different test cases
  // Known cypress bug : https://github.com/cypress-io/cypress/issues/944
  // The suggested solution from cypress is not working, and will likely not be fixed
  // Best workaround: using different tests

  // Demo1 scenario testing for mobility event on the demo app output
  it('Deploy DEMO1 scenario - mobility event #1', function() {

    let scenario = 'demo1';

    openDefaultMeepUrl("");

    // Deploy demo scenario
    deployScenario(scenario);

    // Test events
    //testMobilityEvent();
    //testMobility event was split so that each mobility event can be tested independantly
    cy.log('Create 1st Mobility event');
    createMobilityEvent('ue1', 'zone1-poa2');

  });

  it('DEMO1 scenario - validate mobility event in demo-app frontend', function() {

    cy.wait(2000);
    openDefaultMeepUrl("31111");

    // Already deployed scenario

    // Test events
    validateDemoAppFrontEndUpdateOnMobilityEvent('ue1', 'zone1-poa2');

  });

  // Demo1 scenario testing  mobility event #2
  it('DEMO1 scenario - mobility event #2', function() {

    openDefaultMeepUrl("");
    click(meep.MEEP_TAB_EXEC);

    // Test events
    cy.log('Create 2nd Mobility event');
    createMobilityEvent('ue1', 'zone2-poa1');

  });

  it('DEMO1 scenario - validate mobility event in demo-app frontend', function() {

    cy.wait(2000);
    openDefaultMeepUrl("31111");

    // Test events
    validateDemoAppFrontEndUpdateOnMobilityEvent('ue1', 'zone2-poa1');

  });

  // Demo1 scenario testing  mobility event #3
  it('DEMO1 scenario - mobility event #3', function() {

    openDefaultMeepUrl("");
    click(meep.MEEP_TAB_EXEC);

    // Test events
    cy.log('Create 3rd Mobility event');
    createMobilityEvent('ue1', 'zone1-poa1');

  });

  it('DEMO1 scenario - validate mobility event in demo-app frontend', function() {

    cy.wait(2000);
    openDefaultMeepUrl("31111");

    // Test events
    validateDemoAppFrontEndUpdateOnMobilityEvent('ue1', 'zone1-poa1');

  });

  // Demo1 scenario testing  mobility event #4
  it('DEMO1 scenario - mobility event #4', function() {

    openDefaultMeepUrl("");
    click(meep.MEEP_TAB_EXEC);

    // Test events
    cy.log('Create 4th Mobility event');
    createMobilityEvent('ue2-ext', 'zone1-poa2');

  });

  it('DEMO1 scenario - validate mobility event in demo-app frontend', function() {

    cy.wait(2000);
    openDefaultMeepUrl("31111");

    // Test events
    validateDemoAppFrontEndUpdateOnMobilityEvent('ue2-ext', 'zone1-poa2');

  });

  // Demo1 scenario testing  mobility event #5
  it('DEMO1 scenario - mobility event #5', function() {

    openDefaultMeepUrl("");
    click(meep.MEEP_TAB_EXEC);

    // Test events
    cy.log('Create 5th Mobility event');
    createMobilityEvent('ue2-ext', 'zone2-poa1');

  });

  it('DEMO1 scenario - validate mobility event in demo-app frontend', function() {

    cy.wait(2000);
    openDefaultMeepUrl("31111");

    // Test events
    validateDemoAppFrontEndUpdateOnMobilityEvent('ue2-ext', 'zone2-poa1');

  });

  // Demo1 scenario testing  mobility event #6
  it('DEMO1 scenario - mobility event #6', function() {

    openDefaultMeepUrl("");
    click(meep.MEEP_TAB_EXEC);

    // Test events
    cy.log('Create 6th Mobility event');
    createMobilityEvent('ue2-ext', 'zone1-poa1');

  });

  it('DEMO1 scenario - validate mobility event in demo-app frontend', function() {

    cy.wait(2000);
    openDefaultMeepUrl("31111");

    // Test events
    validateDemoAppFrontEndUpdateOnMobilityEvent('ue2-ext', 'zone1-poa1');

  });

  // Demo1 scenario terminate
  it('DEMO1 scenario - terminate', function() {

    let scenario = 'demo1';

    openDefaultMeepUrl("");
    click(meep.MEEP_TAB_EXEC);

    // Test events
    // Terminate demo scenario
    terminateScenario(scenario);

  });

  // ------------------------------
  //          FUNCTIONS
  // ------------------------------

  // Opens the default meepUrl page
  function openDefaultMeepUrl(port) {
    var meepUrl = Cypress.env('meep_url');
    if (meepUrl == null) {
      meepUrl = 'http://127.0.0.1:30000';
    }
    if (port != "") {
      meepUrl = meepUrl.replace("30000", port);
    }
 
    cy.viewport(1920, 1080);
    cy.visit(meepUrl);

  }

  // Deploy scenario with provided name
  function deployScenario(name) {
    // Go to execution page
    cy.log('Go to execution page');
    click(meep.MEEP_TAB_EXEC);
    verify(meep.MEEP_LBL_SCENARIO_NAME, 'contain', defaultScenario);
    verifyEnabled(meep.EXEC_BTN_DEPLOY, true);
    verifyEnabled(meep.EXEC_BTN_TERMINATE, false);
    verifyEnabled(meep.EXEC_BTN_REFRESH, false);
    verifyEnabled(meep.EXEC_BTN_EVENT, false);

    // Deploy scenario
    cy.log('Deploy scenario: ' + name);
    click(meep.EXEC_BTN_DEPLOY);
    cy.wait(1000);
    select(meep.MEEP_DLG_DEPLOY_SCENARIO_SELECT, name);
    click(meep.MEEP_DLG_DEPLOY_SCENARIO, 'Ok');
    cy.wait(1000);
    verifyEnabled(meep.EXEC_BTN_EVENT, true, 30000);
    verifyEnabled(meep.EXEC_BTN_DEPLOY, false);
    verifyEnabled(meep.EXEC_BTN_TERMINATE, true);
    verifyEnabled(meep.EXEC_BTN_REFRESH, true);
    verify(meep.MEEP_LBL_SCENARIO_NAME, 'contain', name);
  }

  // Terminate deployed scenario
  function terminateScenario(name) {
    cy.log('Terminate Scenario: ' + name);
    click(meep.EXEC_BTN_TERMINATE);
    click(meep.MEEP_DLG_TERMINATE_SCENARIO, 'Ok');
    cy.wait(1000);
    verifyEnabled(meep.EXEC_BTN_DEPLOY, true, 120000);
    verifyEnabled(meep.EXEC_BTN_TERMINATE, false);
    verifyEnabled(meep.EXEC_BTN_REFRESH, false);
    verifyEnabled(meep.EXEC_BTN_EVENT, false);
    verify(meep.MEEP_LBL_SCENARIO_NAME, 'contain', defaultScenario);
  }

  // Cancel Event creation
  function testCancelEvent() {
    cy.log('Cancel event creation');
    click(meep.EXEC_BTN_EVENT);
    verifyForm(meep.EXEC_EVT_TYPE, true);
    verifyEnabled(meep.MEEP_BTN_CANCEL, true);
    // verifyEnabled(meep.MEEP_BTN_APPLY, false)
    click(meep.MEEP_BTN_CANCEL);
    cy.wait(1000);
  }

  // Create Network Characteristic events
  function testNetCharEvent(scenario) {
    cy.log('Create & Validate Network Characteristic event');
    createNetCharEvent('SCENARIO', scenario, 60, 5, 1, 200000);
    createNetCharEvent('DOMAIN', 'operator1', 10, 3, 2, 90000);
    createNetCharEvent('ZONE-INTER-EDGE', 'zone1', 5, 0, 1, 80000);
    createNetCharEvent('ZONE-INTER-FOG', 'zone1', 3, 2, 1, 75000);
    createNetCharEvent('ZONE-EDGE-FOG', 'zone1', 6, 2, 1, 70000);
    createNetCharEvent('ZONE-INTER-EDGE', 'zone2', 5, 0, 1, 80000);
    createNetCharEvent('ZONE-INTER-FOG', 'zone2', 3, 2, 1, 75000);
    createNetCharEvent('ZONE-EDGE-FOG', 'zone2', 6, 2, 1, 70000);
    createNetCharEvent('POA', 'zone1-poa1', 2, 3, 4, 10000);
    createNetCharEvent('POA', 'zone1-poa2', 40, 5, 2, 20000);
    createNetCharEvent('POA', 'zone2-poa1', 0, 0, 1, 15000);
  }

  // Create a Mobility event
  function createMobilityEvent(ue, dest) {
    cy.log('Moving ' + ue + ' --> ' + dest);
    click(meep.EXEC_BTN_EVENT);
    select(meep.EXEC_EVT_TYPE, states.UE_MOBILITY_EVENT);
    select(meep.EXEC_EVT_MOB_TARGET, ue);
    select(meep.EXEC_EVT_MOB_DEST, dest);
    click(meep.MEEP_BTN_APPLY);

    // Validate event
    cy.wait(1000);
    validateMobilityEvent(ue, dest);
  }

  // Create a Network Characteristic event
  function createNetCharEvent(elemType, name, l, lv, pl, tp) {
    cy.log('Setting Net Char for type[' + elemType + '] name[' + name + '] latency[' + l +
      '] variation[' + lv + '] packetLoss[' + pl + '] throughput[' + tp + ']');
    click(meep.EXEC_BTN_EVENT);
    select(meep.EXEC_EVT_TYPE, states.NETWORK_CHARACTERISTICS_EVENT);
    select(meep.EXEC_EVT_NC_TYPE, elemType);
    select(meep.EXEC_EVT_NC_NAME, name);
    cy.wait(1000);
    type(meep.CFG_ELEM_LATENCY, l);
    type(meep.CFG_ELEM_LATENCY_VAR, lv);
    type(meep.CFG_ELEM_PKT_LOSS, pl);
    type(meep.CFG_ELEM_THROUGHPUT, tp);
    click(meep.MEEP_BTN_APPLY);

    // Validate event
    cy.wait(1000);
    validateNetCharEvent(elemType, name, l, lv, pl, tp);
  }

  // Retrieve Element entry from Application table
  function getEntry(entries, name) {
    if (entries) {
      for (var i = 0; i < entries.length; i++) {
        if (getElemFieldVal(entries[i], FIELD_NAME) == name) {
          return entries[i];
        }
      }
    }
    return null;
  }

  // Validate that new UE parent matches destination
  function validateMobilityEvent(ue, dest) {
    cy.window().then((win) => {
      var entry = getEntry(win.meepStore.getState().exec.table.entries, ue);
      assert.isNotNull(entry);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), dest);
    });
  }

  // Validate that network characteristics were correctly applied
  function validateNetCharEvent(elemType, name, l, lv, pl, tp) {
    cy.window().then((win) => {
      var entry = getEntry(win.meepStore.getState().exec.table.entries, name);
      assert.isNotNull(entry);

      switch (elemType) {
      case 'SCENARIO':
        assert.equal(getElemFieldVal(entry, FIELD_INT_DOM_LATENCY), l);
        assert.equal(getElemFieldVal(entry, FIELD_INT_DOM_LATENCY_VAR), lv);
        assert.equal(getElemFieldVal(entry, FIELD_INT_DOM_PKT_LOSS), pl);
        assert.equal(getElemFieldVal(entry, FIELD_INT_DOM_THROUGPUT), tp);
        break;
      case 'DOMAIN':
        assert.equal(getElemFieldVal(entry, FIELD_INT_ZONE_LATENCY), l);
        assert.equal(getElemFieldVal(entry, FIELD_INT_ZONE_LATENCY_VAR), lv);
        assert.equal(getElemFieldVal(entry, FIELD_INT_ZONE_PKT_LOSS), pl);
        assert.equal(getElemFieldVal(entry, FIELD_INT_ZONE_THROUGPUT), tp);
        break;
      case 'ZONE-INTER-EDGE':
        assert.equal(getElemFieldVal(entry, FIELD_INT_EDGE_LATENCY), l);
        assert.equal(getElemFieldVal(entry, FIELD_INT_EDGE_LATENCY_VAR), lv);
        assert.equal(getElemFieldVal(entry, FIELD_INT_EDGE_PKT_LOSS), pl);
        assert.equal(getElemFieldVal(entry, FIELD_INT_EDGE_THROUGPUT), tp);
        break;
      case 'ZONE-INTER-FOG':
        assert.equal(getElemFieldVal(entry, FIELD_INT_FOG_LATENCY), l);
        assert.equal(getElemFieldVal(entry, FIELD_INT_FOG_LATENCY_VAR), lv);
        assert.equal(getElemFieldVal(entry, FIELD_INT_FOG_PKT_LOSS), pl);
        assert.equal(getElemFieldVal(entry, FIELD_INT_FOG_THROUGPUT), tp);
        break;
      case 'ZONE-EDGE-FOG':
        assert.equal(getElemFieldVal(entry, FIELD_EDGE_FOG_LATENCY), l);
        assert.equal(getElemFieldVal(entry, FIELD_EDGE_FOG_LATENCY_VAR), lv);
        assert.equal(getElemFieldVal(entry, FIELD_EDGE_FOG_PKT_LOSS), pl);
        assert.equal(getElemFieldVal(entry, FIELD_EDGE_FOG_THROUGPUT), tp);
        break;
      case 'POA':
        assert.equal(getElemFieldVal(entry, FIELD_LINK_LATENCY), l);
        assert.equal(getElemFieldVal(entry, FIELD_LINK_LATENCY_VAR), lv);
        assert.equal(getElemFieldVal(entry, FIELD_LINK_PKT_LOSS), pl);
        assert.equal(getElemFieldVal(entry, FIELD_LINK_THROUGPUT), tp);
        break;
      default:
        assert.isOk(false, 'Unsupported element type');
      }
    });
  }

  // Validate that demo app front-end update matches the mobility event that occured
  function validateDemoAppFrontEndUpdateOnMobilityEvent(ue, dest) {

    var ueAppName;
    var ueAppLocation;
    var ueAddressLabel;
    var ueLocationLabel;

    if (ue == 'ue1') {
      ueAppName = 'ue1-iperf';
      ueAddressLabel = 'span[id="demo-svc-loc-serv-address-1"]';
      ueLocationLabel = 'span[id="demo-svc-loc-serv-location-1"]';
    } else {
      if (ue == 'ue2-ext') {
        ueAppName = 'ue2-svc';
        ueAddressLabel = 'span[id="demo-svc-loc-serv-address-2"]';
        ueLocationLabel = 'span[id="demo-svc-loc-serv-location-2"]';
      }
    }
    if (dest == 'zone1-poa1') {
      ueAppLocation = 'zone1 / zone1-poa1';
    } else {
      if (dest == 'zone1-poa2') {
        ueAppLocation = 'zone1 / zone1-poa2';
      } else {
        if (dest == 'zone2-poa1') {
          ueAppLocation = 'zone2 / zone2-poa1';
        }
      }
    }
    cy.contains(ueLocationLabel, ueAppLocation)
  }

});


