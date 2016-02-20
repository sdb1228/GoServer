/**
 * Copyright © 2015-2016 SoccerLC. All rights reserved.
 *
 */
import React, { Component, PropTypes } from 'react';
var Autocomplete = require('react-autocomplete');
var Schedule = require('./schedule');
var axios = require('axios');
var styles = {
  item: {
    padding: '2px 6px',
    cursor: 'default'
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

class App extends Component {

findTeam (team, value) {
  return (
    team.name.toLowerCase().indexOf(value.toLowerCase()) !== -1
  )
}

sortTeams (a, b, value) {
  return (
    a.name.toLowerCase().indexOf(value.toLowerCase()) >
    b.name.toLowerCase().indexOf(value.toLowerCase()) ? 1 : -1
  )
}

fakeRequest (value, cb) {
    if (value === '')
        return this.state.fullTeams
    var items = this.state.fullTeams.filter((team) => {
        return this.findTeam(team, value)
    })
    setTimeout(() => {
    cb(items)
    }, 500)
}
onSelect(value, item){
  this.setState({teams: this.state.fullTeams, loading: false, team: item});
}

  constructor(props) {
    super(props);
    this.selectedTeam = "lets_play";
    this.state = {count: props.initialCount};
    axios.get('http://soccerlc.com/api/v1/teams/6')
      .then(function (response) {
        this.success(response.data);
      }.bind(this))
      .catch(function (response) {
        console.log(response);
      });
  }
  success(data){
    this.setState({teams: data, loading: false, selected: this.selectedTeam, fullTeams: data});
  }
renderItems (items) {
  return items.map((item, index) => {
    var text = item.props.children[0].props.children.charAt(0)
    if (index === 0 || items[index - 1].props.children[0].props.children.charAt(0) !== text) {
      var style = {
        background: '#eee',
        color: '#454545',
        padding: '2px 6px',
        fontWeight: 'bold'
      }
      return [<div style={style}>{text}</div>, item]
    }
    else {
      return item
    }
  })
}
facilityChoose(facility, event){
  event.target.parentElement.className="active";
  $("#" + this.state.selected).toggleClass("active");
  var url = "http://soccerlc.com/api/v1/teams/" + event.target.id
  this.selectedTeam = event.target.parentElement.id
  axios.get(url)
  .then(function (response) {
    this.setState({teams: response.data, loading: false, selected: this.selectedTeam, fullTeams:response.data});
  }.bind(this))
  .catch(function (response) {
    console.log(response);
  });
}
  static contextTypes = {
    onSetTitle: PropTypes.func.isRequired,
  };
  render() {
    return (
      <div>
        <div className="root">
          <ul className="nav nav-pills">
              <li id="lets_play" role="presentation" className="active" ><a id="6" href="#" onClick={this.facilityChoose.bind(this, 'lets_play')}>Lets Play Soccer</a></li>
              <li id="soccer_city" role="presentation"><a id="5" href="#" onClick={this.facilityChoose.bind(this, 'soccer_city')}>Soccer City</a></li>
              <li id="utah_soccer" role="presentation"><a id="1" href="#" onClick={this.facilityChoose.bind(this, 'utah_soccer')}>Utah Soccer</a></li>
              <li id="uysa_boys" role="presentation"><a id="4" href="#" onClick={this.facilityChoose.bind(this, 'uysa_boys')}>UYSA Boys</a></li>
              <li id="uysa_girls" role="presentation"><a id="3" href="#" onClick={this.facilityChoose.bind(this, 'uysa_girls')}>UYSA Girls</a></li>
            </ul>
            <Autocomplete
              items={this.state.teams}
              getItemValue={(item) => item.name}
              inputProps={{"className": "form-control input-lg", "style": {"border": "solid 1px #ccc", "width": "900px", "marginTop": "20px"}}}
              onSelect={this.onSelect.bind(this)}
              onChange={(event, value) => {
                this.setState({loading: true, team: null})
                this.fakeRequest(value, (items) => {
                  this.setState({ teams: items, loading: false })
                })
              }}
              renderItem={(item, isHighlighted) => (
                <div
                  style={isHighlighted ? styles.highlightedItem : styles.item}
                  key={item.Id}
                  id={item.Id}
                  location={item.facility}
                ><h1>{item.name}</h1><div>{item.division}</div></div>
              )}
              renderMenu={(items, value, style) => (
                <div style={{...styles.menu, ...style}}>
                  {value === '' ? (
                    <div style={{padding: 6}}>Type the name of a Team</div>
                  ) : this.state.loading ? (
                    <div style={{padding: 6}}>Loading...</div>
                  ) : items.length === 0 ? (
                    <div style={{padding: 6}}>No matches for {value}</div>
                  ) : this.renderItems(items)}
                </div>
              )} />
            </div>
            <Schedule
              team={this.state.team}
              loading={this.state.team ? true : false}/>
          </div>
    );
  }

}

React.render(<App/>, document.getElementById('app'))


