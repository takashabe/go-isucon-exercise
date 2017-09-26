import React from 'react';
import {BrowserRouter, Route, Link, Redirect} from 'react-router-dom';

import Login from './login.js';
import Dashboard from './dashboard.js';
import Enqueue from './enqueue.js';

export default class Routes extends React.Component {
  constructor() {
    super();

    const sess = sessionStorage.getItem('portal-session');
    this.state = {
      sessionId: sess,
    };

    this.isAuthentication = this.isAuthentication.bind(this);
    this.authenticate = this.authenticate.bind(this);
  }

  isAuthentication() {
    return this.state.sessionId !== null;
  }

  authenticate(id, password) {
    console.log('From Auth: ID=' + id + ', Pass=' + password);
    // dummy
    const sessionId = id;

    sessionStorage.setItem('portal-session', sessionId);
    this.setState({
      sessionId: sessionId,
    });
  }

  render() {
    return (
      <BrowserRouter>
        <div>
          <Header />
          <PropsRoute
            path="/login"
            component={Login}
            auth={this.authenticate}
          />
          <PrivateRoute
            exact
            path="/"
            component={Dashboard}
            auth={this.isAuthentication}
          />
          <PrivateRoute
            path="/enqueue"
            component={Enqueue}
            auth={this.isAuthentication}
          />
        </div>
      </BrowserRouter>
    );
  }
}

const Header = () => (
  <div>
    <p>Header</p>
    <ul>
      <li>
        <Link to="/">Dashboard</Link>
      </li>
      <li>
        <Link to="/login">Login</Link>
      </li>
      <li>
        <Link to="/enqueue">Enqueue</Link>
      </li>
    </ul>
  </div>
);

const PrivateRoute = ({component: Component, auth, ...rest}) => {
  return (
    <Route
      {...rest}
      render={props =>
        auth() === true ? (
          <Component {...props} />
        ) : (
          <Redirect to={{pathname: '/login', state: {from: props.location}}} />
        )}
    />
  );
};

const renderMergedProps = (component, ...rest) => {
  const finalProps = Object.assign({}, ...rest);
  return React.createElement(component, finalProps);
};

const PropsRoute = ({component, ...rest}) => {
  return (
    <Route
      {...rest}
      render={props => {
        return renderMergedProps(component, props, rest);
      }}
    />
  );
};
