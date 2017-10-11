import React from 'react';
import axios from 'axios';
import Typography from 'material-ui/Typography';
import PropTypes from 'prop-types';
import Table, {
  TableBody,
  TableCell,
  TableHead,
  TableRow,
} from 'material-ui/Table';
import Paper from 'material-ui/Paper';
import {withStyles} from 'material-ui/styles';

import BenchDetail from './bench_detail.js';

const styles = theme => ({
  root: {
    width: '90%',
    margin: 'auto',
  },
  paper: {
    width: '100%',
    margin: 'auto',
  },
});

const Score = props => {
  const {classes, data, detail, detailOpen, detailClose} = props;
  const histories = data.map(x => {
    const timestamp = new Date(x.submitted_at * 1000);
    return {
      id: x.id,
      summary: x.summary,
      score: x.score,
      timestamp: timestamp,
    };
  });

  return (
    <div className={classes.root}>
      <Typography type="display1">Scores</Typography>
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
                <TableRow key={n.id} hover onClick={() => detailOpen(n.id)}>
                  <TableCell>{n.summary}</TableCell>
                  <TableCell numeric>{n.score}</TableCell>
                  <TableCell>{n.timestamp}</TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
        <BenchDetail
          detail={detail}
          detailOpen={detailOpen}
          detailClose={detailClose}
        />
      </Paper>
    </div>
  );
};

Score.prototype = {
  classes: PropTypes.object.isRequired,
  data: PropTypes.array.isRequired,
  detail: PropTypes.object.isRequired,
  detailOpen: PropTypes.func.isRequired,
  detailClose: PropTypes.func.isRequired,
};

export default withStyles(styles)(Score);
