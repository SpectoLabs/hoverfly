import React, {PropTypes} from 'react'
import {Navbar, Nav, NavItem, NavDropdown, MenuItem, Grid, Row} from 'react-bootstrap'
import '../../styles/core.scss'

// Note: Stateless/function components *will not* hot reload!
// react-transform *only* works on component classes.
//
// Since layouts rarely change, they are a good place to
// leverage React's new Stateless Functions:
// https://facebook.github.io/react/docs/reusable-components.html#stateless-functions
//
// CoreLayout is a pure function of its props, so we can
// define it with a plain javascript function...
function CoreLayout ({children}) {
  const navbarInstance = (
    <Navbar inverse>
      <Navbar.Header>
        <Navbar.Brand>
          <a href='/'>Hoverfly</a>
        </Navbar.Brand>
        <Navbar.Toggle />
      </Navbar.Header>
      <Navbar.Collapse>
        <Nav>
          <NavItem eventKey={1} href='/modes'>Modes</NavItem>
          <NavItem eventKey={2} href='/records'>Records</NavItem>
          <NavDropdown eventKey={3} title='Dropdown' id='basic-nav-dropdown'>
            <MenuItem eventKey={3.1}>Action</MenuItem>
            <MenuItem eventKey={3.2}>Another action</MenuItem>
            <MenuItem eventKey={3.3}>Something else here</MenuItem>
            <MenuItem divider/>
            <MenuItem eventKey={3.3}>Separated link</MenuItem>
          </NavDropdown>
        </Nav>
        <Nav pullRight>
          <NavItem eventKey={2} href='/logout'>Logout</NavItem>
        </Nav>
      </Navbar.Collapse>
    </Navbar>
  )

  const gridInstance = (
    <Grid>
      <Row>
        {children}
      </Row>
    </Grid>
  )

  return (
    <div className='page-container'>
      {navbarInstance}
      {gridInstance}
    </div>
  )
}

CoreLayout.propTypes = {
  children: PropTypes.element
}

export default CoreLayout
