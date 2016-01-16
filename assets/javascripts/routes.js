import React from 'react'
import { Route, IndexRoute } from 'react-router'

import App from './components/App'
import Home from './components/Home'
import Jobs from './components/Jobs'
import Builds from './components/Builds'

let routes = (
  <Route path="/" component={App}>
    <IndexRoute component={Home} />
    <Route path="jobs" component={Jobs} />
    <Route path="jobs/:job/builds" component={Builds} />
  </Route>
)
export default routes