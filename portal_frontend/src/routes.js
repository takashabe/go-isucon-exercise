import React from 'react';
import {BrowserRouter, Route, Link, Redirect} from 'react-router-dom';
import axios from 'axios';

import Login from './login.js';
import Dashboard from './dashboard.js';
import Enqueue from './enqueue.js';

export default class Routes extends React.Component {
  constructor() {
    super();

    this.state = {
      authed: sessionStorage.getItem('authed_session') !== null,
    };
    this.isAuthentication = this.isAuthentication.bind(this);
    this.authenticate = this.authenticate.bind(this);
    this.updateAuthentication = this.updateAuthentication.bind(this);
  }

  isAuthentication() {
    if (this.state.authed === false) {
      return false;
    }

    // expire session
    axios.get('/api/team', {withCredentials: true}).catch(e => {
      if (e.status === 401) {
        this.updateAuthentication(false);
      }
    });
    return this.state.authed;
  }

  updateAuthentication(isAuth) {
    if (isAuth) {
      sessionStorage.setItem('authed_session', true);
    } else {
      sessionStorage.removeItem('authed_session');
    }
    this.setState({
      authed: isAuth,
    });
  }

  authenticate(id, password) {
    let params = new URLSearchParams();
    params.append('email', id);
    params.append('password', password);

    axios
      .post('/api/login', params, {withCredentials: true})
      .then(res => {
        console.log(res.status + ': ' + JSON.stringify(res.data));
        this.updateAuthentication(true);
      })
      .catch(e => {
        console.log(e.response.status + ': ' + JSON.stringify(e.response.data));
        this.updateAuthentication(false);
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
