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
    // AMD. Register as an anonymous module.
    define(['ApiClient'], factory);
  } else if (typeof module === 'object' && module.exports) {
    // CommonJS-like environments that support module.exports, like Node.
    module.exports = factory(require('../ApiClient'));
  } else {
    // Browser globals (root is window)
    if (!root.AdvantEdgeSandboxControllerRestApi) {
      root.AdvantEdgeSandboxControllerRestApi = {};
    }
    root.AdvantEdgeSandboxControllerRestApi.CellularDomainConfig = factory(root.AdvantEdgeSandboxControllerRestApi.ApiClient);
  }
}(this, function(ApiClient) {
  'use strict';

  /**
   * The CellularDomainConfig model module.
   * @module model/CellularDomainConfig
   * @version 1.0.0
   */

  /**
   * Constructs a new <code>CellularDomainConfig</code>.
   * Cellular domain configuration information
   * @alias module:model/CellularDomainConfig
   * @class
   */
  var exports = function() {
  };

  /**
   * Constructs a <code>CellularDomainConfig</code> from a plain JavaScript object, optionally creating a new instance.
   * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
   * @param {Object} data The plain JavaScript object bearing properties of interest.
   * @param {module:model/CellularDomainConfig} obj Optional instance to populate.
   * @return {module:model/CellularDomainConfig} The populated <code>CellularDomainConfig</code> instance.
   */
  exports.constructFromObject = function(data, obj) {
    if (data) {
      obj = obj || new exports();
      if (data.hasOwnProperty('mnc'))
        obj.mnc = ApiClient.convertToType(data['mnc'], 'String');
      if (data.hasOwnProperty('mcc'))
        obj.mcc = ApiClient.convertToType(data['mcc'], 'String');
      if (data.hasOwnProperty('defaultCellId'))
        obj.defaultCellId = ApiClient.convertToType(data['defaultCellId'], 'String');
    }
    return obj;
  }

  /**
   * Mobile Network Code part of PLMN identity as defined in ETSI TS 136 413
   * @member {String} mnc
   */
  exports.prototype.mnc = undefined;

  /**
   * Mobile Country Code part of PLMN identity as defined in ETSI TS 136 413
   * @member {String} mcc
   */
  exports.prototype.mcc = undefined;

  /**
   * The E-UTRAN Cell Identity as defined in ETSI TS 136 413 if no cellId is defined for the cell or if not applicable
   * @member {String} defaultCellId
   */
  exports.prototype.defaultCellId = undefined;

  return exports;

}));
