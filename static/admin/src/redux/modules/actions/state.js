/* @flow */
import fetch from 'isomorphic-fetch'
import {push} from 'react-router-redux'
import {loginUserFailure} from './auth'
import {checkHttpStatus, parseJSON} from '../../../utils'

import {
  GET_MODE,
  REQUEST_STATE,
  RECEIVE_STATE,
  REQUEST_RECORDS_COUNT,
  RECEIVE_RECORDS_COUNT,
  REQUEST_STATS,
  RECEIVE_STATS,
  SET_REFRESH_ID,
  CLEAR_REFRESH_ID
} from '../../../constants'

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
    })
      .then(checkHttpStatus)
      .then(parseJSON)
      .then((response) => dispatch(receiveState(response))
      )
      .catch((error) => {
        if (error.response.status === 401) {
          dispatch(loginUserFailure(error))
          dispatch(push('/login'))
        }
      })
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
  localStorage.setItem('refreshID', id)
  return {
    type: SET_REFRESH_ID,
    payload: id,
    receivedAt: Date.now()
  }
}

export function clearRefreshID () {
  let id = localStorage.getItem('refreshID')
  clearInterval(id)
  localStorage.removeItem('refreshID')
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
        dispatch(clearRefreshID)
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
