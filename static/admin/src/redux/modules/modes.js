/* @flow */
import fetch from 'isomorphic-fetch'
// ------------------------------------
// Constants
// ------------------------------------
export const SET_MODE = 'SET_MODE'
export const GET_MODE = 'GET_MODE'

export const REQUEST_STATE = 'REQUEST_STATE'
export const RECEIVE_STATE = 'RECEIVE_STATE'

// ------------------------------------
// Actions
// ------------------------------------
// NOTE: "Action" is a Flow interface defined in https://github.com/TechnologyAdvice/flow-interfaces
// If you're unfamiliar with Flow, you are completely welcome to avoid annotating your code, but
// if you'd like to learn more you can check out: flowtype.org.
// DOUBLE NOTE: there is currently a bug with babel-eslint where a `space-infix-ops` error is
// incorrectly thrown when using arrow functions, hence the oddity.
export function setMode (mode):Action {
  return function (dispatch) {
    dispatch(requestState())
    console.log('setting mode')
    return fetch('/api/state', {
      method: 'POST',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        mode: mode
      })
    }).then((response) => response.json())
      .then((json) => dispatch(receiveState(json))
      )
  }
}

export function getMode (mode):Action {
  return {
    type: GET_MODE,
    payload: mode
  }
}

export function requestState () {
  return {
    type: REQUEST_STATE
  }
}

export function receiveState (json) {
  console.log(json)
  return {
    type: RECEIVE_STATE,
    payload: json.mode,
    receivedAt: Date.now()
  }
}

export function fetchState () {
  return function (dispatch) {
    dispatch(requestState())
    console.log('fetching state')
    return fetch('/api/state')
      .then((response) => response.json())
      .then((json) => dispatch(receiveState(json))
      )
    // TODO: should also catch any error in the network call.
  }
}

export const actions = {
  setMode,
  getMode,
  requestState,
  receiveState,
  fetchState
}
// ------------------------------------
// Action Handlers
// ------------------------------------
const ACTION_HANDLERS = {
  [SET_MODE]: (state:mode, action:{payload: mode}):mode => action.payload,
  [RECEIVE_STATE]: (state:mode, action:{payload: mode}):mode => action.payload
}

// ------------------------------------
// Reducer
// ------------------------------------
let initialState = 'fetching data..'
export default function modeReducer (state:mode = initialState, action:Action):mode {
  const handler = ACTION_HANDLERS[action.type]

  return handler ? handler(state, action) : state
}
