import {
  REQUEST_STATE,
  RECEIVE_STATE,
  setMode
  default as stateReducer
} from 'redux/modules/actions/state'

describe('(Redux Module) State', function () {
  it('Should export a constant COUNTER_INCREMENT.', function () {
    expect(REQUEST_STATE).to.equal('REQUEST_STATE')
  })

  describe('(Reducer)', function () {
    it('Should be a function.', function () {
      expect(stateReducer).to.be.a('function')
    })
  })
})
