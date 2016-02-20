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
var Schedule = React.createClass({

  getInitialState () {
    return {loading: true}
  },
  componentWillReceiveProps(nextProps){
    if (nextProps.team == null) {
      return;
    }
    var url = 'http://soccerlc.com/api/v1/games/' + nextProps.team.teamid
    axios.get(url)
      .then(function (response) {
        this.setState({games: response.data, loading: false});
      }.bind(this))
      .catch(function (response) {
        console.log(response);
      });
  },

  render () {
    var content = [];
    if (this.props.team) {
      if(this.state.loading){
        content = [<tr><td>loading</td></tr>]  
      }
      else{
        content = [<tr><th>When</th><th>Where</th><th>Home Team</th><th>Home Team Score</th><th>Away Team Score</th><th>Away Team</th></tr>]
        for (var i =0; i < this.state.games.length; i++) {
          content.push(this.renderItem(this.state.games[i]));
        };
      }
    }
    else{
      content = <tr></tr>
    }
    return( <div style={styles.tableContainer}><table className="table table-striped">
            <tbody>{content}</tbody>
            </table></div>);
  },

  renderItem (item) {
    var date = new Date(item.gamesdatetime)
    var hours = date.getUTCHours()
    var suffix = (hours >= 12)? 'PM' : 'AM';
    var parsedHour = ((hours + 11) % 12 + 1);
    var stringDate = weekday[date.getDay()] + " " + (date.getMonth() + 1) + "-" + (date.getDate()) + "-" + date.getFullYear() + " " + parsedHour + ":" + date.getMinutes() + suffix
    return (<tr><td>{stringDate}</td><td>{item.field}</td><td>{item.hometeam}</td><td>{item.hometeamscore}</td><td>{item.awayteamscore}</td><td>{item.awayteam}</td></tr>)
  }
})

module.exports = Schedule;

