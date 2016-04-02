/* @flow */
import React, {PropTypes} from 'react'
import {connect} from 'react-redux'
import {setMode, fetchState} from '../../redux/modules/modes'
import classes from './ModeView.scss'

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
  }

  virtualize () {
    console.log('virtualize clicked')
    this.props.setMode('virtualize')
  }

  capture () {
    console.log('capture clicked')
    this.props.setMode('capture')
  }

  componentWillMount () {
    console.log('mounted')
    this.props.fetchState()
  }

  render () {
    console.log(this.props)
    return (
      <div className='container text-center'>
        <h1>Modes</h1>
        <h2>
          Current:
          {' '}
          <span className={classes['counter--green']}>{this.props.mode}</span>
        </h2>
        <button className='btn btn-default' onClick={this.virtualize}>
          Virtualize
        </button>
        {' '}
        <button className='btn btn-default' onClick={this.capture}>
          Capture
        </button>
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
