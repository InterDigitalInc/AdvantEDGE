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