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

import React, { Component } from 'react';
import { Grid, GridCell } from '@rmwc/grid';

import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  DialogButton
} from '@rmwc/dialog';

class IDDialog extends Component {
  constructor(props) {
    super(props);

    this.state = {
      sandboxName: '',
      err: null
      
    };
  }

  render() {
    return (
      <Dialog
        style={{zIndex: 10000}}
        open={this.props.open}
        onClose={this.props.onClose}
        data-cy={this.props.cydata}
      >
        <DialogTitle style={styles.title}>
          <Grid>
            <GridCell span={12}>{this.props.title}</GridCell>
          </Grid>
        </DialogTitle>

        <DialogContent style={styles.content}>
          <Grid>
            <GridCell span={12}>{this.props.children}</GridCell>
          </Grid>
        </DialogContent>

        <DialogActions style={styles.actions}>
          <Grid>
            <GridCell span={8}></GridCell>
            <GridCell span={2}>
              <DialogButton
                style={styles.button}
                action="close"
                onClick={this.props.onClose}
              >
                Cancel
              </DialogButton>
            </GridCell>

            <GridCell span={2}>
              <DialogButton
                style={styles.button}
                action="accept"
                isDefaultAction
                onClick={this.props.onSubmit}
                disabled={this.props.okDisabled}
              >
                Ok
              </DialogButton>
            </GridCell>
          </Grid>
        </DialogActions>
      </Dialog>
    );
  }
}

const styles = {
  title: {
    paddingLeft: 25,
    paddingRight: 25
  },
  content: {
    paddingLeft: 25,
    paddingRight: 30
  },
  actions: {
    marginTop: 20
  },
  button: {
    margin: 5
  }
};

export default IDDialog;
