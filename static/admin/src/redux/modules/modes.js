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
  return {
    type: SET_MODE,
    payload: mode
  }
}

export function getMode (mode):Action {
  return {
    type: GET_MODE,
    payload: mode
  }
}

export function requestState (hoverfly) {
  return {
    type: REQUEST_STATE,
    hoverfly
  }
}

export function receiveState (hoverfly, json) {
  console.log(json)
  return {
    type: RECEIVE_STATE,
    hoverfly,
    // mode: json.data.children.map((child) => child.data),
    payload: json.mode,
    receivedAt: Date.now()
  }
}

export function fetchState (hoverfly) {
  // Thunk middleware knows how to handle functions.
  // It passes the dispatch method as an argument to the function,
  // thus making it able to dispatch actions itself.

  return function (dispatch) {
    // First dispatch: the app state is updated to inform
    // that the API call is starting.

    dispatch(requestState(hoverfly))

    // The function called by the thunk middleware can return a value,
    // that is passed on as the return value of the dispatch method.

    // In this case, we return a promise to wait for.
    // This is not required by thunk middleware, but it is convenient for us.
    console.log(hoverfly)
    return fetch(`http://${hoverfly}/state`)
      .then((response) => response.json())
      .then((json) =>

        // We can dispatch many times!
        // Here, we update the app state with the results of the API call.

        dispatch(receiveState(hoverfly, json))
      )

    // In a real world app, you also want to
    // catch any error in the network call.
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
// const initialState = 'initial'
let initialState = 'init'
export default function modeReducer (state:mode = initialState, action:Action):mode {
  const handler = ACTION_HANDLERS[action.type]

  return handler ? handler(state, action) : state
}
