import page from 'page'
import React from 'react'
import App from './main'
import ReactDOM from 'react-dom'

const element = document.querySelector('#app')

page('/', function(ctx) {
    ReactDOM.render(<App/>, element)
})

//page('/standings/{team}', function(ctx) {
// team
    //ReactDOM.render(<standingsComponent team={ctx.params.team}/>, element)
//})

//page('/', function(ctx) {
    //ReactDOM.render(HomeComponent, element)
//})

page.start()
