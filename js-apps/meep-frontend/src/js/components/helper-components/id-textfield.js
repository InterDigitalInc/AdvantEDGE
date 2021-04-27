/*
 * Copyright (c) 2021  InterDigital Communications, Inc
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
import { TextField, TextFieldIcon, TextFieldHelperText } from '@rmwc/textfield';
import { GridCell } from '@rmwc/grid';

import {
  getElemFieldErr
} from '../../util/elem-utils';

// COMPONENTS
export const IDTextField = props => {
  var err = props.element[props.fieldName]
    ? props.element[props.fieldName].err
    : null;
  return (
    <>
      <TextField
        outlined
        style={{ width: '100%', marginBottom: 0 }}
        label={props.label}
        withLeadingIcon={!props.icon ? null : 
          <TextFieldIcon
            tabIndex="0"
            icon={props.icon}
            onClick={props.onIconClick}
          />
        }
        type={props.type}
        onChange={event => {
          var err = props.validate ? props.validate(event.target.value) : null;
          var val =
            event.target.value && props.isNumber && !err
              ? Number(event.target.value)
              : event.target.value;
          props.onUpdate(props.fieldName, val, err);
        }}
        invalid={err}
        value={
          props.element[props.fieldName]
            ? props.element[props.fieldName].val
            : ''
        }
        disabled={props.disabled}
        data-cy={props.cydata}
      />
      <TextFieldHelperText validationMsg={true}>
        <span>{getElemFieldErr(props.element, props.fieldName)}</span>
      </TextFieldHelperText>
    </>
  );
};

export const IDTextFieldCell = props => {
  return (
    <GridCell span={props.span}>
      <IDTextField {...props} />
    </GridCell>
  );
};
