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

// Import Test utility functions
import { selector, click, check, type, select, verify, verifyEnabled } from '../util/util';

// Scenario Execution Tests
describe('MEEP Settings', function() {

  // Test Setup
  beforeEach(() => {
    var meepUrl = Cypress.env('meep_url');
    if (meepUrl == null) {
      meepUrl = 'http://127.0.0.1:30000';
    }

    cy.viewport(1920, 1080);
    cy.visit(meepUrl);
  });

  it('Execution Settings', function() {
    let refreshInterval = '10000';

    // Go to monitoring page
    cy.log('Go to settings page');
    click(meep.MEEP_TAB_SET);
    verify(meep.SET_EXEC_REFRESH_CHECKBOX, 'not.be.checked');
    verifyEnabled(meep.SET_EXEC_REFRESH_INT, false);

    // Enable refresh interval
    check(meep.SET_EXEC_REFRESH_CHECKBOX, true);
    verifyEnabled(meep.SET_EXEC_REFRESH_INT, true);
    type(meep.SET_EXEC_REFRESH_INT, refreshInterval);
    verify(meep.SET_EXEC_REFRESH_CHECKBOX, 'be.checked');
    // verify(meep.SET_EXEC_REFRESH_INT, 'contain', refreshInterval)

    // Disable refresh interval
    check(meep.SET_EXEC_REFRESH_CHECKBOX, false);
    verifyEnabled(meep.SET_EXEC_REFRESH_INT, false);
    verify(meep.SET_EXEC_REFRESH_CHECKBOX, 'not.be.checked');
  });

  it('Development Settings', function() {
    // Go to monitoring page
    cy.log('Go to settings page');
    click(meep.MEEP_TAB_SET);
    verify(meep.SET_DEV_MODE_CHECKBOX, 'not.be.checked');

    // Enable dev mode
    check(meep.SET_DEV_MODE_CHECKBOX, true);
    verify(meep.SET_DEV_MODE_CHECKBOX, 'be.checked');

    // Disable dev mode
    check(meep.SET_DEV_MODE_CHECKBOX, false);
    verify(meep.SET_DEV_MODE_CHECKBOX, 'not.be.checked');
  });

});


