const initialState = {
  debug: false
};

const CHANGE_SETTINGS = 'CHANGE_SETTINGS';
export function changeSettings(settings) {
  return {
    type: CHANGE_SETTINGS,
    payload: settings
  };
}

export default function settingsReducer(state = initialState, action) {
  switch (action.type) {
  case CHANGE_SETTINGS:
    return action.payload;
  default:
    return state;
  }
}