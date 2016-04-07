/**
 * Created by karolisrusenas on 06/04/2016.
 */

import React, {PropTypes} from 'react/addons'
import {bindActionCreators} from 'redux'
import {connect} from 'react-redux'
import reactMixin from 'react-mixin'
import * as actionCreators from '../../redux/modules/auth'
import {Col} from 'react-bootstrap'

import Card from 'material-ui/lib/card/card'
// import CardActions from 'material-ui/lib/card/card-actions'
// import CardHeader from 'material-ui/lib/card/card-header'
import CardTitle from 'material-ui/lib/card/card-title'
import RaisedButton from 'material-ui/lib/raised-button'
import CardText from 'material-ui/lib/card/card-text'

export class LoginView extends React.Component {
  static propTypes = {
    location: PropTypes.object.isRequired,
    actions: PropTypes.object.isRequired,
    statusText: PropTypes.string.isRequired,
    isAuthenticating: PropTypes.bool.isRequired
  };

  constructor (props) {
    super(props)
    const redirectRoute = this.props.location.query.redirect || '/login'
    this.state = {
      email: '',
      password: '',
      redirectTo: redirectRoute
    }
  }

  login (e) {
    e.preventDefault()
    this.props.actions.loginUser(this.state.email, this.state.password, this.state.redirectTo)
  }

  render () {
    return (
      <div>
        <Col md={3}/>
        <Col md={6}>
          <Card>
            <CardTitle title='Login required!' subtitle='Hint: hf/hf'/>
            <CardText>
              {this.props.statusText ? <div className='alert alert-info'>{this.props.statusText}</div> : ''}
              <form role='form'>
                <div className='form-group'>
                  <input
                    type='text'
                    className='form-control input-lg'
                    valueLink={this.linkState('email')}
                    placeholder='Username'/>
                </div>
                <div className='form-group'>
                  <input
                    type='password'
                    className='form-control input-lg'
                    valueLink={this.linkState('password')}
                    placeholder='Password'/>
                </div>
                <RaisedButton
                  type='submit'
                  label='Submit'
                  onClick={this.login.bind(this)}
                  disabled={this.props.isAuthenticating}
                  primary/>
              </form>
            </CardText>
          </Card>
        </Col>
      </div>
    )
  }
}

reactMixin(LoginView.prototype, React.addons.LinkedStateMixin)

const mapStateToProps = (state) => ({
  isAuthenticating: state.auth.isAuthenticating,
  statusText: state.auth.statusText
})

const mapDispatchToProps = (dispatch) => ({
  actions: bindActionCreators(actionCreators, dispatch)
})

export default connect(mapStateToProps, mapDispatchToProps)(LoginView)
