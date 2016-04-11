import { combineReducers } from 'redux'
import { routerReducer as router } from 'react-router-redux'
import counter from './modules/counter'
import info from './modules/state'
import navigation from './modules/navigation'
import auth from './modules/authReducers/auth'

export default combineReducers({
  auth,
  counter,
  info,
  navigation,
  router
})
