import React from 'react';
import ReactDOM from 'react-dom';

export default class Login extends React.Component {
  constructor() {
    super();
    this.state = {
      inputId: null,
      inputPass: null,
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
    this.props.auth(this.state.inputId, this.state.inputPass);
  }

  render() {
    return (
      <div className="login">
        Hello Login
        <form className="loingForm" onSubmit={e => this.handleSubmit(e)}>
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
}
