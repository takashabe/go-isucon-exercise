import React from 'react';
import Typography from 'material-ui/Typography';
import PropTypes from 'prop-types';
import Table, {
  TableBody,
  TableCell,
  TableHead,
  TableRow,
} from 'material-ui/Table';
import {withStyles} from 'material-ui/styles';

const styles = theme => ({
  root: {
    width: '90%',
    marginTop: theme.spacing.unit * 3,
    marginLeft: 'auto',
    marginRight: 'auto',
  },
});

class Summary extends React.Component {
  constructor() {
    super();

    this.createScoreRow = this.createScoreRow.bind(this);
  }

  createScoreRow(label, data) {
    if (!data) {
      return (
        <TableRow>
          <TableCell padding="dense">{label}</TableCell>
          <TableCell>{'failed to receive score data'}</TableCell>
        </TableRow>
      );
    }
    return (
      <TableRow>
        <TableCell padding="dense">{label}</TableCell>
        <TableCell>{data.summary}</TableCell>
        <TableCell>{data.score}</TableCell>
        <TableCell>{data.submitted_at}</TableCell>
      </TableRow>
    );
  }

  render() {
    let highScore = null;
    let latestScore = null;
    for (const v of this.props.data) {
      if (highScore === null || highScore.score < v.score) {
        highScore = v;
      }
      if (latestScore === null || latestScore.score < v.score) {
        latestScore = v;
      }
    }

    return (
      <div className={this.props.classes.root}>
        <Typography type="display1">Summary</Typography>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell />
              <TableCell>Summary</TableCell>
              <TableCell>Score</TableCell>
              <TableCell>Timestamp</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {this.createScoreRow('High score', highScore)}
            {this.createScoreRow('Latest score', latestScore)}
          </TableBody>
        </Table>
      </div>
    );
  }
}

Summary.propTypes = {
  data: PropTypes.array.isRequired,
  team: PropTypes.object,
};

export default withStyles(styles)(Summary);
