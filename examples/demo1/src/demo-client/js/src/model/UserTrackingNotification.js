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
    // AMD. Register as an anonymous module.
    define(['ApiClient', 'model/TimeStamp', 'model/UserEventType', 'model/UserInfo'], factory);
  } else if (typeof module === 'object' && module.exports) {
    // CommonJS-like environments that support module.exports, like Node.
    module.exports = factory(require('../ApiClient'), require('./TimeStamp'), require('./UserEventType'), require('./UserInfo'));
  } else {
    // Browser globals (root is window)
    if (!root.MeepDemoAppApi) {
      root.MeepDemoAppApi = {};
    }
    root.MeepDemoAppApi.UserTrackingNotification = factory(root.MeepDemoAppApi.ApiClient, root.MeepDemoAppApi.TimeStamp, root.MeepDemoAppApi.UserEventType, root.MeepDemoAppApi.UserInfo);
  }
}(this, function(ApiClient, TimeStamp, UserEventType, UserInfo) {
  'use strict';

  /**
   * The UserTrackingNotification model module.
   * @module model/UserTrackingNotification
   * @version 0.0.1
   */

  /**
   * Constructs a new <code>UserTrackingNotification</code>.
   * User tracking notification - callback generated toward an ME app with a user tracking subscription
   * @alias module:model/UserTrackingNotification
   * @class
   * @param callbackData {String} CallBackData if passed by the application during the associated Subscription (Zone or User Tracking) operation
   * @param userInfo {module:model/UserInfo} 
   * @param timeStamp {module:model/TimeStamp} 
   */
  var exports = function(callbackData, userInfo, timeStamp) {
    this.callbackData = callbackData;
    this.userInfo = userInfo;
    this.timeStamp = timeStamp;
  };

  /**
   * Constructs a <code>UserTrackingNotification</code> from a plain JavaScript object, optionally creating a new instance.
   * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
   * @param {Object} data The plain JavaScript object bearing properties of interest.
   * @param {module:model/UserTrackingNotification} obj Optional instance to populate.
   * @return {module:model/UserTrackingNotification} The populated <code>UserTrackingNotification</code> instance.
   */
  exports.constructFromObject = function(data, obj) {
    if (data) {
      obj = obj || new exports();
      if (data.hasOwnProperty('callbackData'))
        obj.callbackData = ApiClient.convertToType(data['callbackData'], 'String');
      if (data.hasOwnProperty('userInfo'))
        obj.userInfo = UserInfo.constructFromObject(data['userInfo']);
      if (data.hasOwnProperty('timeStamp'))
        obj.timeStamp = TimeStamp.constructFromObject(data['timeStamp']);
      if (data.hasOwnProperty('userEventType'))
        obj.userEventType = UserEventType.constructFromObject(data['userEventType']);
    }
    return obj;
  }

  /**
   * CallBackData if passed by the application during the associated Subscription (Zone or User Tracking) operation
   * @member {String} callbackData
   */
  exports.prototype.callbackData = undefined;

  /**
   * @member {module:model/UserInfo} userInfo
   */
  exports.prototype.userInfo = undefined;

  /**
   * @member {module:model/TimeStamp} timeStamp
   */
  exports.prototype.timeStamp = undefined;

  /**
   * @member {module:model/UserEventType} userEventType
   */
  exports.prototype.userEventType = undefined;

  return exports;

}));
