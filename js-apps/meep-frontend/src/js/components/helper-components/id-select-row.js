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
import IDSelect from './id-select';

const IDSelectRow = (props) => {
  return (
    <Grid style={{marginBottom: 10}}>
      <IDSelect {...props} />
      <GridCell span={12 - props.span} />
    </Grid>
  );
};

export default IDSelectRow;
