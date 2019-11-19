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
  ELEMENT_TYPE_ZONE,
  ELEMENT_TYPE_POA,
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

  // Layout type
  MEEP_COMPONENT_SINGLE_COLUMN_LAYOUT
} from '../../meep-constants';

import {
  // Field Names
  FIELD_NAME,
  FIELD_TYPE,
  FIELD_INT_DOM_LATENCY,
  FIELD_INT_DOM_LATENCY_VAR,
  FIELD_INT_DOM_THROUGPUT,
  FIELD_INT_DOM_PKT_LOSS,
  FIELD_INT_ZONE_LATENCY,
  FIELD_INT_ZONE_LATENCY_VAR,
  FIELD_INT_ZONE_THROUGPUT,
  FIELD_INT_ZONE_PKT_LOSS,
  FIELD_INTRA_ZONE_LATENCY,
  FIELD_INTRA_ZONE_LATENCY_VAR,
  FIELD_INTRA_ZONE_THROUGPUT,
  FIELD_INTRA_ZONE_PKT_LOSS,
  FIELD_TERM_LINK_LATENCY,
  FIELD_TERM_LINK_LATENCY_VAR,
  FIELD_TERM_LINK_THROUGPUT,
  FIELD_TERM_LINK_PKT_LOSS,
  FIELD_LINK_LATENCY,
  FIELD_LINK_LATENCY_VAR,
  FIELD_LINK_THROUGPUT,
  FIELD_LINK_PKT_LOSS,
  FIELD_APP_LATENCY,
  FIELD_APP_LATENCY_VAR,
  FIELD_APP_THROUGPUT,
  FIELD_APP_PKT_LOSS,
  getElemFieldVal,
  setElemFieldVal,
  setElemFieldErr
} from '../../util/elem-utils';

