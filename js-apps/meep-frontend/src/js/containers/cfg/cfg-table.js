/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
import PropTypes from 'prop-types';
import { Elevation } from '@rmwc/elevation';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TablePagination from '@material-ui/core/TablePagination';
import TableRow from '@material-ui/core/TableRow';
import TableSortLabel from '@material-ui/core/TableSortLabel';
import Paper from '@material-ui/core/Paper';
import Tooltip from '@material-ui/core/Tooltip';
import { Grid, GridCell } from '@rmwc/grid';
import { withStyles } from '@material-ui/core/styles';
import { connect } from 'react-redux';
import React, { Component }  from 'react';

import { updateObject } from '../../util/object-util';

import {
  cfgChangeTable
} from '../../state/cfg';

import {
  getSortingByField,
  handleRequestSort,
  handleChangePage,
  handleChangeRowsPerPage,
  isRowSelected
} from '../../util/table-utils';

import {
  FIELD_TYPE,
  FIELD_NAME,
  FIELD_PARENT,
  getElemFieldVal
} from '../../util/elem-utils';

// Network Element Cfg Styles & Table
const cfgTableStyles = theme => ({
  root: {
    width: '100%',
    marginTop: theme.spacing.unit * 3
  },
  table: {
    minWidth: 1020
  },
  tableWrapper: {
    overflowX: 'auto'
  },
  tableHead: {
    'background-color': '#379DD8'
  },
  tableHeadColor: {
    color: '#FFFFFF'
  }
});

const cfgTableColumnData = [
  { id: FIELD_NAME, numeric: false, disablePadding: false, label: 'NAME' },
  { id: FIELD_TYPE, numeric: false, disablePadding: false, label: 'TYPE' },
  { id: FIELD_PARENT, numeric: false, disablePadding: false, label: 'PARENT NODE' }
];

class CfgTable extends Component {

  constructor(props) {
    super(props);
    this.state = {
      dismissibleOpen: true
    };
    this.classes = props.classes;
  }

  onRequestSort(event, property) {
    var table = updateObject({}, this.props.table);
    handleRequestSort(table, event, property);
    this.props.changeTable(table);
  }

  onClick(/*event, name*/) {
    // var table = updateObject({}, this.props.table);
    // handleClick(table, event, name);
    // this.props.changeTable(table);
  }

  onChangePage(event, page) {
    var table = updateObject({}, this.props.table);
    handleChangePage(table, event, page);
    this.props.changeTable(table);
  }

  onChangeRowsPerPage(event) {
    var table = updateObject({}, this.props.table);
    handleChangeRowsPerPage(table, event);
    this.props.changeTable(table);
  }

  render() {
    const classes = this.classes;
    const table = this.props.table;
    const data = table.entries || [];
    const order = table.order;
    const orderBy = table.orderBy;
    const rowsPerPage = table.rowsPerPage;
    const page = table.page;
    const emptyRows = rowsPerPage - Math.min(rowsPerPage, data.length - page * rowsPerPage);

    if (!data || data.length < 1) {return null;}

    return (
        <>
          <Grid>
            <GridCell span={12}>
              <Elevation className="component-style" z={2}  style={styles.cfgTable}>
                <div>
                  <span className="mdc-typography--headline6">Network Elements </span>
                </div>
                <Paper className={classes.root}>
                  <div className={classes.tableWrapper}>
                    <Table className={classes.table} aria-labelledby="tableTitle" style={{width: '100%'}}>
                      <TableHead className={classes.tableHead}>
                        <TableRow>
                          {cfgTableColumnData.map(column => {
                            return (
                              <TableCell
                                key={column.id}
                                numeric={column.numeric}
                                padding={column.disablePadding ? 'none' : 'default'}
                                sortDirection={orderBy === column.id ? order : false}
                                className={classes.tableHeadColor}
                              >
                                <Tooltip
                                  title="Sort"
                                  placement={column.numeric ? 'bottom-end' : 'bottom-start'}
                                  enterDelay={300}
                                >
                                  <TableSortLabel
                                    active={orderBy === column.id}
                                    direction={order}
                                    onClick={(event) => this.onRequestSort(event, column.id)}
                                    className={classes.tableHeadColor}
                                  >
                                    {column.label}
                                  </TableSortLabel>
                                </Tooltip>
                              </TableCell>
                            );
                          }, this)}
                        </TableRow>
                      </TableHead>
                      <TableBody>
                        {data.sort(getSortingByField(order, orderBy))
                          .slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage)
                          .map(elem => {
                            const name = getElemFieldVal(elem, FIELD_NAME);
                            const type = getElemFieldVal(elem, FIELD_TYPE);
                            const parent = getElemFieldVal(elem, FIELD_PARENT);
                            const isSelected = isRowSelected(table, name);
                            return (
                              <TableRow
                                hover
                                onClick={event => this.onClick(event, name)}
                                role="checkbox"
                                aria-checked={isSelected}
                                tabIndex={-1}
                                key={name}
                                selected={isSelected}
                              >
                                <TableCell component="th" scope="row">{name}</TableCell>
                                <TableCell>{type}</TableCell>
                                <TableCell>{parent}</TableCell>
                              </TableRow>
                            );
                          })}
                        {emptyRows > 0 && (
                          <TableRow style={{ height: 49 * emptyRows }}>
                            <TableCell colSpan={6} />
                          </TableRow>
                        )}
                      </TableBody>
                    </Table>
                  </div>
                  <TablePagination
                    component="div"
                    count={data.length}
                    rowsPerPage={rowsPerPage}
                    page={page}
                    backIconButtonProps={{'aria-label': 'Previous Page'}}
                    nextIconButtonProps={{'aria-label': 'Next Page'}}
                    onChangePage={(event, page) => this.onChangePage(event, page)}
                    onChangeRowsPerPage={event => this.onChangeRowsPerPage(event)}
                  />
                </Paper>
              </Elevation>
            </GridCell>
          </Grid>
      </>
    );
  }
}

const styles = {
  cfgTable: {
    marginTop: 20,
    padding: 10
  }
};

CfgTable.propTypes = {
  classes: PropTypes.object.isRequired
};

const mapStateToProps = (state) => {
  return {
    table: state.cfg.table
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeTable: (table) => dispatch(cfgChangeTable(table))
  };
};

export default withStyles(cfgTableStyles)(connect(mapStateToProps, mapDispatchToProps)(CfgTable));
