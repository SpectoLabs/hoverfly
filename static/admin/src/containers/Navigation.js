/**
 * Created by karolisrusenas on 04/04/2016.
 */
import React, {PropTypes} from 'react'
// import { push } from 'react-router-redux'
import AppBar from 'material-ui/lib/app-bar'
import LeftNav from 'material-ui/lib/left-nav'
import Menu from 'material-ui/lib/menus/menu'
import MenuItem from 'material-ui/lib/menus/menu-item'
import {Tabs, Tab} from 'material-ui'

export default class Navigation extends React.Component<void, Props, void> {
  static propTypes = {
    store: PropTypes.object.isRequired
  };

  constructor (props) {
    super(props)
    this.state = {open: false}
  }

  handleToggle = () => this.setState({open: !this.state.open});

  handleClose = () => this.setState({open: false});

  onRequestChange = (open) => this.setState({open})

  render () {
    let myTabs = (
      <Tabs>
        <Tab label='Logout' route='/logout'/>
      </Tabs>
    )

    return (
      <div>
        <AppBar title='Hoverfly' onLeftIconButtonTouchTap={this.handleToggle} iconElementRight={myTabs}/>
        <LeftNav
          ref='leftNav'
          docked={false}
          open={this.state.open}
          onRequestChange={this.onRequestChange}
        >
          <Menu>
            <MenuItem primaryText='State'/>
            <MenuItem disabled primaryText='Records'/>
            <MenuItem disabled primaryText='Middlewares'/>
          </Menu>
        </LeftNav>
      </div>
    )
  }
}

// Since this is not a <Route> component, we add History to the context
Navigation.contextTypes = {
  history: React.PropTypes.object
}

export default Navigation
