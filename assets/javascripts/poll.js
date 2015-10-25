import React, { Component, PropTypes } from 'react';
var Autocomplete = require('react-autocomplete');
var Schedule = require('./schedule');
var blah = 1231;

let styles = {
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

function findTeam (team, value) {
  return (
    team.name.toLowerCase().indexOf(value.toLowerCase()) !== -1
  )
}

function sortTeams (a, b, value) {
  return (
    a.name.toLowerCase().indexOf(value.toLowerCase()) >
    b.name.toLowerCase().indexOf(value.toLowerCase()) ? 1 : -1
  )
}

let App = React.createClass({

  getInitialState () {
        $.ajax({
          type: "GET", 
          url: "https://api.parse.com:443/1/classes/Teams",
          headers: 
		  { 
		    'X-Parse-Application-Id': 'UnWG5wrHS2fIl7xpzxHqStks4ei4sc6p0plxUOGv',
		    'X-Parse-REST-API-Key': 'g7Cj2NeORxfnKRXCHVv3ZcxxjRNpPU1RVuUxX19b'
		  },
          data: { "limit": "1000"},
          dataType: "json",
          success: function(response) {
          	this.setState({teams: response.results, loading: false});
       	  }.bind(this),
          error: function(xhr, ajaxOptions, thrownError) { alert(xhr.responseText); }
        });
    return {
      teams: [],
      team: null,
      loading: true
    }
  },
	fakeRequest (value, cb) {
		if (value === '')
			return this.state.teams
		var items = this.state.teams.filter((team) => {
			return findTeam(team, value)
		})
		setTimeout(() => {
		cb(items)
		}, 500)
	},
	onSelect(value, item){
		this.setState({team: item }) 
	},


  render () {
    return (
    <div>
      <div>
        <Autocomplete
          items={this.state.teams}
          getItemValue={(item) => item.name}
          sortItems={sortTeams}
		  inputProps={{"className": "input-lg", "style": {"border": "solid 1px #ccc", "width": "900px"}}}
          onSelect={this.onSelect}
          onChange={(event, value) => {
            this.setState({loading: true, team: null})
            this.fakeRequest(value, (items) => {
              this.setState({ teams: items, loading: false })
            })
          }}
          renderItem={(item, isHighlighted) => (
            <div
              style={isHighlighted ? styles.highlightedItem : styles.item}
              key={item.teamId}
              id={item.teamId}
              location={"Lets Play"}
            ><h1>{item.name}</h1><h2>{item.division}</h2></div>
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
          )}
        />
      </div>
      <Schedule
      team={this.state.team}
      loading={this.state.team ? true : false}
      />
      </div>
    )
  },

  renderItems (items) {
    return items.map((item, index) => {
      var text = item._store.originalProps.location
      if (index === 0 || items[index - 1]._store.originalProps.location !== text) {
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
})

React.render(<App/>, document.getElementById('app'))

