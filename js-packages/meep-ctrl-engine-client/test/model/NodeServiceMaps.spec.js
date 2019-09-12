/**
 * MEEP Controller REST API
 * Copyright (c) 2019 InterDigital Communications, Inc. All rights reserved. The information provided herein is the proprietary and confidential information of InterDigital Communications, Inc. 
 *
 * OpenAPI spec version: 1.0.0
 *
 * NOTE: This class is auto generated by the swagger code generator program.
 * https://github.com/swagger-api/swagger-codegen.git
 *
 * Swagger Codegen version: 2.3.1
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
    factory(root.expect, root.MeepControllerRestApi);
  }
}(this, function(expect, MeepControllerRestApi) {
  'use strict';

  var instance;

  beforeEach(function() {
    instance = new MeepControllerRestApi.NodeServiceMaps();
  });

  var getProperty = function(object, getter, property) {
    // Use getter method if present; otherwise, get the property directly.
    if (typeof object[getter] === 'function')
      return object[getter]();
    else
      return object[property];
  }

  var setProperty = function(object, setter, property, value) {
    // Use setter method if present; otherwise, set the property directly.
    if (typeof object[setter] === 'function')
      object[setter](value);
    else
      object[property] = value;
  }

  describe('NodeServiceMaps', function() {
    it('should create an instance of NodeServiceMaps', function() {
      // uncomment below and update the code to test NodeServiceMaps
      //var instane = new MeepControllerRestApi.NodeServiceMaps();
      //expect(instance).to.be.a(MeepControllerRestApi.NodeServiceMaps);
    });

    it('should have the property node (base name: "node")', function() {
      // uncomment below and update the code to test the property node
      //var instane = new MeepControllerRestApi.NodeServiceMaps();
      //expect(instance).to.be();
    });

    it('should have the property ingressServiceMap (base name: "ingressServiceMap")', function() {
      // uncomment below and update the code to test the property ingressServiceMap
      //var instane = new MeepControllerRestApi.NodeServiceMaps();
      //expect(instance).to.be();
    });

    it('should have the property egressServiceMap (base name: "egressServiceMap")', function() {
      // uncomment below and update the code to test the property egressServiceMap
      //var instane = new MeepControllerRestApi.NodeServiceMaps();
      //expect(instance).to.be();
    });

  });

}));