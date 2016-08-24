import React from 'react'
import { Route, IndexRoute } from 'react-router'
import { requireAuthentication } from '../utils'

// NOTE: here we're making use of the `resolve.root` configuration
// option in webpack, which allows us to specify import paths as if
// they were from the root of the ~/src directory. This makes it
// very easy to navigate to files regardless of how deeply nested
// your current file is.
import CoreLayout from 'layouts/CoreLayout/CoreLayout'
import StateView from 'views/StateView/StateView'
import LoginView from 'views/Auth/LoginView'
import LogoutView from 'views/Auth/LogoutView'

export default (store) => (
  <Route path='/' component={CoreLayout}>
    <IndexRoute component={requireAuthentication(StateView)} />
    <Route path='login' component={LoginView} />
    <Route path='logout' component={requireAuthentication(LogoutView)} />
    <Route path='state' component={requireAuthentication(StateView)} />
  </Route>
)
