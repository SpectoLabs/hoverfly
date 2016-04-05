/* @flow */
import React, {PropTypes} from 'react'
import {connect} from 'react-redux'
import {setMode, fetchState} from '../../redux/modules/modes'

import Card from 'material-ui/lib/card/card'
import CardActions from 'material-ui/lib/card/card-actions'
// import CardHeader from 'material-ui/lib/card/card-header'
import CardTitle from 'material-ui/lib/card/card-title'
import RaisedButton from 'material-ui/lib/raised-button'
import CardText from 'material-ui/lib/card/card-text'

// import {Tabs, Tab, Button} from 'react-bootstrap'
// We can use Flow (http://flowtype.org/) to type our component's props
// and state. For convenience we've included both regular propTypes and
// Flow types, but if you want to try just using Flow you'll want to
// disable the eslint rule `react/prop-types`.
// NOTE: You can run `npm run flow:check` to check for any errors in your
// code, or `npm i -g flow-bin` to have access to the binary globally.
// Sorry Windows users :(.
type Props = {
  mode: string,
  setMode: Function
};

const virtualizeMode = 'virtualize'
const captureMode = 'capture'
const modifyMode = 'modify'
const synthesizeMode = 'synthesize'

export class ModeInfoComponent extends React.Component<void, Props, void> {
  static propTypes = {
    mode: PropTypes.string.isRequired
  };

  render () {
    let mode = this.props.mode
    if (mode === virtualizeMode) {
      return (
        <div>
          <p>This mode enables service virtualization. Hoverfly uses captured requests and their unique
                        identifiers (such as a query, a method, etc.) to find the best response. In this mode,
                        middleware
                        will be applied to matched responses.
          </p>
        </div>
      )
    } else if (mode === captureMode) {
      return (
        <div>
          <p>
            When capture mode is active, Hoverfly intercepts requests and then makes them on behalf of the
            client.
            In this mode, middleware is applied to outgoing traffic. Requests and responses are stored in
            embedded database as JSON structures.
          </p>
        </div>
      )
    } else if (mode === synthesizeMode) {
      return (
        <div>
          <p>
            Synthesize mode enables completely synthetic, virtual services. Middleware is required for this
            mode to work. The JSON payload with the incoming request information is supplied to the
            middleware.
            The middleware must then supply data to be returned in the response. More about this in project
            readme.
          </p>
        </div>
      )
    } else if (mode === modifyMode) {
      return (
        <div>
          <p>
            Modify mode applies middleware to both outbound and inbound HTTP/HTTPS traffic, allowing you to
            modify requests
            and responses on the fly. Hoverfly doesn't record anything when modify mode is enabled.
          </p>
        </div>
      )
    } else {
      return (
        <div></div>
      )
    }
  }
}

// We avoid using the `@connect` decorator on the class definition so
// that we can export the undecorated component for testing.
// See: http://rackt.github.io/redux/docs/recipes/WritingTests.html
export class ModeView extends React.Component<void, Props, void> {
  static propTypes = {
    mode: PropTypes.string.isRequired,
    setMode: PropTypes.func.isRequired,
    fetchState: PropTypes.func.isRequired
  };

  constructor (props) {
    super(props)
    this.virtualize = this.virtualize.bind(this)
    this.capture = this.capture.bind(this)
    this.modify = this.modify.bind(this)
    this.synthesize = this.synthesize.bind(this)
  }

  virtualize () {
    this.props.setMode('virtualize')
  }

  capture () {
    this.props.setMode('capture')
  }

  modify () {
    this.props.setMode('modify')
  }

  synthesize () {
    this.props.setMode('synthesize')
  }

  componentWillMount () {
    this.props.fetchState()
  }

  render () {
    let currentMode = this.props.mode
    let currentModeInfo = 'Current mode: ' + currentMode

    let virtualizeButton
    let captureButton
    let modifyButton
    let synthesizeButton

    if (currentMode === virtualizeMode) {
      virtualizeButton = <RaisedButton label='Virtualize' onClick={this.virtualize} primary/>
    } else {
      virtualizeButton = <RaisedButton label='Virtualize' onClick={this.virtualize}/>
    }

    if (currentMode === captureMode) {
      captureButton = <RaisedButton label='Capture' onClick={this.capture} primary/>
    } else {
      captureButton = <RaisedButton label='Capture' onClick={this.capture}/>
    }

    if (currentMode === modifyMode) {
      modifyButton = <RaisedButton label='Modify' onClick={this.modify} primary/>
    } else {
      modifyButton = <RaisedButton label='Modify' onClick={this.modify}/>
    }

    if (currentMode === synthesizeMode) {
      synthesizeButton = <RaisedButton label='Synthesize' onClick={this.synthesize} primary/>
    } else {
      synthesizeButton = <RaisedButton label='Synthesize' onClick={this.synthesize}/>
    }

    const modeInfo = (
      <Card>
        <CardTitle title={currentModeInfo} subtitle='You can change proxy behaviour here'/>
        <CardText>
          <ModeInfoComponent mode={this.props.mode}/>
        </CardText>
        <CardActions>
          {virtualizeButton}
          {captureButton}
          {modifyButton}
          {synthesizeButton}
        </CardActions>
      </Card>
    )
    return (
      <div>
        {modeInfo}
      </div>
    )
  }
}

const mapStateToProps = (state) => ({
  mode: state.modes
})

export default connect((mapStateToProps), {
  setMode,
  fetchState
})(ModeView)
