/**
 * Created by karolisrusenas on 06/04/2016.
 */

import React, {PropTypes} from 'react'
import {bindActionCreators} from 'redux'
import {connect} from 'react-redux'
import * as actionCreators from '../../redux/modules/actions/auth'
import {Col} from 'react-bootstrap'

import Card from 'material-ui/Card/Card'
import CardTitle from 'material-ui/Card/CardTitle'
import RaisedButton from 'material-ui/RaisedButton'
import CardText from 'material-ui/Card/CardText'
import TextField from 'material-ui/TextField'

export class LoginView extends React.Component {
  static propTypes = {
    location: PropTypes.object.isRequired,
    actions: PropTypes.object.isRequired,
    statusText: PropTypes.string,
    isAuthenticating: PropTypes.bool.isRequired
  };

  constructor (props) {
    super(props)
    const redirectRoute = this.props.location.query.redirect || '/'
    this.state = {
      email: '',
      password: '',
      redirectTo: redirectRoute
    }
    this.login = this.login.bind(this)
  }

  componentDidMount () {
    let token = localStorage.getItem('token')
    if (token !== null) {
      this.props.actions.loginWithTokenAndRedirect(token, this.state.redirectTo)
    }
  }

  login (e) {
    e.preventDefault()
    this.props.actions.loginUser(this.refs.username.getValue(), this.refs.password.getValue(), this.state.redirectTo)
  }

  render () {
    return (
      <div>
        <Col md={3} />
        <Col md={6}>
          <Card>
            <CardTitle
              title='Login required!'
              subtitle='Hint: if auth is disabled - use any username/password combination' />
            <CardText>
              {this.props.statusText ? <div className='alert alert-info'>{this.props.statusText}</div> : ''}
              <form role='form'>
                <div className='form-group'>
                  <TextField
                    hintText='Username'
                    ref='username'
                    type='text'
                  />
                  <br />
                  <TextField
                    hintText='Password'
                    type='password'
                    ref='password'
                  />
                </div>
                <RaisedButton
                  type='submit'
                  label='Submit'
                  onClick={this.login}
                  disabled={this.props.isAuthenticating}
                  primary />
              </form>
            </CardText>
          </Card>
        </Col>
      </div>
    )
  }
}

const mapStateToProps = (state) => ({
  isAuthenticating: state.auth.isAuthenticating,
  statusText: state.auth.statusText
})

const mapDispatchToProps = (dispatch) => ({
  actions: bindActionCreators(actionCreators, dispatch)
})

export default connect(mapStateToProps, mapDispatchToProps)(LoginView)
