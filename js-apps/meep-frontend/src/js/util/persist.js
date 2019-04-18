// Will persist the app state between browser refresh

const STATE_KEY = 'IDC-meep-frontend:state';

export function saveState(state) {
  try {
    let serializedState = JSON.stringify(state);
    localStorage.setItem(STATE_KEY, serializedState);
  } catch(e) {
    // TODO: consider showing an alert.
    // console.log('Error while saving app state: ', e);
  }
}

export function loadState() {
  try {
    let serializedState = localStorage.getItem(STATE_KEY);

    if (serializedState === null) {
      return this.initializeState();
    }

    return JSON.parse(serializedState);
  }
  catch (err) {
    return null;
  }
}