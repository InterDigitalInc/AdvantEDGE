import _ from 'lodash';
import React, { Component }  from 'react';
import { Select } from '@rmwc/select';
import { Grid, GridCell } from '@rmwc/grid';
import CancelApplyPair from '../../components/helper-components/cancel-apply-pair';

import {
  EXEC_EVT_MOB_TARGET,
  EXEC_EVT_MOB_DEST
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

  triggerEvent(e) {
    e.preventDefault();
    var meepEvent = {
      'name': 'name',
      'type': this.props.currentEvent,
      'eventUeMobility': {
        'ue': this.values.eventTarget,
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

    return (
      <div>
        <Grid style={styles.field}>
          <GridCell span="8">
            <Select
              style={styles.select}
              label="Target"
              outlined
              options={_.map(this.props.UEs, elem => getElemFieldVal(elem, FIELD_NAME))}
              onChange={(event)=>{this.values['eventTarget'] = event.target.value;}}
              data-cy={EXEC_EVT_MOB_TARGET}
            />
          </GridCell> 
          <GridCell span="4">
          </GridCell>
        </Grid>
        <Grid style={styles.block}>
          <GridCell span="8">
            <Select
              style= {styles.select}
              label="Destination"
              outlined
              options={_.map(this.props.POAs, elem => getElemFieldVal(elem, FIELD_NAME))}
              onChange={(event)=>{this.values['eventDestination'] = event.target.value;}}
              data-cy={EXEC_EVT_MOB_DEST}
            />
          </GridCell>
          <GridCell span="4">
          </GridCell>   
        </Grid>

        <CancelApplyPair
          cancelText="Close"
          applyText="Submit"
          onCancel={this.props.onClose}
          onApply={(e) => this.triggerEvent(e)}
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

export default MobilityEventPane;


