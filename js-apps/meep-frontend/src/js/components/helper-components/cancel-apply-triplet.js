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
import { Grid, GridInner, GridCell } from '@rmwc/grid';
import { Button } from '@rmwc/button';
import { MEEP_BTN_CANCEL, MEEP_BTN_APPLY, MEEP_BTN_DUPLICATE } from '../../meep-constants';

const buttonStyles = {
  marginRight: 10
};

const CancelApplyTriplet = props => {
  return (
    <Grid style={{ marginTop: 10 }}>
      <GridInner align={'right'}>
        <GridCell span={12}>
          <Button
            outlined
            style={buttonStyles}
            onClick={props.onCancel}
            data-cy={MEEP_BTN_CANCEL}
          >
            {props.cancelText ? props.cancelText : 'Cancel'}
          </Button>
          <Button
            outlined
            style={buttonStyles}
            onClick={props.onApply}
            disabled={props.saveDisabled}
            data-cy={MEEP_BTN_APPLY}
          >
            {props.applyText ? props.applyText : 'Apply'}
          </Button>
          <Button
            outlined
            style={buttonStyles}
            onClick={props.onDuplicate}
            disabled={props.duplicateDisabled}
            data-cy={MEEP_BTN_DUPLICATE}
          >
            {props.duplicateText ? props.duplicateText : 'Duplicate'}
          </Button>
        </GridCell>
      </GridInner>
    </Grid>
  );
};

export default CancelApplyTriplet;
