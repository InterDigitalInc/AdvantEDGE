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

import _ from 'lodash';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import autoBind from 'react-autobind';
import { Select } from '@rmwc/select';
import { Grid, GridCell } from '@rmwc/grid';

import CancelApplyPair from '@/js/components/helper-components/cancel-apply-pair';
import { IDTextFieldCell } from '@/js/components/helper-components/id-textfield';

import {
  PDU_SESSION_ACTION_ADD,
  PDU_SESSION_ACTION_REMOVE,
  EXEC_EVT_PDU_SESSION_ACTION,
  EXEC_EVT_PDU_SESSION_UE,
  EXEC_EVT_PDU_SESSION_ID,
  EXEC_EVT_PDU_SESSION_DNN
} from '@/js/meep-constants';

import {
  uiExecChangePduSessionEvent,
  uiExecChangeEventStatus
} from '@/js/state/ui';

import {
  validateVariableName
} from '@/js/util/validate';

import {
  FIELD_NAME,
  FIELD_DN_NAME,
  setElemFieldVal,
  setElemFieldErr,
  getElemFieldVal,
  validElem
} from '@/js/util/elem-utils';

import {
  updateObject
} from '@/js/util/object-util';

const FIELD_PDU_SESSION_ACTION =  'action';
const FIELD_PDU_SESSION_UE =  'ue';
const FIELD_PDU_SESSION_ID =  'id';
const FIELD_PDU_SESSION_DNN =  'dnn';

const actionTypes = [
  PDU_SESSION_ACTION_ADD,
  PDU_SESSION_ACTION_REMOVE
];

class PduSessionEventPane extends Component {
  constructor(props) {
    super(props);
    autoBind(this);
  }

  onClose(e) {
    e.preventDefault();
    this.props.changePduSessionEvent({});
    this.props.onClose(e);
  }

  // Event update handler
  onUpdateEvent(name, val, err) {
    var updatedEvent = updateObject({}, this.props.pduSessionEvent);
    setElemFieldVal(updatedEvent, name, val);
    setElemFieldErr(updatedEvent, name, err);
    this.props.changePduSessionEvent(updatedEvent);
  }

  /**
   * Callback function to receive the result of the sendEvent operation.
   * @callback module:api/EventsApi~sendEventCallback
   * @param {String} error Error message, if any.
   * @param data This operation does not return a value.
   * @param {String} response The complete HTTP response.
   */
  sendEventCb(error, _, response) {
    var status = '';
    if (error) {
      status = '[' + response.statusCode + '] ' + response.statusText + ': ' + response.text;
      this.props.changeEventStatus(status);
      return;
    }

    status = '[' + response.statusCode + '] ' + response.statusText;
    this.props.changeEventStatus(status);
    this.props.onSuccess();
  }

  sendEvent() {
    if (!validElem(this.props.pduSessionEvent)) {
      this.props.changeEventStatus('Error in params');
      return;
    }

    var pduSessionEvent = {
      name: 'name',
      type: this.props.currentEvent,
      eventPduSession: {
        action: getElemFieldVal(this.props.pduSessionEvent, FIELD_PDU_SESSION_ACTION),
        pduSession: {
          ue: getElemFieldVal(this.props.pduSessionEvent, FIELD_PDU_SESSION_UE),
          id: getElemFieldVal(this.props.pduSessionEvent, FIELD_PDU_SESSION_ID),
          info: {
            dnn: getElemFieldVal(this.props.pduSessionEvent, FIELD_PDU_SESSION_DNN)
          }
        }
      }
    };

    this.props.api.sendEvent(this.props.currentEvent, pduSessionEvent, (error, data, response) => {
      this.sendEventCb(error, data, response);
    });
  }

  render() {
    var action = getElemFieldVal(this.props.pduSessionEvent, FIELD_PDU_SESSION_ACTION) || '';
    var ueList = _.map(this.props.UEs, elem => getElemFieldVal(elem, FIELD_NAME));
    var dnnList = _.uniq(_.map(this.props.DNs, elem => getElemFieldVal(elem, FIELD_DN_NAME)));

    return (
      <div>
        <Grid>
          <GridCell span={6}>
            <Select
              style={styles.select}
              label='Action Type'
              outlined
              options={actionTypes}
              onChange={e => { this.onUpdateEvent(FIELD_PDU_SESSION_ACTION, e.target.value, null); }}
              value={action}
              data-cy={EXEC_EVT_PDU_SESSION_ACTION}
            />
          </GridCell>
          <GridCell span={6}/>

          { action === PDU_SESSION_ACTION_ADD || action === PDU_SESSION_ACTION_REMOVE ?
            <>
              <GridCell span={8}>
                <Select
                  style={styles.select}
                  label='Terminal'
                  outlined
                  options={ueList}
                  onChange={e => { this.onUpdateEvent(FIELD_PDU_SESSION_UE, e.target.value, null); }}
                  value={getElemFieldVal(this.props.pduSessionEvent, FIELD_PDU_SESSION_UE) || ''}
                  data-cy={EXEC_EVT_PDU_SESSION_UE}
                />
              </GridCell>
              <GridCell span={4}/>
                <IDTextFieldCell
                  span={8}
                  onUpdate={this.onUpdateEvent}
                  element={this.props.pduSessionEvent}
                  validate={validateVariableName}
                  label='PDU Session ID'
                  fieldName={FIELD_PDU_SESSION_ID}
                  cydata={EXEC_EVT_PDU_SESSION_ID}
                />
              <GridCell span={4}/>
            </> : null
          }

          { action === PDU_SESSION_ACTION_ADD ?
            <>
              <GridCell span={8}>
                <Select
                  style={styles.select}
                  label='Data Network Name'
                  outlined
                  options={dnnList}
                  onChange={e => { this.onUpdateEvent(FIELD_PDU_SESSION_DNN, e.target.value, null); }}
                  value={getElemFieldVal(this.props.pduSessionEvent, FIELD_PDU_SESSION_DNN) || ''}
                  data-cy={EXEC_EVT_PDU_SESSION_DNN}
                />
              </GridCell>
              <GridCell span={4}/>
            </> : null
          }
        </Grid>

        <CancelApplyPair
          cancelText="Close"
          applyText="Submit"
          onCancel={e => this.onClose(e)}
          onApply={() => this.sendEvent()}
          removeCyCancel={true}
        />
      </div>
    );
  }
}

const styles = {
  field: {
    marginBottom: 10,
    width: '100%'
  },
  select: {
    width: '100%',
    marginBottom: 10
  },
  page: {
    height: '100%',
    marginBottom: 10,
    width: '100%'
  }
};

const mapStateToProps = state => {
  return {
    pduSessionEvent: state.ui.pduSessionEvent
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changePduSessionEvent: event => dispatch(uiExecChangePduSessionEvent(event)),
    changeEventStatus: status => dispatch(uiExecChangeEventStatus(status))
  };
};

const ConnectedPduSessionEventPane = connect(
  mapStateToProps,
  mapDispatchToProps
)(PduSessionEventPane);

export default ConnectedPduSessionEventPane;
