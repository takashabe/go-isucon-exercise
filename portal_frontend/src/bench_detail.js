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

const BenchDetail = props => {
  const {detail, detailOpen, detailClose} = props;
  let violations = (
    <TableRow>
      <TableCell>{'Empty violations. Your request is valid.'}</TableCell>
    </TableRow>
  );

  let detailContent = '';
  if (detail.data != null) {
    if (detail.data.violations.length > 0) {
      violations = detail.data.violations.map(x => {
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
              <TableCell>{detail.data.response.success}</TableCell>
              <TableCell>{detail.data.response.redirect}</TableCell>
              <TableCell>{detail.data.response.client_error}</TableCell>
              <TableCell>
                {detail.data.response.server_error +
                  detail.data.response.exception}
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
    <Dialog open={detail.open} onRequestClose={() => detailClose()}>
      <DialogTitle>{'Detail'}</DialogTitle>
      <DialogContent>
        <DialogContentText>{detail.message}</DialogContentText>
        {detailContent}
      </DialogContent>
    </Dialog>
  );
};

BenchDetail.propTypes = {
  detail: PropTypes.object,
  detailOpen: PropTypes.func.isRequired,
  detailClose: PropTypes.func.isRequired,
};

export default BenchDetail;
