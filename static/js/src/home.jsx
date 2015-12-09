import React from 'react';
import ReactDOM from 'react-dom';
import request from 'superagent';

let StateChangeButton = React.createClass({
    displayName: "StateChangeButton",

    getInitialState() {
        return {"state": null}
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
        console.log("getting current state");
        this.setState({state: "virtualize"});
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
      return (
          <button className="button-primary">Virtualize</button>
      )
    }
});

ReactDOM.render(
    <StateChangeButton />,
    document.getElementById("app")
);