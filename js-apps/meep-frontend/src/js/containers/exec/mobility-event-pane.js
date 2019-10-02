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
import React, { Component }  from 'react';
import { Select } from '@rmwc/select';
import { Grid, GridCell } from '@rmwc/grid';
import CancelApplyPair from '../../components/helper-components/cancel-apply-pair';

import {
  EXEC_EVT_MOB_TARGET,
  EXEC_EVT_MOB_DEST,
  EVENT_CREATION_PANE_TABLE_LAYOUT,
  EVENT_CREATION_PANE_LINE_LAYOUT
} from '../../meep-constants';

import {
  getElemFieldVal,
  FIELD_NAME
} from '../../util/elem-utils';

class MobilityEventPane extends Component {

  constructor(props) {
    super(props);
    this.values = {
      eventType: '',
      eventTarget: '',
      eventDestination: ''
    };

    this.state = {
    };
  }

  shouldComponentUpdate(nextProps) {

    /**
     * element={props.element}
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
     */
    return this.props.api !== nextProps.api
      || this.props.element !== nextProps.element
      || this.props.api !== nextProps.api
      || this.props.currentEvent !== nextProps.currentEvent
      || this.props.UEs !== nextProps.UEs
      || this.props.POAs !== nextProps.POAs
      || this.props.EDGEs !== nextProps.EDGEs
      || this.props.FOGs !== nextProps.FOGs
      || this.props.ZONEs !== nextProps.ZONEs
      || this.props.MobTypes !== nextProps.MobTypes
      || this.props.FogEdges !== nextProps.FogEdges
      || this.props.EdgeApps !== nextProps.EdgeApps;
      
  }

  triggerEvent(e) {
    e.preventDefault();
    var meepEvent = {
      'name': 'name',
      'type': this.props.currentEvent,
      'eventMobility': {
        'elementName': this.values.eventTarget,
        'dest': this.values.eventDestination
      }
    };

    // trigger event with this.props.api
    this.props.api.sendEvent(this.props.currentEvent, meepEvent, (error) => {
      if (!error) {
        this.props.onSuccess();
      }
    });
  }

  render() {

    //let found = this.props.UEs.find(element => element.label == this.values.eventTarget);
    //find if its the selection was a UE, otherwise (in order) EDGE, FOG, EDGE-APP, UE-APP
    var target = this.values.eventTarget;
    var found = this.props.UEs.find(function(element) {
      return element.label === target;
    });

    var populateDestination;
    if (found !== undefined) {
      populateDestination = this.props.POAs;
    } else {
      found = this.props.EDGEs.find(function(element) {
        return element.label === target;
      });

      if (found !== undefined) {
        populateDestination = this.props.ZONEs;
      } else {
        found = this.props.FOGs.find(function(element) {
          return element.label === target;
        });
        if (found !== undefined) {
          populateDestination = this.props.POAs;
        } else {
          found = this.props.EdgeApps.find(function(element) {
            return element.label === target;
          });
          if (found !== undefined) {
            populateDestination = this.props.FogEdges;
          }
        }
      }
    }

    const TargetSelectComponent = () => (
      <>
      <Select
        style={styles.select}
        label="Target"
        outlined
        options={_.map(this.props.MobTypes, elem => getElemFieldVal(elem, FIELD_NAME))}
        onChange={(event)=>{this.values['eventTarget'] = event.target.value;}}
        data-cy={EXEC_EVT_MOB_TARGET}
      />
      </>
    );

    const DestinationSelectComponent = () => (
      <>
      <Select
        style= {styles.select}
        label="Destination"
        outlined
        options={_.map(populateDestination, elem => getElemFieldVal(elem, FIELD_NAME))}
        onChange={(event)=>{this.values['eventDestination'] = event.target.value;}}
        data-cy={EXEC_EVT_MOB_DEST}
      />
      </>
    );

    const CancelApplyPairComponent = () => (
      <>
      <CancelApplyPair
        cancelText="Close"
        applyText="Submit"
        onCancel={this.props.onClose}
        onApply={(e) => this.triggerEvent(e)}
      />
      </>
    );

    let Layout = () => (
      <>
      <Grid style={styles.field}>
        <GridCell span="8">
          <TargetSelectComponent />
        </GridCell>
        <GridCell span="4">
        </GridCell>
      </Grid>
      <Grid style={styles.block}>
        <GridCell span="8">
          <DestinationSelectComponent />
        </GridCell>
        <GridCell span="4">
        </GridCell>
      </Grid>
      <CancelApplyPairComponent />
      </>
    );
  
    return (
      <div>
        <Layout />
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

export default MobilityEventPane;
