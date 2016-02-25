import React from 'react';
import Modal from 'boron/DropModal'

//var Modal = require('boron/DropModal');

let AddRecordComponent = React.createClass({
    displayName: "AddNewRecordComponent",

    showModal: function(){
        this.refs.modal.show();
    },
    hideModal: function(){
        this.refs.modal.hide();
    },
    render: function() {
        return (
            <div>
                <button onClick={this.showModal}>Add record</button>
                <Modal ref="modal">
                    <h2>I am a dialog</h2>
                    <button onClick={this.hideModal}>Close</button>
                </Modal>
            </div>
        );
    }


});


module.exports = AddRecordComponent;