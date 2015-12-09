import React from 'react';
import ReactDOM from 'react-dom';
import request from 'superagent';

let StateChangeButton = React.createClass({
    displayName: "StateChangeButton",

    getInitialState() {
        return {"state": null}
    },

    componentWillMount() {
        console.log("getting current state");
        this.setState({state: "virtualize"});
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