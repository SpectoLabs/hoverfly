/**
 * Created by karolisrusenas on 06/04/2016.
 */

import {checkHttpStatus, parseJSON} from '../../../utils'
import {push} from 'react-router-redux'
import fetch from 'isomorphic-fetch'

// ------------------------------------
// Constants
// ------------------------------------
import {
  LOGIN_USER_REQUEST,
  LOGIN_USER_SUCCESS,
  LOGIN_USER_FAILURE,
  LOGOUT_USER
} from '../../../constants'

export function loginUserSuccess (token) {
  localStorage.setItem('token', token)
  return {
    type: LOGIN_USER_SUCCESS,
    payload: {
      token: token
    }
  }
}

export function loginUserFailure (error) {
  localStorage.removeItem('token')
  return {
    type: LOGIN_USER_FAILURE,
    payload: {
      status: error.response.status,
      statusText: error.response.statusText
    }
  }
}

export function loginUserRequest () {
  return {
    type: LOGIN_USER_REQUEST
  }
}

export function logout () {
  localStorage.removeItem('token')
  return {
    type: LOGOUT_USER
  }
}

export function logoutAndRedirect () {
  console.log('logoutAndRedirect action fired')
  return (dispatch, state) => {
    dispatch(logout())
    dispatch(push('/login'))
  }
}

export function loginUser (email, password, redirect = '/') {
  return function (dispatch) {
    dispatch(loginUserRequest())
    return fetch('/api/token-auth', {
      method: 'post',
      credentials: 'include',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({username: email, password: password})
    })
      .then(checkHttpStatus)
      .then(parseJSON)
      .then((response) => {
        try {
          dispatch(loginUserSuccess(response.token))
          console.log(`loginUserSucess dispatched, redirecting to ${redirect}`)
          dispatch(push(redirect))
          console.log('login success')
        } catch (e) {
          console.log('login failed')
          dispatch(loginUserFailure({
            response: {
              status: 403,
              statusText: 'Invalid token'
            }
          }))
        }
      })
      .catch((error) => {
        dispatch(loginUserFailure(error))
      })
  }
}
