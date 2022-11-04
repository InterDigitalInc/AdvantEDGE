/*
 * Copyright (c) 2022  The AdvantEDGE Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * AdvantEDGE Sandbox Controller REST API
 * This API is the main Sandbox Controller API for scenario deployment & event injection <p>**Micro-service**<br>[meep-sandbox-ctrl](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/go-apps/meep-sandbox-ctrl) <p>**Type & Usage**<br>Platform runtime interface to manage active scenarios and inject events in AdvantEDGE platform <p>**Details**<br>API details available at _your-AdvantEDGE-ip-address/api_
 *
 * OpenAPI spec version: 1.0.0
 * Contact: AdvantEDGE@InterDigital.com
 *
 * NOTE: This class is auto generated by the swagger code generator program.
 * https://github.com/swagger-api/swagger-codegen.git
 *
 * Swagger Codegen version: 2.4.9
 *
 * Do not edit the class manually.
 *
 */

(function(root, factory) {
  if (typeof define === 'function' && define.amd) {
    // AMD.
    define(['expect.js', '../../src/index'], factory);
  } else if (typeof module === 'object' && module.exports) {
    // CommonJS-like environments that support module.exports, like Node.
    factory(require('expect.js'), require('../../src/index'));
  } else {
    // Browser globals (root is window)
    factory(root.expect, root.AdvantEdgeSandboxControllerRestApi);
  }
}(this, function(expect, AdvantEdgeSandboxControllerRestApi) {
  'use strict';

  var instance;

  beforeEach(function() {
    instance = new AdvantEdgeSandboxControllerRestApi.ApplicationsApi();
  });

  describe('(package)', function() {
    describe('ApplicationsApi', function() {
      describe('applicationsAppInstanceIdDELETE', function() {
        it('should call applicationsAppInstanceIdDELETE successfully', function(done) {
          // TODO: uncomment, update parameter values for applicationsAppInstanceIdDELETE call
          /*
          var appInstanceId = "appInstanceId_example";

          instance.applicationsAppInstanceIdDELETE(appInstanceId, function(error, data, response) {
            if (error) {
              done(error);
              return;
            }

            done();
          });
          */
          // TODO: uncomment and complete method invocation above, then delete this line and the next:
          done();
        });
      });
      describe('applicationsAppInstanceIdGET', function() {
        it('should call applicationsAppInstanceIdGET successfully', function(done) {
          // TODO: uncomment, update parameter values for applicationsAppInstanceIdGET call and complete the assertions
          /*
          var appInstanceId = "appInstanceId_example";

          instance.applicationsAppInstanceIdGET(appInstanceId, function(error, data, response) {
            if (error) {
              done(error);
              return;
            }
            // TODO: update response assertions
            expect(data).to.be.a(AdvantEdgeSandboxControllerRestApi.ApplicationInfo);
            expect(data.id).to.be.a('string');
            expect(data.id).to.be("");
            expect(data.name).to.be.a('string');
            expect(data.name).to.be("");
            expect(data.nodeName).to.be.a('string');
            expect(data.nodeName).to.be("");
            expect(data.type).to.be.a('string');
            expect(data.type).to.be("USER");
            expect(data.persist).to.be.a('boolean');
            expect(data.persist).to.be(false);

            done();
          });
          */
          // TODO: uncomment and complete method invocation above, then delete this line and the next:
          done();
        });
      });
      describe('applicationsAppInstanceIdPUT', function() {
        it('should call applicationsAppInstanceIdPUT successfully', function(done) {
          // TODO: uncomment, update parameter values for applicationsAppInstanceIdPUT call and complete the assertions
          /*
          var appInstanceId = "appInstanceId_example";
          var applicationInfo = new AdvantEdgeSandboxControllerRestApi.ApplicationInfo();
          applicationInfo.id = "";
          applicationInfo.name = "";
          applicationInfo.nodeName = "";
          applicationInfo.type = "USER";
          applicationInfo.persist = false;

          instance.applicationsAppInstanceIdPUT(appInstanceId, applicationInfo, function(error, data, response) {
            if (error) {
              done(error);
              return;
            }
            // TODO: update response assertions
            expect(data).to.be.a(AdvantEdgeSandboxControllerRestApi.ApplicationInfo);
            expect(data.id).to.be.a('string');
            expect(data.id).to.be("");
            expect(data.name).to.be.a('string');
            expect(data.name).to.be("");
            expect(data.nodeName).to.be.a('string');
            expect(data.nodeName).to.be("");
            expect(data.type).to.be.a('string');
            expect(data.type).to.be("USER");
            expect(data.persist).to.be.a('boolean');
            expect(data.persist).to.be(false);

            done();
          });
          */
          // TODO: uncomment and complete method invocation above, then delete this line and the next:
          done();
        });
      });
      describe('applicationsGET', function() {
        it('should call applicationsGET successfully', function(done) {
          // TODO: uncomment, update parameter values for applicationsGET call and complete the assertions
          /*
          var opts = {};
          opts.app = "app_example";
          opts.type = "type_example";
          opts.nodeName = "nodeName_example";

          instance.applicationsGET(opts, function(error, data, response) {
            if (error) {
              done(error);
              return;
            }
            // TODO: update response assertions
            let dataCtr = data;
            expect(dataCtr).to.be.an(Array);
            expect(dataCtr).to.not.be.empty();
            for (let p in dataCtr) {
              let data = dataCtr[p];
              expect(data).to.be.a(AdvantEdgeSandboxControllerRestApi.ApplicationInfo);
              expect(data.id).to.be.a('string');
              expect(data.id).to.be("");
              expect(data.name).to.be.a('string');
              expect(data.name).to.be("");
              expect(data.nodeName).to.be.a('string');
              expect(data.nodeName).to.be("");
              expect(data.type).to.be.a('string');
              expect(data.type).to.be("USER");
              expect(data.persist).to.be.a('boolean');
              expect(data.persist).to.be(false);
            }

            done();
          });
          */
          // TODO: uncomment and complete method invocation above, then delete this line and the next:
          done();
        });
      });
      describe('applicationsPOST', function() {
        it('should call applicationsPOST successfully', function(done) {
          // TODO: uncomment, update parameter values for applicationsPOST call and complete the assertions
          /*
          var applicationInfo = new AdvantEdgeSandboxControllerRestApi.ApplicationInfo();
          applicationInfo.id = "";
          applicationInfo.name = "";
          applicationInfo.nodeName = "";
          applicationInfo.type = "USER";
          applicationInfo.persist = false;

          instance.applicationsPOST(applicationInfo, function(error, data, response) {
            if (error) {
              done(error);
              return;
            }
            // TODO: update response assertions
            expect(data).to.be.a(AdvantEdgeSandboxControllerRestApi.ApplicationInfo);
            expect(data.id).to.be.a('string');
            expect(data.id).to.be("");
            expect(data.name).to.be.a('string');
            expect(data.name).to.be("");
            expect(data.nodeName).to.be.a('string');
            expect(data.nodeName).to.be("");
            expect(data.type).to.be.a('string');
            expect(data.type).to.be("USER");
            expect(data.persist).to.be.a('boolean');
            expect(data.persist).to.be(false);

            done();
          });
          */
          // TODO: uncomment and complete method invocation above, then delete this line and the next:
          done();
        });
      });
    });
  });

}));
