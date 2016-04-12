/**
 * Created by karolisrusenas on 11/04/2016.
 */
import React, {PropTypes} from 'react'
import {connect} from 'react-redux'

import Card from 'material-ui/lib/card/card'
import CardActions from 'material-ui/lib/card/card-actions'
import CardTitle from 'material-ui/lib/card/card-title'
import RaisedButton from 'material-ui/lib/raised-button'
import CardText from 'material-ui/lib/card/card-text'

import Table from 'material-ui/lib/table/table'
import TableHeaderColumn from 'material-ui/lib/table/table-header-column'
import TableRow from 'material-ui/lib/table/table-row'
import TableHeader from 'material-ui/lib/table/table-header'
import TableRowColumn from 'material-ui/lib/table/table-row-column'
import TableBody from 'material-ui/lib/table/table-body'

import {Row, Col} from 'react-bootstrap'

import {fetchRecordsCount, fetchStats} from '../../redux/modules/state'

export class RowWrapper extends React.Component<void, Props, void> {
  render () {
    return (
      <TableRow key={name}>
        <TableRowColumn>{this.props.name}</TableRowColumn>
        <TableRowColumn>{this.props.val}</TableRowColumn>
      </TableRow>
    )
  }
}

RowWrapper.propTypes = {
  name: PropTypes.string,
  val: PropTypes.number
}

export class StatsComponent extends React.Component<void, Props, void> {
  componentDidMount () {
    this._fetchRecordsCount()
    this._fetchStats()
  }

  _fetchRecordsCount () {
    this.props.fetchRecordsCount(this.props.token)
  }

  _fetchStats () {
    this.props.fetchStats(this.props.token)
  }

  getCounterRows () {
    let counters = this.props.info.stats.counters
    if (counters !== 'undefined') {
      let rows = []
      for (var key in counters) {
        if (counters.hasOwnProperty(key)) {
          rows.push(<RowWrapper key={key} name={key} val={counters[key]}/>)
        }
      }
      return rows
    }
  }

  render () {
    let recordsCountInfo = 'Captured request count: ' + this.props.info.recordsCount
    let rows = null
    if (this.props.info.stats !== 'undefined' && this.props.info.stats !== null) {
      rows = this.getCounterRows()
    }

    const counterTable = (
      <Col md={4}>
        <Table>
          <TableHeader displaySelectAll={false} adjustForCheckbox={false} enableSelectAll={false}>
            <TableRow>
              <TableHeaderColumn
                colSpan='2'
                tooltip='Basic statistics for requests that are passing through proxy'
                style={{textAlign: 'center'}}>
                Counters
              </TableHeaderColumn>
            </TableRow>
            <TableRow>
              <TableHeaderColumn tooltip='Name'>Name</TableHeaderColumn>
              <TableHeaderColumn tooltip='Value of this parameter'>Value</TableHeaderColumn>
            </TableRow>
          </TableHeader>
          <TableBody>
            {rows}
          </TableBody>
        </Table>
      </Col>
    )

    const statsInfo = (
      <Card>
        <CardTitle title={recordsCountInfo} subtitle='This section provides real-time Hoverfly metrics'/>
        <CardText>
          <Row>
            {counterTable}
          </Row>
        </CardText>
        <CardActions>
          <RaisedButton label='Wipe Records' secondary/>
        </CardActions>
      </Card>
    )
    return (
      <div>
        {statsInfo}
      </div>
    )
  }
}

StatsComponent.propTypes = {
  fetchRecordsCount: PropTypes.func.isRequired,
  fetchStats: PropTypes.func.isRequired,
  token: PropTypes.string,
  info: PropTypes.object.isRequired
}

export default StatsComponent

const mapStateToProps = (state) => ({
  info: state.info
})

export default connect(mapStateToProps, {
  fetchRecordsCount,
  fetchStats
})(StatsComponent)
