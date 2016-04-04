import React, {PropTypes} from 'react'
// import {Navbar, Nav, NavItem, NavDropdown, MenuItem, Grid, Row} from 'react-bootstrap'
import '../../styles/core.scss'

// import {deepOrange500} from 'material-ui/lib/styles/colors'
import getMuiTheme from 'material-ui/lib/styles/getMuiTheme'
import MuiThemeProvider from 'material-ui/lib/MuiThemeProvider'

import {AppBar, Tabs, Tab} from 'material-ui'

// Note: Stateless/function components *will not* hot reload!
// react-transform *only* works on component classes.
//
// Since layouts rarely change, they are a good place to
// leverage React's new Stateless Functions:
// https://facebook.github.io/react/docs/reusable-components.html#stateless-functions
//
// CoreLayout is a pure function of its props, so we can
// define it with a plain javascript function...

const styles = {
  container: {
    textAlign: 'center'
    // paddingTop: 200
  }
}

const muiTheme = getMuiTheme({
  // palette: {
  //   accent1Color: deepOrange500
  // }
})

function CoreLayout ({children}) {
  var myTabs = (
    <Tabs>
      <Tab label='Logout' route='/logout'/>
    </Tabs>
  )

  return (
    <MuiThemeProvider muiTheme={muiTheme}>
      <div style={styles.container}>
        <AppBar title='Hoverfly' iconElementRight={myTabs}/>
        {children}
      </div>
    </MuiThemeProvider>
  )
}

CoreLayout.propTypes = {
  children: PropTypes.element
}

export default CoreLayout
