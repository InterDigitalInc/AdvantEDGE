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
import { selector, click, select, verify } from '../util/util';

// Scenario Execution Tests
describe('Scenario Monitoring', function() {

  // Test Setup
  beforeEach(() => {
    var meepUrl = Cypress.env('meep_url');
    if (meepUrl == null) {
      meepUrl = 'http://127.0.0.1:30000';
    }

    cy.viewport(1920, 1080);
    cy.visit(meepUrl);
  });

  it('Monitor Scenario', function() {
    // let latencyDashboard = 'Latency Dashboard';
    // let demoSvcIntUeDashboard = 'Demo Service Internal UE (ue1)';
    // let demoSvcExtUeDashboard = 'Demo Service External UE (ue2-ext)';
    let metricsDashboard = 'Metrics Dashboard';

    // Go to monitoring page
    cy.log('Go to monitoring page');
    click(meep.MEEP_TAB_MON);

    // Verify available dashboards
    cy.log('Verify available dashboards');
    // verify(meep.MON_DASHBOARD_SELECT, 'contain', latencyDashboard);
    // verify(meep.MON_DASHBOARD_SELECT, 'contain', demoSvcIntUeDashboard);
    // verify(meep.MON_DASHBOARD_SELECT, 'contain', demoSvcExtUeDashboard);
    verify(meep.MON_DASHBOARD_SELECT, 'contain', metricsDashboard);

    // // Open Latency Dashboard
    // select(meep.MON_DASHBOARD_SELECT, latencyDashboard);
    // verifyIframe(meep.MON_DASHBOARD_IFRAME, 'have.attr', 'src');

    // // Open Demo Service Internal UE Dashboard
    // select(meep.MON_DASHBOARD_SELECT, demoSvcIntUeDashboard);
    // verifyIframe(meep.MON_DASHBOARD_IFRAME, 'have.attr', 'src');

    // // Open Demo Service External UE Dashboard
    // select(meep.MON_DASHBOARD_SELECT, demoSvcExtUeDashboard);
    // verifyIframe(meep.MON_DASHBOARD_IFRAME, 'have.attr', 'src');

    // Open Metrics Dashboard
    select(meep.MON_DASHBOARD_SELECT, metricsDashboard);
    verifyIframe(meep.MON_DASHBOARD_IFRAME, 'have.attr', 'src');
  });

  function verifyIframe(name, options, value) {
    cy.get(selector(name)).children().first().should(options, value);
  }

});


