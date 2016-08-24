/**
 * Created by karolisrusenas on 06/04/2016.
 */
import {UserAuthWrapper} from 'redux-auth-wrapper'
import {push} from 'react-router-redux'

export function createConstants (...constants) {
  return constants.reduce((acc, constant) => {
    acc[constant] = constant
    return acc
  }, {})
}

export function createReducer (initialState, reducerMap) {
  return (state = initialState, action) => {
    const reducer = reducerMap[action.type]

    return reducer
      ? reducer(state, action.payload)
      : state
  }
}

export function checkHttpStatus (response) {
  if (response.status >= 200 && response.status < 300) {
    return response
  } else {
    var error = new Error(response.statusText)
    error.response = response
    throw error
  }
}

export function parseJSON (response) {
  return response.json()
}

export const requireAuthentication = UserAuthWrapper({
  authSelector: (state) => state.auth,
  predicate: (auth) => auth.isAuthenticated,
  redirectAction: push,
  wrapperDisplayName: 'UserIsJWTAuthenticated'
})
