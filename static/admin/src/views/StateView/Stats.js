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

import {fetchRecordsCount} from '../../redux/modules/state'

export class StatsComponent extends React.Component<void, Props, void> {
  componentWillMount () {
    this.props.fetchRecordsCount(this.props.token)
  }

  render () {
    let recordsCountInfo = 'Captured request count: ' + this.props.info.recordsCount
    const statsInfo = (
      <Card>
        <CardTitle title={recordsCountInfo} subtitle='This section provides real-time Hoverfly metrics'/>
        <CardText>
          <p> Stats here</p>
        </CardText>
        <CardActions>
          <RaisedButton label='Wipe Records' primary/>
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
  token: PropTypes.string,
  info: PropTypes.object.isRequired
}

export default StatsComponent

const mapStateToProps = (state) => ({
  info: state.info
})

export default connect(mapStateToProps, {
  fetchRecordsCount
})(StatsComponent)
