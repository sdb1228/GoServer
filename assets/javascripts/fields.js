/**
 * Copyright © 2015-2016 SoccerLC. All rights reserved.
 *
 */
import React, { Component, PropTypes } from 'react';
var axios = require('axios');

class App extends Component {

  constructor(props) {
    super(props);
    axios.get('/api/v1/fields/correction')
      .then(function (response) {
        this.success(response.data);
      }.bind(this))
      .catch(function (response) {
        console.log(response);
      });
  }
  renderItem (item) {
    var boundClick = this.handleClick.bind(this, item.id);
    return ([<div class="form-group"><label for="fieldName">{item.name}:</label><input placeholder="Address" type="text" className="form-control"></input><input placeholder="City" type="text" className="form-control"></input><input placeholder="Zip" type="text" className="form-control"></input><button onClick={this.handleClick.bind(this)} type="button" id={item.id}>Submit</button></div>]);
  }
  handleClick (event) {
    var id = event.currentTarget.id
    var address = $(event.currentTarget.parentElement.children[1]).val()
    var city = $(event.currentTarget.parentElement.children[2]).val()
    var zip = $(event.currentTarget.parentElement.children[3]).val()
    if (!address || !city || !zip ) {
      alert("missing a field!");
    }
    else{
      axios.post('/api/v1/fields/postCorrection', {
        id: id,
        address: address,
        city: city,
        zip: zip
      }).then(function (response) {
          $(".form-inline")[0].reset();
          this.success(response.data);
        }.bind(this))
        .catch(function (response) {
          console.log(response);
        });
    }
  }
  success(data){
    this.setState({fields: data, loading: false});
  }

  render() {
    var content = [];
    if (!this.state) {
      content = [<h1>Loading</h1>]
    }
    else{
      for (var i = 0; i < this.state.fields.length; i++) {
        content.push(this.renderItem(this.state.fields[i]));
      }
    }
    return (
    <div><form className="form-inline">{content}</form></div>
    );
  }

}

React.render(<App/>, document.getElementById('fields'))


