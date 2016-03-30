/* @flow */
// ------------------------------------
// Constants
// ------------------------------------
export const SET_MODE = 'SET_MODE'

// ------------------------------------
// Actions
// ------------------------------------
// NOTE: "Action" is a Flow interface defined in https://github.com/TechnologyAdvice/flow-interfaces
// If you're unfamiliar with Flow, you are completely welcome to avoid annotating your code, but
// if you'd like to learn more you can check out: flowtype.org.
// DOUBLE NOTE: there is currently a bug with babel-eslint where a `space-infix-ops` error is
// incorrectly thrown when using arrow functions, hence the oddity.
export function setMode (mode): Action {
  return {
    type: SET_MODE,
    payload: mode
  }
}

export const actions = {
  setMode
}

// ------------------------------------
// Action Handlers
// ------------------------------------
const ACTION_HANDLERS = {
  [SET_MODE]: (state: mode, action: {payload: mode}): mode => action.payload
}

// ------------------------------------
// Reducer
// ------------------------------------
const initialState = 'initial'
export default function modeReducer (state: mode = initialState, action: Action): mode {
  const handler = ACTION_HANDLERS[action.type]

  return handler ? handler(state, action) : state
}
