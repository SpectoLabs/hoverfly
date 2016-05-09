/**
 * Created by karolisrusenas on 07/04/2016.
 */
import React, {PropTypes} from 'react'
import {bindActionCreators} from 'redux'
import {connect} from 'react-redux'
import * as actionCreators from '../../redux/modules/actions/auth'
import {Col} from 'react-bootstrap'

import Card from 'material-ui/Card/Card'
import CardTitle from 'material-ui/Card/CardTitle'
import CardText from 'material-ui/Card/CardText'
import CircularProgress from 'material-ui/CircularProgress'

export class LogoutView extends React.Component {
  static propTypes = {
    actions: PropTypes.object.isRequired,
    statusText: PropTypes.string
  };

  componentDidMount () {
    setTimeout(() => this.props.actions.logoutAndRedirect(), 2000)
  }

  render () {
    return (
      <div>
        <Col md={3} />
        <Col md={6}>
          <Card>
            <CardTitle title='Logged out!' />
            <CardText>
              {this.props.statusText ? <div className='alert alert-info'>{this.props.statusText}</div> : ''}
              <div>
                <CircularProgress size={1.5} />
              </div>
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

export default connect(mapStateToProps, mapDispatchToProps)(LogoutView)
