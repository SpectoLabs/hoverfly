import React from 'react';
import Modal from 'boron/DropModal'


// Style object
var modalStyle = {
    width: '60%'
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
                <div className="row">
                    <div className="four columns">
                        <button onClick={this.showModal}>Add record</button>
                    </div>

                    <div className="eight columns">
                        <div>
                            <p>
                                Add new record. This functionality is currently limited (you cannot add custom response headers)
                                and will be replaced when the new UI comes out. This form serves as a shortcut to quickly
                                insert requests/responses and bypass any encoding issues which might come up. Capturing HTTP requests and
                                responses is still a go-to way of API simulation.
                            </p>
                        </div>
                    </div>

                </div>


                <Modal ref="modal" modalStyle={modalStyle}>
                    <div className="row container">

                        <form className="form-signin" role="form" method="post" action="/add">

                            <h4 className="form-signin-heading">Define Request</h4>

                            <div className="row">
                                <div className="three columns">
                                    <label htmlFor="inputDestination" className="sr-only">Destination</label>
                                    <input id="inputDestination" name="inputDestination"
                                           className="form-control u-full-width"
                                           placeholder="www.google.com" required autofocus/>{" "}
                                </div>

                                <div className="three columns">
                                    <label htmlFor="inputPath" className="sr-only">Path</label>
                                    <input id="inputPath" name="inputPath" className="form-control u-full-width"
                                           placeholder="/api/something"/>
                                </div>

                                <div className="three columns">
                                    <label htmlFor="inputQuery" className="sr-only">Query</label>
                                    <input id="inputQuery" name="inputQuery" className="form-control"
                                           placeholder="?search=forthis"/>
                                </div>

                            </div>

                            <label htmlFor="inputMethod" className="u-full-width">Method</label>
                            <input id="inputMethod" name="inputMethod" className="form-control"
                                   placeholder="GET"
                                   required/>

                            <label htmlFor="inputRequestBody">Request body</label>
                            <textarea className="u-full-width"
                                      name="inputRequestBody"
                                      id="inputRequestBody"></textarea>

                            <h4 className="form-signin-heading">Define Response</h4>

                            <label htmlFor="inputContentType">Content Type</label>
                            <select className="u-full-width" id="inputContentType" name="inputContentType" value="json"
                                    required>
                                <option value="xml">application/xml</option>
                                <option value="json">application/json</option>
                                <option value="text">text/html</option>
                            </select>


                            <label htmlFor="inputResponseStatusCode" className="sr-only">Status Code</label>
                            <input type="number" id="inputResponseStatusCode"
                                   name="inputResponseStatusCode" className="form-control" placeholder="200"
                                   required/>

                            <label htmlFor="inputResponseBody">Response body</label>
                            <textarea className="u-full-width"
                                      name="inputResponseBody"
                                      id="inputResponseBody"></textarea>

                            <div className="row">
                                <button className="button-primary" type="submit">Add</button>
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