const ncApplicableTypes = [
  ELEMENT_TYPE_SCENARIO,
  ELEMENT_TYPE_OPERATOR,
  ELEMENT_TYPE_ZONE,
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
      currentElementType: ''
    };
  }

  triggerEvent(e) {
    e.preventDefault();
    var element = this.props.element;

    // Verify that no field is in error
    var fieldsInError = 0;
    _.forOwn(
      element,
      val => (fieldsInError = val.err ? fieldsInError + 1 : fieldsInError)
    );
    if (fieldsInError) {
      return;
    }

    var neType =
      this.state.currentElementType === 'DOMAIN'
        ? 'OPERATOR'
        : this.state.currentElementType;
    var ncEvent = {
      name: 'name',
      type: this.props.currentEvent,
      eventNetworkCharacteristicsUpdate: {
        elementName: getElemFieldVal(element, FIELD_NAME),
        elementType: neType
      }
    };

    // Retrieve and set net characteristics from element
    this.setNetCharFromElem(ncEvent.eventNetworkCharacteristicsUpdate, element);

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
        if (type === 'DOMAIN' || type === 'OPERATOR') {
          return elemType === 'OPERATOR' || elemType === 'DOMAIN';
        }
        if (elemType === 'ZONE') {
          return type.startsWith(elemType);
        }
        return type === elemType;
      })
      .value();

    return elements.length ? elements[0] : null;
  }

  currentPrefix() {
    switch (this.state.currentElementType) {
    case ELEMENT_TYPE_SCENARIO:
      return PREFIX_INT_DOM;
    case ELEMENT_TYPE_OPERATOR:
      return PREFIX_INT_ZONE;
    case ELEMENT_TYPE_ZONE:
      return PREFIX_INTRA_ZONE;
    case ELEMENT_TYPE_POA:
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
    var throughputFieldName = null;
    var packetLossFieldName = null;
    switch (this.currentPrefix()) {
    case PREFIX_INT_DOM:
      latencyFieldName = FIELD_INT_DOM_LATENCY;
      latencyVarFieldName = FIELD_INT_DOM_LATENCY_VAR;
      throughputFieldName = FIELD_INT_DOM_THROUGPUT;
      packetLossFieldName = FIELD_INT_DOM_PKT_LOSS;
      break;
    case PREFIX_INT_ZONE:
      latencyFieldName = FIELD_INT_ZONE_LATENCY;
      latencyVarFieldName = FIELD_INT_ZONE_LATENCY_VAR;
      throughputFieldName = FIELD_INT_ZONE_THROUGPUT;
      packetLossFieldName = FIELD_INT_ZONE_PKT_LOSS;
      break;
    case PREFIX_INTRA_ZONE:
      latencyFieldName = FIELD_INTRA_ZONE_LATENCY;
      latencyVarFieldName = FIELD_INTRA_ZONE_LATENCY_VAR;
      throughputFieldName = FIELD_INTRA_ZONE_THROUGPUT;
      packetLossFieldName = FIELD_INTRA_ZONE_PKT_LOSS;
      break;
    case PREFIX_TERM_LINK:
      latencyFieldName = FIELD_TERM_LINK_LATENCY;
      latencyVarFieldName = FIELD_TERM_LINK_LATENCY_VAR;
      throughputFieldName = FIELD_TERM_LINK_THROUGPUT;
      packetLossFieldName = FIELD_TERM_LINK_PKT_LOSS;
      break;
    case PREFIX_LINK:
      latencyFieldName = FIELD_LINK_LATENCY;
      latencyVarFieldName = FIELD_LINK_LATENCY_VAR;
      throughputFieldName = FIELD_LINK_THROUGPUT;
      packetLossFieldName = FIELD_LINK_PKT_LOSS;
      break;
    case PREFIX_APP:
      latencyFieldName = FIELD_APP_LATENCY;
      latencyVarFieldName = FIELD_APP_LATENCY_VAR;
      throughputFieldName = FIELD_APP_THROUGPUT;
      packetLossFieldName = FIELD_APP_PKT_LOSS;
      break;
    default:
      return null;
    }

    // Update net characteristics
    netChar.latency = getElemFieldVal(element, latencyFieldName);
    netChar.latencyVariation = getElemFieldVal(element, latencyVarFieldName);
    netChar.throughput = getElemFieldVal(element, throughputFieldName);
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
    for (var i = 0; i < elements.length; i++) {
      if (getElemFieldVal(elements[i], FIELD_NAME) === name) {
        return elements[i];
      }
    }
    return null;
  }

  render() {
    var element = this.props.element;
    var type = getElemFieldVal(element, FIELD_TYPE);
    var nbErrors = _.reduce(
      element,
      (result, value) => {
        return value.err ? (result = result + 1) : result;
      },
      0
    );

    var elements = _.chain(this.props.networkElements)
      .filter(e => {
        var elemType = getElemFieldVal(e, FIELD_TYPE);
        if (type === 'DOMAIN' || type === 'OPERATOR') {
          return elemType === 'OPERATOR' || elemType === 'DOMAIN';
        }
        if (elemType === 'ZONE') {
          return this.state.currentElementType.startsWith(elemType);
        }
        return this.state.currentElementType === elemType;
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
              options={ncApplicableTypes}
              onChange={event => {
                var elem = this.firstElementMatchingType(event.target.value);
                this.props.updateElement(elem);
                this.setState({ currentElementType: event.target.value });
              }}
              data-cy={EXEC_EVT_NC_TYPE}
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

        <NCGroup
          layout={MEEP_COMPONENT_SINGLE_COLUMN_LAYOUT}
          onUpdate={(name, val, err) => {
            this.onUpdateElement(name, val, err);
          }}
          parent={this}
          element={element}
          prefix={this.currentPrefix()}
        />

        <CancelApplyPair
          cancelText="Close"
          applyText="Submit"
          onCancel={() => {
            this.props.onClose();
          }}
          onApply={e => this.triggerEvent(e)}
          saveDisabled={
            !element.elementType || !this.props.element.name || nbErrors
          }
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
