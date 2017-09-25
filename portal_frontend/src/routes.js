import React from 'react';
import {BrowserRouter, Route, Link} from 'react-router-dom';

import Login from './login.js';
import Dashboard from './dashboard.js';

export default class Routes extends React.Component {
  render() {
    return (
      <BrowserRouter>
        <div>
          <Route exact path="/" component={Dashboard} />
          <Route path="/login" component={Login} />
        </div>
      </BrowserRouter>
    );
  }
}
