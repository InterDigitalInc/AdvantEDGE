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

import _ from 'lodash';
import PropTypes from 'prop-types';
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
import { Elevation } from '@rmwc/elevation';
import { withStyles } from '@material-ui/core/styles';
import { connect } from 'react-redux';
import React, { Component } from 'react';

import { updateObject } from '../../util/object-util';
import { podsWithServiceMaps } from '../../state/exec';

import { execChangeTable } from '../../state/exec';

import {
  getSorting,
  handleRequestSort,
  handleChangePage,
  handleChangeRowsPerPage,
  isRowSelected
} from '../../util/table-utils';

const IngressServiceMapRow = props => {
  return (
    <Grid style={{ marginBottom: 10, marginTop: 10, marginLeft: -10 }}>
      <GridCell span={12}>
        <span>
          {' '}
          I: {props.entry.name}: {props.entry.externalPort}{' '}
        </span>
      </GridCell>
    </Grid>
  );
};

const EgressServiceMapRow = props => {
  return (
    <Grid style={{ marginBottom: 10, marginTop: 10, marginLeft: -10 }}>
      <GridCell span={12}>
        <span>
          {' '}
          E: {props.entry.name}: {props.entry.meSvcName},{props.entry.ip},
          {props.entry.port},{props.entry.protocol}
        </span>
      </GridCell>
    </Grid>
  );
};

// Network Element Execution Styles & Table
const execTableStyles = theme => ({
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
  },
  statusRunning: {
    color: '#388E3C',
    'font-weight': 'bold'
  },
  statusPending: {
    color: '#FFA500',
    'font-weight': 'bold'
  }
});

const execTableColumnData = [
  { id: 'name', numeric: false, disablePadding: false, label: 'NAME' },
  {
    id: 'logicalState',
    numeric: false,
    disablePadding: false,
    label: 'STATUS'
  },
  {
    id: 'serviceMaps',
    numeric: false,
    disablePadding: false,
    label: 'SERVICE MAPS'
  }
];

class ExecTable extends Component {
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
    const data = this.props.podsWithServiceMaps;
    const order = table.order;
    const orderBy = table.orderBy;
    const rowsPerPage = table.rowsPerPage;
    const page = table.page;
    const emptyRows =
      rowsPerPage - Math.min(rowsPerPage, data.length - page * rowsPerPage);

    if (!data || !data.length) {
      return null;
    }

    return (
      <Grid>
        <GridCell span={12}>
          <Elevation className="component-style" z={2} style={styles.execTable}>
            <div>
              <span className="mdc-typography--headline6">
                Network Elements{' '}
              </span>
            </div>
            <Paper className={classes.root}>
              <div className={classes.tableWrapper}>
                <Table
                  className={classes.table}
                  aria-labelledby="tableTitle"
                  style={{ width: '100%' }}
                >
                  <TableHead className={classes.tableHead}>
                    <TableRow>
                      {execTableColumnData.map(column => {
                        return (
                          <TableCell
                            key={column.id}
                            numeric={column.numeric}
                            padding={column.disablePadding ? 'none' : 'default'}
                            sortDirection={
                              orderBy === column.id ? order : false
                            }
                            className={classes.tableHeadColor}
                          >
                            <Tooltip
                              title="Sort"
                              placement={
                                column.numeric ? 'bottom-end' : 'bottom-start'
                              }
                              enterDelay={300}
                            >
                              <TableSortLabel
                                active={orderBy === column.id}
                                direction={order}
                                onClick={event =>
                                  this.onRequestSort(event, column.id)
                                }
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
                    {data
                      .sort(getSorting(order, orderBy))
                      .slice(
                        page * rowsPerPage,
                        page * rowsPerPage + rowsPerPage
                      )
                      .map(n => {
                        const isSelected = isRowSelected(table, n.name);
                        return (
                          <TableRow
                            hover
                            onClick={event => this.onClick(event, n.name)}
                            role="checkbox"
                            aria-checked={isSelected}
                            tabIndex={-1}
                            key={n.name}
                            selected={isSelected}
                          >
                            <TableCell component="th" scope="row">
                              {n.name}
                            </TableCell>
                            <TableCell
                              className={
                                n.logicalState === 'Running'
                                  ? classes.statusRunning
                                  : classes.statusPending
                              }
                            >
                              {n.logicalState}
                            </TableCell>
                            <TableCell>
                              {n.ingressServiceMap
                                ? _.map(n.ingressServiceMap, sm => {
                                  return (
                                    <IngressServiceMapRow
                                      entry={sm}
                                      key={sm.externalPort}
                                    />
                                  );
                                })
                                : ''}
                              {n.egressServiceMap
                                ? _.map(n.egressServiceMap, sm => {
                                  return (
                                    <EgressServiceMapRow
                                      entry={sm}
                                      key={sm.name}
                                    />
                                  );
                                })
                                : ''}
                            </TableCell>
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
                count={data.length ? data.length : 0}
                rowsPerPage={rowsPerPage}
                page={page}
                backIconButtonProps={{ 'aria-label': 'Previous Page' }}
                nextIconButtonProps={{ 'aria-label': 'Next Page' }}
                onChangePage={(event, page) => this.onChangePage(event, page)}
                onChangeRowsPerPage={event => this.onChangeRowsPerPage(event)}
              />
            </Paper>
          </Elevation>
        </GridCell>
      </Grid>
    );
  }
}

const styles = {
  execTable: {
    marginTop: 20,
    padding: 10
  }
};

ExecTable.propTypes = {
  classes: PropTypes.object.isRequired
};

const mapStateToProps = state => {
  return {
    table: state.exec.table,
    podsWithServiceMaps: podsWithServiceMaps(state)
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeTable: table => dispatch(execChangeTable(table))
  };
};

export default withStyles(execTableStyles)(
  connect(
    mapStateToProps,
    mapDispatchToProps
  )(ExecTable)
);
