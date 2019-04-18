import { getElemFieldVal } from './elem-utils';

function getSortingByField(order, orderBy) {
  return order === 'desc'
    ? (a, b) => (getElemFieldVal(b, orderBy) < getElemFieldVal(a, orderBy) ? -1 : 1)
    : (a, b) => (getElemFieldVal(a, orderBy) < getElemFieldVal(b, orderBy) ? -1 : 1);
}

function getSorting(order, orderBy) {
  return order === 'desc'
    ? (a, b) => (b[orderBy] < a[orderBy] ? -1 : 1)
    : (a, b) => (a[orderBy] < b[orderBy] ? -1 : 1);
}

function handleRequestSort(table, event, property) {
  const orderBy = property;
  let order = 'desc';

  if (table.orderBy === property && table.order === 'desc') {
    order = 'asc';
  }

  table.order = order;
  table.orderBy = orderBy;
}

function handleClick(table, event, name) {
  const selected = table.selected;
  const selectedIndex = selected.indexOf(name);
  let newSelected = [];

  if (selectedIndex === -1) {
    newSelected = newSelected.concat(selected, name);
  } else if (selectedIndex === 0) {
    newSelected = newSelected.concat(selected.slice(1));
  } else if (selectedIndex === selected.length - 1) {
    newSelected = newSelected.concat(selected.slice(0, -1));
  } else if (selectedIndex > 0) {
    newSelected = newSelected.concat(selected.slice(0, selectedIndex), selected.slice(selectedIndex + 1));
  }

  table.selected = newSelected;
}

function handleChangePage(table, event, page) {
  table.page = page;
}

function handleChangeRowsPerPage(table, event) {
  table.rowsPerPage = event.target.value;
}

function isRowSelected(table, name) {
  return table.selected.indexOf(name) !== -1;
}

export {
  getSortingByField,
  getSorting,
  handleRequestSort,
  handleClick,
  handleChangePage,
  handleChangeRowsPerPage,
  isRowSelected
};