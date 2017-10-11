import React from 'react';
import ReactDOM from 'react-dom';
import {Redirect} from 'react-router-dom';
import axios from 'axios';
import Input, {InputLabel} from 'material-ui/Input';
import {FormControl, FormHelperText} from 'material-ui/Form';
import PropTypes from 'prop-types';
import {withStyles} from 'material-ui/styles';
import Typography from 'material-ui/Typography';
import Button from 'material-ui/Button';

const styles = theme => ({
  root: {
    margin: theme.spacing.unit * 3,
  },
  formControl: {
    margin: theme.spacing.unit,
  },
});

class Login extends React.Component {
  constructor() {
    super();
    this.state = {
      inputId: null,
      inputPass: null,
      authed: false,
      failMessage: null,
    };

    this.handleClick = this.handleClick.bind(this);
  }

  handleChange(event) {
    const target = event.target;
    const value = target.value;
    const id = target.id;

    this.setState({
      [id]: value,
    });
  }

  handleClick(event) {
    event.preventDefault();

    let params = new URLSearchParams();
    params.append('email', this.state.inputId);
    params.append('password', this.state.inputPass);

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
          failMessage: 'Incorrect team id or password',
        });
      });
  }

  render() {
    const {classes, updateSession} = this.props;
    const message =
      this.state.failMessage !== null ? (
        <Typography color="error" type="subheading">
          {this.state.failMessage}
        </Typography>
      ) : (
        ''
      );

    const component =
      this.state.authed === true ? (
        <Redirect push to={'/'} />
      ) : (
        <div className={classes.root}>
          {message}
          <div className={classes.container}>
            <FormControl className={classes.formControl}>
              <InputLabel shrink>Email</InputLabel>
              <Input id="inputId" onChange={e => this.handleChange(e)} />
            </FormControl>
            <FormControl className={classes.formControl}>
              <InputLabel shrink>Password</InputLabel>
              <Input
                id="inputPass"
                type="password"
                onChange={e => this.handleChange(e)}
              />
            </FormControl>
            <Button raised color="primary" onClick={e => this.handleClick(e)}>
              Enqueue
            </Button>
          </div>
        </div>
      );

    return <div>{component}</div>;
  }
}

Login.propTypes = {
  classes: PropTypes.object.isRequired,
  updateSession: PropTypes.func.isRequired,
};

export default withStyles(styles)(Login);
