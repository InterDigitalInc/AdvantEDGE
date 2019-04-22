/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
 import React from 'react';
import { Grid, GridInner, GridCell } from '@rmwc/grid';
import { Button } from '@rmwc/button';
import {
  MEEP_BTN_CANCEL,
  MEEP_BTN_APPLY
} from '../../meep-constants';

const buttonStyles = {
  marginRight: 10
};

const CancelApplyPair = (props) => {
  return (
    <Grid style={{marginTop: 10}}>
      <GridInner align={'right'}>
        <GridCell span={12}>
          <Button outlined style={buttonStyles} onClick={props.onCancel} data-cy={MEEP_BTN_CANCEL}>
            {props.cancelText ? props.cancelText : 'Cancel'}
          </Button>
          <Button outlined style={buttonStyles} onClick={props.onApply} disabled={props.saveDisabled} data-cy={MEEP_BTN_APPLY}>
            {props.applyText ? props.applyText : 'Apply'}
          </Button>
        </GridCell>
      </GridInner>
    </Grid>
  );
};

export default CancelApplyPair;
