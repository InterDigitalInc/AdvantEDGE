import { updateObject } from '../../util/update';

const initialState = {
  data: [],
  selected: [],
  order: 'asc',
  orderBy: 'name',
  rowsPerPage: 10,
  page: 0
};

const CFG_CHANGE_TABLE = 'CFG_CHANGE_TABLE';
function cfgChangeTable(table) {
  return {
    type: CFG_CHANGE_TABLE,
    payload: table
  };
}

export { cfgChangeTable, CFG_CHANGE_TABLE };

export function cfgTableReducer(state = initialState, action) {
  switch (action.type) {
  case CFG_CHANGE_TABLE:
    var ret = updateObject({}, action.payload);
    return ret;
  default:
    return state;
  }
}