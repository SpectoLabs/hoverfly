import React from 'react';
import Modal from 'boron/DropModal'


// Style object
var modalStyle = {
    width: '60%'
};


var backdropStyle = {
    backgroundColor: 'red'
};

var contentStyle = {
    backgroundColor: 'blue',
    height: '100%'
};

let AddRecordComponent = React.createClass({
    displayName: "AddNewRecordComponent",

    showModal: function () {
        this.refs.modal.show();
    },
    hideModal: function () {
        this.refs.modal.hide();
    },
    render: function () {
        return (
            <div>
                <button onClick={this.showModal}>Add record</button>
                <Modal ref="modal" modalStyle={modalStyle}>
                    <div className="row container">

                        <form className="form-signin" role="form" method="post" action="/add">

                            <h4 className="form-signin-heading">Request info</h4>
                            <label for="inputDestination" className="sr-only">Destination</label>
                            <input id="inputDestination" name="inputDestination" className="form-control"
                                   placeholder="www.google.com" required autofocus/>

                            <label for="inputMethod" className="sr-only">Method</label>
                            <input id="inputMethod" name="inputMethod" className="form-control" placeholder="GET"
                                   required/>

                            <label for="inputPath" className="sr-only">Path</label>
                            <input id="inputPath" name="inputPath" className="form-control"
                                   placeholder="/api/something"/>

                            <label for="inputQuery" className="sr-only">Query</label>
                            <input id="inputQuery" name="inputQuery" className="form-control"
                                   placeholder="?search=forthis"/>


                            <label for="inputRequestBody">Request body</label>
                            <textarea className="u-full-width" value="<Request> value here </Request"
                                      id="inputRequestBody"> </textarea>


                            <h4 className="form-signin-heading">Response info</h4>

                            <label for="inputContentType">Content Type</label>
                            <select className="u-full-width" id="inputContentType" name="inputContentType" required>
                                <option value="xml">application/xml</option>
                                <option value="json">application/json</option>
                                <option value="text">text/html</option>
                            </select>


                            <label for="inputResponseStatusCode" className="sr-only">Status Code</label>
                            <input type="number" id="inputResponseStatusCode"
                                   name="inputResponseStatusCode" className="form-control" placeholder="200"
                                   required/>

                            <label for="inputResponseBody">Response body</label>
                            <textarea className="u-full-width" value="<Response> value here </Response"
                                      name="inputResponseBody"
                                      id="inputResponseBody"> </textarea>

                            <div className="row">
                                <button className="btn btn-lg btn-primary btn-block" type="submit">Add</button>
                                <button onClick={this.hideModal}>Close</button>
                            </div>

                        </form>

                    </div>

                </Modal>
            </div>

        );
    }


});


module.exports = AddRecordComponent;