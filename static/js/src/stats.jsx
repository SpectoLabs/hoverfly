import React from 'react';
import request from 'superagent';
var $ = require('jquery');


let WipeRecordsComponent = React.createClass({
    displayName: "WipeRecordsComponent",

    handleClick(){
        let that = this;
        request
            .del('/records')
            .end(function (err, res) {
                that.props.parent.fetchData()
            });
    },

    render() {
        return (
            <button className="button" onClick={this.handleClick}>Wipe Records</button>
        )
    }
});

let RowWrapper = React.createClass({
    render() {
        var name = this.props.name;
        var val = this.props.val;
        return (
            <tr>
                <td>{name}</td>
                <td>{val}</td>
            </tr>
        )
    }
});

let MetricsComponent = React.createClass({
    displayName: "MetricsComponent",

    getRows() {
        let rows = [];
        let metrics = this.props.counters;

        $.each(metrics, function (key, value) {
            rows.push(
                <RowWrapper key={key} name={key} val={value} />)
        });

        return rows
    },

    render() {

        let rows = this.getRows();

        return (
            <table className="u-full-width">
                <thead>
                <tr>
                    <th>Name</th>
                    <th>Value</th>
                </tr>
                </thead>
                <tbody>
                {rows}
                </tbody>
            </table>
        )

    }

});
let StatsComponent = React.createClass({
    displayName: "StatsComponent",

    getInitialState() {
        return {
            "records": null,
            "interval": 1000
        }
    },

    fetchData() {
        var url = '/count';
        var that = this;
        request
            .get(url)
            .end(function (err, res) {
                if (err) throw err;
                if (that.isMounted()) {
                    // checking whether there are any records
                    that.setState({
                        'records': res.body.count
                    });
                }
            });
    },

    greetWebsocket() {
        if (this.state.ws.readyState != this.state.ws.open) {
            this.state.ws.send("hi");
        }
    },
    componentDidMount() {
        setInterval(this.fetchData, parseInt(this.state.interval));
    },

    render() {
        let msg = "Fetching data...";
        if (this.state.records == 0) {
            msg = "No records available.";
        } else if (this.state.records == 1) {
            msg = "Currently there is 1 record."
        } else if (this.state.records > 1) {
            msg = "Currently there are " + this.state.records + " records."
        }

        return (
            <div>
                <div className="two-thirds column">
                    <WipeRecordsComponent parent={this}/>
                </div>
                <div className="one-third column">
                    {msg}
                </div>


            </div>
        )
    }
});

module.exports = StatsComponent;