/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
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

const IDDialog = (props) => {
  return (
    <Dialog
      open={props.open}
      onClose={() => {
        props.onClose();
      }}
      data-cy={props.cydata}
    >
      <DialogTitle style={styles.title}>
        <Grid>
          <GridCell span={12}>
            {props.title}
          </GridCell>
        </Grid>
      </DialogTitle>

      <DialogContent style={styles.content}>
        <Grid>
          <GridCell span={12}>
            {props.children}
          </GridCell>
        </Grid>
      </DialogContent>

      <DialogActions style={styles.actions}>
        <Grid>
          <GridCell span={8}>
          </GridCell>
          <GridCell span={2}>
            <DialogButton style={styles.button}
              action="close"
              onClick={props.onClose}
            >
                                Cancel
            </DialogButton>
          </GridCell>

          <GridCell span={2}>
            <DialogButton style={styles.button}
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
    paddingLeft:25,
    paddingRight:25
  },
  content: {
    paddingLeft:25,
    paddingRight:30
  },
  actions: {
    marginTop: 20
  },
  button: {
    margin: 5
  }
};

export default IDDialog;
