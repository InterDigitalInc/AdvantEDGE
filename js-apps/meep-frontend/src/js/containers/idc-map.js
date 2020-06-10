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

import _ from 'lodash';
import { connect } from 'react-redux';
import React, { Component, createRef } from 'react';
import ReactDOM from 'react-dom';
import L from 'leaflet';
import 'mapbox-gl';
import 'mapbox-gl-leaflet';
import { updateObject } from '../util/object-util';
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

    // Get Map div instance
    var domNode = ReactDOM.findDOMNode(this);
    var map = L.map(domNode, {
      center: [43.73752,7.42892],
      zoom: 15,
      minZoom: 15,
      maxZoom: 18,
      drawControl: true
    });
    console.log(map);

    // Creating GL Baselayers
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

    // Layer management
    var markerOverlay = L.layerGroup();
    var movingOverlay = L.layerGroup();
    var circleOverlay = L.layerGroup();
    var overlays = {
      'marker-PoC': markerOverlay,
      'circle-PoC': circleOverlay,
      'moving-PoC': movingOverlay
    };
      
    // Layer Controls
    var layerCtrl = L.control.layers(baselayers, overlays);

    // Initialize map & layers
    layerCtrl.addTo(map);
    markerOverlay.addTo(map); // this to have markers checked by default
    positronBaselayer.addTo(map); // this to select default baselayer
    // movingMarker.addTo(movingOverlay);
    // map.on('click', HandleClickMap);
  }

  render() {
    console.log(this.props.execMap.ueList);
    return (
      <div style={{ height: '100%' }}
        id="map"
      >
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
