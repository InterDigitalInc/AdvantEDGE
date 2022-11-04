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
 * AdvantEDGE Platform Controller REST API
 * This API is the main Platform Controller API for scenario configuration & sandbox management <p>**Micro-service**<br>[meep-pfm-ctrl](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/go-apps/meep-platform-ctrl) <p>**Type & Usage**<br>Platform main interface used by controller software to configure scenarios and manage sandboxes in the AdvantEDGE platform <p>**Details**<br>API details available at _your-AdvantEDGE-ip-address/api_
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
    factory(root.expect, root.AdvantEdgePlatformControllerRestApi);
  }
}(this, function(expect, AdvantEdgePlatformControllerRestApi) {
  'use strict';

  var instance;

  describe('(package)', function() {
    describe('GeoData', function() {
      beforeEach(function() {
        instance = new AdvantEdgePlatformControllerRestApi.GeoData();
      });

      it('should create an instance of GeoData', function() {
        // TODO: update the code to test GeoData
        expect(instance).to.be.a(AdvantEdgePlatformControllerRestApi.GeoData);
      });

      it('should have the property location (base name: "location")', function() {
        // TODO: update the code to test the property location
        expect(instance).to.have.property('location');
        // expect(instance.location).to.be(expectedValueLiteral);
      });

      it('should have the property radius (base name: "radius")', function() {
        // TODO: update the code to test the property radius
        expect(instance).to.have.property('radius');
        // expect(instance.radius).to.be(expectedValueLiteral);
      });

      it('should have the property path (base name: "path")', function() {
        // TODO: update the code to test the property path
        expect(instance).to.have.property('path');
        // expect(instance.path).to.be(expectedValueLiteral);
      });

      it('should have the property eopMode (base name: "eopMode")', function() {
        // TODO: update the code to test the property eopMode
        expect(instance).to.have.property('eopMode');
        // expect(instance.eopMode).to.be(expectedValueLiteral);
      });

      it('should have the property velocity (base name: "velocity")', function() {
        // TODO: update the code to test the property velocity
        expect(instance).to.have.property('velocity');
        // expect(instance.velocity).to.be(expectedValueLiteral);
      });

      it('should have the property d2dInRange (base name: "d2dInRange")', function() {
        // TODO: update the code to test the property d2dInRange
        expect(instance).to.have.property('d2dInRange');
        // expect(instance.d2dInRange).to.be(expectedValueLiteral);
      });

      it('should have the property poaInRange (base name: "poaInRange")', function() {
        // TODO: update the code to test the property poaInRange
        expect(instance).to.have.property('poaInRange');
        // expect(instance.poaInRange).to.be(expectedValueLiteral);
      });

    });
  });

}));
