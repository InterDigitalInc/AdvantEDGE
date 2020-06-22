/*
 * Copyright (c) 2020  InterDigital Communications, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { connect } from 'react-redux';
import React, { Component, createRef } from 'react';
import ReactDOM from 'react-dom';
import L from 'leaflet';
import 'mapbox-gl';
import 'mapbox-gl-leaflet';
import '@geoman-io/leaflet-geoman-free';
import deepEqual from 'deep-equal';
import {
  updateObject,
  deepCopy
} from '../util/object-util';
import { execChangeTable } from '../state/exec';
import {
  cfgChangeTable,
  cfgElemUpdate,
  cfgElemEdit
} from '../state/cfg';
import {
  uiCfgChangeMapCfg,
  uiExecChangeSandboxCfg
} from '../state/ui';
import {
  TYPE_CFG,
  TYPE_EXEC,
  HOST_PATH,
  DEFAULT_MAP_LATITUDE,
  DEFAULT_MAP_LONGITUDE,
  DEFAULT_MAP_ZOOM,
  ELEMENT_TYPE_UE,
  ELEMENT_TYPE_POA,
  ELEMENT_TYPE_POA_CELL,
  ELEMENT_TYPE_FOG,
  ELEMENT_TYPE_EDGE,
  ELEMENT_TYPE_DC
} from '../meep-constants';
import {
  FIELD_NAME,
  FIELD_TYPE,
  FIELD_GEO_LOCATION,
  FIELD_GEO_PATH,
  FIELD_GEO_RADIUS,
  getElemFieldVal,
  getElemFieldErr,
  setElemFieldVal,
  setElemFieldErr
} from '../util/elem-utils';

class IDCMap extends Component {
  constructor(props) {
    super(props);
    this.state = {};
    this.thisRef = createRef();
    this.configRef = createRef();
  }

  componentDidMount() {
    this.createMap();
  }

  componentWillUnmount() {
    this.destroyMap();
  }

  componentDidUpdate(prevProps) {
    if (prevProps.sandbox !== this.props.sandbox) {
      this.destroyMap();
      this.createMap();
      this.updateMarkers();
    }
  }

  shouldComponentUpdate(nextProps) {
    // Map size update
    let width = this.thisRef.current.offsetWidth;
    let height = this.thisRef.current.offsetHeight;
    if ((width && this.width !== width) || (height && this.height !== height)) {
      this.width = width;
      this.height = height;
      // console.log('Map view resized to: ' + width + 'x' + height);
      this.map.invalidateSize();
      return true;
    }

    // Target element change
    if (nextProps.configuredElement !== this.props.configuredElement) {
      return true;
    }

    // Sandbox changed
    if (nextProps.sandbox !== this.props.sandbox) {
      return true;
    }

    // Map changed
    if (!deepEqual(this.getMap(nextProps), this.getMap(this.props))) {
      return true;
    }

    return false;
  }

  getMap(props) {
    return (this.props.type === TYPE_CFG) ? props.cfgPageMap : props.execPageMap;
  }

  getCfg() {
    return (this.props.type === TYPE_CFG) ? this.props.mapCfg : this.props.sandboxCfg[this.props.sandboxName];
  }

  updateCfg(cfg) {
    switch (this.props.type) {
    case TYPE_CFG:
      this.props.changeMapCfg(updateObject(this.getCfg(), cfg));
      break;
    case TYPE_EXEC:
      var sandboxCfg = updateObject({}, this.props.sandboxCfg);
      if (sandboxCfg[this.props.sandboxName]) {
        sandboxCfg[this.props.sandboxName] = updateObject(sandboxCfg[this.props.sandboxName], cfg);
        this.props.changeSandboxCfg(sandboxCfg);
      }
      break;
    default:
      break;
    }
  }

  changeTable(table) {
    switch (this.props.type) {
    case TYPE_CFG:
      this.props.changeCfgTable(table);
      break;
    case TYPE_EXEC:
      this.props.changeExecTable(table);
      break;
    default:
      break;
    }
  }

  editElement(name) {
    // Update selected nodes in table
    const table = updateObject({}, this.props.cfgTable);
    const elem = this.getElementByName(table.entries, name);
    table.selected = elem ? [elem.id] : [];
    this.changeTable(table);

    // Open selected element in element configuration pane
    if (this.props.type === TYPE_CFG) {
      this.props.onEditElement(elem);

      // Update target element name & reset controls on target change
      if (name !== this.targetElemName) {
        this.map.pm.disableDraw('Marker');
        this.map.pm.disableDraw('Line');
        this.map.pm.disableGlobalEditMode();
        if (this.map.pm.globalDragModeEnabled()) {
          this.map.pm.toggleGlobalDragMode();
        }
        if (this.map.pm.globalRemovalEnabled()) {
          this.map.pm.toggleGlobalRemovalMode();
        }
      }
      this.targetElemName = name;
    }
  }

  getElementByName(entries, name) {
    for (var i = 0; i < entries.length; i++) {
      if (getElemFieldVal(entries[i], FIELD_NAME) === name) {
        return entries[i];
      }
    }
    return null;
  }

  createMap() {
    // Get stored configuration
    var cfg = this.getCfg();
    var lat = cfg.center ? cfg.center.lat : DEFAULT_MAP_LATITUDE;
    var lng = cfg.center ? cfg.center.lng : DEFAULT_MAP_LONGITUDE;
    var zoom = cfg.zoom ? cfg.zoom : DEFAULT_MAP_ZOOM;
    var baselayerName = cfg.baselayerName ? cfg.baselayerName : 'Positron';

    // Create Map instance
    var domNode = ReactDOM.findDOMNode(this);
    this.map = L.map(domNode, {
      center: [lat,lng],
      zoom: zoom,
      minZoom: 15,
      maxZoom: 18,
      drawControl: true
    });

    // Create GL Baselayers
    var positronBaselayer = L.mapboxGL({style: HOST_PATH + '/styles/positron/style.json'});
    var darkBaselayer = L.mapboxGL({style: HOST_PATH + '/styles/dark-matter/style.json'});
    var klokBaselayer = L.mapboxGL({style: HOST_PATH + '/styles/klokantech-basic/style.json'});
    var osmBaselayer = L.mapboxGL({style: HOST_PATH + '/styles/osm-bright/style.json'});
    var baselayers = {
      'Positron': positronBaselayer,
      'Black Matter': darkBaselayer,
      'Klokantech': klokBaselayer,
      'OSM Bright': osmBaselayer
    };

    // Create Layer Group Overlays
    this.ueOverlay = L.layerGroup();
    this.uePathOverlay = L.layerGroup();
    this.poaOverlay = L.layerGroup();
    this.poaRangeOverlay = L.layerGroup();
    this.computeOverlay = L.layerGroup();
    var overlays = {
      'terminal': this.ueOverlay,
      'terminal-path': this.uePathOverlay,
      'poa': this.poaOverlay,
      'poa-coverage': this.poaRangeOverlay,
      'compute': this.computeOverlay
    };

    // Create Layer Controller
    this.layerCtrl = L.control.layers(baselayers, overlays);

    // Create popup
    this.popup = L.popup();

    // Initialize map & layers
    this.layerCtrl.addTo(this.map);
    this.ueOverlay.addTo(this.map);
    this.uePathOverlay.addTo(this.map);
    this.poaOverlay.addTo(this.map);
    this.poaRangeOverlay.addTo(this.map);
    this.computeOverlay.addTo(this.map);

    // Set default base layer
    var baselayer = baselayers[baselayerName] ? baselayers[baselayerName] : positronBaselayer;
    baselayer.addTo(this.map);    

    // Handlers
    var _this = this;
    this.map.on('zoomend', function() {_this.setZoom(this);});
    this.map.on('moveend', function() {_this.setCenter(this);});
    this.map.on('baselayerchange', function(e) {_this.setBaseLayer(e);});

    // Add asset markers
    this.updateMarkers();

    if (this.props.type === TYPE_CFG) {
      // Draw Controls -- add leaflet-geoman controls with some options to the map
      this.map.pm.addControls({
        position: 'topleft',
        drawMarker: false, // adds button to draw markers
        drawCircleMarker: false, // adds button to draw circle markers
        drawPolyline: false, // adds button to draw rectangle
        drawRectangle: false, // adds button to draw rectangle
        drawPolygon: false, // adds button to draw polygon
        drawCircle: false, // adds button to draw circle
        editMode: false, // adds button to toggle edit mode for all layers
        dragMode: false, // adds button to toggle drag mode for all layers
        cutPolygon: false, // adds button to cut a hole in a polygon
        removalMode: false, // adds a button to remove layers
        pinningOption: false,	// $$ adds a button to toggle the Pinning Option
        snappingOption:	false	// $$ adds a button to toggle the Snapping Option
      });

      // Map handlers
      this.map.on('pm:globaleditmodetoggled', e => this.onEditModeToggle(e));
      this.map.on('pm:globaldragmodetoggled', e => this.onDragModeToggle(e));
      this.map.on('pm:globalremovalmodetoggled', e => this.onRemovalModeToggle(e));
      this.map.on('pm:create', e => this.onLayerCreated(e));
      this.map.on('pm:remove', e => this.onLayerRemoved(e));
    }
  }

  destroyMap() {
    if (this.map) {
      this.map.remove();
    }
  }

  setZoom(map) {
    this.updateCfg({zoom: map.getZoom()});
  }

  setCenter(map) {
    this.updateCfg({center: map.getCenter()});
  }

  setBaseLayer(event) {
    this.updateCfg({baselayerName: event.name});
  }

  setUeMarker(ue) {
    var latlng = L.latLng(L.GeoJSON.coordsToLatLng(ue.location.coordinates));
    var pathLatLngs = ue.path ? L.GeoJSON.coordsToLatLngs(ue.path.coordinates) : null;

    // Find existing UE marker
    var existingMarker;
    this.ueOverlay.eachLayer((marker) => {
      if (marker.options.meep.ue.id === ue.assetName){
        existingMarker = marker;
        return;
      }
    });

    if (existingMarker === undefined) {
      // Create path, if any
      // var p = ue.path ? L.GeoJSON.geometryToLayer(ue.path) : null;
      var p = pathLatLngs ? L.polyline(pathLatLngs) : null;
      console.log('New path', p);

      // Create new UE marker
      var m = L.marker(latlng, {
        meep: {
          ue: {
            id: ue.assetName,
            path: p,
            eopMode: ue.eopMode,
            velocity: ue.velocity
          }
        },
        draggable: true,
        pmIgnore: false
      });
      m.bindTooltip(ue.assetName).openTooltip();

      // Handlers
      var _this = this;
      m.on('click', function() {_this.clickUeMarker(this);});
      // m.on('pm:edit', e => console.log(e));
      // m.on('pm:update', e => console.log(e));
      // if (p) {
      //   p.on('pm:edit', e => console.log(e));
      //   p.on('pm:update', e => console.log(e));
      // }

      // Add to map overlay
      m.addTo(this.ueOverlay);
      if (p) {
        p.addTo(this.uePathOverlay);
      }
      // console.log('UE ' + id + ' added @ ' + latlng.toString());
    } else {
      // Update UE position & path
      existingMarker.setLatLng(latlng);

      // console.log('pathLatLngs', pathLatLngs);
      // console.log('Existing: updated path', existingMarker.options.meep.ue.path);
      // if (pathLatLngs) {
      //   if (existingMarker.options.meep.ue.path) {
      //     existingMarker.options.meep.ue.path.setLatLngs(pathLatLngs);
      //     console.log('Existing: updated path', existingMarker.options.meep.ue.path);
      //   }
      // }

      // Update, create or remove path
      if (pathLatLngs) {
        if (existingMarker.options.meep.ue.path) {
          existingMarker.options.meep.ue.path.setLatLngs(pathLatLngs);
          console.log('Existing: updated path', existingMarker.options.meep.ue.path);
        } else {
          var path = L.polyline(pathLatLngs);
          existingMarker.options.meep.ue.path = path;
          path.addTo(this.uePathOverlay);
          console.log('Existing: New polyline', path);
        }
      } else {
        if (existingMarker.options.meep.ue.path) {
          console.log('Existing: removing path', existingMarker.options.meep.ue.path);
          existingMarker.options.meep.ue.path.removeFrom(this.uePathOverlay);
          existingMarker.options.meep.ue.path = null;
        }
      }
    }
  }

  setPoaMarker(poa) {
    var latlng = L.latLng(L.GeoJSON.coordsToLatLng(poa.location.coordinates));

    // Find existing POA marker
    var existingMarker;
    this.poaOverlay.eachLayer((marker) => {
      if (marker.options.meep.poa.id === poa.assetName){
        existingMarker = marker;
        return;
      }
    });

    if (existingMarker === undefined) {
      // Create new POA marker & circle
      var c = L.circle(latlng, {
        meep: {
          range: {
            id: poa.assetName
          }
        },
        radius: poa.radius,
        opacity: '0.5',
        pmIgnore: true
      });
      var m = L.marker(latlng, {
        meep: {
          poa: {
            id: poa.assetName,
            range: c
          }
        },
        opacity: '0.5',
        draggable: true,
        pmIgnore: false
      });
      m.bindTooltip(poa.assetName).openTooltip();

      // Handlers
      var _this = this;
      m.on('click', function() {_this.clickPoaMarker(this);});
      m.on('drag', e => _this.onPoaMoved(e));
      m.on('dragend', e => _this.onPoaMoved(e));

      // Add to map overlay
      m.addTo(this.poaOverlay);
      c.addTo(this.poaRangeOverlay);
      // console.log('PoA ' + id + ' added @ ' + latlng.toString());
    } else {
      // Update POA position & range
      existingMarker.setLatLng(latlng);
      existingMarker.options.meep.poa.range.setLatLng(latlng);
      if (Number.isInteger(poa.radius) && poa.radius >= 0) {
        existingMarker.options.meep.poa.range.setRadius(poa.radius);
      }
    }
  }

  setComputeMarker(compute) {
    var latlng = L.latLng(L.GeoJSON.coordsToLatLng(compute.location.coordinates));

    // Find existing COMPUTE marker
    var existingMarker;
    this.computeOverlay.eachLayer((marker) => {
      if (marker.options.meep.compute.id === compute.assetName){
        existingMarker = marker;
        return;
      }
    });

    if (existingMarker === undefined) {
      // Creating new COMPUTE marker
      var m = L.marker(latlng, {
        meep: {
          compute: {
            id: compute.assetName
          }
        },
        opacity: '0.5',
        draggable: true,
        pmIgnore: false
      });
      m.bindTooltip(compute.assetName).openTooltip();

      // Handlers
      var _this = this;
      m.on('click', function() {_this.clickComputeMarker(this);});

      // Add to map overlay
      m.addTo(this.computeOverlay);
      // console.log('Compute ' + id + ' added @ ' + latlng.toString());
    } else {
      // Update COMPUTE position
      existingMarker.setLatLng(latlng);
    }
  }

  // UE Marker Event Handler
  clickUeMarker(marker) {
    this.editElement(marker.options.meep.ue.id);
    if (this.props.type === TYPE_EXEC) {
      var latlng = marker.getLatLng();
      var msg = '<b>id: ' + marker.options.meep.ue.id + '</b><br>';
      msg += 'path-mode: ' + marker.options.meep.ue.eopMode + '<br>';
      msg += 'velocity: ' + marker.options.meep.ue.velocity + ' m/s<br>';
      msg += latlng.toString();
      this.showPopup(latlng, msg);
    }
  }

  // POA Marker Event Handler
  clickPoaMarker(marker) {
    this.editElement(marker.options.meep.poa.id);
    if (this.props.type === TYPE_EXEC) {
      var latlng = marker.getLatLng();
      var msg = '<b>id: ' + marker.options.meep.poa.id + '</b><br>';
      msg += 'radius: ' + marker.options.meep.poa.range.options.radius + ' m<br>';
      msg += latlng.toString();
      this.showPopup(latlng, msg);
    }
  }

  // UE Marker Event Handler
  clickComputeMarker(marker) {
    this.editElement(marker.options.meep.compute.id);
    if (this.props.type === TYPE_EXEC) {
      var latlng = marker.getLatLng();
      var msg = '<b>id: ' + marker.options.meep.compute.id + '</b><br>';
      msg += latlng.toString();
      this.showPopup(latlng, msg);
    }
  }

  // Show position popup
  showPopup(latlng, msg) {
    // console.log(msg);
    this.popup
      .setLatLng(latlng)
      .setContent(msg)
      .openOn(this.map);
  }

  updateTargetMarker(map) {
    const target = this.props.configuredElement;
    if (!target) {
      return;
    }

    const location = getElemFieldVal(target, FIELD_GEO_LOCATION);
    const locationErr = getElemFieldErr(target, FIELD_GEO_LOCATION);
    if (location && !locationErr) {
      var name = getElemFieldVal(target, FIELD_NAME);
      var type = getElemFieldVal(target, FIELD_TYPE);
      var geoDataAsset;

      switch (type) {
      case ELEMENT_TYPE_UE:
        for (let i = 0; i < map.ueList.length; i++) {
          if (map.ueList[i].assetName === name) {
            geoDataAsset = map.ueList[i];
            break;
          }
        }
        if (!geoDataAsset) {
          geoDataAsset = {assetName: name, assetType: 'UE', subType: type};
          map.ueList.push(geoDataAsset);
        }
        geoDataAsset.location = {type: 'Point', coordinates: JSON.parse(location)};

        var path = getElemFieldVal(target, FIELD_GEO_PATH);
        var pathErr = getElemFieldErr(target, FIELD_GEO_PATH);
        geoDataAsset.path = (pathErr || !path) ? null : {type: 'LineString', coordinates: JSON.parse(path)};
        break;

      case ELEMENT_TYPE_POA:
      case ELEMENT_TYPE_POA_CELL:
        for (let i = 0; i < map.poaList.length; i++) {
          if (map.poaList[i].assetName === name) {
            geoDataAsset = map.poaList[i];
            break;
          }
        }
        if (!geoDataAsset) {
          geoDataAsset = {assetName: name, assetType: 'POA', subType: type};
          map.poaList.push(geoDataAsset);
        }
        geoDataAsset.location = {type: 'Point', coordinates: JSON.parse(location)};
        geoDataAsset.radius = getElemFieldVal(target, FIELD_GEO_RADIUS);
        break;

      case ELEMENT_TYPE_FOG:
      case ELEMENT_TYPE_EDGE:
      case ELEMENT_TYPE_DC:
        for (let i = 0; i < map.computeList.length; i++) {
          if (map.computeList[i].assetName === name) {
            geoDataAsset = map.computeList[i];
            break;
          }
        }
        if (!geoDataAsset) {
          geoDataAsset = {assetName: name, assetType: 'COMPUTE', subType: type};
          map.computeList.push(geoDataAsset);
        }
        geoDataAsset.location = {type: 'Point', coordinates: JSON.parse(location)};
        break;

      default:
        break;
      }
    }
  }

  updateMarkers() {
    if (!this.map) {
      return;
    }

    // Get copy of map data 
    var map = deepCopy(this.getMap(this.props));
    if (!map) {
      return;
    }

    // Update target marker geodata using configured element geodata, if any
    this.updateTargetMarker(map);

    // Set UE markers
    var ueMap = {};
    if (map.ueList) {
      for (let i = 0; i < map.ueList.length; i++) {
        let ue = map.ueList[i];
        this.setUeMarker(ue);
        ueMap[ue.assetName] = true;
      }
    }

    // Remove old UE markers
    this.ueOverlay.eachLayer((marker) => {
      if (!ueMap[marker.options.meep.ue.id]) {
        if (marker.options.meep.ue.path) {
          marker.options.meep.ue.path.removeFrom(this.uePathOverlay);
        }
        marker.removeFrom(this.ueOverlay);
      }
    });

    // Set POA markers
    var poaMap = {};
    if (map.poaList) {
      for (let i = 0; i < map.poaList.length; i++) {
        let poa = map.poaList[i];
        this.setPoaMarker(poa);
        poaMap[poa.assetName] = true;
      }
    }

    // Remove old POA markers
    this.poaOverlay.eachLayer((marker) => {
      if (!poaMap[marker.options.meep.poa.id]) {
        marker.options.meep.poa.range.removeFrom(this.poaRangeOverlay);
        marker.removeFrom(this.poaOverlay);
      }
    });

    // Set COMPUTE markers
    var computeMap = {};
    if (map.computeList) {
      for (let i = 0; i < map.computeList.length; i++) {
        let compute = map.computeList[i];
        this.setComputeMarker(compute);
        computeMap[compute.assetName] = true;
      }
    }

    // Remove old COMPUTE markers
    this.computeOverlay.eachLayer((marker) => {
      if (!computeMap[marker.options.meep.compute.id]) {
        marker.removeFrom(this.computeOverlay);
      }
    });
  }

  onEditModeToggle(e) {
    console.log('onEditModeToggle', this.map);
    console.log(e);
    var targetElemName = getElemFieldVal(this.props.configuredElement, FIELD_NAME);
    if (e.enabled) {
      this.setTarget(targetElemName);
    } else {
      this.updateTargetGeoData(targetElemName);
    }
  }

  onDragModeToggle(e) {
    var targetElemName = getElemFieldVal(this.props.configuredElement, FIELD_NAME);
    if (e.enabled) {
      this.setTarget(targetElemName);
    } else {
      this.updateTargetGeoData(targetElemName);
    }
  }

  onRemovalModeToggle(e) {
    console.log('onRemovalModeToggle');
    var targetElemName = getElemFieldVal(this.props.configuredElement, FIELD_NAME);
    if (e.enabled) {
      this.setTarget(targetElemName);
    } else {
      this.updateTargetGeoData(targetElemName);
    }
  }

  // onDraw(e) {
  //   console.log('onDraw');
  //   console.log(e);
  //   var marker = this.findMarker(this.props.configuredElement);
  //   if (marker) {
  //     if (e.shape === 'Marker') {
  //       // Disable marker drawing if it already exists
  //       this.map.pm.disableDraw('Marker');
  //     } else if (e.shape === 'Line') {
  //       var meepOptions = marker.options.meep;
  //       if (meepOptions && meepOptions.ue && meepOptions.ue.path) {
  //         // Disable path drawing if it already exists
  //         this.map.pm.disableDraw('Line');
  //       }
  //     }
  //   }
  // }

  onLayerCreated(e) {
    console.log('onLayerCreated');
    console.log(e);
  }

  onLayerRemoved(e) {
    console.log('onLayerRemoved');
    console.log(e);
    this.removeMarker(this.props.configuredElement);
    // var marker = this.findMarker(this.props.configuredElement);
  }

  onPoaMoved(e) {
    e.target.options.meep.poa.range.setLatLng(e.target.getLatLng());
  }

  updateConfiguredElement(name, val, err) {
    var updatedElem = updateObject({}, this.props.configuredElement);
    setElemFieldVal(updatedElem, name, val);
    setElemFieldErr(updatedElem, name, err);
    console.log(updatedElem);
    this.props.cfgElemUpdate(updatedElem);
  }

  updateTargetGeoData(targetElemName) {
    if (!targetElemName) {
      return;
    }
    console.log('Updating geodata for: ', targetElemName);
    var location = '';
    var path = '';

    // Get latest geoData from map, if any
    var markerInfo = this.getMarkerInfo(targetElemName);
    if (markerInfo && markerInfo.marker) {
      location = JSON.stringify(L.GeoJSON.latLngToCoords(markerInfo.marker.getLatLng()));
      if (markerInfo.type === 'UE' && markerInfo.marker.options.meep.ue.path) {
        path = JSON.stringify(L.GeoJSON.latLngsToCoords(markerInfo.marker.options.meep.ue.path.getLatLngs()));
      }
    }

    // Update configured element with map geodata
    this.updateConfiguredElement(FIELD_GEO_LOCATION, location, null);
    this.updateConfiguredElement(FIELD_GEO_PATH, path, null);
  }

  getMarkerInfo(name) {
    var marker;
    for (marker of this.ueOverlay.getLayers()) {
      if (marker.options.meep && (marker.options.meep.ue.id === name)) {
        return {marker: marker, type: 'UE'};
      }
    }
    for (marker of this.poaOverlay.getLayers()) {
      if (marker.options.meep && (marker.options.meep.poa.id === name)) {
        return {marker: marker, type: 'POA'};
      }
    }
    for (marker of this.computeOverlay.getLayers()) {
      if (marker.options.meep && (marker.options.meep.compute.id === name)) {
        return {marker: marker, type: 'COMPUTE'};
      }
    }
    return null;
  }

  removeMarker(name) {
    var marker;
    for (marker of this.ueOverlay.getLayers()) {
      if (marker.options.meep && (marker.options.meep.ue.id === name)) {
        marker.removeFrom(this.ueOverlay);
        return;
      }
    }
    for (marker of this.poaOverlay.getLayers()) {
      if (marker.options.meep && (marker.options.meep.poa.id === name)) {
        marker.removeFrom(this.poaOverlay);
        return;
      }
    }
    for (marker of this.computeOverlay.getLayers()) {
      if (marker.options.meep && (marker.options.meep.compute.id === name)) {
        marker.removeFrom(this.computeOverlay);
        return;
      }
    }
  }

  setTarget(name) {
    // Disable changes on all markers except target
    this.ueOverlay.eachLayer((marker) => {
      if (marker.pm && (!name || marker.options.meep.ue.id !== name)) {
        marker.pm.disable();
        var path = marker.options.meep.ue.path;
        if (path && path.pm) {
          path.pm.disable();
        }
      }
    });
    this.poaOverlay.eachLayer((marker) => {
      if (marker.pm && (!name || marker.options.meep.poa.id !== name)) {
        marker.pm.disable();
      }
    });
    this.computeOverlay.eachLayer((marker) => {
      if (marker.pm && (!name || marker.options.meep.compute.id !== name)) {
        marker.pm.disable();
      }
    });
  }

  updateEditControls() {
    if (this.props.type !== TYPE_CFG || !this.map) {
      return;
    }

    var drawMarkerEnabled = false;
    var drawPolylineEnabled = false;
    var editModeEnabled = false;
    var dragModeEnabled = false;
    var removalModeEnabled = false;

    // Update target element name & reset controls on target change
    var targetElemName = getElemFieldVal(this.props.configuredElement, FIELD_NAME);

    // Determine which controls to enable
    if (targetElemName) {
      var markerInfo = this.getMarkerInfo(targetElemName);
      if (markerInfo && markerInfo.marker) {
        // Enable path create/edit for UE only
        if (markerInfo.type === 'UE') {
          if (!markerInfo.marker.options.meep.ue.path) {
            drawPolylineEnabled = true;
          }
          editModeEnabled = true;
        }
        dragModeEnabled = true;
        removalModeEnabled = true;
      } else {
        // Enable marker creation
        drawMarkerEnabled = true;
      }
    }

    // Enable necessary controls
    this.map.pm.addControls({
      drawMarker: drawMarkerEnabled,
      drawPolyline: drawPolylineEnabled,
      editMode: editModeEnabled,
      dragMode: dragModeEnabled,
      removalMode: removalModeEnabled
    });

    // If in drawMarker mode, enable it by default
    if (drawMarkerEnabled) {
      this.map.pm.enableDraw('Marker');
    }

    // Set target element & disable edit on all other markers
    this.setTarget(targetElemName);
  }

  render() {
    this.updateMarkers();
    this.updateEditControls();
    return (
      <div ref={this.thisRef} style={{ height: '100%' }}>
        Map Component
      </div>
    );
  }
}

const mapStateToProps = state => {
  return {
    cfgPageMap: state.cfg.map,
    execPageMap: state.exec.map,
    sandbox: state.ui.sandbox,
    sandboxCfg: state.ui.sandboxCfg,
    mapCfg: state.ui.mapCfg,
    cfgTable: state.cfg.table,
    execTable: state.exec.table,
    configuredElement: state.cfg.elementConfiguration.configuredElement
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeMapCfg: cfg => dispatch(uiCfgChangeMapCfg(cfg)),
    changeSandboxCfg: cfg => dispatch(uiExecChangeSandboxCfg(cfg)),
    changeExecTable: table => dispatch(execChangeTable(table)),
    changeCfgTable: table => dispatch(cfgChangeTable(table)),
    cfgElemUpdate: element => dispatch(cfgElemUpdate(element)),
    changeCfgElement: element => dispatch(cfgElemEdit(element))
  };
};

const ConnectedIDCMap = connect(
  mapStateToProps,
  mapDispatchToProps
)(IDCMap);

export default ConnectedIDCMap;
