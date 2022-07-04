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

(function(factory) {
  if (typeof define === 'function' && define.amd) {
    // AMD. Register as an anonymous module.
    define(['ApiClient', 'model/AppTerminationSubscription', 'model/AppTerminationSubscriptionLinks', 'model/ApplicationContextState', 'model/ApplicationInstance', 'model/ApplicationInstanceAmsLinkListSubscription', 'model/ApplicationInstanceAppTerminationSubscription', 'model/ApplicationInstanceDiscoveredServices', 'model/ApplicationInstanceOfferedService', 'model/ApplicationInstanceSerAvailabilitySubscription', 'model/ApplicationInstanceSubscriptions', 'model/AssociateId', 'model/CommunicationInterface', 'model/LinkType', 'model/LocalityType', 'model/MobilityProcedureNotification', 'model/MobilityProcedureNotificationTargetAppInfo', 'model/SerInstanceId', 'model/SerName', 'model/ServiceAvailabilityNotification', 'model/ServiceAvailabilityNotificationServiceReferences', 'model/ServiceState', 'model/Subscription', 'model/TimeStamp', 'api/FrontendApi', 'api/NotificationApi'], factory);
  } else if (typeof module === 'object' && module.exports) {
    // CommonJS-like environments that support module.exports, like Node.
    module.exports = factory(require('./ApiClient'), require('./model/AppTerminationSubscription'), require('./model/AppTerminationSubscriptionLinks'), require('./model/ApplicationContextState'), require('./model/ApplicationInstance'), require('./model/ApplicationInstanceAmsLinkListSubscription'), require('./model/ApplicationInstanceAppTerminationSubscription'), require('./model/ApplicationInstanceDiscoveredServices'), require('./model/ApplicationInstanceOfferedService'), require('./model/ApplicationInstanceSerAvailabilitySubscription'), require('./model/ApplicationInstanceSubscriptions'), require('./model/AssociateId'), require('./model/CommunicationInterface'), require('./model/LinkType'), require('./model/LocalityType'), require('./model/MobilityProcedureNotification'), require('./model/MobilityProcedureNotificationTargetAppInfo'), require('./model/SerInstanceId'), require('./model/SerName'), require('./model/ServiceAvailabilityNotification'), require('./model/ServiceAvailabilityNotificationServiceReferences'), require('./model/ServiceState'), require('./model/Subscription'), require('./model/TimeStamp'), require('./api/FrontendApi'), require('./api/NotificationApi'));
  }
}(function(ApiClient, AppTerminationSubscription, AppTerminationSubscriptionLinks, ApplicationContextState, ApplicationInstance, ApplicationInstanceAmsLinkListSubscription, ApplicationInstanceAppTerminationSubscription, ApplicationInstanceDiscoveredServices, ApplicationInstanceOfferedService, ApplicationInstanceSerAvailabilitySubscription, ApplicationInstanceSubscriptions, AssociateId, CommunicationInterface, LinkType, LocalityType, MobilityProcedureNotification, MobilityProcedureNotificationTargetAppInfo, SerInstanceId, SerName, ServiceAvailabilityNotification, ServiceAvailabilityNotificationServiceReferences, ServiceState, Subscription, TimeStamp, FrontendApi, NotificationApi) {
  'use strict';

  /**
   * Demo_3_is_an_edge_application_that_can_be_used_with_AdvantEDGE_or_ETSI_MEC_Sandbox_to_demonstrate_MEC011_and_MEC021_usage.<br>
   * The <code>index</code> module provides access to constructors for all the classes which comprise the public API.
   * <p>
   * An AMD (recommended!) or CommonJS application will generally do something equivalent to the following:
   * <pre>
   * var MecDemo3Api = require('index'); // See note below*.
   * var xxxSvc = new MecDemo3Api.XxxApi(); // Allocate the API class we're going to use.
   * var yyyModel = new MecDemo3Api.Yyy(); // Construct a model instance.
   * yyyModel.someProperty = 'someValue';
   * ...
   * var zzz = xxxSvc.doSomething(yyyModel); // Invoke the service.
   * ...
   * </pre>
   * <em>*NOTE: For a top-level AMD script, use require(['index'], function(){...})
   * and put the application logic within the callback function.</em>
   * </p>
   * <p>
   * A non-AMD browser application (discouraged) might do something like this:
   * <pre>
   * var xxxSvc = new MecDemo3Api.XxxApi(); // Allocate the API class we're going to use.
   * var yyy = new MecDemo3Api.Yyy(); // Construct a model instance.
   * yyyModel.someProperty = 'someValue';
   * ...
   * var zzz = xxxSvc.doSomething(yyyModel); // Invoke the service.
   * ...
   * </pre>
   * </p>
   * @module index
   * @version 0.0.1
   */
  var exports = {
    /**
     * The ApiClient constructor.
     * @property {module:ApiClient}
     */
    ApiClient: ApiClient,
    /**
     * The AppTerminationSubscription model constructor.
     * @property {module:model/AppTerminationSubscription}
     */
    AppTerminationSubscription: AppTerminationSubscription,
    /**
     * The AppTerminationSubscriptionLinks model constructor.
     * @property {module:model/AppTerminationSubscriptionLinks}
     */
    AppTerminationSubscriptionLinks: AppTerminationSubscriptionLinks,
    /**
     * The ApplicationContextState model constructor.
     * @property {module:model/ApplicationContextState}
     */
    ApplicationContextState: ApplicationContextState,
    /**
     * The ApplicationInstance model constructor.
     * @property {module:model/ApplicationInstance}
     */
    ApplicationInstance: ApplicationInstance,
    /**
     * The ApplicationInstanceAmsLinkListSubscription model constructor.
     * @property {module:model/ApplicationInstanceAmsLinkListSubscription}
     */
    ApplicationInstanceAmsLinkListSubscription: ApplicationInstanceAmsLinkListSubscription,
    /**
     * The ApplicationInstanceAppTerminationSubscription model constructor.
     * @property {module:model/ApplicationInstanceAppTerminationSubscription}
     */
    ApplicationInstanceAppTerminationSubscription: ApplicationInstanceAppTerminationSubscription,
    /**
     * The ApplicationInstanceDiscoveredServices model constructor.
     * @property {module:model/ApplicationInstanceDiscoveredServices}
     */
    ApplicationInstanceDiscoveredServices: ApplicationInstanceDiscoveredServices,
    /**
     * The ApplicationInstanceOfferedService model constructor.
     * @property {module:model/ApplicationInstanceOfferedService}
     */
    ApplicationInstanceOfferedService: ApplicationInstanceOfferedService,
    /**
     * The ApplicationInstanceSerAvailabilitySubscription model constructor.
     * @property {module:model/ApplicationInstanceSerAvailabilitySubscription}
     */
    ApplicationInstanceSerAvailabilitySubscription: ApplicationInstanceSerAvailabilitySubscription,
    /**
     * The ApplicationInstanceSubscriptions model constructor.
     * @property {module:model/ApplicationInstanceSubscriptions}
     */
    ApplicationInstanceSubscriptions: ApplicationInstanceSubscriptions,
    /**
     * The AssociateId model constructor.
     * @property {module:model/AssociateId}
     */
    AssociateId: AssociateId,
    /**
     * The CommunicationInterface model constructor.
     * @property {module:model/CommunicationInterface}
     */
    CommunicationInterface: CommunicationInterface,
    /**
     * The LinkType model constructor.
     * @property {module:model/LinkType}
     */
    LinkType: LinkType,
    /**
     * The LocalityType model constructor.
     * @property {module:model/LocalityType}
     */
    LocalityType: LocalityType,
    /**
     * The MobilityProcedureNotification model constructor.
     * @property {module:model/MobilityProcedureNotification}
     */
    MobilityProcedureNotification: MobilityProcedureNotification,
    /**
     * The MobilityProcedureNotificationTargetAppInfo model constructor.
     * @property {module:model/MobilityProcedureNotificationTargetAppInfo}
     */
    MobilityProcedureNotificationTargetAppInfo: MobilityProcedureNotificationTargetAppInfo,
    /**
     * The SerInstanceId model constructor.
     * @property {module:model/SerInstanceId}
     */
    SerInstanceId: SerInstanceId,
    /**
     * The SerName model constructor.
     * @property {module:model/SerName}
     */
    SerName: SerName,
    /**
     * The ServiceAvailabilityNotification model constructor.
     * @property {module:model/ServiceAvailabilityNotification}
     */
    ServiceAvailabilityNotification: ServiceAvailabilityNotification,
    /**
     * The ServiceAvailabilityNotificationServiceReferences model constructor.
     * @property {module:model/ServiceAvailabilityNotificationServiceReferences}
     */
    ServiceAvailabilityNotificationServiceReferences: ServiceAvailabilityNotificationServiceReferences,
    /**
     * The ServiceState model constructor.
     * @property {module:model/ServiceState}
     */
    ServiceState: ServiceState,
    /**
     * The Subscription model constructor.
     * @property {module:model/Subscription}
     */
    Subscription: Subscription,
    /**
     * The TimeStamp model constructor.
     * @property {module:model/TimeStamp}
     */
    TimeStamp: TimeStamp,
    /**
     * The FrontendApi service constructor.
     * @property {module:api/FrontendApi}
     */
    FrontendApi: FrontendApi,
    /**
     * The NotificationApi service constructor.
     * @property {module:api/NotificationApi}
     */
    NotificationApi: NotificationApi
  };

  return exports;
}));