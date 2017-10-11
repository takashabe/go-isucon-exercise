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

class BenchDetail extends React.Component {
  constructor() {
    super();
    this.state = {
      open: false,
      message: '',
      data: null,
    };

    this.hanldeOnRequestClose = this.hanldeOnRequestClose.bind(this);
  }

  componentWillReceiveProps(nextProps) {
    this.setState({
      open: nextProps.detail.open,
      message: nextProps.detail.message,
      data: nextProps.detail.data,
    });
  }

  hanldeOnRequestClose() {
    this.setState({open: false});
  }

  render() {
    const detail = this.state.data;
    let violations = (
      <TableRow>
        <TableCell>{'Empty violations. Your request is valid.'}</TableCell>
      </TableRow>
    );

    let detailContent = '';
    if (detail != null) {
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
      <Dialog
        maxWidth="md"
        open={this.state.open}
        onRequestClose={() => this.hanldeOnRequestClose()}>
        <DialogTitle>{'Detail'}</DialogTitle>
        <DialogContent>
          <DialogContentText>{this.state.message}</DialogContentText>
          {detailContent}
        </DialogContent>
      </Dialog>
    );
  }
}

BenchDetail.propTypes = {
  detail: PropTypes.object.isRequired,
};

export default BenchDetail;
