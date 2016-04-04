/* @flow */
import React, {PropTypes} from 'react'
import {connect} from 'react-redux'
import {setMode, fetchState} from '../../redux/modules/modes'
import classes from './ModeView.scss'
import {Tabs, Tab, Button} from 'react-bootstrap'
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
    return (
      <div>
        <h1>System State</h1>
        <div className='row'>
          <Tabs defaultActiveKey={1}>
            <Tab eventKey={1} title='Modes'>
              <h3>
                Current:
                {' '}
                <span className={classes['counter--green']}>{this.props.mode}</span>
              </h3>
              <Button onClick={this.virtualize}>Virtualize</Button>
              {' '}
              <Button onClick={this.capture}>Capture</Button>
              {' '}
              <Button onClick={this.modify}>Modify</Button>
              {' '}
              <Button onClick={this.synthesize}>Synthesize</Button>
            </Tab>
          </Tabs>
        </div>
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
