/**
 * Created by karolisrusenas on 05/04/2016.
 */
/* @flow */
// import {browserHistory} from 'react-router'
import {push} from 'react-router-redux'

// ------------------------------------
// Constants
// ------------------------------------
export const NAVIGATE_TO = 'NAVIGATE_TO'

// ------------------------------------
// Actions
// ------------------------------------
export function navigateTo (path:string = '/'):Action {
  return {
    type: NAVIGATE_TO,
    path: path
  }
}

export function pushPath (path) {
  return function (dispatch) {
    dispatch(navigateTo(path))
    dispatch(push(path))
  }
}

export const actions = {
  navigateTo,
  pushPath
}

// ------------------------------------
// Action Handlers
// ------------------------------------
const ACTION_HANDLERS = {
  [NAVIGATE_TO]: (state:string, action:{path:string}):string => action.path
}

// ------------------------------------
// Reducer
// ------------------------------------
const initialState = '/'
export default function navigationReducer (state:string = initialState, action:Action):string {
  const handler = ACTION_HANDLERS[action.type]

  return handler ? handler(state, action) : state
}
