/**
 * Created by karolisrusenas on 04/04/2016.
 */
import React, {PropTypes} from 'react'
import AppBar from 'material-ui/lib/app-bar'
import {Tabs, Tab} from 'material-ui'

import {connect} from 'react-redux'
import {pushPath} from '../redux/modules/navigation'

type Props = {
  path: string,
  pushPath: Function
};

export default class Navigation extends React.Component<void, Props, void> {
  static propTypes = {
    // path: PropTypes.string.isRequired,
    pushPath: PropTypes.func.isRequired
  };

  constructor (props) {
    super(props)
    this.state = {open: false}
    this.handleActive = this.handleActive.bind(this)
    // this.goToPage = this.goToPage.bind(this)
  }

  handleToggle = () => this.setState({open: !this.state.open});

  handleClose = () => this.setState({open: false});

  onRequestChange = (open) => this.setState({open})

  handleActive (tab) {
    this.props.pushPath(tab.props.route)
  }

  render () {
    let myLeftTabs = (
      <Tabs>
        <Tab label='State' route='state' onActive={this.handleActive}/>
        <Tab label='Records' route='records' onActive={this.handleActive}/>
        <Tab label='Middleware' route='middleware' onActive={this.handleActive}/>
      </Tabs>
    )

    let myRightTabs = (
      <Tabs>
        <Tab label='Logout' route='/logout'/>
      </Tabs>
    )

    return (
      <div>
        <AppBar
          onLeftIconButtonTouchTap={this.handleToggle}
          iconElementLeft={myLeftTabs} iconElementRight={myRightTabs}/>
      </div>
    )
  }
}

// Since this is not a <Route> component, we add History to the context
Navigation.contextTypes = {
  history: React.PropTypes.object
}

const mapStateToProps = (state) => ({})

export default connect((mapStateToProps), {
  pushPath
})(Navigation)
