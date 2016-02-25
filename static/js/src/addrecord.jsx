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
                <button onClick={this.showModal}>Add record</button>
                <Modal ref="modal" modalStyle={modalStyle}>
                    <div className="row container">

                        <form className="form-signin" role="form" method="post" action="/add">

                            <h4 className="form-signin-heading">Request info</h4>
                            <label htmlFor="inputDestination" className="sr-only">Destination</label>
                            <input id="inputDestination" name="inputDestination" className="form-control"
                                   placeholder="www.google.com" required autofocus/>

                            <label htmlFor="inputMethod" className="sr-only">Method</label>
                            <input id="inputMethod" name="inputMethod" className="form-control" placeholder="GET"
                                   required/>

                            <label htmlFor="inputPath" className="sr-only">Path</label>
                            <input id="inputPath" name="inputPath" className="form-control"
                                   placeholder="/api/something"/>

                            <label htmlFor="inputQuery" className="sr-only">Query</label>
                            <input id="inputQuery" name="inputQuery" className="form-control"
                                   placeholder="?search=forthis"/>


                            <label htmlFor="inputRequestBody">Request body</label>
                            <textarea className="u-full-width"
                                      name="inputRequestBody"
                                      id="inputRequestBody"> </textarea>


                            <h4 className="form-signin-heading">Response info</h4>

                            <label htmlFor="inputContentType">Content Type</label>
                            <select className="u-full-width" id="inputContentType" name="inputContentType" value="json" required>
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
                                      id="inputResponseBody"> </textarea>

                            <div className="row">
                                <button className="button-primary" type="submit">Add</button>  <button onClick={this.hideModal}>Close</button>
                            </div>

                        </form>
                    </div>

                </Modal>
            </div>

        );
    }


});


module.exports = AddRecordComponent;