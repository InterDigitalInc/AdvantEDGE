/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
import _ from 'lodash';
import React, { Component }  from 'react';
import { Select } from '@rmwc/select';
import { Grid, GridCell } from '@rmwc/grid';
import { updateObject } from '../../util/update';
import NCGroup from '../../components/helper-components/nc-group';
import CancelApplyPair from '../../components/helper-components/cancel-apply-pair';
import IDSelect from '../../components/helper-components/id-select';

import {
  camelCasePrefix,
  firstLetterUpper

} from '../../util/stringManipulation';

import {
  EXEC_EVT_NC_TYPE,
  EXEC_EVT_NC_NAME,

  // NC Group Prefixes
  PREFIX_INT_DOM,
  PREFIX_INT_ZONE,
  PREFIX_INT_EDGE,
  PREFIX_INT_FOG,
  PREFIX_EDGE_FOG,
  PREFIX_TERM_LINK,

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
  FIELD_INT_EDGE_LATENCY,
  FIELD_INT_EDGE_LATENCY_VAR,
  FIELD_INT_EDGE_THROUGPUT,
  FIELD_INT_EDGE_PKT_LOSS,
  FIELD_INT_FOG_LATENCY,
  FIELD_INT_FOG_LATENCY_VAR,
  FIELD_INT_FOG_THROUGPUT,
  FIELD_INT_FOG_PKT_LOSS,
  FIELD_EDGE_FOG_LATENCY,
  FIELD_EDGE_FOG_LATENCY_VAR,
  FIELD_EDGE_FOG_THROUGPUT,
  FIELD_EDGE_FOG_PKT_LOSS,
  FIELD_LINK_LATENCY,
  FIELD_LINK_LATENCY_VAR,
  FIELD_LINK_THROUGPUT,
  FIELD_LINK_PKT_LOSS,

  getElemFieldVal,
  setElemFieldVal,
  setElemFieldErr
} from '../../util/elem-utils';

const ncApplicableTypes = ['SCENARIO', 'DOMAIN', 'ZONE-INTER-EDGE', 'ZONE-INTER-FOG', 'ZONE-EDGE-FOG', 'POA'];

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
    var fieldsInError=0;
    _.forOwn(element, (val) => fieldsInError = val.err ? fieldsInError+1 : fieldsInError);
    if (fieldsInError) {
      return;
    }

    var neType = (this.state.currentElementType == 'DOMAIN') ? 'OPERATOR' : this.state.currentElementType;
    var ncEvent = {
      name: 'name',
      type: this.props.currentEvent,
      eventNetworkCharacteristicsUpdate: {
        elementName: getElemFieldVal(element, FIELD_NAME),
        elementType: neType,
      }
    };

    // Retrieve and set net characteristics from element
    this.setNetCharFromElem(ncEvent.eventNetworkCharacteristicsUpdate, element);

    // trigger event with this.props.api
    this.props.api.sendEvent(this.props.currentEvent, ncEvent, (error) => {
      if (!error) {
        this.setState({dialogOpen: true});
        this.props.onSuccess();
      }
    });
  }

  firstElementMatchingType(type) {
    var elements = _.chain(this.props.networkElements)
      .filter((e) => {
        var elemType = getElemFieldVal(e, FIELD_TYPE);
        if (type === 'DOMAIN' || type === 'OPERATOR') {
          return elemType === 'OPERATOR' || elemType === 'DOMAIN';
        }
        return type.startsWith(elemType);
      })
      .value();

    return elements.length ? elements[0] : null;
  }

  currentPrefix() {
    switch(this.state.currentElementType) {
    case 'SCENARIO':
      return PREFIX_INT_DOM;
    case 'DOMAIN':
      return PREFIX_INT_ZONE;
    case 'ZONE-INTER-EDGE':
      return PREFIX_INT_EDGE;
    case 'ZONE-INTER-FOG':
      return PREFIX_INT_FOG;
    case 'ZONE-EDGE-FOG':
      return PREFIX_EDGE_FOG;
    case 'POA':
      return PREFIX_TERM_LINK;
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
    case PREFIX_INT_EDGE:
      latencyFieldName = FIELD_INT_EDGE_LATENCY;
      latencyVarFieldName = FIELD_INT_EDGE_LATENCY_VAR;
      throughputFieldName = FIELD_INT_EDGE_THROUGPUT;
      packetLossFieldName = FIELD_INT_EDGE_PKT_LOSS;
      break;
    case PREFIX_INT_FOG:
      latencyFieldName = FIELD_INT_FOG_LATENCY;
      latencyVarFieldName = FIELD_INT_FOG_LATENCY_VAR;
      throughputFieldName = FIELD_INT_FOG_THROUGPUT;
      packetLossFieldName = FIELD_INT_FOG_PKT_LOSS;
      break;
    case PREFIX_EDGE_FOG:
      latencyFieldName = FIELD_EDGE_FOG_LATENCY;
      latencyVarFieldName = FIELD_EDGE_FOG_LATENCY_VAR;
      throughputFieldName = FIELD_EDGE_FOG_THROUGPUT;
      packetLossFieldName = FIELD_EDGE_FOG_PKT_LOSS;
      break;
    case PREFIX_TERM_LINK:
      latencyFieldName = FIELD_LINK_LATENCY;
      latencyVarFieldName = FIELD_LINK_LATENCY_VAR;
      throughputFieldName = FIELD_LINK_THROUGPUT;
      packetLossFieldName = FIELD_LINK_PKT_LOSS;
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
    var name = (camelCasePrefix(prefix) + firstLetterUpper(genericFieldName)).replace(/\s/g,'');
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
      if (getElemFieldVal(elements[i], FIELD_NAME) == name) {
        return elements[i];
      }
    }
    return null;
  }

  render() {
    var element = this.props.element;
    var type = getElemFieldVal(element, FIELD_TYPE);
    var nbErrors = _.reduce(element, (result, value) => {
      return (value.err) ? result = result + 1 : result;
    }, 0);

    var elements = _.chain(this.props.networkElements)
      .filter((e) => {
        var elemType = getElemFieldVal(e, FIELD_TYPE);
        if (type === 'DOMAIN' || type === 'OPERATOR') {
          return elemType === 'OPERATOR' || elemType === 'DOMAIN';
        }
        return this.state.currentElementType.startsWith(elemType);
      })
      .map((e) => {
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
              onChange={(event) => {
                var elem = this.firstElementMatchingType(event.target.value);
                this.props.updateElement(elem);
                this.setState({currentElementType: event.target.value});
              }}
              data-cy={EXEC_EVT_NC_TYPE}
            />
          </GridCell>
          <GridCell span="4">
          </GridCell>
        </Grid>

        <Grid>
          <GridCell span="8">
            <IDSelect
              span="8"
              label="Network Element"
              value={element ? getElemFieldVal(element, FIELD_NAME) || '' : ''}
              options={elements}
              onChange={(event)=>{
                this.props.updateElement(this.getElementByName(event.target.value));
              }}
              cydata={EXEC_EVT_NC_NAME}
            />
          </GridCell>
          <GridCell span="4">
          </GridCell>
        </Grid>

        <NCGroup
          onUpdate={(name, val, err) => {this.onUpdateElement(name, val, err);}}
          parent={this}
          element={element}
          prefix={this.currentPrefix()}
        />

        <CancelApplyPair
          cancelText="Close"
          applyText="Submit"
          onCancel={() => {this.props.onClose();}}
          onApply={(e) => this.triggerEvent(e)}
          saveDisabled={!element.elementType || !this.props.element.name || nbErrors}
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
