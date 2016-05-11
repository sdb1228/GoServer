/**
 * Copyright © 2014-2016 SoccerLC. All rights reserved.
 *
 */

import React, { Component, PropTypes } from 'react';
var axios = require('axios');

let styles = {
  div: {
    margin: '0 auto'
  },
  tableContainer: {
    marginTop: '20px'
  },
  linkDiv: {
    float: "left",
  },
  linkContainer: {
    margin: '0 auto'
  }
}
let weekday = new Array(7);
weekday[0]=  "Sunday";
weekday[1] = "Monday";
weekday[2] = "Tuesday";
weekday[3] = "Wednesday";
weekday[4] = "Thursday";
weekday[5] = "Friday";
weekday[6] = "Saturday";
var Standings = React.createClass({

  componentWillMount(){
    const url = "http://soccerlc.com/api/v1/standings/" + this.props.division;
    axios.get(url)
      .then(function (response) {
        this.success(response.data);
      }.bind(this))
      .catch(function (response) {
        console.log(response);
      });
  },

  printClicked (event){
    var divToPrint=document.getElementById("teamSchedule");
    var newWin=window.open("");
    newWin.document.write(divToPrint.outerHTML);
    newWin.print();
    newWin.close();
  },
  success(data){
    this.setState({standings: data, loading: false});
  },

  getInitialState () {
    return {loading: true}
  },

  render () {
    var content = [];
    var tableClass= "table table-striped";
    var printButton = [];
    var searchLink = [];
    if (this.state.standings) {
      searchLink = [<a style={styles.linkDiv} href="/">Back to Search</a>];
      printButton = [<button className="hvr-grow-shadow btn btn-primary btn-lg" onClick={this.printClicked}>Print me</button>];
      tableClass = tableClass + " shadow";
      if(this.state.loading){
        content = [<tr><td>loading</td></tr>]
      }
      else{
        content = [<tr><th>Position</th><th>Name</th><th>GP</th><th>GF</th><th>GA</th><th>Points</th></tr>]
        for (var i =0; i < this.state.standings.length; i++) {
          content.push(this.renderItem(this.state.standings[i], i+1));
        };
      }
    }
    else{
      content = <tr></tr>
    }
    return( <div className="tableSize" style={styles.tableContainer}><div style={styles.linkContainer}>{searchLink}</div><table id="teamSchedule" className={tableClass}>
            <tbody>{content}</tbody>
            </table>{printButton}</div>);
  },

  renderItem (item, position) {
    return (<tr><td>{position}</td><td>{item.teamname}</td><td>{item.gamesplayed}</td><td>{item.goalsfor}</td><td>{item.goalsagainst}</td><td>{item.points}</td></tr>)
  }
})

module.exports = Standings;

