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
    define(['ApiClient', 'model/ActivationInfo', 'model/Domains', 'model/NetworkLocations', 'model/NodeServiceMaps', 'model/PhysicalLocations', 'model/Processes', 'model/Scenario', 'model/Zones'], factory);
  } else if (typeof module === 'object' && module.exports) {
    // CommonJS-like environments that support module.exports, like Node.
    module.exports = factory(require('../ApiClient'), require('../model/ActivationInfo'), require('../model/Domains'), require('../model/NetworkLocations'), require('../model/NodeServiceMaps'), require('../model/PhysicalLocations'), require('../model/Processes'), require('../model/Scenario'), require('../model/Zones'));
  } else {
    // Browser globals (root is window)
    if (!root.AdvantEdgeSandboxControllerRestApi) {
      root.AdvantEdgeSandboxControllerRestApi = {};
    }
    root.AdvantEdgeSandboxControllerRestApi.ActiveScenarioApi = factory(root.AdvantEdgeSandboxControllerRestApi.ApiClient, root.AdvantEdgeSandboxControllerRestApi.ActivationInfo, root.AdvantEdgeSandboxControllerRestApi.Domains, root.AdvantEdgeSandboxControllerRestApi.NetworkLocations, root.AdvantEdgeSandboxControllerRestApi.NodeServiceMaps, root.AdvantEdgeSandboxControllerRestApi.PhysicalLocations, root.AdvantEdgeSandboxControllerRestApi.Processes, root.AdvantEdgeSandboxControllerRestApi.Scenario, root.AdvantEdgeSandboxControllerRestApi.Zones);
  }
}(this, function(ApiClient, ActivationInfo, Domains, NetworkLocations, NodeServiceMaps, PhysicalLocations, Processes, Scenario, Zones) {
  'use strict';

  /**
   * ActiveScenario service.
   * @module api/ActiveScenarioApi
   * @version 1.0.0
   */

  /**
   * Constructs a new ActiveScenarioApi. 
   * @alias module:api/ActiveScenarioApi
   * @class
   * @param {module:ApiClient} [apiClient] Optional API client implementation to use,
   * default to {@link module:ApiClient#instance} if unspecified.
   */
  var exports = function(apiClient) {
    this.apiClient = apiClient || ApiClient.instance;


    /**
     * Callback function to receive the result of the activateScenario operation.
     * @callback module:api/ActiveScenarioApi~activateScenarioCallback
     * @param {String} error Error message, if any.
     * @param data This operation does not return a value.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Deploy a scenario
     * Deploy a scenario present in the platform scenario store
     * @param {String} name Scenario name
     * @param {Object} opts Optional parameters
     * @param {module:model/ActivationInfo} opts.activationInfo Activation information
     * @param {module:api/ActiveScenarioApi~activateScenarioCallback} callback The callback function, accepting three arguments: error, data, response
     */
    this.activateScenario = function(name, opts, callback) {
      opts = opts || {};
      var postBody = opts['activationInfo'];

      // verify the required parameter 'name' is set
      if (name === undefined || name === null) {
        throw new Error("Missing the required parameter 'name' when calling activateScenario");
      }


      var pathParams = {
        'name': name
      };
      var queryParams = {
      };
      var collectionQueryParams = {
      };
      var headerParams = {
      };
      var formParams = {
      };

      var authNames = [];
      var contentTypes = ['application/json'];
      var accepts = ['application/json'];
      var returnType = null;

      return this.apiClient.callApi(
        '/active/{name}', 'POST',
        pathParams, queryParams, collectionQueryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, callback
      );
    }

    /**
     * Callback function to receive the result of the getActiveNodeServiceMaps operation.
     * @callback module:api/ActiveScenarioApi~getActiveNodeServiceMapsCallback
     * @param {String} error Error message, if any.
     * @param {Array.<module:model/NodeServiceMaps>} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Get deployed scenario's port mapping
     * Returns the deployed scenario's port mapping<p> <li>Ports are used by external nodes to access services internal to the platform <li>Port mapping concept for external nodes is available [here](https://github.com/InterDigitalInc/AdvantEDGE/wiki/external-ue#port-mapping)
     * @param {Object} opts Optional parameters
     * @param {String} opts.node Unique node identifier
     * @param {String} opts.type Exposed service type (ingress or egress)
     * @param {String} opts.service Exposed service name
     * @param {module:api/ActiveScenarioApi~getActiveNodeServiceMapsCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link Array.<module:model/NodeServiceMaps>}
     */
    this.getActiveNodeServiceMaps = function(opts, callback) {
      opts = opts || {};
      var postBody = null;


      var pathParams = {
      };
      var queryParams = {
        'node': opts['node'],
        'type': opts['type'],
        'service': opts['service'],
      };
      var collectionQueryParams = {
      };
      var headerParams = {
      };
      var formParams = {
      };

      var authNames = [];
      var contentTypes = ['application/json'];
      var accepts = ['application/json'];
      var returnType = [NodeServiceMaps];

      return this.apiClient.callApi(
        '/active/serviceMaps', 'GET',
        pathParams, queryParams, collectionQueryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, callback
      );
    }

    /**
     * Callback function to receive the result of the getActiveScenario operation.
     * @callback module:api/ActiveScenarioApi~getActiveScenarioCallback
     * @param {String} error Error message, if any.
     * @param {module:model/Scenario} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Get the deployed scenario
     * Get the scenario currently deployed on the platform
     * @param {Object} opts Optional parameters
     * @param {Boolean} opts.minimize Return minimized scenario element content
     * @param {module:api/ActiveScenarioApi~getActiveScenarioCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link module:model/Scenario}
     */
    this.getActiveScenario = function(opts, callback) {
      opts = opts || {};
      var postBody = null;


      var pathParams = {
      };
      var queryParams = {
        'minimize': opts['minimize'],
      };
      var collectionQueryParams = {
      };
      var headerParams = {
      };
      var formParams = {
      };

      var authNames = [];
      var contentTypes = ['application/json'];
      var accepts = ['application/json'];
      var returnType = Scenario;

      return this.apiClient.callApi(
        '/active', 'GET',
        pathParams, queryParams, collectionQueryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, callback
      );
    }

    /**
     * Callback function to receive the result of the getActiveScenarioDomain operation.
     * @callback module:api/ActiveScenarioApi~getActiveScenarioDomainCallback
     * @param {String} error Error message, if any.
     * @param {module:model/Domains} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Get domain elements from the deployed scenario
     * Returns a filtered list of domain elements from the deployed scenario using the provided query parameters
     * @param {Object} opts Optional parameters
     * @param {String} opts.domain Domain name
     * @param {module:model/String} opts.domainType Domain type
     * @param {String} opts.zone Zone name
     * @param {String} opts.networkLocation Network Location name
     * @param {module:model/String} opts.networkLocationType Network Location type
     * @param {String} opts.physicalLocation Physical Location name
     * @param {module:model/String} opts.physicalLocationType Physical Location type
     * @param {String} opts.process Process name
     * @param {module:model/String} opts.processType Process type
     * @param {Boolean} opts.excludeChildren Include child elements in response
     * @param {Boolean} opts.minimize Return minimized scenario element content
     * @param {module:api/ActiveScenarioApi~getActiveScenarioDomainCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link module:model/Domains}
     */
    this.getActiveScenarioDomain = function(opts, callback) {
      opts = opts || {};
      var postBody = null;


      var pathParams = {
      };
      var queryParams = {
        'domain': opts['domain'],
        'domainType': opts['domainType'],
        'zone': opts['zone'],
        'networkLocation': opts['networkLocation'],
        'networkLocationType': opts['networkLocationType'],
        'physicalLocation': opts['physicalLocation'],
        'physicalLocationType': opts['physicalLocationType'],
        'process': opts['process'],
        'processType': opts['processType'],
        'excludeChildren': opts['excludeChildren'],
        'minimize': opts['minimize'],
      };
      var collectionQueryParams = {
      };
      var headerParams = {
      };
      var formParams = {
      };

      var authNames = [];
      var contentTypes = ['application/json'];
      var accepts = ['application/json'];
      var returnType = Domains;

      return this.apiClient.callApi(
        '/active/domains', 'GET',
        pathParams, queryParams, collectionQueryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, callback
      );
    }

    /**
     * Callback function to receive the result of the getActiveScenarioNetworkLocation operation.
     * @callback module:api/ActiveScenarioApi~getActiveScenarioNetworkLocationCallback
     * @param {String} error Error message, if any.
     * @param {module:model/NetworkLocations} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Get network location elements from the deployed scenario
     * Returns a filtered list of network location elements from the deployed scenario using the provided query parameters
     * @param {Object} opts Optional parameters
     * @param {String} opts.domain Domain name
     * @param {module:model/String} opts.domainType Domain type
     * @param {String} opts.zone Zone name
     * @param {String} opts.networkLocation Network Location name
     * @param {module:model/String} opts.networkLocationType Network Location type
     * @param {String} opts.physicalLocation Physical Location name
     * @param {module:model/String} opts.physicalLocationType Physical Location type
     * @param {String} opts.process Process name
     * @param {module:model/String} opts.processType Process type
     * @param {Boolean} opts.excludeChildren Include child elements in response
     * @param {Boolean} opts.minimize Return minimized scenario element content
     * @param {module:api/ActiveScenarioApi~getActiveScenarioNetworkLocationCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link module:model/NetworkLocations}
     */
    this.getActiveScenarioNetworkLocation = function(opts, callback) {
      opts = opts || {};
      var postBody = null;


      var pathParams = {
      };
      var queryParams = {
        'domain': opts['domain'],
        'domainType': opts['domainType'],
        'zone': opts['zone'],
        'networkLocation': opts['networkLocation'],
        'networkLocationType': opts['networkLocationType'],
        'physicalLocation': opts['physicalLocation'],
        'physicalLocationType': opts['physicalLocationType'],
        'process': opts['process'],
        'processType': opts['processType'],
        'excludeChildren': opts['excludeChildren'],
        'minimize': opts['minimize'],
      };
      var collectionQueryParams = {
      };
      var headerParams = {
      };
      var formParams = {
      };

      var authNames = [];
      var contentTypes = ['application/json'];
      var accepts = ['application/json'];
      var returnType = NetworkLocations;

      return this.apiClient.callApi(
        '/active/networkLocations', 'GET',
        pathParams, queryParams, collectionQueryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, callback
      );
    }

    /**
     * Callback function to receive the result of the getActiveScenarioPhysicalLocation operation.
     * @callback module:api/ActiveScenarioApi~getActiveScenarioPhysicalLocationCallback
     * @param {String} error Error message, if any.
     * @param {module:model/PhysicalLocations} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Get physical location elements from the deployed scenario
     * Returns a filtered list of physical location elements from the deployed scenario using the provided query parameters
     * @param {Object} opts Optional parameters
     * @param {String} opts.domain Domain name
     * @param {module:model/String} opts.domainType Domain type
     * @param {String} opts.zone Zone name
     * @param {String} opts.networkLocation Network Location name
     * @param {module:model/String} opts.networkLocationType Network Location type
     * @param {String} opts.physicalLocation Physical Location name
     * @param {module:model/String} opts.physicalLocationType Physical Location type
     * @param {String} opts.process Process name
     * @param {module:model/String} opts.processType Process type
     * @param {Boolean} opts.excludeChildren Include child elements in response
     * @param {Boolean} opts.minimize Return minimized scenario element content
     * @param {module:api/ActiveScenarioApi~getActiveScenarioPhysicalLocationCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link module:model/PhysicalLocations}
     */
    this.getActiveScenarioPhysicalLocation = function(opts, callback) {
      opts = opts || {};
      var postBody = null;


      var pathParams = {
      };
      var queryParams = {
        'domain': opts['domain'],
        'domainType': opts['domainType'],
        'zone': opts['zone'],
        'networkLocation': opts['networkLocation'],
        'networkLocationType': opts['networkLocationType'],
        'physicalLocation': opts['physicalLocation'],
        'physicalLocationType': opts['physicalLocationType'],
        'process': opts['process'],
        'processType': opts['processType'],
        'excludeChildren': opts['excludeChildren'],
        'minimize': opts['minimize'],
      };
      var collectionQueryParams = {
      };
      var headerParams = {
      };
      var formParams = {
      };

      var authNames = [];
      var contentTypes = ['application/json'];
      var accepts = ['application/json'];
      var returnType = PhysicalLocations;

      return this.apiClient.callApi(
        '/active/physicalLocations', 'GET',
        pathParams, queryParams, collectionQueryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, callback
      );
    }

    /**
     * Callback function to receive the result of the getActiveScenarioProcess operation.
     * @callback module:api/ActiveScenarioApi~getActiveScenarioProcessCallback
     * @param {String} error Error message, if any.
     * @param {module:model/Processes} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Get process elements from the deployed scenario
     * Returns a filtered list of process elements from the deployed scenario using the provided query parameters
     * @param {Object} opts Optional parameters
     * @param {String} opts.domain Domain name
     * @param {module:model/String} opts.domainType Domain type
     * @param {String} opts.zone Zone name
     * @param {String} opts.networkLocation Network Location name
     * @param {module:model/String} opts.networkLocationType Network Location type
     * @param {String} opts.physicalLocation Physical Location name
     * @param {module:model/String} opts.physicalLocationType Physical Location type
     * @param {String} opts.process Process name
     * @param {module:model/String} opts.processType Process type
     * @param {Boolean} opts.excludeChildren Include child elements in response
     * @param {Boolean} opts.minimize Return minimized scenario element content
     * @param {module:api/ActiveScenarioApi~getActiveScenarioProcessCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link module:model/Processes}
     */
    this.getActiveScenarioProcess = function(opts, callback) {
      opts = opts || {};
      var postBody = null;


      var pathParams = {
      };
      var queryParams = {
        'domain': opts['domain'],
        'domainType': opts['domainType'],
        'zone': opts['zone'],
        'networkLocation': opts['networkLocation'],
        'networkLocationType': opts['networkLocationType'],
        'physicalLocation': opts['physicalLocation'],
        'physicalLocationType': opts['physicalLocationType'],
        'process': opts['process'],
        'processType': opts['processType'],
        'excludeChildren': opts['excludeChildren'],
        'minimize': opts['minimize'],
      };
      var collectionQueryParams = {
      };
      var headerParams = {
      };
      var formParams = {
      };

      var authNames = [];
      var contentTypes = ['application/json'];
      var accepts = ['application/json'];
      var returnType = Processes;

      return this.apiClient.callApi(
        '/active/processes', 'GET',
        pathParams, queryParams, collectionQueryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, callback
      );
    }

    /**
     * Callback function to receive the result of the getActiveScenarioZone operation.
     * @callback module:api/ActiveScenarioApi~getActiveScenarioZoneCallback
     * @param {String} error Error message, if any.
     * @param {module:model/Zones} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Get zone elements from the deployed scenario
     * Returns a filtered list of zone elements from the deployed scenario using the provided query parameters
     * @param {Object} opts Optional parameters
     * @param {String} opts.domain Domain name
     * @param {module:model/String} opts.domainType Domain type
     * @param {String} opts.zone Zone name
     * @param {String} opts.networkLocation Network Location name
     * @param {module:model/String} opts.networkLocationType Network Location type
     * @param {String} opts.physicalLocation Physical Location name
     * @param {module:model/String} opts.physicalLocationType Physical Location type
     * @param {String} opts.process Process name
     * @param {module:model/String} opts.processType Process type
     * @param {Boolean} opts.excludeChildren Include child elements in response
     * @param {Boolean} opts.minimize Return minimized scenario element content
     * @param {module:api/ActiveScenarioApi~getActiveScenarioZoneCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link module:model/Zones}
     */
    this.getActiveScenarioZone = function(opts, callback) {
      opts = opts || {};
      var postBody = null;


      var pathParams = {
      };
      var queryParams = {
        'domain': opts['domain'],
        'domainType': opts['domainType'],
        'zone': opts['zone'],
        'networkLocation': opts['networkLocation'],
        'networkLocationType': opts['networkLocationType'],
        'physicalLocation': opts['physicalLocation'],
        'physicalLocationType': opts['physicalLocationType'],
        'process': opts['process'],
        'processType': opts['processType'],
        'excludeChildren': opts['excludeChildren'],
        'minimize': opts['minimize'],
      };
      var collectionQueryParams = {
      };
      var headerParams = {
      };
      var formParams = {
      };

      var authNames = [];
      var contentTypes = ['application/json'];
      var accepts = ['application/json'];
      var returnType = Zones;

      return this.apiClient.callApi(
        '/active/zones', 'GET',
        pathParams, queryParams, collectionQueryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, callback
      );
    }

    /**
     * Callback function to receive the result of the terminateScenario operation.
     * @callback module:api/ActiveScenarioApi~terminateScenarioCallback
     * @param {String} error Error message, if any.
     * @param data This operation does not return a value.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Terminate the deployed scenario
     * Terminate the scenario currently deployed on the platform
     * @param {module:api/ActiveScenarioApi~terminateScenarioCallback} callback The callback function, accepting three arguments: error, data, response
     */
    this.terminateScenario = function(callback) {
      var postBody = null;


      var pathParams = {
      };
      var queryParams = {
      };
      var collectionQueryParams = {
      };
      var headerParams = {
      };
      var formParams = {
      };

      var authNames = [];
      var contentTypes = ['application/json'];
      var accepts = ['application/json'];
      var returnType = null;

      return this.apiClient.callApi(
        '/active', 'DELETE',
        pathParams, queryParams, collectionQueryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, callback
      );
    }
  };

  return exports;
}));
