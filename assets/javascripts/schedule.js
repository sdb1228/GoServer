import React, { Component, PropTypes } from 'react';

let styles = {
  item: {
    padding: '2px 6px',
    cursor: 'default'
  },
  div: {
    margin: '0 auto',
  },
  highlightedItem: {
    color: 'white',
    background: 'hsl(200, 50%, 50%)',
    padding: '2px 6px',
    cursor: 'default'
  },
  divStyle: {
    "background-color": 'black'
  }, 
  menu: {
    borderRadius: '3px',
    boxShadow: '0 2px 12px rgba(0, 0, 0, 0.1)',
    background: 'rgba(255, 255, 255, 0.9)',
    padding: '2px 0',
    fontSize: '90%',
    position: 'fixed',
    overflow: 'auto',
    maxHeight: '50%',
  }
}
var Schedule = React.createClass({

  getInitialState () {
    return {loading: true}
  },
  componentWillReceiveProps(nextProps){
    if (nextProps.team == null) {
      return;
    }
    $.ajax({
          type: "GET", 
          url: "https://api.parse.com:443/1/classes/Games",
          headers: 
      { 
        'X-Parse-Application-Id': 'UnWG5wrHS2fIl7xpzxHqStks4ei4sc6p0plxUOGv',
        'X-Parse-REST-API-Key': 'g7Cj2NeORxfnKRXCHVv3ZcxxjRNpPU1RVuUxX19b'
      },
          data: {"where": {"$or":[{"awayTeam": nextProps.team.teamId},{"homeTeam":nextProps.team.teamId}]}},
          dataType: "json",
          success: function(response) {
            this.getTeamNames(response.results)
            debugger
            this.setState({teams: response.results, loading: false});
          }.bind(this),
          error: function(xhr, ajaxOptions, thrownError) { alert(xhr.responseText); }
    });
  },
  getTeamNames(games){
    for (var i = 0; i < games.length; i++) {
      var homeTeam = games[i].homeTeam; 
      var awayTeam = games[i].awayTeam;
      if (homeTeam == this.props.team.teamId) {
        $.ajax({
              type: "GET", 
              url: "https://api.parse.com/1/classes/Teams",
              headers: 
          { 
            'X-Parse-Application-Id': 'UnWG5wrHS2fIl7xpzxHqStks4ei4sc6p0plxUOGv',
            'X-Parse-REST-API-Key': 'g7Cj2NeORxfnKRXCHVv3ZcxxjRNpPU1RVuUxX19b'
          },
              data: {"where": {"$relatedTo":{"object":{"__type": "Pointer", "className": "Games", "objectId": games[i].objectId}, "key": "awayTeamPointer"}}},
              dataType: "json",
              success: function(response) {
                var elementPos = this.state.teams.map(function(x) {return x.awayTeam; }).indexOf(response.results[0].teamId);
                this.state.teams[elementPos].awayTeamName = response.results[0].name
              }.bind(this),
              error: function(xhr, ajaxOptions, thrownError) { alert(xhr.responseText); }
        });
      }
      else{
        //home team needed
        awayTeam = this.props.team.name;
      }
    }
  },

  render () {
    var content = [];
    if (this.props.team) {
      if(this.state.loading){
        content = [<tr><td>loading</td></tr>]  
      }
      else{
        content = [<tr><td>When</td><td>Where</td><td>Home Team</td><td>Home Team Score</td><td>Away Team Score</td><td>Away Team</td></tr>]
        for (var i =0; i < this.state.teams.length; i++) {
          content.push(this.renderItem(this.state.teams[i]));
        };
      }
    }
    else{
      content = <tr><td></td></tr>
    }
    return( <div><table style={styles.div}>
            <tbody>{content}</tbody>
            </table></div>);
  },

  renderItem (item) {
    var homeTeam = item.homeTeam;
    var awayTeam = item.awayTeam;
    if (item.homeTeam == this.props.team.teamId) {
      homeTeam = this.props.team.name;
    }
    else{
      awayTeam = this.props.team.name;
    }
    return (<tr><td>{item.date}</td><td>{item.field}</td><td>{homeTeam}</td><td>{item.homeTeamScore}</td><td>{item.awayTeamScore}</td><td>{awayTeam}</td></tr>)
  }
})

module.exports = Schedule;

