import React from 'react';
import ReactDOM from 'react-dom';
import request from 'superagent';

const VirtualizeMode = "virtualize";
const CaptureMode = "capture";
const SynthesizeMode = "synthesize";
const ModifyMode = "modify";

let ModeInfoComponent = React.createClass({
    displayName: "ModeInfoComponent",

    getInitialState() {
        return {
            "mode": this.props.data
        }
    },

    render() {
        let mode = this.props.data.mode;

        if (mode == VirtualizeMode) {
            return (
                <div>
                    <p>This mode enables service virtualization. Hoverfly uses captured requests and their unique
                       identifiers (such as a query, a method, etc.) to find the best response. In this mode, middleware
                       will be applied to matched responses.
                    </p>
                </div>
            )
        } else if (mode == CaptureMode) {
            return (
                <div>
                    <p>
                        When capture mode is active, Hoverfly intercepts requests and then makes them on behalf of the client.
                        In this mode, middleware is applied to outgoing traffic. Requests and responses are stored in Redis as JSON structures.
                    </p>
                </div>
            )
        } else if (mode == SynthesizeMode) {
            return (
                <div>
                    <p>
                        Synthesize mode enables completely synthetic, virtual services. Middleware is required for this
                        mode to work. The JSON payload with the incoming request information is supplied to the middleware.
                        The middleware must then supply data to be returned in the response. More about this in project readme.
                    </p>
                </div>
            )
        } else if (mode == ModifyMode) {
            return (
                <div>
                    <p>
                        Modify mode applies middleware to both outbound and inbound HTTP/HTTPS traffic, allowing you to modify requests
                        and responses on the fly. Hoverfly doesn't record anything when modify mode is enabled.
                    </p>
                </div>
            )
        } else {
            return (
                <div></div>
            )
        }
    }

});

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
        var url = '/records';
        var that = this;
        request
            .get(url)
            .end(function (err, res) {
                if (err) throw err;
                if (that.isMounted()) {
                    // checking whether there are any records
                    if (res.body.data == null) {
                        that.setState({
                            'records': 0
                        });
                    } else {
                        that.setState({
                            'records': res.body.data.length
                        });
                    }
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

let StateChangeComponent = React.createClass({
    displayName: "StateChangeComponent",

    getInitialState() {
        return {"mode": null}
    },

    getCurrentMode() {
        var url = '/state';
        var that = this;
        request
            .get(url)
            .end(function (err, res) {
                if (err) throw err;
                if (that.isMounted()) {
                    that.setState({
                        'mode': res.body.mode
                    });
                }
            });
    },

    componentWillMount() {
        this.getCurrentMode();
    },

    changeMode(e){
        var url = '/state';
        var that = this;
        request
            .post(url)
            .send({mode: e.target.value})
            .end(function (err, res) {
                if (err) throw err;
                if (that.isMounted()) {
                    that.setState({
                        'mode': res.body.mode
                    });
                }
            });
    },

    render() {
        let defaultBtn = "button";
        let primaryBtn = "button-primary";
        // deciding states
        let virtualizeClass = defaultBtn;
        let modifyClass = defaultBtn;
        let captureClass = defaultBtn;
        let synthesizeClass = defaultBtn;


        if (this.state.mode == VirtualizeMode) {
            virtualizeClass = primaryBtn;
        } else if (this.state.mode == ModifyMode) {
            modifyClass = primaryBtn;
        } else if (this.state.mode == CaptureMode) {
            captureClass = primaryBtn;
        } else if (this.state.mode == SynthesizeMode) {
            synthesizeClass = primaryBtn;
        }

        let data = {
            "mode": this.state.mode
        };

        return (
            <div>
                <hr/>
                <div className="row">
                    <div className="two-thirds column">
                        <button className={virtualizeClass} onClick={this.changeMode} value="virtualize">Virtualize
                        </button>
                        {' '}
                        <button className={modifyClass} onClick={this.changeMode} value="modify">Modify</button>
                        {' '}
                        <button className={captureClass} onClick={this.changeMode} value="capture">Capture</button>
                        {' '}
                        <button className={synthesizeClass} onClick={this.changeMode} value="synthesize">Synthesize
                        </button>
                    </div>
                    <div className="one-third column">
                        <ModeInfoComponent data={data}/>
                    </div>
                </div>
                    <hr/>
                    <div className="row">
                        <StatsComponent />
                    </div>

            </div>
        )
    }
});

ReactDOM.render(
    <StateChangeComponent />,
    document.getElementById("app")
);
