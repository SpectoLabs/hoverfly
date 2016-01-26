import React from 'react';
import request from 'superagent';

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