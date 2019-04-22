/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
import React from 'react';
import { Select } from '@rmwc/select';
import { GridCell } from '@rmwc/grid';


const IDSelect = (props) => {
  return (
    <GridCell span={props.span}>
      <Select
        style={{width: '100%'}}
        label={props.label}
        outlined
        options={props.options}
        onChange={props.onChange}
        disabled={props.disabled}
        value={props.value}
        data-cy={props.cydata}
      />
    </GridCell>
  );
};

export default IDSelect;
