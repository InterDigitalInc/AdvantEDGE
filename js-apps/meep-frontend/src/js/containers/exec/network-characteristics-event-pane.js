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
import React, { Component } from 'react';
import { Select } from '@rmwc/select';
import { Grid, GridCell } from '@rmwc/grid';
import { updateObject } from '../../util/object-util';
import NCGroup from '../../components/helper-components/nc-group';
import CancelApplyPair from '../../components/helper-components/cancel-apply-pair';
import IDSelect from '../../components/helper-components/id-select';

import {
  camelCasePrefix,
  firstLetterUpper
} from '../../util/string-manipulation';

import {
  EXEC_EVT_NC_TYPE,
  EXEC_EVT_NC_NAME,

  // Network element types
  ELEMENT_TYPE_SCENARIO,
  ELEMENT_TYPE_OPERATOR,
  ELEMENT_TYPE_OPERATOR_GENERIC,
  ELEMENT_TYPE_OPERATOR_CELL,
  ELEMENT_TYPE_ZONE,
  ELEMENT_TYPE_POA,
  ELEMENT_TYPE_POA_GENERIC,
  ELEMENT_TYPE_POA_4G,
  ELEMENT_TYPE_POA_5G,
  ELEMENT_TYPE_POA_WIFI,
  ELEMENT_TYPE_DC,
  //ELEMENT_TYPE_CN,
  ELEMENT_TYPE_EDGE,
  ELEMENT_TYPE_FOG,
  ELEMENT_TYPE_UE,
  //ELEMENT_TYPE_MECSVC,
  ELEMENT_TYPE_UE_APP,
  //ELEMENT_TYPE_EXT_UE_APP,
  ELEMENT_TYPE_EDGE_APP,
  ELEMENT_TYPE_CLOUD_APP,

  // NC Group Prefixes
  PREFIX_INT_DOM,
  PREFIX_INT_ZONE,
  PREFIX_INTRA_ZONE,
  PREFIX_TERM_LINK,
  PREFIX_LINK,
  PREFIX_APP,

  DOMAIN_TYPE_STR,
  DOMAIN_CELL_TYPE_STR,
  POA_TYPE_STR,
  POA_4G_TYPE_STR,
  POA_5G_TYPE_STR,
  POA_WIFI_TYPE_STR,
  DC_TYPE_STR,
  UE_APP_TYPE_STR,
  EDGE_APP_TYPE_STR,
  CLOUD_APP_TYPE_STR
} from '../../meep-constants';

import {
  // Field Names
  FIELD_NAME,
  FIELD_TYPE,
  FIELD_INT_DOM_LATENCY,
  FIELD_INT_DOM_LATENCY_VAR,
  FIELD_INT_DOM_LATENCY_DIST,
  FIELD_INT_DOM_THROUGHPUT_DL,
  FIELD_INT_DOM_THROUGHPUT_UL,
  FIELD_INT_DOM_PKT_LOSS,
  FIELD_INT_ZONE_LATENCY,
  FIELD_INT_ZONE_LATENCY_VAR,
  FIELD_INT_ZONE_THROUGHPUT_DL,
  FIELD_INT_ZONE_THROUGHPUT_UL,
  FIELD_INT_ZONE_PKT_LOSS,
  FIELD_INTRA_ZONE_LATENCY,
  FIELD_INTRA_ZONE_LATENCY_VAR,
  FIELD_INTRA_ZONE_THROUGHPUT_DL,
  FIELD_INTRA_ZONE_THROUGHPUT_UL,
  FIELD_INTRA_ZONE_PKT_LOSS,
  FIELD_TERM_LINK_LATENCY,
  FIELD_TERM_LINK_LATENCY_VAR,
  FIELD_TERM_LINK_THROUGHPUT_DL,
  FIELD_TERM_LINK_THROUGHPUT_UL,
  FIELD_TERM_LINK_PKT_LOSS,
  FIELD_LINK_LATENCY,
  FIELD_LINK_LATENCY_VAR,
  FIELD_LINK_THROUGHPUT_DL,
  FIELD_LINK_THROUGHPUT_UL,
  FIELD_LINK_PKT_LOSS,
  FIELD_APP_LATENCY,
  FIELD_APP_LATENCY_VAR,
  FIELD_APP_THROUGHPUT_DL,
  FIELD_APP_THROUGHPUT_UL,
  FIELD_APP_PKT_LOSS,
  getElemFieldVal,
  setElemFieldVal,
  setElemFieldErr
} from '../../util/elem-utils';

