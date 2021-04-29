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
import { Typography } from '@rmwc/typography';

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
      style={{ zIndex: 10000 }}
      open={props.open}
      onClose={(evt) => {
        var closeOnSubmit = (evt.detail.action === 'accept' ? true : false);
        props.onClose(closeOnSubmit);
      }}
      className={props.className ? props.className : ''}
      data-cy={props.cydata}
    >

      <DialogTitle theme="primary" style={styles.title}>
        <Typography use="headline5">{props.title}</Typography>
      </DialogTitle>

      <DialogContent style={styles.content}>
        {props.children}
      </DialogContent>

      <DialogActions style={styles.actions}>
        {props.onClose && (
          <DialogButton
            style={styles.button}
            action="close"
            onClick={props.onClose}
          >
            {(props.closeLabel === undefined) ? 'Cancel' : props.closeLabel}
          </DialogButton>
        )}
        {props.onSubmit && (
          <DialogButton
            style={styles.button}
            action="accept"
            isDefaultAction
            onClick={() => props.onSubmit()}
            disabled={props.okDisabled}
          >
            {(props.submitLabel === undefined) ? 'Ok' : props.submitLabel}
          </DialogButton>
        )}
      </DialogActions>
    </Dialog>
  );
};

const styles = {
  title: {
    paddingTop: 10,
    paddingBottom: 15
  },
  actions: {
    paddingBottom: 10,
    paddingRight: 24
  },
  button: {
    marginBottom: 5,
    marginLeft: 10
  }
};

export default IDDialog;
