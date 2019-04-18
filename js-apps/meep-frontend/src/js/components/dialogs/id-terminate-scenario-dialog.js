import React, { Component }  from 'react';

import IDDialog from './id-dialog';
import { MEEP_DLG_TERMINATE_SCENARIO } from '../../meep-constants';

class IDTerminateScenarioDialog extends Component {

  constructor(props) {
    super(props);
    this.state={};
  }

  render() {
    return (
      <IDDialog
        title='Terminate Scenario'
        open={this.props.open}
        onClose={this.props.onClose}
        onSubmit={() => this.props.onSubmit()}
        cydata={MEEP_DLG_TERMINATE_SCENARIO}
      >
        <span style={styles.text}>{'Are you sure you want to terminate the deployed scenario?'}</span>
      </IDDialog>
    );
  }
}

const styles = {
  text: {
    color: 'gray'
  }
};

export default IDTerminateScenarioDialog;