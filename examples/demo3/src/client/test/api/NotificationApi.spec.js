/*
 * Copyright (c) 2020  InterDigital Communications, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the \"License\");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an \"AS IS\" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * MEC Demo 3 API
 * Demo 3 is an edge application that can be used with AdvantEDGE or ETSI MEC Sandbox to demonstrate MEC011 and MEC021 usage
 *
 * OpenAPI spec version: 0.0.1
 *
 * NOTE: This class is auto generated by the swagger code generator program.
 * https://github.com/swagger-api/swagger-codegen.git
 *
 * Swagger Codegen version: 3.0.29
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
    factory(root.expect, root.MecDemo3Api);
  }
}(this, function(expect, MecDemo3Api) {
  'use strict';

  var instance;

  beforeEach(function() {
    instance = new MecDemo3Api.NotificationApi();
  });

  describe('(package)', function() {
    describe('NotificationApi', function() {
      describe('amsNotificationCallback', function() {
        it('should call amsNotificationCallback successfully', function(done) {
          // TODO: uncomment, update parameter values for amsNotificationCallback call
          /*

          instance.amsNotificationCallback(body, function(error, data, response) {
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
      describe('appTerminationNotificationCallback', function() {
        it('should call appTerminationNotificationCallback successfully', function(done) {
          // TODO: uncomment, update parameter values for appTerminationNotificationCallback call
          /*

          instance.appTerminationNotificationCallback(body, function(error, data, response) {
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
      describe('contextTransferNotificationCallback', function() {
        it('should call contextTransferNotificationCallback successfully', function(done) {
          // TODO: uncomment, update parameter values for contextTransferNotificationCallback call
          /*

          instance.contextTransferNotificationCallback(body, function(error, data, response) {
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
      describe('serviceAvailNotificationCallback', function() {
        it('should call serviceAvailNotificationCallback successfully', function(done) {
          // TODO: uncomment, update parameter values for serviceAvailNotificationCallback call
          /*

          instance.serviceAvailNotificationCallback(body, function(error, data, response) {
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
    });
  });

}));
