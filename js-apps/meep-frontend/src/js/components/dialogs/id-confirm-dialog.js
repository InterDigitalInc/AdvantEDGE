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
import { MEEP_DLG_CONFIRM } from '../../meep-constants';

class IDConfirmDialog extends Component {

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
        onSubmit={() => this.props.onSubmit()}
        cydata={MEEP_DLG_CONFIRM}
      >
        <span style={styles.text}>{`Are you sure you want to ${this.props.title.toLowerCase()}?`}</span>
      </IDDialog>
    );
  }
}

const styles = {
  text: {
    color: 'gray'
  }
};

export default IDConfirmDialog;
