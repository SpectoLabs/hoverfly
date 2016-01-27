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
    displayName: "RowWrapper",

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
            "ws": null,
            "records": null,
            "counters": null,
            "interval": 1000
        }
    },

    componentWillMount() {
        if ("WebSocket" in window) {
            this.state.ws = new WebSocket("ws:/" + window.location.host + "/statsws");

            this.state.ws.onclose = function () {
                console.log("Connection is closed, fetching manually");
                this.state.ws = null;
                // starting to fetch manually
                setInterval(this.fetchFromHTTP, parseInt(this.state.interval));
            }.bind(this);

        } else {
            console.log("WebSocket not supported by your browser.");
            // starting to fetch manually
            setInterval(this.fetchFromHTTP, parseInt(this.state.interval));
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

    // Make the function wait until the connection is made...
    waitForSocketConnection(socket, callback){
        setTimeout(
            function () {
                if (socket.readyState === 1) {
                    if (callback != null) {
                        callback();
                    }
                } else {
                    this.waitForSocketConnection(socket, callback);
                }
            }.bind(this), 5); // wait 5 ms for the connection...
    },

    // fetch counter stats, captured request count through HTTP API
    fetchFromHTTP() {
        var url = '/stats';
        request
            .get(url)
            .end(function (err, res) {
                if (err) throw err;
                if (this.isMounted()) {
                    this.setState({
                        "records": res.body.recordsCount,
                        "counters": res.body.stats.counters
                    });
                }
            }.bind(this));
    },

    componentDidMount() {
        if (this.state.ws != null) {
            this.waitForSocketConnection(this.state.ws, this.greetWebsocket);

            // getting response with data
            this.state.ws.onmessage = function (response) {
                let parsedData = JSON.parse(response.data);
                this.setState({
                    "records": parsedData.recordsCount,
                    "counters": parsedData.stats.counters
                });
            }.bind(this);

        } else {
            console.log("fetching data manually:(");
             setInterval(this.fetchFromHTTP, parseInt(this.state.interval));
        }
    },

    render() {
        let msg = "Fetching data...";
        if (this.state.records == 0) {
            msg = "No records available.";
        } else if (this.state.records == 1) {
            msg = "Currently there is 1 captured request."
        } else if (this.state.records > 1) {
            msg = "Currently there are " + this.state.records + " captured requests."
        }

        return (
            <div>
                <div className="one-third column">
                    <WipeRecordsComponent parent={this}/>
                </div>
                <div className="one-third column">
                    {msg}
                </div>
                <div className="one-third column">
                    <MetricsComponent counters={this.state.counters} />
                </div>


            </div>
        )
    }
});


module.exports = StatsComponent;