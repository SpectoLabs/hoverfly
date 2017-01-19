import React from 'react'
import TestUtils from 'react-addons-test-utils'
import {bindActionCreators} from 'redux'
import {StateView} from 'views/StateView/StateView'
import {mount} from 'enzyme'

function shallowRender (component) {
  const renderer = TestUtils.createRenderer()

  renderer.render(component)
  return renderer.getRenderOutput()
}

function renderWithProps (props = {}) {
  return TestUtils.renderIntoDocument(<StateView {...props} />)
}

function shallowRenderWithProps (props = {}) {
  return shallowRender(<StateView {...props} />)
}

describe('(View) State', function () {
  let _component, _rendered, _props, _spies

  let initialAuthState = {
    token: null,
    userName: null,
    isAuthenticated: false,
    isAuthenticating: false,
    statusText: null
  }

  let initialState = {
    recordsCount: null,
    mode: null,
    stats: null,
    refreshID: null
  }

  beforeEach(function () {
    _spies = {}
    _props = {
      authData: initialAuthState,
      info: initialState,
      ...bindActionCreators({
        setMode: (_spies.setMode = sinon.spy()),
        fetchState: (_spies.fetchState = sinon.spy())
      }, _spies.dispatch = sinon.spy())
    }

    _component = shallowRenderWithProps(_props)
    _rendered = renderWithProps(_props)
  })
})
