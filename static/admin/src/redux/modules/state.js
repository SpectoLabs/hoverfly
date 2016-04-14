/* @flow */
import fetch from 'isomorphic-fetch'
import {push} from 'react-router-redux'
import {loginUserFailure} from './actions/auth'
import {checkHttpStatus, parseJSON, createReducer} from '../../utils'
// ------------------------------------
// Constants
// ------------------------------------
export const SET_MODE = 'SET_MODE'
export const GET_MODE = 'GET_MODE'

export const REQUEST_STATE = 'REQUEST_STATE'
export const RECEIVE_STATE = 'RECEIVE_STATE'

export const REQUEST_RECORDS_COUNT = 'REQUEST_RECORDS_COUNT'
export const RECEIVE_RECORDS_COUNT = 'RECEIVE_RECORDS_COUNT'

export const REQUEST_STATS = 'REQUEST_STATS'
export const RECEIVE_STATS = 'RECEIVE_STATS'

export const SET_REFRESH_ID = 'SET_REFRESH_ID'
export const CLEAR_REFRESH_ID = 'CLEAR_REFRESH_ID'

// ------------------------------------
// Actions
// ------------------------------------
// NOTE: "Action" is a Flow interface defined in https://github.com/TechnologyAdvice/flow-interfaces
// If you're unfamiliar with Flow, you are completely welcome to avoid annotating your code, but
// if you'd like to learn more you can check out: flowtype.org.
// DOUBLE NOTE: there is currently a bug with babel-eslint where a `space-infix-ops` error is
// incorrectly thrown when using arrow functions, hence the oddity.
export function setMode (mode, token):Action {
  return function (dispatch) {
    dispatch(requestState())
    return fetch('/api/state', {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
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

export function requestStats () {
  return {
    type: REQUEST_STATS
  }
}

export function receiveStats (json) {
  return {
    type: RECEIVE_STATS,
    payload: json,
    receivedAt: Date.now()
  }
}

export function requestRecordsCount () {
  return {
    type: REQUEST_RECORDS_COUNT
  }
}

export function receiveRecordsCount (json) {
  return {
    type: RECEIVE_RECORDS_COUNT,
    payload: json.count,
    receivedAt: Date.now()
  }
}

export function receiveState (json) {
  return {
    type: RECEIVE_STATE,
    payload: json.mode,
    receivedAt: Date.now()
  }
}

export function setRefreshID (id) {
  return {
    type: SET_REFRESH_ID,
    payload: id,
    receivedAt: Date.now()
  }
}

export function clearRefreshID (id) {
  clearInterval(id)
  return {
    type: CLEAR_REFRESH_ID,
    receivedAt: Date.now()
  }
}

export function fetchState (token) {
  if (typeof token === 'undefined') {
    token = ''
  }
  return function (dispatch) {
    dispatch(requestState())
    return fetch('/api/state', {
      credentials: 'include',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
      .then(checkHttpStatus)
      .then(parseJSON)
      .then((response) => {
        dispatch(receiveState(response))
      })
      .catch((error) => {
        if (error.response.status === 401) {
          dispatch(loginUserFailure(error))
          dispatch(push('/login'))
        }
      })
  }
}

export function fetchStats (token) {
  if (typeof token === 'undefined') {
    token = ''
  }
  return function (dispatch) {
    dispatch(requestStats())
    return fetch('/api/stats', {
      credentials: 'include',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
      .then(checkHttpStatus)
      .then(parseJSON)
      .then((response) => {
        dispatch(receiveStats(response))
      })
      .catch((error) => {
        if (error.response.status === 401) {
          dispatch(loginUserFailure(error))
          dispatch(push('/login'))
        }
      })
  }
}

export function fetchRecordsCount (token) {
  if (typeof token === 'undefined') {
    token = ''
  }
  return function (dispatch) {
    dispatch(requestRecordsCount())
    return fetch('/api/count', {
      credentials: 'include',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
      .then(checkHttpStatus)
      .then(parseJSON)
      .then((response) => {
        dispatch(receiveRecordsCount(response))
      })
      .catch((error) => {
        if (error.response.status === 401) {
          dispatch(loginUserFailure(error))
          dispatch(push('/login'))
        }
      })
  }
}

export function wipeRecords (token) {
  if (typeof token === 'undefined') {
    token = ''
  }
  return function (dispatch) {
    dispatch(requestRecordsCount())
    return fetch('/api/records', {
      method: 'DELETE',
      credentials: 'include',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
      .then(checkHttpStatus)
      .catch((error) => {
        if (error.response.status === 401) {
          dispatch(loginUserFailure(error))
          dispatch(push('/login'))
        }
      })
  }
}

export const actions = {
  setMode,
  getMode,
  requestState,
  receiveState,
  fetchState,
  requestRecordsCount,
  receiveRecordsCount,
  fetchRecordsCount,
  requestStats,
  receiveStats,
  fetchStats,
  wipeRecords,
  setRefreshID,
  clearRefreshID
}
// ------------------------------------
// Action Handlers
// ------------------------------------

const initialState = {
  recordsCount: null,
  mode: null,
  stats: null,
  refreshID: null
}

export default createReducer(initialState, {
  [SET_MODE]: (state, payload) => {
    return Object.assign({}, state, {
      'mode': payload,
      'recordsCount': state.recordsCount
    })
  },
  [RECEIVE_STATE]: (state, payload) => {
    return Object.assign({}, state, {
      'mode': payload,
      'recordsCount': state.recordsCount
    })
  },
  [RECEIVE_RECORDS_COUNT]: (state, payload) => {
    return Object.assign({}, state, {
      'mode': state.mode,
      'recordsCount': payload
    })
  },
  [RECEIVE_STATS]: (state, payload) => {
    return Object.assign({}, state, {
      'recordsCount': payload.recordsCount,
      'stats': payload.stats
    })
  },
  [SET_REFRESH_ID]: (state, payload) => {
    return Object.assign({}, state, {
      'refreshID': payload
    })
  },
  [CLEAR_REFRESH_ID]: (state, payload) => {
    return Object.assign({}, state, {
      'refreshID': null
    })
  }
})
