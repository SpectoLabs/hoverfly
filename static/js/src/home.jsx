import React from 'react';
import ReactDOM from 'react-dom';
import request from 'superagent';

let StateChangeComponent = React.createClass({
    displayName: "StateChangeButton",

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
                    console.log(res.body);
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
      console.log(e.target.value);
        var url = '/state';
        var that = this;
        request
            .post(url)
            .send({ mode: e.target.value })
            .end(function (err, res) {
                if (err) throw err;
                if (that.isMounted()) {
                    console.log(res.body);
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


        if (this.state.mode == "virtualize") {
            virtualizeClass = primaryBtn;
        } else if (this.state.mode == "modify") {
            modifyClass = primaryBtn;
        } else if (this.state.mode == "capture") {
            captureClass = primaryBtn;
        } else if (this.state.mode == "synthesize") {
            synthesizeClass = primaryBtn;
        }

        return (
            <div>
                <button className={virtualizeClass} onClick={this.changeMode} value="virtualize">Virtualize</button>
                {' '}
                <button className={modifyClass} onClick={this.changeMode} value="modify">Modify</button>
                {' '}
                <button className={captureClass} onClick={this.changeMode} value="capture">Capture</button>
                {' '}
                <button className={synthesizeClass} onClick={this.changeMode} value="synthesize">Synthesize</button>
            </div>
        )
    }
});

ReactDOM.render(
    <StateChangeComponent />,
    document.getElementById("app")
);