const ncApplicableTypes = [
  ELEMENT_TYPE_SCENARIO,
  ELEMENT_TYPE_OPERATOR_GENERIC,
  ELEMENT_TYPE_OPERATOR_CELL,
  ELEMENT_TYPE_OPERATOR,
  ELEMENT_TYPE_ZONE,
  ELEMENT_TYPE_POA_GENERIC,
  ELEMENT_TYPE_POA_4G,
  ELEMENT_TYPE_POA_5G,
  ELEMENT_TYPE_POA_WIFI,
  ELEMENT_TYPE_POA,
  ELEMENT_TYPE_DC,
  ELEMENT_TYPE_EDGE,
  ELEMENT_TYPE_FOG,
  ELEMENT_TYPE_UE,
  ELEMENT_TYPE_UE_APP,
  ELEMENT_TYPE_EDGE_APP,
  ELEMENT_TYPE_CLOUD_APP
];

class NetworkCharacteristicsEventPane extends Component {
  constructor(props) {
    super(props);

    this.state = {
      dialogOpen: false,
      ncTypes: []
    };
  }

  componentDidMount() {
    let ncTypes = ncApplicableTypes.filter(e => {
      for (const key in this.props.networkElements) {
        if (e === getElemFieldVal(this.props.networkElements[key], FIELD_TYPE)) {
          return true;
        }
      }
      return false;
    });

    this.setState({ ncTypes });
  }

  onNetworkCharacPaneClose(e) {
    e.preventDefault();
    var updatedElem = updateObject({}, this.props.element);
    setElemFieldVal(updatedElem, FIELD_NAME, '');
    setElemFieldVal(updatedElem, FIELD_TYPE, '');
    this.props.updateElement(updatedElem);
    this.props.onClose(e);
  }

  triggerEvent(e) {
    e.preventDefault();
    var element = this.props.element;
    var type = getElemFieldVal(element, FIELD_TYPE);

    // Verify that no field is in error
    var fieldsInError = 0;
    _.forOwn(
      element,
      val => (fieldsInError = val.err ? fieldsInError + 1 : fieldsInError)
    );
    if (fieldsInError) {
      return;
    }

    var neType = '';
    switch(type) {
    case ELEMENT_TYPE_OPERATOR_GENERIC:
      neType = DOMAIN_TYPE_STR;
      break;
    case ELEMENT_TYPE_OPERATOR_CELL:
      neType = DOMAIN_CELL_TYPE_STR;
      break;
    case ELEMENT_TYPE_POA_GENERIC:
      neType = POA_TYPE_STR;
      break;
    case ELEMENT_TYPE_POA_4G:
      neType = POA_4G_TYPE_STR;
      break;
    case ELEMENT_TYPE_POA_5G:
      neType = POA_5G_TYPE_STR;
      break;
    case ELEMENT_TYPE_POA_WIFI:
      neType = POA_WIFI_TYPE_STR;
      break;
    case ELEMENT_TYPE_DC:
      neType = DC_TYPE_STR;
      break;
    case ELEMENT_TYPE_UE_APP:
      neType = UE_APP_TYPE_STR;
      break;
    case ELEMENT_TYPE_EDGE_APP:
      neType = EDGE_APP_TYPE_STR;
      break;
    case ELEMENT_TYPE_CLOUD_APP:
      neType = CLOUD_APP_TYPE_STR;
      break;
    default:
      neType = type;
    }

    var ncEvent = {
      name: 'name',
      type: this.props.currentEvent,
      eventNetworkCharacteristicsUpdate: {
        elementName: getElemFieldVal(element, FIELD_NAME),
        elementType: neType,
        netChar: {}
      }
    };

    // Retrieve and set net characteristics from element
    this.setNetCharFromElem(ncEvent.eventNetworkCharacteristicsUpdate.netChar, element);

    // trigger event with this.props.api
    this.props.api.sendEvent(this.props.currentEvent, ncEvent, error => {
      if (!error) {
        this.setState({ dialogOpen: true });
        this.props.onSuccess();
      }
    });
  }

  firstElementMatchingType(type) {
    var elements = _.chain(this.props.networkElements)
      .filter(e => {
        var elemType = getElemFieldVal(e, FIELD_TYPE);
        if (type === ELEMENT_TYPE_OPERATOR_GENERIC) {
          return elemType === ELEMENT_TYPE_OPERATOR;
        } else if (type === ELEMENT_TYPE_POA_GENERIC) {
          return elemType === ELEMENT_TYPE_POA;
        } else if (elemType === ELEMENT_TYPE_ZONE) {
          return type.startsWith(elemType);
        } 
        return type === elemType;
      })
      .value();

    return elements.length ? elements[0] : null;
  }

