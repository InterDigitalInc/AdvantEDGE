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
import deepEqual from 'deep-equal';
import { updateObject } from '../util/object-util';
import {
  uiCfgChangeMapCfg,
  uiExecChangeSandboxCfg
} from '../state/ui';
import {
  TYPE_CFG,
  HOST_PATH,
  DEFAULT_MAP_LATITUDE,
  DEFAULT_MAP_LONGITUDE,
  DEFAULT_MAP_ZOOM
} from '../meep-constants';

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
    // Check for map size update
    let width = this.thisRef.current.offsetWidth;
    let height = this.thisRef.current.offsetHeight;
    if ((width && this.width !== width) || (height && this.height !== height)) {
      this.width = width;
      this.height = height;
      // console.log('Map view resized to: ' + width + 'x' + height);
      this.map.invalidateSize();
      return true;
    }

    // Update if sandbox changed
    if (nextProps.sandbox !== this.props.sandbox) {
      return true;
    }

    // Update if map changed
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
    this.map.on('zoomend', function() {_this.updateZoom(this);});
    this.map.on('moveend', function() {_this.updateCenter(this);});
    this.map.on('baselayerchange', function(e) {_this.updateBaseLayer(e);});
  }

  destroyMap() {
    if (this.map) {
      this.map.remove();
    }
  }

  updateCfg(cfg) {
    if (this.props.type === TYPE_CFG) {
      this.props.changeMapCfg(updateObject(this.getCfg(), cfg));
    } else {
      var sandboxCfg = updateObject({}, this.props.sandboxCfg);
      if (sandboxCfg[this.props.sandboxName]) {
        sandboxCfg[this.props.sandboxName] = updateObject(sandboxCfg[this.props.sandboxName], cfg);
        this.props.changeSandboxCfg(sandboxCfg);
      }
    }
  }

  updateZoom(map) {
    this.updateCfg({zoom: map.getZoom()});
  }

  updateCenter(map) {
    this.updateCfg({center: map.getCenter()});
  }

  updateBaseLayer(event) {
    this.updateCfg({baselayerName: event.name});
  }

  setUeMarker(ue) {
    var latlng = L.latLng(L.GeoJSON.coordsToLatLng(ue.location.coordinates));
    var p = ue.path ? L.GeoJSON.geometryToLayer(ue.path) : null;

    // Find existing UE marker
    var existingMarker;
    this.ueOverlay.eachLayer((marker) => {
      if (marker.options.meep.ue.id === ue.assetName){
        existingMarker = marker;
        return;
      }
    });

    if (existingMarker === undefined) {
      // Create new UE marker & path
      var m = L.marker(latlng, {
        meep: {
          ue: {
            id: ue.assetName,
            path: p,
            eopMode: ue.eopMode,
            velocity: ue.velocity
          }
        },
        draggable: false
      });
      m.bindTooltip(ue.assetName).openTooltip();

      // Click handler
      var _this = this;
      m.on('click', function() {_this.clickUeMarker(this);});

      // Add to map overlay
      m.addTo(this.ueOverlay);
      if (p) {
        p.addTo(this.uePathOverlay);
      }
      // console.log('UE ' + id + ' added @ ' + latlng.toString());
    } else {
      // Update UE position & path
      existingMarker.setLatLng(latlng);

      // Update path
      if (existingMarker.options.meep.ue.path) {
        existingMarker.options.meep.ue.path.removeFrom(this.uePathOverlay);
      }
      if (p) {
        existingMarker.options.meep.ue.path = p;
        p.addTo(this.uePathOverlay);
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
        opacity: '0.5'
      });
      var m = L.marker(latlng, {
        meep: {
          poa: {
            id: poa.assetName,
            range: c
          }
        },
        opacity: '0.5',
        draggable: false
      });
      m.bindTooltip(poa.assetName).openTooltip();

      // Click handler
      var _this = this;
      m.on('click', function(e) {_this.clickPoaMarker(this,e);});

      // Add to map overlay
      m.addTo(this.poaOverlay);
      c.addTo(this.poaRangeOverlay);
      // console.log('PoA ' + id + ' added @ ' + latlng.toString());
    } else {
      // Update POA position
      existingMarker.setLatLng(latlng);
      existingMarker.options.meep.poa.range.setLatLng(latlng);
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
        draggable: false
      });
      m.bindTooltip(compute.assetName).openTooltip();

      // Click handler
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
    var latlng = marker.getLatLng();
    var msg = '<b>id: ' + marker.options.meep.ue.id + '</b><br>';
    msg += 'path-mode: ' + marker.options.meep.ue.eopMode + '<br>';
    msg += 'velocity: ' + marker.options.meep.ue.velocity + ' m/s<br>';
    msg += latlng.toString();
    this.showPopup(latlng, msg);
  }

  // POA Marker Event Handler
  clickPoaMarker(marker) {
    var latlng = marker.getLatLng();
    var msg = '<b>id: ' + marker.options.meep.poa.id + '</b><br>';
    msg += 'radius: ' + marker.options.meep.poa.range.options.radius + ' m<br>';
    msg += latlng.toString();
    this.showPopup(latlng, msg);
  }

  // UE Marker Event Handler
  clickComputeMarker(marker) {
    var latlng = marker.getLatLng();
    var msg = '<b>id: ' + marker.options.meep.compute.id + '</b><br>';
    msg += latlng.toString();
    this.showPopup(latlng, msg);
  }

  // Show position popup
  showPopup(latlng, msg) {
    // console.log(msg);
    this.popup
      .setLatLng(latlng)
      .setContent(msg)
      .openOn(this.map);
  }

  updateMarkers() {
    if (!this.map) {
      return;
    }
    var map = this.getMap(this.props);
    if (!map) {
      return;
    }

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

  render() {
    this.updateMarkers();
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
    mapCfg: state.ui.mapCfg
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeMapCfg: cfg => dispatch(uiCfgChangeMapCfg(cfg)),
    changeSandboxCfg: cfg => dispatch(uiExecChangeSandboxCfg(cfg))
  };
};

const ConnectedIDCMap = connect(
  mapStateToProps,
  mapDispatchToProps
)(IDCMap);

export default ConnectedIDCMap;
