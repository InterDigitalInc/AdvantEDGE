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
 * MEEP Demo App API
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * OpenAPI spec version: 0.0.1
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
    factory(root.expect, root.MeepDemoAppApi);
  }
}(this, function(expect, MeepDemoAppApi) {
  'use strict';

  var instance;

  beforeEach(function() {
    instance = new MeepDemoAppApi.EdgeAppInfoApi();
  });

  describe('(package)', function() {
    describe('EdgeAppInfoApi', function() {
      describe('getEdgeInfo', function() {
        it('should call getEdgeInfo successfully', function(done) {
          // TODO: uncomment getEdgeInfo call and complete the assertions
          /*

          instance.getEdgeInfo(function(error, data, response) {
            if (error) {
              done(error);
              return;
            }
            // TODO: update response assertions
            expect(data).to.be.a(MeepDemoAppApi.EdgeInfo);
            expect(data.svc).to.be.a('string');
            expect(data.svc).to.be("");
            expect(data.name).to.be.a('string');
            expect(data.name).to.be("");
            expect(data.ip).to.be.a('string');
            expect(data.ip).to.be("");

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
