import React from 'react';
import PropTypes from 'prop-types';
import axios from 'axios';
import {withStyles} from 'material-ui/styles';
import Table, {
  TableBody,
  TableCell,
  TableHead,
  TableRow,
} from 'material-ui/Table';
import Paper from 'material-ui/Paper';
import Dialog, {
  DialogContent,
  DialogContentText,
  DialogTitle,
} from 'material-ui/Dialog';

const styles = theme => ({
  paper: {
    width: '100%',
    marginTop: theme.spacing.unit,
    marginLeft: 'auto',
    marginRight: 'auto',
    overflowX: 'auto',
  },
});

class HistoryTable extends React.Component {
  constructor() {
    super();
    this.state = {
      open: false,
      message: null,
      detail: null,
    };
  }

  handleOnRequestClose() {
    this.setState({open: false});
  }

  handleOnClick(id) {
    axios
      .get('/api/bench_detail/' + id, {withCredentials: true})
      .then(res => {
        this.setState({
          open: true,
          detail: JSON.parse(res.data.detail),
          message: '',
        });
      })
      .catch(e => {
        this.setState({
          open: true,
          message: 'Failed to receive detail score',
        });
      });
  }

  render() {
    const classes = this.props.classes;
    const data = this.props.data;

    const detail = this.state.detail;
    let violations = (
      <TableRow>
        <TableCell>{'Empty violations. Your request is valid.'}</TableCell>
      </TableRow>
    );
    let detailContent = '';

    if (detail !== null) {
      if (detail.violations.length > 0) {
        violations = detail.violations.map(x => {
          return (
            <TableRow key={x.request_type}>
              <TableCell>{x.request_type}</TableCell>
              <TableCell>{x.num}</TableCell>
              <TableCell>{x.description}</TableCell>
            </TableRow>
          );
        });
      }

      detailContent = (
        <div>
          <DialogContentText>Summary</DialogContentText>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Request count</TableCell>
                <TableCell>Elapsed time</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              <TableRow>
                <TableCell>{detail.request_count}</TableCell>
                <TableCell>{detail.elapsed_time}</TableCell>
              </TableRow>
            </TableBody>
          </Table>

          <DialogContentText>Count of response code</DialogContentText>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>2xx</TableCell>
                <TableCell>3xx</TableCell>
                <TableCell>4xx</TableCell>
                <TableCell>5xx/timeout</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              <TableRow>
                <TableCell>{detail.response.success}</TableCell>
                <TableCell>{detail.response.redirect}</TableCell>
                <TableCell>{detail.response.client_error}</TableCell>
                <TableCell>
                  {detail.response.server_error + detail.response.exception}
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>

          <DialogContentText>Violations</DialogContentText>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Request type</TableCell>
                <TableCell>Num</TableCell>
                <TableCell>Description</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>{violations}</TableBody>
          </Table>
        </div>
      );
    }

    return (
      <Paper className={classes.paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Summary</TableCell>
              <TableCell numeric>Score</TableCell>
              <TableCell>Timestamp</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {data.map(n => {
              return (
                <TableRow
                  key={n.id}
                  hover
                  onClick={() => this.handleOnClick(n.id)}>
                  <TableCell>{n.summary}</TableCell>
                  <TableCell numeric>{n.score}</TableCell>
                  <TableCell>{n.timestamp}</TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
        <Dialog
          open={this.state.open}
          onRequestClose={() => this.handleOnRequestClose()}>
          <DialogTitle>{'Detail'}</DialogTitle>
          <DialogContent>
            <DialogContentText>{this.state.message}</DialogContentText>
            {detailContent}
          </DialogContent>
        </Dialog>
      </Paper>
    );
  }
}

HistoryTable.propTypes = {
  classes: PropTypes.object.isRequired,
  data: PropTypes.array.isRequired,
};

export default withStyles(styles)(HistoryTable);
