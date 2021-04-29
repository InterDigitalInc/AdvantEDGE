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
import React, { Component } from 'react';
import { Select } from '@rmwc/select';
import { Grid, GridCell } from '@rmwc/grid';
import { Typography } from '@rmwc/typography';
import { updateObject } from '@/js/util/object-util';
import MobilityEventPane from './mobility-event-pane';
import NetworkCharacteristicsEventPane from './network-characteristics-event-pane';
import ScenarioUpdateEventPane from './scenario-update-event-pane';
import PduSessionEventPane from './pdu-session-event-pane';

import CancelApplyPair from '@/js/components/helper-components/cancel-apply-pair';

import {
  uiExecChangeCurrentEvent,
  uiExecChangeEventStatus
} from '@/js/state/ui';

import {
  execChangeSelectedScenarioElement,
  execResetSelectedScenarioElement,
  execUEs,
  execPOAs,
  execMobTypes,
  execEdges,
  execEdgeApps,
  execFogs,
  execFogEdges,
  execZones,
  execDNs
} from '@/js/state/exec';

import {
  MOBILITY_EVENT,
  NETWORK_CHARACTERISTICS_EVENT,
  SCENARIO_UPDATE_EVENT,
  PDU_SESSION_EVENT,
  EXEC_EVT_TYPE,
  PAGE_EXECUTE
} from '@/js/meep-constants';

const EventTypeSelect = props => {
  return (
    <Grid style={styles.field}>
      <GridCell span={12}>
        <Select
          style={styles.select}
          label="Event Type"
          fullwidth="true"
          outlined
          options={props.eventTypes}
          onChange={props.onChange}
          data-cy={EXEC_EVT_TYPE}
          value={props.value}
        />
      </GridCell>
    </Grid>
  );
};

const EventCreationFields = props => {
  switch (props.currentEvent) {
  case MOBILITY_EVENT:
    return (
      <MobilityEventPane
        element={props.element}
        eventTypes={props.eventTypes}
        api={props.api}
        onSuccess={props.onSuccess}
        onClose={props.onClose}
        currentEvent={props.currentEvent}
        UEs={props.UEs}
        POAs={props.POAs}
        EDGEs={props.EDGEs}
        FOGs={props.FOGs}
        ZONEs={props.ZONEs}
        MobTypes={props.MobTypes}
        FogEdges={props.FogEdges}
        EdgeApps={props.EdgeApps}
      />
    );
  case NETWORK_CHARACTERISTICS_EVENT:
    return (
      <NetworkCharacteristicsEventPane
        element={props.element}
        updateElement={props.updateElement}
        api={props.api}
        onSuccess={props.onSuccess}
        onClose={props.onClose}
        currentEvent={props.currentEvent}
        table={props.table}
        networkElements={props.networkElements}
      />
    );
  case SCENARIO_UPDATE_EVENT:
    return (
      <ScenarioUpdateEventPane
        currentEvent={props.currentEvent}
        onSuccess={props.onSuccess}
        onClose={props.onClose}
        api={props.api}
      />
    );
  case PDU_SESSION_EVENT:
    return (
      <PduSessionEventPane
        currentEvent={props.currentEvent}
        onSuccess={props.onSuccess}
        onClose={props.onClose}
        api={props.api}
        UEs={props.UEs}
        DNs={props.DNs}
      />
    );
  default:
    return <div></div>;
  }
};

class EventCreationPane extends Component {
  constructor(props) {
    super(props);
    this.state = {};

    if (!this.props.currentEvent) {
      this.props.changeEvent('');
    }
  }
  
  onEventPaneClose(e) {
    e.preventDefault();
    this.props.changeEvent('');
    this.props.changeEventStatus('');
    this.props.onClose(e);
  }

  updateElement(values) {
    if (values === null) {
      this.props.resetSelectedScenarioElement();
    } else {
      var element = updateObject({}, this.props.selectedScenarioElement);
      element = updateObject(element, values);
      this.props.changeSelectedScenarioElement(element);
    }
  }

  render() {
    if (this.props.page !== PAGE_EXECUTE || this.props.hide) {
      return null;
    }
    
    const statusColor = (this.props.eventStatus && this.props.eventStatus.startsWith('[200]')) ? 'green' : 'red';

    return (
      <div style={{ padding: 10 }}>
        <div style={styles.block}>
          <Typography use="headline6">Trigger Event</Typography>
        </div>
        <EventTypeSelect
          eventTypes={this.props.eventTypes}
          onChange={event => {
            this.props.changeEvent(event.target.value);
            this.props.changeEventStatus('');
          }}
          value={this.props.currentEvent}
        />
        <EventCreationFields
          element={this.props.selectedScenarioElement}
          currentEvent={this.props.currentEvent}
          api={this.props.api}
          updateElement={element => {
            this.updateElement(element);
          }}
          onSuccess={this.props.onSuccess}
          onClose={e => this.onEventPaneClose(e)}
          UEs={this.props.UEs}
          POAs={this.props.POAs}
          EDGEs={this.props.EDGEs}
          FOGs={this.props.FOGs}
          ZONEs={this.props.ZONEs}
          MobTypes={this.props.MobTypes}
          EdgeApps={this.props.EdgeApps}
          FogEdges={this.props.FogEdges}
          DNs={this.props.DNs}
          table={this.props.table}
          networkElements={this.props.networkElements}
        />

        <div hidden={this.props.currentEvent !== ''}>
          <CancelApplyPair
            cancelText="Close"
            applyText="Submit"
            onCancel={e => this.onEventPaneClose(e)}
            saveDisabled={true}
            removeCyApply={true}
          />
        </div>

        {
          (this.props.eventStatus) ?
          // <Grid style={{ marginTop: 20, border: '1px solid #e4e4e4' }}>
            <Grid style={{ marginTop: 20 }}>
              <GridCell span={12} style={{ padding: 5 }}>
                <Typography use="body1">Status:</Typography>
                <Typography use="body2" style={{ marginLeft: 5, color: statusColor}}>{this.props.eventStatus}</Typography>
              </GridCell>
            </Grid> : null
        }
      </div>
    );
  }
}

const styles = {
  block: {
    marginBottom: 20
  },
  field: {
    marginBottom: 20
  },
  select: {
    width: '100%'
  }
};

const mapStateToProps = state => {
  return {
    currentEvent: state.ui.execCurrentEvent,
    selectedScenarioElement: state.exec.selectedScenarioElement,
    page: state.ui.page,
    UEs: execUEs(state),
    POAs: execPOAs(state),
    EDGEs: execEdges(state),
    FOGs: execFogs(state),
    ZONEs: execZones(state),
    DNs: execDNs(state),
    MobTypes: execMobTypes(state),
    EdgeApps: execEdgeApps(state),
    FogEdges: execFogEdges(state),
    table: state.exec.table,
    networkElements: state.exec.table.entries,
    eventStatus: state.ui.eventStatus
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeEvent: event => dispatch(uiExecChangeCurrentEvent(event)),
    changeSelectedScenarioElement: element =>
      dispatch(execChangeSelectedScenarioElement(element)),
    resetSelectedScenarioElement: () => dispatch(execResetSelectedScenarioElement()),
    changeEventStatus: status => dispatch(uiExecChangeEventStatus(status))
  };
};

const ConnectedEventCreationPane = connect(
  mapStateToProps,
  mapDispatchToProps
)(EventCreationPane);

export default ConnectedEventCreationPane;
