/**
 * Created by karolisrusenas on 13/04/2016.
 */
import React, {PropTypes} from 'react'
import ReactHighcharts from 'react-highcharts'
import {connect} from 'react-redux'

export default class CountersPie extends React.Component<void, Props, void> {
  constructor (props) {
    super(props)
    this.getOptions = this.getOptions.bind(this)
    this._prepareChartData = this._prepareChartData.bind(this)
  }

  _prepareChartData (counters) {
    let data = []
    for (var key in counters) {
      if (counters.hasOwnProperty(key)) {
        data.push({'name': key, 'y': counters[key]})
      }
    }
    return data
  }

  getOptions () {
    let data = this._prepareChartData(this.props.stats.counters)
    return {
      chart: {
        plotBackgroundColor: null,
        plotBorderWidth: null,
        plotShadow: false,
        type: 'pie'
      },
      credits: {
        enabled: false
      },
      title: {
        text: 'Requests statistics'
      },
      tooltip: {
        pointFormat: '{series.name}: <b>{point.percentage:.1f}%</b>'
      },
      plotOptions: {
        pie: {
          allowPointSelect: true,
          cursor: 'pointer',
          dataLabels: {
            enabled: false
          },
          showInLegend: true
        }
      },
      series: [{
        name: 'Type',
        colorByPoint: true,
        data: data
      }]
    }
  }

  componentWillMount () {
    this.getOptions()
  }

  render () {
    let options = this.getOptions()
    return (
      <ReactHighcharts config={options} />
    )
  }
}

CountersPie.propTypes = {
  stats: PropTypes.object.isRequired
}

const mapStateToProps = (state) => ({
  stats: state.info.stats
})

export default connect(mapStateToProps, {})(CountersPie)

