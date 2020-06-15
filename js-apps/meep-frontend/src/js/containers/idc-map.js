/*
 * Copyright (c) 2019  InterDigital Communications, Inc
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
  execChangeMap
} from '../state/exec';
import {
  uiExecChangeSandboxCfg
} from '../state/ui';

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
    if (!deepEqual(nextProps.execMap, this.props.execMap)) {
      return true;
    }

    return false;
  }

  createMap() {
    // Get stored configuration
    var cfg = this.props.sandboxCfg[this.props.sandboxName];
    var lat = cfg.center ? cfg.center.lat : 43.73752;
    var lng = cfg.center ? cfg.center.lng : 7.42892;
    var zoom = cfg.zoom ? cfg.zoom : 15;
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
    var positronBaselayer = L.mapboxGL({style: 'http://10.3.16.105:30080/styles/positron/style.json'});
    var darkBaselayer = L.mapboxGL({style: 'http://10.3.16.105:30080/styles/dark-matter/style.json'});
    var klokBaselayer = L.mapboxGL({style: 'http://10.3.16.105:30080/styles/klokantech-basic/style.json'});
    var osmBaselayer = L.mapboxGL({style: 'http://10.3.16.105:30080/styles/osm-bright/style.json'});
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
      'UE': this.ueOverlay,
      'UE Path': this.uePathOverlay,
      'POA': this.poaOverlay,
      'POA Range': this.poaRangeOverlay,
      'Compute': this.computeOverlay
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

  updateZoom(map) {
    var zoom = map.getZoom();
    var sandboxCfg = updateObject({}, this.props.sandboxCfg);
    sandboxCfg[this.props.sandboxName]['zoom'] = zoom;
    this.props.changeSandboxCfg(sandboxCfg);
  }

  updateCenter(map) {
    var center = map.getCenter();
    var sandboxCfg = updateObject({}, this.props.sandboxCfg);
    sandboxCfg[this.props.sandboxName]['center'] = center;
    this.props.changeSandboxCfg(sandboxCfg);
  }

  updateBaseLayer(event) {
    var sandboxCfg = updateObject({}, this.props.sandboxCfg);
    sandboxCfg[this.props.sandboxName]['baselayerName'] = event.name;
    this.props.changeSandboxCfg(sandboxCfg);
  }

  setUeMarker(ue) {
    var latlng = L.latLng(L.GeoJSON.coordsToLatLng(ue.location.coordinates));
    var p = ue.path ? L.GeoJSON.geometryToLayer(ue.path) : null;

    // Find existing UE marker
    var existingMarker;
    this.ueOverlay.eachLayer((marker) => {
      if (marker.options.myId === ue.assetName){
        existingMarker = marker;
        return;
      }
    });

    if (existingMarker === undefined) {
      // Create new UE marker & path
      var m = L.marker(latlng, {
        myId: ue.assetName,
        path: p,
        eopMode: ue.eopMode,
        velocity: ue.velocity,
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
      if (existingMarker.options.path) {
        existingMarker.options.path.removeFrom(this.uePathOverlay);
      }
      if (p) {
        existingMarker.options.path = p;
        p.addTo(this.uePathOverlay);
      }
    }
  }

  setPoaMarker(poa) {
    var latlng = L.latLng(L.GeoJSON.coordsToLatLng(poa.location.coordinates));

    // Find existing POA marker
    var existingMarker;
    this.poaOverlay.eachLayer((marker) => {
      if (marker.options.myId === poa.assetName){
        existingMarker = marker;
        return;
      }
    });

    if (existingMarker === undefined) {
      // Create new POA marker & circle
      var c = L.circle(latlng, {
        myId: poa.assetName,
        radius: poa.radius,
        opacity: '0.5'
      });
      var m = L.marker(latlng, {
        myId: poa.assetName,
        myCircle: c,
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
      existingMarker.options.myCircle.setLatLng(latlng);
    }
  }

  setComputeMarker(compute) {
    var latlng = L.latLng(L.GeoJSON.coordsToLatLng(compute.location.coordinates));

    // Find existing COMPUTE marker
    var existingMarker;
    this.computeOverlay.eachLayer((marker) => {
      if (marker.options.myId === compute.assetName){
        existingMarker = marker;
        return;
      }
    });

    if (existingMarker === undefined) {
      // Creating new COMPUTE marker
      var m = L.marker(latlng, {
        myId: compute.assetName,
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
    var msg = '<b>UE: ' + marker.options.myId + '</b><br>';
    msg += 'eopMode: ' + marker.options.eopMode + '<br>';
    msg += 'velocity: ' + marker.options.velocity + ' m/s<br>';
    msg += latlng.toString();
    this.showPopup(latlng, msg);
  }

  // POA Marker Event Handler
  clickPoaMarker(marker) {
    var latlng = marker.getLatLng();
    var msg = '<b>POA: ' + marker.options.myId + '</b><br>';
    msg += 'radius: ' + marker.options.myCircle.options.radius + ' m<br>';
    msg += latlng.toString();
    this.showPopup(latlng, msg);
  }

  // UE Marker Event Handler
  clickComputeMarker(marker) {
    var latlng = marker.getLatLng();
    var msg = '<b>Compute: ' + marker.options.myId + '</b><br>';
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

    // Set UE markers
    let ueList = this.props.execMap.ueList;
    var ueMap = {};
    for (let i = 0; i < ueList.length; i++) {
      let ue = ueList[i];
      this.setUeMarker(ue);
      ueMap[ue.assetName] = true;
    }

    // Remove old UE markers
    this.ueOverlay.eachLayer((marker) => {
      if (!ueMap[marker.options.myId]) {
        if (marker.options.path) {
          marker.options.path.removeFrom(this.uePathOverlay);
        }
        marker.removeFrom(this.ueOverlay);
      }
    });

    // Set POA markers
    let poaList = this.props.execMap.poaList;
    var poaMap = {};
    for (let i = 0; i < poaList.length; i++) {
      let poa = poaList[i];
      this.setPoaMarker(poa);
      poaMap[poa.assetName] = true;
    }

    // Remove old POA markers
    this.poaOverlay.eachLayer((marker) => {
      if (!poaMap[marker.options.myId]) {
        marker.options.myCircle.removeFrom(this.poaRangeOverlay);
        marker.removeFrom(this.poaOverlay);
      }
    });

    // Set COMPUTE markers
    let computeList = this.props.execMap.computeList;
    var computeMap = {};
    for (let i = 0; i < computeList.length; i++) {
      let compute = computeList[i];
      this.setComputeMarker(compute);
      computeMap[compute.assetName] = true;
    }

    // Remove old COMPUTE markers
    this.computeOverlay.eachLayer((marker) => {
      if (!computeMap[marker.options.myId]) {
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
    execMap: state.exec.map,
    sandbox: state.ui.sandbox,
    sandboxCfg: state.ui.sandboxCfg
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeExecMap: map => dispatch(execChangeMap(map)),
    changeSandboxCfg: cfg => dispatch(uiExecChangeSandboxCfg(cfg))
  };
};

const ConnectedIDCMap = connect(
  mapStateToProps,
  mapDispatchToProps
)(IDCMap);

export default ConnectedIDCMap;