  currentPrefix() {
    var type = getElemFieldVal(this.props.element, FIELD_TYPE);
    switch (type) {
    case ELEMENT_TYPE_SCENARIO:
      return PREFIX_INT_DOM;
    case ELEMENT_TYPE_OPERATOR:
    case ELEMENT_TYPE_OPERATOR_GENERIC:
    case ELEMENT_TYPE_OPERATOR_CELL:
      return PREFIX_INT_ZONE;
    case ELEMENT_TYPE_ZONE:
      return PREFIX_INTRA_ZONE;
    case ELEMENT_TYPE_POA:
    case ELEMENT_TYPE_POA_GENERIC:
    case ELEMENT_TYPE_POA_4G:
    case ELEMENT_TYPE_POA_5G:
    case ELEMENT_TYPE_POA_WIFI:
      return PREFIX_TERM_LINK;
    case ELEMENT_TYPE_EDGE:
      return PREFIX_LINK;
    case ELEMENT_TYPE_FOG:
      return PREFIX_LINK;
    case ELEMENT_TYPE_DC:
      return PREFIX_LINK;
    case ELEMENT_TYPE_UE:
      return PREFIX_LINK;
    case ELEMENT_TYPE_UE_APP:
      return PREFIX_APP;
    case ELEMENT_TYPE_EDGE_APP:
      return PREFIX_APP;
    case ELEMENT_TYPE_CLOUD_APP:
      return PREFIX_APP;
    default:
      return '';
    }
  }

  setNetCharFromElem(netChar, element) {
    // Retrieve field names
    var latencyFieldName = null;
    var latencyVarFieldName = null;
    var latencyDistFieldName = null;
    var throughputDlFieldName = null;
    var throughputUlFieldName = null;
    var packetLossFieldName = null;
    switch (this.currentPrefix()) {
    case PREFIX_INT_DOM:
      latencyFieldName = FIELD_INT_DOM_LATENCY;
      latencyVarFieldName = FIELD_INT_DOM_LATENCY_VAR;
      latencyDistFieldName = FIELD_INT_DOM_LATENCY_DIST;
      throughputDlFieldName = FIELD_INT_DOM_THROUGHPUT_DL;
      throughputUlFieldName = FIELD_INT_DOM_THROUGHPUT_UL;
      packetLossFieldName = FIELD_INT_DOM_PKT_LOSS;
      break;
    case PREFIX_INT_ZONE:
      latencyFieldName = FIELD_INT_ZONE_LATENCY;
      latencyVarFieldName = FIELD_INT_ZONE_LATENCY_VAR;
      throughputDlFieldName = FIELD_INT_ZONE_THROUGHPUT_DL;
      throughputUlFieldName = FIELD_INT_ZONE_THROUGHPUT_UL;
      packetLossFieldName = FIELD_INT_ZONE_PKT_LOSS;
      break;
    case PREFIX_INTRA_ZONE:
      latencyFieldName = FIELD_INTRA_ZONE_LATENCY;
      latencyVarFieldName = FIELD_INTRA_ZONE_LATENCY_VAR;
      throughputDlFieldName = FIELD_INTRA_ZONE_THROUGHPUT_DL;
      throughputUlFieldName = FIELD_INTRA_ZONE_THROUGHPUT_UL;
      packetLossFieldName = FIELD_INTRA_ZONE_PKT_LOSS;
      break;
    case PREFIX_TERM_LINK:
      latencyFieldName = FIELD_TERM_LINK_LATENCY;
      latencyVarFieldName = FIELD_TERM_LINK_LATENCY_VAR;
      throughputDlFieldName = FIELD_TERM_LINK_THROUGHPUT_DL;
      throughputUlFieldName = FIELD_TERM_LINK_THROUGHPUT_UL;
      packetLossFieldName = FIELD_TERM_LINK_PKT_LOSS;
      break;
    case PREFIX_LINK:
      latencyFieldName = FIELD_LINK_LATENCY;
      latencyVarFieldName = FIELD_LINK_LATENCY_VAR;
      throughputDlFieldName = FIELD_LINK_THROUGHPUT_DL;
      throughputUlFieldName = FIELD_LINK_THROUGHPUT_UL;
      packetLossFieldName = FIELD_LINK_PKT_LOSS;
      break;
    case PREFIX_APP:
      latencyFieldName = FIELD_APP_LATENCY;
      latencyVarFieldName = FIELD_APP_LATENCY_VAR;
      throughputDlFieldName = FIELD_APP_THROUGHPUT_DL;
      throughputUlFieldName = FIELD_APP_THROUGHPUT_UL;
      packetLossFieldName = FIELD_APP_PKT_LOSS;
      break;
    default:
      return null;
    }

    // Update net characteristics
    netChar.latency = getElemFieldVal(element, latencyFieldName);
    netChar.latencyVariation = getElemFieldVal(element, latencyVarFieldName);
    if (latencyDistFieldName) {
      netChar.latencyDistribution = getElemFieldVal(element, latencyDistFieldName);
    }
    netChar.throughputDl = getElemFieldVal(element, throughputDlFieldName);
    netChar.throughputUl = getElemFieldVal(element, throughputUlFieldName);
    netChar.packetLoss = getElemFieldVal(element, packetLossFieldName);
  }

