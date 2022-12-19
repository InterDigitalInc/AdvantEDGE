/*
 * Copyright (c) 2022  The AdvantEDGE Authors
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
import tinycolor from 'tinycolor2';
import _ from 'lodash';
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
  ELEMENT_TYPE_UE,
  ELEMENT_TYPE_POA,
  ELEMENT_TYPE_POA_4G,
  ELEMENT_TYPE_POA_5G,
  ELEMENT_TYPE_POA_WIFI,
  ELEMENT_TYPE_FOG,
  ELEMENT_TYPE_EDGE,
  ELEMENT_TYPE_DC
} from '../meep-constants';
import {
  FIELD_NAME,
  FIELD_TYPE,
  FIELD_PARENT,
  FIELD_CELL_ID,
  FIELD_NR_CELL_ID,
  FIELD_MAC_ID,
  FIELD_UE_MAC_ID,
  FIELD_GEO_LOCATION,
  FIELD_GEO_PATH,
  FIELD_GEO_RADIUS,
  FIELD_CONNECTED,
  FIELD_WIRELESS_TYPE,
  FIELD_META_DISPLAY_MAP_COLOR,
  FIELD_META_DISPLAY_MAP_ICON,
  getElemFieldVal,
  getElemFieldErr,
  setElemFieldVal,
  setElemFieldErr,
  FIELD_DN_NAME,
  FIELD_DN_ECSP,
  FIELD_DN_LADN,
  FIELD_D2D_RADIUS
} from '../util/elem-utils';

import 'leaflet/dist/images/marker-shadow.png';

const ZONE_COLOR_LIST = [
  'blueviolet',
  'darkorange',
  'darkred',
  'limegreen',
  'blue',
  'purple',
  'gold',
  'darkturquoise'
];
const DISCONNECTED_COLOR = 'red';

const TYPE_UE = 'UE';
const TYPE_POA = 'POA';
const TYPE_COMPUTE = 'COMPUTE';

const UE_ICON = 'ion-iphone';
const UE_COLOR_DEFAULT = '#00ccff';
const UE_PATH_COLOR = '#008fb3';
const UE_OPACITY = 1.0;
const UE_OPACITY_BACKGROUND = 0.3;
const UE_PATH_OPACITY = 0.5;
const UE_PATH_OPACITY_BACKGROUND = 0.3;
const UE_RANGE_OPACITY = 0.05;

const POA_ICON = 'ion-connection-bars';
const POA_ICON_WIFI = 'ion-wifi';
const POA_COLOR_DEFAULT = '#696969';
const POA_OPACITY = 1.0;
const POA_OPACITY_BACKGROUND = 0.35;
const POA_RANGE_OPACITY = 0.05;
const POA_RANGE_OPACITY_BACKGROUND = 0.05;

const COMPUTE_ICON = 'ion-android-cloud';
const COMPUTE_COLOR_DEFAULT = '#0a50f2';
const COMPUTE_OPACITY = 1.0;
const COMPUTE_OPACITY_BACKGROUND = 0.35;

const OPACITY_TARGET = 1;

const LOCATION_PRECISION = 6;

const DEFAULT_MAP_STYLE = 'Positron';
const DEFAULT_MAP_LATITUDE = 0;
const DEFAULT_MAP_LONGITUDE = 0;
const DEFAULT_MAP_ZOOM = 2;

class IDCMap extends Component {
  constructor(props) {
    super(props);
    this.state = {};
    this.thisRef = createRef();
    this.configRef = createRef();
    this.rendering = false;
    this.zoneColorMap = {};
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

    // In CFG mode, only update when necessary
    if (this.props.type === TYPE_CFG) {
      // Target element update
      if (nextProps.configuredElement !== this.props.configuredElement) {
        return true;
      }
      // Scenario change
      if (nextProps.cfgScenarioName !== this.props.cfgScenarioName) {
        return true;
      }
      // Sandbox update
      if (nextProps.cfgView !== this.props.cfgView) {
        return true;
      }
      // Map asset change
      if (!deepEqual(this.getMap(nextProps), this.getMap(this.props))) {
        return true;
      }
      return false;
    }

    // Always update in EXEC mode
    return true;
  }

  getMap(props) {
    return (this.props.type === TYPE_CFG) ? props.cfgPageMap : props.execPageMap;
  }

  getCfg() {
    return (this.props.type === TYPE_CFG) ? this.props.mapCfg : this.props.sandboxCfg[this.props.sandboxName];
  }

  getTable() {
    return (this.props.type === TYPE_CFG) ? this.props.cfgTable : this.props.execTable;
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
    const table = updateObject({}, this.getTable());
    const elem = this.getElementByName(table.entries, name);
    table.selected = elem ? [elem.id] : [];
    this.changeTable(table);

    // Open selected element in element configuration pane
    if (this.props.type === TYPE_CFG) {
      this.props.onEditElement(elem ? elem : this.props.configuredElement);

      // Update target element name & reset controls on target change
      if (name !== this.targetElemName) {
        this.map.pm.disableDraw('Marker');
        this.map.pm.disableDraw('Line');
        if (this.map.pm.globalEditEnabled()) {
          this.map.pm.disableGlobalEditMode();
        }
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
    var element = entries[name];
    return element ? element : null;
  }

  createMap() {
    // Get stored configuration
    var cfg = this.getCfg();
    var lat = (cfg && cfg.center) ? cfg.center.lat : DEFAULT_MAP_LATITUDE;
    var lng = (cfg && cfg.center) ? cfg.center.lng : DEFAULT_MAP_LONGITUDE;
    var zoom = (cfg && cfg.zoom) ? cfg.zoom : DEFAULT_MAP_ZOOM;
    var baselayerName = (cfg && cfg.baselayerName) ? cfg.baselayerName : DEFAULT_MAP_STYLE;

    // Map bounds
    const corner1 = L.latLng(-90, -180);
    const corner2 = L.latLng(90, 180);
    const bounds = L.latLngBounds(corner1, corner2);

    // Create Map instance
    var domNode = ReactDOM.findDOMNode(this);
    this.map = L.map(domNode, {
      center: [lat,lng],
      zoom: zoom,
      minZoom: 2,
      maxZoom: 20,
      drawControl: true,
      maxBounds: bounds,
      maxBoundsViscosity: 1.0
    });
    this.map.attributionControl.addAttribution('<a href="https://www.maptiler.com/copyright/?_ga=2.45788834.742970109.1593090041-1523068243.1593090041" target="_blank">© MapTiler</a>');
    this.map.attributionControl.addAttribution('<a href="https://www.openstreetmap.org/copyright" target="_blank">© OpenStreetMap contributors</a>');

    // Create GL Baselayers
    var positronBaselayer = L.mapboxGL({style: HOST_PATH + '/map/styles/positron/style.json'});
    var darkBaselayer = L.mapboxGL({style: HOST_PATH + '/map/styles/dark-matter/style.json'});
    var klokBaselayer = L.mapboxGL({style: HOST_PATH + '/map/styles/klokantech-basic/style.json'});
    var osmBaselayer = L.mapboxGL({style: HOST_PATH + '/map/styles/osm-bright/style.json'});
    var baselayers = {
      'Positron': positronBaselayer,
      'Black Matter': darkBaselayer,
      'Klokantech': klokBaselayer,
      'OSM Bright': osmBaselayer
    };

    // Create Layer Group Overlays
    this.ueOverlay = L.layerGroup();
    this.ueRangeOverlay = L.layerGroup();
    this.uePathOverlay = L.layerGroup();
    this.poaOverlay = L.layerGroup();
    this.poaRangeOverlay = L.layerGroup();
    this.computeOverlay = L.layerGroup();
    var overlays = {
      'terminal': this.ueOverlay,
      'terminal-path': this.uePathOverlay,
      'poa': this.poaOverlay,
      'poa-coverage': this.poaRangeOverlay,
      'compute': this.computeOverlay,
      'd2d-coverage': this.ueRangeOverlay
    };

    // Create Layer Controller
    this.layerCtrl = L.control.layers(baselayers, overlays);

    // Create popup
    this.popup = L.popup();

    // Initialize map & layers
    this.layerCtrl.addTo(this.map);
    this.ueOverlay.addTo(this.map);
    this.ueRangeOverlay.addTo(this.map);
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

      // Set control states
      this.updateEditControls();

      // Map handlers
      this.map.on('pm:globaleditmodetoggled', e => this.onEditModeToggle(e));
      this.map.on('pm:globaldragmodetoggled', e => this.onDragModeToggle(e));
      this.map.on('pm:create', e => this.onLayerCreated(e));
    }
  }

  destroyMap() {
    if (this.map) {
      this.map.remove();
    }
  }

  setZoom(map) {
    if (map && !this.rendering) {
      this.updateCfg({zoom: map.getZoom()});
    }
  }

  setCenter(map) {
    if (map && !this.rendering) {
      this.updateCfg({center: map.getCenter()});
    }
  }

  setBaseLayer(event) {
    this.updateCfg({baselayerName: event.name});
  }

  getD2dRadius(scenarioName) {
    var radius = 0;
    var table = this.getTable();
    if (table && table.entries) {
      radius = getElemFieldVal(table.entries[scenarioName], FIELD_D2D_RADIUS);
    }
    return radius;
  }

  isD2dEnabled(name) {
    if (this.getWirelessTypePrio(name) && this.getWirelessTypePrio(name).includes('d2d')) {
      return true;
    }
    return false;
  }
  getUePoa(ue) {
    var poa = null;
    var table = this.getTable();
    if (table && table.entries) {
      poa = getElemFieldVal(table.entries[ue], FIELD_PARENT);
    }
    return poa;
  }

  getUeZone(ue) {
    var zone = null;
    var table = this.getTable();
    if (table && table.entries) {
      var poa = getElemFieldVal(table.entries[ue], FIELD_PARENT);
      zone = poa ? this.getPoaZone(poa) : null;
    }
    return zone;
  }

  getPoaZone(poa) {
    var zone = null;
    var table = this.getTable();
    if (table && table.entries) {
      zone = getElemFieldVal(table.entries[poa], FIELD_PARENT);
    }
    return zone;
  }

  getComputeZone(compute) {
    var zone = null;
    var table = this.getTable();
    if (table && table.entries) {
      var computeType = getElemFieldVal(table.entries[compute], FIELD_TYPE);
      var parent = getElemFieldVal(table.entries[compute], FIELD_PARENT);
      if (computeType === ELEMENT_TYPE_EDGE) {
        zone = parent;
      } else if (computeType === ELEMENT_TYPE_FOG) {
        zone = parent ? this.getPoaZone(parent) : null;
      }
    }
    return zone;
  }

  // Get Colors
  getZoneColor(zone) {
    var color = null;
    if (zone) {
      var table = this.getTable();
      if (table && table.entries) {
        // Get zone color from meta
        color = getElemFieldVal(table.entries[zone], FIELD_META_DISPLAY_MAP_COLOR);
        if (!color) {
          // Get zone color from zone color map
          color = this.zoneColorMap[zone];
          if (!color) {
            // Get a new color for this zone
            color = this.zoneColorMap[zone] = ZONE_COLOR_LIST[Object.keys(this.zoneColorMap).length % ZONE_COLOR_LIST.length];
            // // Generate a random color for this zone
            // color = this.zoneColorMap[zone] = tinycolor.random().toHexString();
          }
        }
      }
    }
    return color;
  }

  getUeColor(ue) {
    var color = undefined;
    var connected = this.isConnected(ue.id) || (ue.d2dInRange && ue.d2dInRange.length > 0);
    if (!connected) {
      color = DISCONNECTED_COLOR;
    }
    return color ? color : UE_COLOR_DEFAULT;
  }

  getPoaColor(poa) {
    var color = this.getZoneColor(this.getPoaZone(poa));
    return color ? color : POA_COLOR_DEFAULT;
  }

  getComputeColor(compute) {
    if (!this.isConnected(compute)) {
      return DISCONNECTED_COLOR;
    }
    return COMPUTE_COLOR_DEFAULT;
  }

  // Get connected status
  isConnected(name) {
    var connected = false;
    var table = this.getTable();
    if (table && table.entries) {
      connected = getElemFieldVal(table.entries[name], FIELD_CONNECTED);
    }
    return connected;
  }

  // Get wireless type Priority
  getWirelessTypePrio(name) {
    var wirelessTypePrio = '';
    var table = this.getTable();
    if (table && table.entries) {
      wirelessTypePrio = getElemFieldVal(table.entries[name], FIELD_WIRELESS_TYPE);
    }
    return wirelessTypePrio;
  }

  // Set Icons
  setUeIcon(iconDiv, ue) {
    var table = this.getTable();
    if (table && table.entries) {
      var metaIcon = getElemFieldVal(table.entries[ue], FIELD_META_DISPLAY_MAP_ICON);
      var icon = metaIcon ? metaIcon : UE_ICON;
      iconDiv.className = 'custom-marker-icon ion ' + icon;
      iconDiv.innerHTML = '';
    }
  }

  setPoaIcon(iconDiv, iconTextDiv, poa) {
    var table = this.getTable();
    if (table && table.entries) {
      var poaType = getElemFieldVal(table.entries[poa], FIELD_TYPE);
      var metaIcon = getElemFieldVal(table.entries[poa], FIELD_META_DISPLAY_MAP_ICON);
      var icon = metaIcon ? metaIcon : (poaType === ELEMENT_TYPE_POA_WIFI) ? POA_ICON_WIFI : POA_ICON;
      iconDiv.className = 'custom-marker-icon ion ' + icon;
      iconDiv.innerHTML = '';

      var innerHTML = '';
      if (!metaIcon) {
        if (poaType === ELEMENT_TYPE_POA_4G) {
          innerHTML = '4G';
        }
        if (poaType === ELEMENT_TYPE_POA_5G) {
          innerHTML = '5G';
        }
      }
      iconTextDiv.innerHTML = innerHTML;
    }
  }

  setComputeIcon(iconDiv, compute) {
    var table = this.getTable();
    if (table && table.entries) {
      var metaIcon = getElemFieldVal(table.entries[compute], FIELD_META_DISPLAY_MAP_ICON);
      var icon = metaIcon ? metaIcon : COMPUTE_ICON;
      iconDiv.className = 'custom-marker-icon ion ' + icon;
      iconDiv.innerHTML = '';
    }
  }

  // Set styles
  setUeMarkerStyle(marker) {
    if (marker._icon) {
      // Set marker color
      var color = tinycolor(this.getUeColor(marker.options.meep.ue));
      var markerStyle = marker._icon.querySelector('.custom-marker-pin').style;
      markerStyle['background'] = color;
      markerStyle['border-color'] = color.darken(10);

      // Set UE range color
      if (this.isD2dEnabled(marker.options.meep.ue.id)) {
        marker.options.meep.ue.range.setStyle({color: color});
      }

      // Set marker icon
      var iconDiv = marker._icon.querySelector('.custom-marker-icon');
      this.setUeIcon(iconDiv, marker.options.meep.ue.id);
    }
  }

  setPoaMarkerStyle(marker) {
    if (marker._icon) {
      // Set marker color
      var color = tinycolor(this.getPoaColor(marker.options.meep.poa.id));
      var markerStyle = marker._icon.querySelector('.custom-marker-pin').style;
      markerStyle['background'] = color;
      markerStyle['border-color'] = color.darken(10);

      // Set POA range color
      marker.options.meep.poa.range.setStyle({color: color});

      // Set marker icon
      var iconDiv = marker._icon.querySelector('.custom-marker-icon');
      var iconTextDiv = marker._icon.querySelector('.custom-marker-icon-text');
      this.setPoaIcon(iconDiv, iconTextDiv, marker.options.meep.poa.id);
    }
  }

  setComputeMarkerStyle(marker) {
    if (marker._icon) {
      // Set marker color
      var color = tinycolor(this.getComputeColor(marker.options.meep.compute.id));
      var markerStyle = marker._icon.querySelector('.custom-marker-pin').style;
      markerStyle['background'] = color;
      markerStyle['border-color'] = color.darken(10);

      // Set marker icon
      var iconDiv = marker._icon.querySelector('.custom-marker-icon');
      this.setComputeIcon(iconDiv, marker.options.meep.compute.id);
    }
  }

  getLocationStr(latlng) {
    return '[' + latlng.lat.toFixed(LOCATION_PRECISION) + ', ' + latlng.lng.toFixed(LOCATION_PRECISION) + ']';
  }

  // UE Marker Event Handler
  updateUePopup(marker) {
    var table = this.getTable();
    var d2dInRange = false;
    var poaInRange = false;
    if (marker && table && table.entries) {
      var latlng = marker.getLatLng();
      var hasPath = (marker.options.meep.ue.path) ? true : false;
      var msg = '<b>id: ' + marker.options.meep.ue.id + '</b><br>';
      var ownMac = getElemFieldVal(table.entries[marker.options.meep.ue.id], FIELD_UE_MAC_ID);
      if (ownMac !== '') {
        msg += 'mac: ' + ownMac + '<br>';
      }
      msg += 'velocity: ' + (hasPath ? marker.options.meep.ue.velocity : '0') + ' m/s<br>';
      if (this.isD2dEnabled(marker.options.meep.ue.id)) {
        if (marker.options.meep.ue.d2dInRange) {
          var d2dConnType = 'd2d: ' + marker.options.meep.ue.d2dInRange + '<br>';
          d2dInRange = true;
        }
      }
      if (this.isConnected(marker.options.meep.ue.id)) {
        var poa = this.getUePoa(marker.options.meep.ue.id);
        var poaType = getElemFieldVal(table.entries[poa], FIELD_TYPE);
        var poaConnType = 'poa: ' + poa + '<br>';
        poaInRange = true;
        switch(poaType) {
        case ELEMENT_TYPE_POA_4G:
          poaConnType += 'cell: ' + getElemFieldVal(table.entries[poa], FIELD_CELL_ID) + '<br>';
          break;
        case ELEMENT_TYPE_POA_5G:
          poaConnType += 'cell: ' + getElemFieldVal(table.entries[poa], FIELD_NR_CELL_ID) + '<br>';
          break;
        case ELEMENT_TYPE_POA_WIFI:
          poaConnType += 'poa mac: ' + getElemFieldVal(table.entries[poa], FIELD_MAC_ID) + '<br>';
          break;
        default:
          break;
        }
        poaConnType += 'zone: ' + this.getUeZone(marker.options.meep.ue.id) + '<br>';
      }

      if (!d2dInRange && !poaInRange) {
        msg += 'state: <b style="color:red;">DISCONNECTED</b><br>';
      } else if (poaInRange && !d2dInRange) {
        msg += 'd2d: none <br>';
        msg += poaConnType;
      } else if (d2dInRange && !poaInRange) {
        msg += d2dConnType;
        msg += 'poa: none <br>';
      } else if (d2dInRange && poaInRange) {
        msg += d2dConnType + poaConnType;
      }
      msg += 'wireless: ' + (this.getWirelessTypePrio(marker.options.meep.ue.id) || 'wifi,5g,4g,other') + '<br>';
      msg += 'location: ' + this.getLocationStr(latlng);
      marker.getPopup().setContent(msg);
    }
  }

  // POA Marker Event Handler
  updatePoaPopup(marker) {
    var latlng = marker.getLatLng();
    var poaType = getElemFieldVal(this.getTable().entries[marker.options.meep.poa.id], FIELD_TYPE);
    var msg = '<b>id: ' + marker.options.meep.poa.id + '</b><br>';
    msg += 'radius: ' + marker.options.meep.poa.range.options.radius + ' m<br>';
    switch (poaType) {
    case ELEMENT_TYPE_POA_4G:
      msg += 'cell: ' + getElemFieldVal(this.getTable().entries[marker.options.meep.poa.id], FIELD_CELL_ID) + '<br>';
      break;
    case ELEMENT_TYPE_POA_5G:
      msg += 'cell: ' + getElemFieldVal(this.getTable().entries[marker.options.meep.poa.id], FIELD_NR_CELL_ID) + '<br>';
      break;
    case ELEMENT_TYPE_POA_WIFI:
      msg += 'mac: ' + getElemFieldVal(this.getTable().entries[marker.options.meep.poa.id], FIELD_MAC_ID) + '<br>';
      break;
    default:
      break;
    }
    msg += 'zone: ' + this.getPoaZone(marker.options.meep.poa.id) + '<br>';
    msg += 'location: ' + this.getLocationStr(latlng);
    marker.getPopup().setContent(msg);
  }

  // UE Marker Event Handler
  updateComputePopup(marker) {
    var table = this.getTable();
    if (marker && table && table.entries) {
      // Retrieve state
      const networkName = getElemFieldVal(table.entries[marker.options.meep.compute.id], FIELD_DN_NAME);
      const edgeProvider = getElemFieldVal(table.entries[marker.options.meep.compute.id], FIELD_DN_ECSP);
      const ladn = getElemFieldVal(table.entries[marker.options.meep.compute.id], FIELD_DN_LADN);
      var appInstanceTable = this.props.appInstanceTable;
      var latlng = marker.getLatLng();
      // Parse mec application state on current popup
      var appInstances = [];
      for (var i = 0; i < appInstanceTable.length ; i++) {
        if (appInstanceTable[i].nodeName === marker.options.meep.compute.id) {
          appInstances.push(appInstanceTable[i]);
        }
      }
      // Sort parsed array of mec app
      var sortedAppInstances = _.sortBy(appInstances, ['name']);
      // Modify render message
      var msg = '<b>id: ' + marker.options.meep.compute.id + '</b><br>';
      if (edgeProvider) {
        msg += 'service-provider: ' + edgeProvider + '<br>';
      }
      if (networkName) {
        msg += 'data-network: ' + networkName;
        if (ladn) {
          msg += ' (LADN)';
        }
        msg += '<br>';
      }
      msg += 'applications: <br>';
      if (appInstances) {
        sortedAppInstances.forEach(elem => {
          msg += '<li>' + elem.name + ' ' + '(id: ' + elem.id.substring(0,8) + '...)' + '<br>';
        });
      }
      msg += 'location: ' + this.getLocationStr(latlng);
      marker.getPopup().setContent(msg);
    }
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
      var p = !pathLatLngs ? null : L.polyline(pathLatLngs, {
        meep: {
          path: {
            id: ue.assetName
          }
        },
        color: UE_PATH_COLOR,
        opacity: UE_PATH_OPACITY,
        pmIgnore: (this.props.type === TYPE_CFG) ? false : true
      });

      var markerIcon = L.divIcon({
        className: '',
        html: '<div class="custom-marker-pin"></div><div class="custom-marker-icon"></div>',
        iconSize: [30, 42],
        iconAnchor: [15, 42],
        popupAnchor: [0, -36]
      });

      // Create new UE marker & circle
      var c = L.circle(latlng, {
        meep: {
          range: {
            id: ue.assetName
          }
        },
        radius: (this.props.type === TYPE_CFG) ?  (this.isD2dEnabled(ue.assetName) ? this.getD2dRadius(this.props.cfgScenarioName): 0) : ue.radius || 0,
        opacity: UE_RANGE_OPACITY,
        pmIgnore: true
      });

      var m = L.marker(latlng, {
        meep: {
          ue: {
            id: ue.assetName,
            path: p,
            eopMode: ue.eopMode,
            velocity: ue.velocity,
            connected: true,
            d2dInRange: ue.d2dInRange,
            range: c
          }
        },
        icon: markerIcon,
        opacity: UE_OPACITY,
        draggable: (this.props.type === TYPE_CFG) ? true : false,
        pmIgnore: (this.props.type === TYPE_CFG) ? false : true
      });
      m.bindTooltip(ue.assetName).openTooltip();

      // Handlers
      var _this = this;
      m.on('add', (e) => _this.setUeMarkerStyle(e.target));
      if (this.props.type === TYPE_CFG) {
        m.on('click', function() {_this.editElement(m.options.meep.ue.id);});
      } else {
        m.bindPopup('').openPopup();
        m.on('popupopen', (e) => _this.updateUePopup(e.target));
      }

      // Add to map overlay
      m.addTo(this.ueOverlay);
      c.addTo(this.ueRangeOverlay);
      if (p) {
        p.addTo(this.uePathOverlay);
      }

    } else {
      // Update UE position, path, mode, velocity, range & d2dInRange
      existingMarker.setLatLng(latlng);
      existingMarker.options.meep.ue.eopMode = ue.eopMode;
      existingMarker.options.meep.ue.velocity = ue.velocity;
      if (Number.isInteger(ue.radius) && ue.radius >= 0) {
        existingMarker.options.meep.ue.range.setLatLng(latlng);
        existingMarker.options.meep.ue.range.setRadius(ue.radius);
        existingMarker.options.meep.ue.d2dInRange = ue.d2dInRange;
      }

      // Update, create or remove path
      if (pathLatLngs) {
        if (existingMarker.options.meep.ue.path) {
          existingMarker.options.meep.ue.path.setLatLngs(pathLatLngs);
        } else {
          var path = L.polyline(pathLatLngs, {
            meep: {
              path: {
                id: ue.assetName
              }
            },
            color: UE_PATH_COLOR,
            opacity: UE_PATH_OPACITY,
            pmIgnore: (this.props.type === TYPE_CFG) ? false : true
          });
          existingMarker.options.meep.ue.path = path;
          path.addTo(this.uePathOverlay);
        }
      } else {
        if (existingMarker.options.meep.ue.path) {
          existingMarker.options.meep.ue.path.removeFrom(this.uePathOverlay);
          existingMarker.options.meep.ue.path = null;
        }
      }

      // Refresh marker style if necessary
      if (this.props.type === TYPE_CFG) {
        this.setUeMarkerStyle(existingMarker);
      } else {
        var connected = this.isConnected(ue.assetName) || (ue.d2dInRange && ue.d2dInRange.length > 0);
        if (existingMarker.options.meep.ue.connected !== connected) {
          this.setUeMarkerStyle(existingMarker);
          existingMarker.options.meep.ue.connected = connected;
        }
      }

      // Refresh popup text & position
      if (this.props.type === TYPE_EXEC) {
        this.updateUePopup(existingMarker);
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
        color: this.getPoaColor(poa.assetName),
        radius: poa.radius || 0,
        opacity: POA_RANGE_OPACITY,
        pmIgnore: true
      });

      var markerIcon = L.divIcon({
        className: '',
        html: '<div class="custom-marker-pin"></div><div class="custom-marker-icon"></div><div class="custom-marker-icon-text"></div>',
        iconSize: [30, 42],
        iconAnchor: [15, 42],
        popupAnchor: [0, -36]
      });

      var m = L.marker(latlng, {
        meep: {
          poa: {
            id: poa.assetName,
            range: c
          }
        },
        icon: markerIcon,
        opacity: POA_OPACITY,
        draggable: (this.props.type === TYPE_CFG) ? true : false,
        pmIgnore: (this.props.type === TYPE_CFG) ? false : true
      });
      m.bindTooltip(poa.assetName).openTooltip();

      // Handlers
      var _this = this;
      m.on('add', (e) => _this.setPoaMarkerStyle(e.target));
      if (this.props.type === TYPE_CFG) {
        m.on('click', function() {_this.editElement(m.options.meep.poa.id);});
        m.on('drag', e => _this.onPoaMoved(e));
        m.on('dragend', e => _this.onPoaMoved(e));
      } else {
        m.bindPopup('').openPopup();
        m.on('popupopen', (e) => _this.updatePoaPopup(e.target));
      }

      // Add to map overlay
      m.addTo(this.poaOverlay);
      c.addTo(this.poaRangeOverlay);

    } else {
      // Update POA position & range
      existingMarker.setLatLng(latlng);
      if (Number.isInteger(poa.radius) && poa.radius >= 0) {
        existingMarker.options.meep.poa.range.setLatLng(latlng);
        existingMarker.options.meep.poa.range.setRadius(poa.radius);
      }

      // Refresh marker style & position
      if (this.props.type === TYPE_CFG) {
        this.setPoaMarkerStyle(existingMarker);
      } else {
        this.updatePoaPopup(existingMarker);
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
      // Create new marker
      var markerIcon = L.divIcon({
        className: '',
        html: '<div class="custom-marker-pin"></div><div class="custom-marker-icon"></div>',
        iconSize: [30, 42],
        iconAnchor: [15, 42],
        popupAnchor: [0, -36]
      });

      // Creating new COMPUTE marker
      var m = L.marker(latlng, {
        meep: {
          compute: {
            id: compute.assetName,
            connected: true
          }
        },
        icon: markerIcon,
        opacity: COMPUTE_OPACITY,
        draggable: (this.props.type === TYPE_CFG) ? true : false,
        pmIgnore: (this.props.type === TYPE_CFG) ? false : true
      });
      m.bindTooltip(compute.assetName).openTooltip();

      // Handlers
      var _this = this;
      m.on('add', (e) => _this.setComputeMarkerStyle(e.target));
      if (this.props.type === TYPE_CFG) {
        m.on('click', function() {_this.editElement(m.options.meep.compute.id);});
      } else {
        m.bindPopup('').openPopup();
        m.on('popupopen', (e) => _this.updateComputePopup(e.target));
      }

      // Add to map overlay
      m.addTo(this.computeOverlay);

    } else {
      // Update COMPUTE position
      existingMarker.setLatLng(latlng);

      // Refresh marker style if necessary
      if (this.props.type === TYPE_CFG) {
        this.setComputeMarkerStyle(existingMarker);
      } else {
        var connected = this.isConnected(compute.assetName);
        if (existingMarker.options.meep.compute.connected !== connected) {
          this.setComputeMarkerStyle(existingMarker);
          existingMarker.options.meep.compute.connected = connected;
        }
      }

      // Refresh popup text & position
      if (this.props.type === TYPE_EXEC) {
        this.updateComputePopup(existingMarker);
      }
    }
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
          geoDataAsset = {assetName: name, assetType: TYPE_UE, subType: type};
          map.ueList.push(geoDataAsset);
        }
        geoDataAsset.location = {type: 'Point', coordinates: JSON.parse(location)};

        var path = getElemFieldVal(target, FIELD_GEO_PATH);
        var pathErr = getElemFieldErr(target, FIELD_GEO_PATH);
        geoDataAsset.path = (pathErr || !path) ? null : {type: 'LineString', coordinates: JSON.parse(path)};
        break;

      case ELEMENT_TYPE_POA:
      case ELEMENT_TYPE_POA_4G:
      case ELEMENT_TYPE_POA_5G:
      case ELEMENT_TYPE_POA_WIFI:
        for (let i = 0; i < map.poaList.length; i++) {
          if (map.poaList[i].assetName === name) {
            geoDataAsset = map.poaList[i];
            break;
          }
        }
        if (!geoDataAsset) {
          geoDataAsset = {assetName: name, assetType: TYPE_POA, subType: type};
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
          geoDataAsset = {assetName: name, assetType: TYPE_COMPUTE, subType: type};
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
    if (this.props.type === TYPE_CFG) {
      this.updateTargetMarker(map);
    }

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
        marker.removeFrom(this.ueRangeOverlay);
      }
    });
  }

  onEditModeToggle(e) {
    var targetElemName = getElemFieldVal(this.props.configuredElement, FIELD_NAME);
    if (e.enabled) {
      this.setTarget(targetElemName);
    } else {
      this.updateTargetGeoData(targetElemName, '', '');
    }
  }

  onDragModeToggle(e) {
    var targetElemName = getElemFieldVal(this.props.configuredElement, FIELD_NAME);
    if (e.enabled) {
      this.setTarget(targetElemName);
    } else {
      this.updateTargetGeoData(targetElemName, '', '');
    }
  }

  onLayerCreated(e) {
    var location = '';
    var path = '';

    // Get marker location or path & remove newly created layer
    if (e.shape === 'Marker') {
      location = JSON.stringify(L.GeoJSON.latLngToCoords(e.marker.getLatLng()));
      e.marker.removeFrom(this.map);
    } else if (e.shape === 'Line') {
      path = JSON.stringify(L.GeoJSON.latLngsToCoords(e.layer.getLatLngs()));
      e.layer.removeFrom(this.map);
    } else {
      return;
    }

    // Update configured element & refresh map to create the new marker or path
    var targetElemName = getElemFieldVal(this.props.configuredElement, FIELD_NAME);
    this.updateTargetGeoData(targetElemName, location, path);
  }

  onPoaMoved(e) {
    e.target.options.meep.poa.range.setLatLng(e.target.getLatLng());
  }

  updateConfiguredElement(name, val, err) {
    var updatedElem = updateObject({}, this.props.configuredElement);
    setElemFieldVal(updatedElem, name, val);
    setElemFieldErr(updatedElem, name, err);
    this.props.cfgElemUpdate(updatedElem);
  }

  updateTargetGeoData(targetElemName, location, path) {
    if (!targetElemName) {
      return;
    }

    // Get latest geoData from map, if any
    if (!location) {
      var markerInfo = this.getMarkerInfo(targetElemName);
      if (markerInfo && markerInfo.marker) {
        location = JSON.stringify(L.GeoJSON.latLngToCoords(markerInfo.marker.getLatLng()));
        if (!path && markerInfo.type === TYPE_UE && markerInfo.marker.options.meep.ue.path) {
          path = JSON.stringify(L.GeoJSON.latLngsToCoords(markerInfo.marker.options.meep.ue.path.getLatLngs()));
        }
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
        return {marker: marker, type: TYPE_UE};
      }
    }
    for (marker of this.poaOverlay.getLayers()) {
      if (marker.options.meep && (marker.options.meep.poa.id === name)) {
        return {marker: marker, type: TYPE_POA};
      }
    }
    for (marker of this.computeOverlay.getLayers()) {
      if (marker.options.meep && (marker.options.meep.compute.id === name)) {
        return {marker: marker, type: TYPE_COMPUTE};
      }
    }
    return null;
  }

  setTarget(target) {
    // Disable changes on all markers except target
    this.ueOverlay.eachLayer((marker) => {
      var path = marker.options.meep.ue.path;
      if (marker.pm && (!target || marker.options.meep.ue.id !== target)) {
        marker.pm.disable();
        marker.setOpacity(target ? UE_OPACITY_BACKGROUND : UE_OPACITY);
        if (path && path.pm) {
          path.pm.disable();
          path.setStyle({opacity: target ? UE_PATH_OPACITY_BACKGROUND : UE_PATH_OPACITY});
        }
      } else {
        marker.setOpacity(OPACITY_TARGET);
        if (path) {
          path.setStyle({opacity: OPACITY_TARGET});
        }
      }
    });
    this.poaOverlay.eachLayer((marker) => {
      if (marker.pm && (!target || marker.options.meep.poa.id !== target)) {
        marker.pm.disable();
        marker.setOpacity(target ? POA_OPACITY_BACKGROUND : POA_OPACITY);
        marker.options.meep.poa.range.setStyle({opacity: target ? POA_RANGE_OPACITY_BACKGROUND : POA_RANGE_OPACITY});
      } else {
        marker.setOpacity(OPACITY_TARGET);
        marker.options.meep.poa.range.setStyle({opacity: OPACITY_TARGET});
      }
    });
    this.computeOverlay.eachLayer((marker) => {
      if (marker.pm && (!target || marker.options.meep.compute.id !== target)) {
        marker.pm.disable();
        marker.setOpacity(target ? COMPUTE_OPACITY_BACKGROUND : COMPUTE_OPACITY);
      } else {
        marker.setOpacity(OPACITY_TARGET);
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
        if (markerInfo.type === TYPE_UE) {
          if (!markerInfo.marker.options.meep.ue.path) {
            drawPolylineEnabled = true;
          }
          editModeEnabled = true;
        }
        dragModeEnabled = true;
        // removalModeEnabled = true;
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

    // Disable draw, edit & drag modes if controls disabled
    if (!drawMarkerEnabled) {
      this.map.pm.disableDraw('Marker');
    }
    if (!drawPolylineEnabled) {
      this.map.pm.disableDraw('Line');
    }
    if (!editModeEnabled) {
      if (this.map.pm.globalEditEnabled()) {
        this.map.pm.disableGlobalEditMode();
      }
    }
    if (!dragModeEnabled) {
      if (this.map.pm.globalDragModeEnabled()) {
        this.map.pm.toggleGlobalDragMode();
      }
    }

    // Set target element & disable edit on all other markers
    this.setTarget(targetElemName);
  }

  render() {
    this.rendering = true;
    this.updateMarkers();
    this.updateEditControls();
    this.rendering = false;
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
    configuredElement: state.cfg.elementConfiguration.configuredElement,
    cfgView: state.ui.cfgView,
    cfgScenarioName: state.cfg.scenario.name,
    appInstanceTable: state.exec.appInstanceTable.data
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
