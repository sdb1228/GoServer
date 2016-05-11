import page from 'page'
import React from 'react'
import App from './main'
import Standings from './standings'
import ReactDOM from 'react-dom'

const element = document.querySelector('#app')

page('/', function(ctx) {
    ReactDOM.render(<App/>, element)
})

page('/standings/:division', function(ctx) {
    ReactDOM.render(<Standings division={ctx.params.division}/>, element)
})

//page('/', function(ctx) {
    //ReactDOM.render(HomeComponent, element)
//})

page.start()
