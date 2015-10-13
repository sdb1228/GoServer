import React, { Component, PropTypes } from 'react';
var Autocomplete = require('react-autocomplete');
var Parse = require('parse').Parse;
// ParseReact sits on top of your Parse singleton
var ParseReact = require('parse-react');
var ParseComponent = ParseReact.Component(React);
Parse.initialize('UnWG5wrHS2fIl7xpzxHqStks4ei4sc6p0plxUOGv', '1ZwTDp8rOpTqbZEv6oEtaUHJNK985Fi9pEJWAfDQ');

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
    border: 'solid 1px #ccc'
  }
}

function getStates() {
  return [
    { abbr: "AL", name: "Alabama"},
    { abbr: "AK", name: "Alaska"},
    { abbr: "AZ", name: "Arizona"},
    { abbr: "AR", name: "Arkansas"},
    { abbr: "CA", name: "California"},
    { abbr: "CO", name: "Colorado"},
    { abbr: "CT", name: "Connecticut"},
    { abbr: "DE", name: "Delaware"},
    { abbr: "FL", name: "Florida"},
    { abbr: "GA", name: "Georgia"},
    { abbr: "HI", name: "Hawaii"},
    { abbr: "ID", name: "Idaho"},
    { abbr: "IL", name: "Illinois"},
    { abbr: "IN", name: "Indiana"},
    { abbr: "IA", name: "Iowa"},
    { abbr: "KS", name: "Kansas"},
    { abbr: "KY", name: "Kentucky"},
    { abbr: "LA", name: "Louisiana"},
    { abbr: "ME", name: "Maine"},
    { abbr: "MD", name: "Maryland"},
    { abbr: "MA", name: "Massachusetts"},
    { abbr: "MI", name: "Michigan"},
    { abbr: "MN", name: "Minnesota"},
    { abbr: "MS", name: "Mississippi"},
    { abbr: "MO", name: "Missouri"},
    { abbr: "MT", name: "Montana"},
    { abbr: "NE", name: "Nebraska"},
    { abbr: "NV", name: "Nevada"},
    { abbr: "NH", name: "New Hampshire"},
    { abbr: "NJ", name: "New Jersey"},
    { abbr: "NM", name: "New Mexico"},
    { abbr: "NY", name: "New York"},
    { abbr: "NC", name: "North Carolina"},
    { abbr: "ND", name: "North Dakota"},
    { abbr: "OH", name: "Ohio"},
    { abbr: "OK", name: "Oklahoma"},
    { abbr: "OR", name: "Oregon"},
    { abbr: "PA", name: "Pennsylvania"},
    { abbr: "RI", name: "Rhode Island"},
    { abbr: "SC", name: "South Carolina"},
    { abbr: "SD", name: "South Dakota"},
    { abbr: "TN", name: "Tennessee"},
    { abbr: "TX", name: "Texas"},
    { abbr: "UT", name: "Utah"},
    { abbr: "VT", name: "Vermont"},
    { abbr: "VA", name: "Virginia"},
    { abbr: "WA", name: "Washington"},
    { abbr: "WV", name: "West Virginia"},
    { abbr: "WI", name: "Wisconsin"},
    { abbr: "WY", name: "Wyoming"}
  ]
}


function matchStateToTerm (state, value) {
  return (
    state.name.toLowerCase().indexOf(value.toLowerCase()) !== -1 ||
    state.abbr.toLowerCase().indexOf(value.toLowerCase()) !== -1
  )
}

function sortStates (a, b, value) {
  return (
    a.name.toLowerCase().indexOf(value.toLowerCase()) >
    b.name.toLowerCase().indexOf(value.toLowerCase()) ? 1 : -1
  )
}
function fakeRequest (value, cb) {
  if (value === '')
    return getStates()
  var items = getStates().filter((state) => {
    return matchStateToTerm(state, value)
  })
  setTimeout(() => {
    cb(items)
  }, 500)
}


class App extends ParseComponent {
  componentWillMount () {
  	// look into using getInitialState... but later.
    this.setState({
      unitedStates: getStates(),
      loading: false
    });
  }

  observe  (props, state) {
    return {
      teams: (new Parse.Query('Teams')).limit(1000)
    };
  }

  render () {
    return (
      <div>
        <Autocomplete
          items={this.state.unitedStates}
          getItemValue={(item) => item.name}
		  inputProps={{"className": "input-lg", "style": {"border": "solid 1px #ccc", "width": "900px"}}}
          onSelect={() => this.setState({ unitedStates: [] }) }
          onChange={(event, value) => {
            this.setState({loading: true})
            fakeRequest(value, (items) => {
              this.setState({ unitedStates: items, loading: false })
            })
          }}
          renderItem={(item, isHighlighted) => (
            <div
              style={isHighlighted ? styles.highlightedItem : styles.item}
              key={item.abbr}
              id={item.abbr}
            ><h1>{item.name}</h1><h2>blah</h2></div>
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
    )
  }

  renderItems (items) {
    return items.map((item, index) => {
      var text = item.props.children[0].props.children
      if (index === 0 || items[index - 1].props.children[0].props.children.charAt(0) !== text.charAt(0)) {
        var style = {
          background: '#eee',
          color: '#454545',
          padding: '2px 6px',
          fontWeight: 'bold'
        }
        return [<div style={style}>{text.charAt(0)}</div>, item]
      }
      else {
        return item
      }
    })
  }
}

React.render(<App/>, document.getElementById('app'))

