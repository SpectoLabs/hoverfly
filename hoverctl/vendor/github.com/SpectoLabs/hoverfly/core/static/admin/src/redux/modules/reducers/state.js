/**
 * Created by karolisrusenas on 14/04/2016.
 */
import {createReducer} from '../../../utils'
import {
  SET_MODE,
  RECEIVE_STATE,
  RECEIVE_RECORDS_COUNT,
  RECEIVE_STATS,
  SET_REFRESH_ID,
  CLEAR_REFRESH_ID
} from '../../../constants'

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
