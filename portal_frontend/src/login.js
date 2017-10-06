import React from 'react';
import ReactDOM from 'react-dom';
import {Redirect} from 'react-router-dom';
import axios from 'axios';

export default class Login extends React.Component {
  constructor() {
    super();
    this.state = {
      inputId: null,
      inputPass: null,
      authed: false,
      message: null,
    };
  }

  handleInput(e) {
    const target = e.target;
    const value = target.value;
    const name = target.name;

    this.setState({
      [name]: value,
    });
  }

  handleSubmit(e) {
    e.preventDefault();

    const id = this.state.inputId;
    const pass = this.state.inputPass;
    if (!id) {
      this.setState({
        message: 'require TeamID',
      });
      return;
    } else if (!pass) {
      this.setState({
        message: 'require Password',
      });
      return;
    }

    let params = new URLSearchParams();
    params.append('email', id);
    params.append('password', pass);

    axios
      .post('/api/login', params, {withCredentials: true})
      .then(res => {
        this.props.updateSession(true);
        this.setState({authed: true});
      })
      .catch(e => {
        this.props.updateSession(false);
        this.setState({
          authed: false,
          message: 'Incorrect team id or password',
        });
      });
  }

  render() {
    let component;
    if (this.state.authed) {
      component = <Redirect push to={'/'} />;
    } else {
      component = (
        <div className="login">
          <p>{this.state.message}</p>
          Login
          <form className="loginForm" onSubmit={e => this.handleSubmit(e)}>
            Team ID:
            <input
              type="text"
              name="inputId"
              onChange={e => this.handleInput(e)}
            />
            Password:
            <input
              type="text"
              name="inputPass"
              onChange={e => this.handleInput(e)}
            />
            <input type="submit" value="Submit" />
          </form>
        </div>
      );
    }

    return <div>{component}</div>;
  }
}
