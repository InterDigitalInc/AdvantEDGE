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
 * AdvantEDGE MEC Application Support API
 * MEC Application Support Service is AdvantEDGE's implementation of [ETSI MEC ISG MEC011 Application Enablement API](https://www.etsi.org/deliver/etsi_gs/MEC/001_099/011/02.01.01_60/gs_MEC011v020101p.pdf) <p>[Copyright (c) ETSI 2017](https://forge.etsi.org/etsi-forge-copyright-notice.txt) <p>**Micro-service**<br>[meep-app-enablement](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/go-apps/meep-app-enablement/server/app-support) <p>**Type & Usage**<br>Edge Service used by edge applications that want to get information about applications in the network <p>**Note**<br>AdvantEDGE supports a selected subset of Application Support API endpoints (see below).
 *
 * OpenAPI spec version: 2.1.1
 * Contact: AdvantEDGE@InterDigital.com
 *
 * NOTE: This class is auto generated by the swagger code generator program.
 * https://github.com/swagger-api/swagger-codegen.git
 *
 * Swagger Codegen version: 3.0.22
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
    if (!root.AdvantEdgeMecApplicationSupportApi) {
      root.AdvantEdgeMecApplicationSupportApi = {};
    }
    root.AdvantEdgeMecApplicationSupportApi.TimingCapsNtpServers = factory(root.AdvantEdgeMecApplicationSupportApi.ApiClient);
  }
}(this, function(ApiClient) {
  'use strict';

  /**
   * The TimingCapsNtpServers model module.
   * @module model/TimingCapsNtpServers
   * @version 2.1.1
   */

  /**
   * Constructs a new <code>TimingCapsNtpServers</code>.
   * NTP server detail.
   * @alias module:model/TimingCapsNtpServers
   * @class
   * @param ntpServerAddrType {module:model/TimingCapsNtpServers.NtpServerAddrTypeEnum} Address type of NTP server
   * @param ntpServerAddr {String} NTP server address
   * @param minPollingInterval {Number} Minimum poll interval for NTP messages, in seconds as a power of two. Range 3...17
   * @param maxPollingInterval {Number} Maximum poll interval for NTP messages, in seconds as a power of two. Range 3...17
   * @param localPriority {Number} NTP server local priority
   * @param authenticationOption {module:model/TimingCapsNtpServers.AuthenticationOptionEnum} NTP authentication option
   * @param authenticationKeyNum {Number} Authentication key number
   */
  var exports = function(ntpServerAddrType, ntpServerAddr, minPollingInterval, maxPollingInterval, localPriority, authenticationOption, authenticationKeyNum) {
    this.ntpServerAddrType = ntpServerAddrType;
    this.ntpServerAddr = ntpServerAddr;
    this.minPollingInterval = minPollingInterval;
    this.maxPollingInterval = maxPollingInterval;
    this.localPriority = localPriority;
    this.authenticationOption = authenticationOption;
    this.authenticationKeyNum = authenticationKeyNum;
  };

  /**
   * Constructs a <code>TimingCapsNtpServers</code> from a plain JavaScript object, optionally creating a new instance.
   * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
   * @param {Object} data The plain JavaScript object bearing properties of interest.
   * @param {module:model/TimingCapsNtpServers} obj Optional instance to populate.
   * @return {module:model/TimingCapsNtpServers} The populated <code>TimingCapsNtpServers</code> instance.
   */
  exports.constructFromObject = function(data, obj) {
    if (data) {
      obj = obj || new exports();
      if (data.hasOwnProperty('ntpServerAddrType'))
        obj.ntpServerAddrType = ApiClient.convertToType(data['ntpServerAddrType'], 'String');
      if (data.hasOwnProperty('ntpServerAddr'))
        obj.ntpServerAddr = ApiClient.convertToType(data['ntpServerAddr'], 'String');
      if (data.hasOwnProperty('minPollingInterval'))
        obj.minPollingInterval = ApiClient.convertToType(data['minPollingInterval'], 'Number');
      if (data.hasOwnProperty('maxPollingInterval'))
        obj.maxPollingInterval = ApiClient.convertToType(data['maxPollingInterval'], 'Number');
      if (data.hasOwnProperty('localPriority'))
        obj.localPriority = ApiClient.convertToType(data['localPriority'], 'Number');
      if (data.hasOwnProperty('authenticationOption'))
        obj.authenticationOption = ApiClient.convertToType(data['authenticationOption'], 'String');
      if (data.hasOwnProperty('authenticationKeyNum'))
        obj.authenticationKeyNum = ApiClient.convertToType(data['authenticationKeyNum'], 'Number');
    }
    return obj;
  }

  /**
   * Address type of NTP server
   * @member {module:model/TimingCapsNtpServers.NtpServerAddrTypeEnum} ntpServerAddrType
   */
  exports.prototype.ntpServerAddrType = undefined;

  /**
   * NTP server address
   * @member {String} ntpServerAddr
   */
  exports.prototype.ntpServerAddr = undefined;

  /**
   * Minimum poll interval for NTP messages, in seconds as a power of two. Range 3...17
   * @member {Number} minPollingInterval
   */
  exports.prototype.minPollingInterval = undefined;

  /**
   * Maximum poll interval for NTP messages, in seconds as a power of two. Range 3...17
   * @member {Number} maxPollingInterval
   */
  exports.prototype.maxPollingInterval = undefined;

  /**
   * NTP server local priority
   * @member {Number} localPriority
   */
  exports.prototype.localPriority = undefined;

  /**
   * NTP authentication option
   * @member {module:model/TimingCapsNtpServers.AuthenticationOptionEnum} authenticationOption
   */
  exports.prototype.authenticationOption = undefined;

  /**
   * Authentication key number
   * @member {Number} authenticationKeyNum
   */
  exports.prototype.authenticationKeyNum = undefined;


  /**
   * Allowed values for the <code>ntpServerAddrType</code> property.
   * @enum {String}
   * @readonly
   */
  exports.NtpServerAddrTypeEnum = {
    /**
     * value: "IP_ADDRESS"
     * @const
     */
    IP_ADDRESS: "IP_ADDRESS",

    /**
     * value: "DNS_NAME"
     * @const
     */
    DNS_NAME: "DNS_NAME"
  };


  /**
   * Allowed values for the <code>authenticationOption</code> property.
   * @enum {String}
   * @readonly
   */
  exports.AuthenticationOptionEnum = {
    /**
     * value: "NONE"
     * @const
     */
    NONE: "NONE",

    /**
     * value: "SYMMETRIC_KEY"
     * @const
     */
    SYMMETRIC_KEY: "SYMMETRIC_KEY",

    /**
     * value: "AUTO_KEY"
     * @const
     */
    AUTO_KEY: "AUTO_KEY"
  };

  return exports;

}));