import { combineReducers } from 'redux'
import { routerReducer as router } from 'react-router-redux'
import info from './modules/reducers/state'
import navigation from './modules/navigation'
import auth from './modules/reducers/auth'

export default combineReducers({
  auth,
  info,
  navigation,
  router
})
