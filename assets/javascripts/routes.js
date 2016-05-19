import page from 'page'
import React from 'react'
import App from './main'
import Field from './fields'
import Standings from './standings'
import ReactDOM from 'react-dom'

const element = document.querySelector('#app')

page('/', function(ctx) {
    ReactDOM.render(<App/>, element)
})

page('/standings/:division', function(ctx) {
    ReactDOM.render(<Standings division={ctx.params.division}/>, element)
})

page('/fields/update', function(ctx) {
    ReactDOM.render(<Field/>, element)
})

page.start()
