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
import {
  execChangeMap
} from '../state/exec';

class IDCMap extends Component {
  constructor(props) {
    super(props);
    this.state = {};
    this.thisRef = createRef();
    this.configRef = createRef();
  }

  componentDidMount() {
    // Create Map instance
    var domNode = ReactDOM.findDOMNode(this);
    this.map = L.map(domNode, {
      center: [43.73752,7.42892],
      zoom: 15,
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
    this.poaOverlay = L.layerGroup();
    this.poaRangeOverlay = L.layerGroup();
    this.computeOverlay = L.layerGroup();
    var overlays = {
      'UE': this.ueOverlay,
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
    this.poaOverlay.addTo(this.map);
    this.poaRangeOverlay.addTo(this.map);
    this.computeOverlay.addTo(this.map);

    // Set default base layer
    positronBaselayer.addTo(this.map);
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
        title: poa.assetName,
        opacity: '0.5',
        draggable: true
      });

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

  setUeMarker(ue) {
    var latlng = L.latLng(L.GeoJSON.coordsToLatLng(ue.location.coordinates));

    // Find existing UE marker
    var existingMarker;
    this.ueOverlay.eachLayer((marker) => {
      if (marker.options.myId === ue.assetName){
        existingMarker = marker;
        return;
      }
    });

    if (existingMarker === undefined) {
      // Creating new UE marker
      var m = L.marker(latlng, {
        myId: ue.assetName,
        eopMode: ue.eopMode,
        velocity: ue.velocity,
        title: ue.assetName,
        opacity: '0.5',
        draggable: true
      });

      // Click handler
      var _this = this;
      m.on('click', function() {_this.clickUeMarker(this);});

      // Add to map overlay
      m.addTo(this.ueOverlay);
      // console.log('UE ' + id + ' added @ ' + latlng.toString());
    } else {
      // Update UE position
      existingMarker.setLatLng(latlng);
    }
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
  clickUeMarker(marker) {
    var latlng = marker.getLatLng();
    var msg = '<b>UE: ' + marker.options.myId + '</b><br>';
    msg += 'eopMode: ' + marker.options.eopMode + '<br>';
    msg += 'velocity: ' + marker.options.velocity + ' m/s<br>';
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

    // Render if asset list or positions updated
    return !deepEqual(nextProps.execMap, this.props.execMap);
  }

  render() {
    if (this.map) {
      // Set POA markers
      let poaList = this.props.execMap.poaList;
      for (let i = 0; i < poaList.length; i++) {
        this.setPoaMarker(poaList[i]);
      }

      // TODO -- Remove old POA markers

      // Set UE markers
      let ueList = this.props.execMap.ueList;
      for (let i = 0; i < ueList.length; i++) {
        this.setUeMarker(ueList[i]);
      }

      // TODO -- Remove old UE markers
    }

    return (
      <div ref={this.thisRef} style={{ height: '100%' }}>
        Map Component
      </div>
    );
  }
}

const mapStateToProps = state => {
  return {
    execMap: state.exec.map
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeExecMap: map => dispatch(execChangeMap(map))
  };
};

const ConnectedIDCMap = connect(
  mapStateToProps,
  mapDispatchToProps
)(IDCMap);

export default ConnectedIDCMap;
