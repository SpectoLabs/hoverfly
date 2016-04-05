import { combineReducers } from 'redux'
import { routerReducer as router } from 'react-router-redux'
import counter from './modules/counter'
import modes from './modules/modes'
import navigation from './modules/navigation'

export default combineReducers({
  counter,
  modes,
  navigation,
  router
})