  fieldName(genericFieldName) {
    const prefix = this.currentPrefix();
    var name = (
      camelCasePrefix(prefix) + firstLetterUpper(genericFieldName)
    ).replace(/\s/g, '');
    return name;
  }

  onUpdateElement(name, val, err) {
    var updatedElem = updateObject({}, this.props.element);
    setElemFieldVal(updatedElem, name, val);
    setElemFieldErr(updatedElem, name, err);
    this.props.updateElement(updatedElem);
  }

  getElementByName(name) {
    var elements = this.props.networkElements;
    var element = elements[name];
    return element ? element : null;
  }

  render() {
    var element = this.props.element;
    var nbErrors = _.reduce(
      element,
      (result, value) => {
        return value.err ? (result = result + 1) : result;
      },
      0
    );

    var elements = _.chain(this.props.networkElements)
      .filter(e => {
        var type = getElemFieldVal(element, FIELD_TYPE);
        var elemType = getElemFieldVal(e, FIELD_TYPE);

        if (type === ELEMENT_TYPE_OPERATOR_GENERIC) {
          return elemType === ELEMENT_TYPE_OPERATOR;
        } else if (type === ELEMENT_TYPE_POA_GENERIC) {
          return elemType === ELEMENT_TYPE_POA;
        } else if (elemType === ELEMENT_TYPE_ZONE) {
          return type.startsWith(elemType);
        } 
        return type === elemType;
      })
      .map(e => {
        return getElemFieldVal(e, FIELD_NAME);
      })
      .value();

    return (
      <div>
        <Grid style={styles.field}>
          <GridCell span="8">
            <Select
              style={styles.select}
              label="Network Element Type"
              outlined
              options={this.state.ncTypes}
              onChange={event => {
                var elem = this.firstElementMatchingType(event.target.value);
                this.props.updateElement(elem);
              }}
              data-cy={EXEC_EVT_NC_TYPE}
              value={element ? getElemFieldVal(element, FIELD_TYPE) || '' : ''}
            />
          </GridCell>
          <GridCell span="4"></GridCell>
        </Grid>

        <Grid>
          <GridCell span="8">
            <IDSelect
              span="8"
              label="Network Element"
              value={element ? getElemFieldVal(element, FIELD_NAME) || '' : ''}
              options={elements}
              onChange={event => {
                this.props.updateElement(
                  this.getElementByName(event.target.value)
                );
              }}
              cydata={EXEC_EVT_NC_NAME}
            />
          </GridCell>
          <GridCell span="4"></GridCell>
        </Grid>

        {elements.length ? <NCGroup
          onUpdate={(name, val, err) => {
            this.onUpdateElement(name, val, err);
          }}
          parent={this}
          element={element}
          prefix={this.currentPrefix()}
        /> : null}

        <CancelApplyPair
          cancelText="Close"
          applyText="Submit"
          onCancel={e => this.onNetworkCharacPaneClose(e)}
          onApply={e => this.triggerEvent(e)}
          saveDisabled={
            !elements.length || !element.elementType || !this.props.element.name || nbErrors
          }
          removeCyCancel={true}
        />
      </div>
    );
  }
}

const styles = {
  block: {
    marginBottom: 20
  },
  field: {
    marginBottom: 10
  },
  select: {
    width: '100%'
  }
};

export default NetworkCharacteristicsEventPane;
