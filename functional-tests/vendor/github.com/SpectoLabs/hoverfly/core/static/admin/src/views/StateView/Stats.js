/**
 * Created by karolisrusenas on 11/04/2016.
 */
import React, {PropTypes} from 'react'
import {connect} from 'react-redux'

import Card from 'material-ui/Card/Card'
import CardActions from 'material-ui/Card/CardActions'
import CardTitle from 'material-ui/Card/CardTitle'
import RaisedButton from 'material-ui/RaisedButton'
import CardText from 'material-ui/Card/CardText'

import {
  Table,
  TableHeaderColumn,
  TableRow,
  TableHeader,
  TableRowColumn,
  TableBody
} from 'material-ui/Table'

import {Row, Col} from 'react-bootstrap'

import CountersPie from '../../containers/StatsCharts'

import {
  fetchRecordsCount,
  fetchStats,
  wipeRecords,
  receiveStats,
  setRefreshID,
  clearRefreshID
} from '../../redux/modules/actions/state'

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
  constructor (props) {
    super(props)
    this._fetchRecordsCount = this._fetchRecordsCount.bind(this)
    this._fetchStats = this._fetchStats.bind(this)
    this.handleWipeRecordsClick = this.handleWipeRecordsClick.bind(this)

    this.greetWebsocket = this.greetWebsocket.bind(this)
    this.waitForSocketConnection = this.waitForSocketConnection.bind(this)
    this._cleanup = this._cleanup.bind(this)

    this.state = {
      ws: null,
      interval: 2000,
      refreshId: null
    }
  }

  greetWebsocket () {
    if (this.state.ws.readyState !== this.state.ws.open) {
      this.state.ws.send('hi')
    }
  }

  // Make the function wait until the connection is made...
  waitForSocketConnection (socket, callback) {
    setTimeout(
      function () {
        if (socket.readyState === 1) {
          if (callback != null) {
            callback()
          }
        } else {
          this.waitForSocketConnection(socket, callback)
        }
      }.bind(this), 5) // wait 5 ms for the connection...
  }

  componentWillMount () {
    if ('WebSocket' in window) {
      this.state.ws = new WebSocket('ws:/' + window.location.host + '/api/statsws')

      // this.state.ws.onclose = function () {
      //   console.log('Connection is closed, fetching manually')
      //   this.state.ws = null
      //   this.state.refreshId = setInterval(this._fetchStats, parseInt(this.state.interval))
      //   this.props.setRefreshID(this.state.refreshId)
      // }.bind(this)
    } else {
      console.log('WebSocket not supported by your browser.')
      this.state.refreshId = setInterval(this._fetchStats, parseInt(this.state.interval))
      this.props.setRefreshID(this.state.refreshId)
    }
  }

  _cleanup () {
    if (this.state.refreshId !== null) {
      this.props.clearRefreshID()
      // clearInterval(this.state.refreshId)
    }
    if (this.state.ws !== null) {
      this.state.ws.close()
    }
  }

  componentWillUnmount () {
    this._cleanup()
  }

  componentDidMount () {
    this._fetchRecordsCount()

    if (this.state.ws != null) {
      this.waitForSocketConnection(this.state.ws, this.greetWebsocket)
      // getting response with data
      this.state.ws.onmessage = function (response) {
        let parsedData = JSON.parse(response.data)
        this.props.receiveStats(parsedData)
      }.bind(this)
    }
  }

  _fetchRecordsCount () {
    this.props.fetchRecordsCount(this.props.token)
  }

  _fetchStats () {
    this.props.fetchStats(this.props.token)
  }

  handleWipeRecordsClick () {
    this.props.wipeRecords(this.props.token)
    this._fetchRecordsCount()
  }

  getCounterRows () {
    let counters = this.props.info.stats.counters
    if (counters !== 'undefined') {
      let rows = []
      for (var key in counters) {
        if (counters.hasOwnProperty(key)) {
          rows.push(<RowWrapper key={key} name={key} val={counters[key]} />)
        }
      }
      return rows
    }
  }

  render () {
    let recordsCountInfo = 'Captured request count: ' + this.props.info.recordsCount
    let rows = null
    let statsChart = <div></div>

    if (this.props.info.stats !== 'undefined' && this.props.info.stats !== null) {
      rows = this.getCounterRows()
      statsChart = (
        <Col md={4}>
          <CountersPie />
        </Col>
      )
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
        <CardTitle title={recordsCountInfo} subtitle='This section provides real-time Hoverfly metrics' />
        <CardText>
          <Row>
            {counterTable}
            {statsChart}
          </Row>
        </CardText>
        <CardActions>
          <RaisedButton label='Wipe Records' onClick={this.handleWipeRecordsClick} secondary />
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
  receiveStats: PropTypes.func.isRequired,
  setRefreshID: PropTypes.func.isRequired,
  clearRefreshID: PropTypes.func.isRequired,
  wipeRecords: PropTypes.func.isRequired,
  token: PropTypes.string,
  info: PropTypes.object.isRequired
}

export default StatsComponent

const mapStateToProps = (state) => ({
  info: state.info
})

export default connect(mapStateToProps, {
  fetchRecordsCount,
  fetchStats,
  setRefreshID,
  clearRefreshID,
  wipeRecords,
  receiveStats
})(StatsComponent)
