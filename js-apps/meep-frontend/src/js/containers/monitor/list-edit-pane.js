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
import _ from 'lodash';
import { Grid, GridCell } from '@rmwc/grid';
import { TextField } from '@rmwc/textfield';
import { Checkbox } from '@rmwc/checkbox';
import { Button } from '@rmwc/button';

export const ListEditPaneRow = ({item, itemLabelLabel, itemValueLabel, updateItemSelection, updateItemValue, updateItemLabel}) => {
  return (
    <Grid>
      <GridCell span={4} style={styles.editListItemCell}>
        <TextField outlined style={{width: '100%'}}
          label={itemLabelLabel}
          value={item.label}
          onChange={(e) => {
            updateItemLabel(item.index, e.target.value);
          }}
        />
      </GridCell>
      <GridCell span={7} style={styles.editListItemCell}>
        <TextField outlined style={{width: '100%'}}
          label={itemValueLabel}
          value={item.value}
          onChange={(e) => {
            updateItemValue(item.index, e.target.value);
          }}
        />
      </GridCell>
      <GridCell span={1} style={{...styles.editListItemCell, paddingTop: 30}}>
        <Checkbox
          checked={item.selected}
          onChange={(e) => {
            updateItemSelection(item.index, e.target.checked);
          }}
        />
      </GridCell>
    </Grid>
  );
};
  
export const ListEditPane = (props) => {
  return (
    <div>
      {_.map(props.items, (item, index) => {
        return (<ListEditPaneRow
          item={item}
          key={index}
          itemLabelLabel={props.itemLabelLabel}
          itemValueLabel={props.itemValueLabel}
          updateItemLabel={props.updateItemLabel}
          updateItemValue={props.updateItemValue}
          updateItemSelection={props.updateItemSelection}
        />);
      })
      }

      <Grid style={{marginTop: 20, marginBottom: 10}}>
        <GridCell span={7}>

        </GridCell>

        <GridCell span={5}>
          <Button raised
            style={styles.button}
            onClick={props.cancelEditMode}
          >
              CANCEL
          </Button>
          <Button raised
            style={styles.button}
            onClick={props.deleteItems}
            disabled={!props.canDelete()}
          >
              DELETE
          </Button>
          <Button raised
            style={styles.button}
            onClick={props.addItem}
          >
              ADD
          </Button>
          <Button raised
            style={styles.button}
            onClick={props.saveItems}
          >
              SAVE
          </Button>
        </GridCell>
      </Grid>
    </div>
  );
};

const styles = {
  button: {
    color: 'white',
    marginRight: 5
  },
  editListItemCell: {
    padding: 5
  }
};