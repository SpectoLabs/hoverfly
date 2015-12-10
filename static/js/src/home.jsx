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
                        identifiers (such as query, method, etc.) to find best response. If used with middleware - it
                        will
                        be applied to matched responses</p>
                </div>
            )
        } else if (mode == CaptureMode) {
            return (
                <div>
                    <p>
                        When capture mode is active - Hoverfly intercepts requests and makes them on behalf of client.
                        This enables Hoverfly to also apply middleware (if user supplied middleware setting) on outgoing
                        traffic. Requests and responses are stored in Redis as a JSON structures.
                    </p>
                </div>
            )
        } else if (mode == SynthesizeMode) {
            return (
                <div>
                    <p>
                        Synthesize mode enforces completely synthetic, virtual services. Middleware is required for this mode
                        work. JSON payload with incoming request information is supplied to middleware and it is expected to
                        provide data that will be used to create response. More about this in project readme.
                    </p>
                </div>
            )
        } else if (mode == ModifyMode) {
            return (
                <div>
                    <p>
                        Modify mode applies middleware to both outbound and inbound HTTP traffic. Hoverfly doesn't record
                        anything when modify mode is enabled.
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
        //console.log(e.target.value);
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
            <div className="row">
                <div className="two-thirds column">
                    <button className={virtualizeClass} onClick={this.changeMode} value="virtualize">Virtualize</button>
                    {' '}
                    <button className={modifyClass} onClick={this.changeMode} value="modify">Modify</button>
                    {' '}
                    <button className={captureClass} onClick={this.changeMode} value="capture">Capture</button>
                    {' '}
                    <button className={synthesizeClass} onClick={this.changeMode} value="synthesize">Synthesize</button>
                </div>
                <div className="one-third column">
                    <ModeInfoComponent data={data} />
                </div>
            </div>
        )
    }
});

ReactDOM.render(
    <StateChangeComponent />,
    document.getElementById("app")
);