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

import React from 'react';
import { Grid, GridCell } from '@rmwc/grid';

import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  DialogButton
} from '@rmwc/dialog';

const IDDialog = props => {
  return (
    <Dialog
      style={{zIndex: 10000}}
      open={props.open}
      onClose={() => {
        props.onClose();
      }}
      data-cy={props.cydata}
    >
      <DialogTitle style={styles.title}>
        <Grid>
          <GridCell span={12}>{props.title}</GridCell>
        </Grid>
      </DialogTitle>

      <DialogContent style={styles.content}>
        <Grid>
          <GridCell span={12}>{props.children}</GridCell>
        </Grid>
      </DialogContent>

      <DialogActions style={styles.actions}>
        <Grid>
          <GridCell span={8}></GridCell>
          <GridCell span={2}>
            <DialogButton
              style={styles.button}
              action="close"
              onClick={props.onClose}
            >
              Cancel
            </DialogButton>
          </GridCell>

          <GridCell span={2}>
            <DialogButton
              style={styles.button}
              action="accept"
              isDefaultAction
              onClick={() => props.onSubmit()}
              disabled={props.okDisabled}
            >
              Ok
            </DialogButton>
          </GridCell>
        </Grid>
      </DialogActions>
    </Dialog>
  );
};

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
