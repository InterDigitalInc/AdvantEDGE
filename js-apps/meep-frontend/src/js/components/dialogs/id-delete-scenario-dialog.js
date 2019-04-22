/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
import React, { Component }  from 'react';

import IDDialog from './id-dialog';
import { MEEP_DLG_DEL_SCENARIO } from '../../meep-constants';

class IDDeleteScenarioDialog extends Component {

  constructor(props) {
    super(props);
    this.state={};
  }

  render() {
    return (
      <IDDialog
        title={this.props.title}
        open={this.props.open}
        onClose={this.props.onClose}
        onSubmit={this.props.deleteScenario}
        cydata={MEEP_DLG_DEL_SCENARIO}
      >
        <span style={styles.text}>{'Are you sure you want to delete the current scenario from the MEEP Controller?'}</span>
      </IDDialog>
    );
  }
}

const styles = {
  text: {
    color: 'gray'
  }
};

export default IDDeleteScenarioDialog;
