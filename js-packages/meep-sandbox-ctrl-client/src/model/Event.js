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
    define(['ApiClient', 'model/EventMobility', 'model/EventNetworkCharacteristicsUpdate', 'model/EventPduSession', 'model/EventPoasInRange', 'model/EventScenarioUpdate'], factory);
  } else if (typeof module === 'object' && module.exports) {
    // CommonJS-like environments that support module.exports, like Node.
    module.exports = factory(require('../ApiClient'), require('./EventMobility'), require('./EventNetworkCharacteristicsUpdate'), require('./EventPduSession'), require('./EventPoasInRange'), require('./EventScenarioUpdate'));
  } else {
    // Browser globals (root is window)
    if (!root.AdvantEdgeSandboxControllerRestApi) {
      root.AdvantEdgeSandboxControllerRestApi = {};
    }
    root.AdvantEdgeSandboxControllerRestApi.Event = factory(root.AdvantEdgeSandboxControllerRestApi.ApiClient, root.AdvantEdgeSandboxControllerRestApi.EventMobility, root.AdvantEdgeSandboxControllerRestApi.EventNetworkCharacteristicsUpdate, root.AdvantEdgeSandboxControllerRestApi.EventPduSession, root.AdvantEdgeSandboxControllerRestApi.EventPoasInRange, root.AdvantEdgeSandboxControllerRestApi.EventScenarioUpdate);
  }
}(this, function(ApiClient, EventMobility, EventNetworkCharacteristicsUpdate, EventPduSession, EventPoasInRange, EventScenarioUpdate) {
  'use strict';

  /**
   * The Event model module.
   * @module model/Event
   * @version 1.0.0
   */

  /**
   * Constructs a new <code>Event</code>.
   * Event object
   * @alias module:model/Event
   * @class
   */
  var exports = function() {
  };

  /**
   * Constructs a <code>Event</code> from a plain JavaScript object, optionally creating a new instance.
   * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
   * @param {Object} data The plain JavaScript object bearing properties of interest.
   * @param {module:model/Event} obj Optional instance to populate.
   * @return {module:model/Event} The populated <code>Event</code> instance.
   */
  exports.constructFromObject = function(data, obj) {
    if (data) {
      obj = obj || new exports();
      if (data.hasOwnProperty('name'))
        obj.name = ApiClient.convertToType(data['name'], 'String');
      if (data.hasOwnProperty('type'))
        obj.type = ApiClient.convertToType(data['type'], 'String');
      if (data.hasOwnProperty('eventMobility'))
        obj.eventMobility = EventMobility.constructFromObject(data['eventMobility']);
      if (data.hasOwnProperty('eventNetworkCharacteristicsUpdate'))
        obj.eventNetworkCharacteristicsUpdate = EventNetworkCharacteristicsUpdate.constructFromObject(data['eventNetworkCharacteristicsUpdate']);
      if (data.hasOwnProperty('eventPoasInRange'))
        obj.eventPoasInRange = EventPoasInRange.constructFromObject(data['eventPoasInRange']);
      if (data.hasOwnProperty('eventScenarioUpdate'))
        obj.eventScenarioUpdate = EventScenarioUpdate.constructFromObject(data['eventScenarioUpdate']);
      if (data.hasOwnProperty('eventPduSession'))
        obj.eventPduSession = EventPduSession.constructFromObject(data['eventPduSession']);
    }
    return obj;
  }

  /**
   * Event name
   * @member {String} name
   */
  exports.prototype.name = undefined;

  /**
   * Event type
   * @member {module:model/Event.TypeEnum} type
   */
  exports.prototype.type = undefined;

  /**
   * @member {module:model/EventMobility} eventMobility
   */
  exports.prototype.eventMobility = undefined;

  /**
   * @member {module:model/EventNetworkCharacteristicsUpdate} eventNetworkCharacteristicsUpdate
   */
  exports.prototype.eventNetworkCharacteristicsUpdate = undefined;

  /**
   * @member {module:model/EventPoasInRange} eventPoasInRange
   */
  exports.prototype.eventPoasInRange = undefined;

  /**
   * @member {module:model/EventScenarioUpdate} eventScenarioUpdate
   */
  exports.prototype.eventScenarioUpdate = undefined;

  /**
   * @member {module:model/EventPduSession} eventPduSession
   */
  exports.prototype.eventPduSession = undefined;


  /**
   * Allowed values for the <code>type</code> property.
   * @enum {String}
   * @readonly
   */
  exports.TypeEnum = {
    /**
     * value: "MOBILITY"
     * @const
     */
    MOBILITY: "MOBILITY",

    /**
     * value: "NETWORK-CHARACTERISTICS-UPDATE"
     * @const
     */
    NETWORK_CHARACTERISTICS_UPDATE: "NETWORK-CHARACTERISTICS-UPDATE",

    /**
     * value: "POAS-IN-RANGE"
     * @const
     */
    POAS_IN_RANGE: "POAS-IN-RANGE",

    /**
     * value: "SCENARIO-UPDATE"
     * @const
     */
    SCENARIO_UPDATE: "SCENARIO-UPDATE",

    /**
     * value: "PDU-SESSION"
     * @const
     */
    PDU_SESSION: "PDU-SESSION"
  };

  return exports;

}));
