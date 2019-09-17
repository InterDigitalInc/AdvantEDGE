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
import React, { Component }  from 'react';
import { Select } from '@rmwc/select';
import { Grid, GridCell } from '@rmwc/grid';
import { Typography } from '@rmwc/typography';
import { updateObject } from '../../util/object-util';
import MobilityEventPane from './mobility-event-pane';
import NetworkCharacteristicsEventPane from './network-characteristics-event-pane';

import { uiExecChangeCurrentEvent, UE_MOBILITY_EVENT, NETWORK_CHARACTERISTICS_EVENT } from '../../state/ui';
import {
  execChangeSelectedScenarioElement,
  execUEs,
  execPOAs
} from '../../state/exec';

import {
  PAGE_EXECUTE
} from '../../state/ui';

import {
  EXEC_EVT_TYPE
} from '../../meep-constants';

const EventTypeSelect = (props) => {
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
        />
      </GridCell>
    </Grid>
  );
};

const EventCreationFields = (props) => {
  switch (props.currentEvent) {
  case UE_MOBILITY_EVENT:
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
  default:
    return (<div>NO EVENT</div>);
  }
};

class EventCreationPane extends Component {

  constructor(props) {
    super(props);
    this.state = {};
  }

  updateElement(values) {
    var element = updateObject({}, this.props.selectedScenarioElement);
    element = updateObject(element, values);
    this.props.changeSelectedScenarioElement(element);
  }

  render() {
    if (this.props.page !== PAGE_EXECUTE) { return null; }

    return (
      <div style={{padding: 10}}>
        <div style={styles.block}>
          <Typography use="headline6">Trigger Event</Typography>
        </div>
        <EventTypeSelect
          eventTypes={this.props.eventTypes}
          onChange={(event)=>{this.props.changeEvent(event.target.value);}}
        />
        <EventCreationFields
          element={this.props.selectedScenarioElement}
          currentEvent={this.props.currentEvent}
          api={this.props.api}
          updateElement={(element) => {this.updateElement(element);}}
          onSuccess={this.props.onSuccess}
          onClose={this.props.onClose}
          UEs={this.props.UEs}
          POAs={this.props.POAs}
          table={this.props.table}
          networkElements={this.props.networkElements}
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

const mapStateToProps = state => {
  return {
    currentEvent: state.ui.execCurrentEvent,
    selectedScenarioElement: state.exec.selectedScenarioElement,
    page: state.ui.page,
    UEs: execUEs(state),
    POAs: execPOAs(state),
    table: state.exec.table,
    networkElements: state.exec.table.entries
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeEvent: (event) => dispatch(uiExecChangeCurrentEvent(event)),
    changeSelectedScenarioElement: (element) => dispatch(execChangeSelectedScenarioElement(element))
  };
};

const ConnectedEventCreationPane = connect(
  mapStateToProps,
  mapDispatchToProps
)(EventCreationPane);

export default ConnectedEventCreationPane;
