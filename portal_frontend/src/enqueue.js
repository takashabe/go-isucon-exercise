import React from 'react';
import axios from 'axios';

import {withStyles} from 'material-ui/styles';
import Button from 'material-ui/Button';
import Dialog, {
  DialogContent,
  DialogContentText,
  DialogTitle,
} from 'material-ui/Dialog';

const styles = theme => ({
  root: {
    marginTop: theme.spacing.unit * 3,
    textAlign: 'center',
  },
});

class Enqueue extends React.Component {
  constructor() {
    super();
    this.state = {
      open: false,
      message: null,
    };

    this.handleOnRequestClose = this.handleOnRequestClose.bind(this);
    this.handleOnClick = this.handleOnClick.bind(this);
  }

  handleOnRequestClose() {
    this.setState({open: false});
  }

  handleOnClick() {
    axios
      .post('/api/enqueue', {withCredentials: true})
      .then(res => {
        this.setState({
          open: true,
          message: 'Success send request queue',
        });
      })
      .catch(e => {
        this.setState({
          open: true,
          message:
            'Failed to send request queue. Receive error message:' +
            JSON.stringify(e.response.data),
        });
      });
  }

  render() {
    return (
      <div className={this.props.classes.root}>
        <Button raised color="primary" onClick={this.handleOnClick}>
          Enqueue
        </Button>
        <Dialog
          open={this.state.open}
          onRequestClose={this.handleOnRequestClose}>
          <DialogTitle>{'Enqueue'}</DialogTitle>
          <DialogContent>
            <DialogContentText>{this.state.message}</DialogContentText>
          </DialogContent>
        </Dialog>
      </div>
    );
  }
}

export default withStyles(styles)(Enqueue);
