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
import Paper from 'material-ui/Paper';

const styles = theme => ({
  root: {
    width: '90%',
    marginTop: theme.spacing.unit * 3,
    marginLeft: 'auto',
    marginRight: 'auto',
  },
  paper: {
    width: '100%',
    marginTop: theme.spacing.unit,
    marginLeft: 'auto',
    marginRight: 'auto',
    overflowX: 'auto',
  },
  paperTable: {
    width: '90%',
    marginLeft: 'auto',
    marginRight: 'auto',
  },
  secondary: {
    marginTop: theme.spacing.unit,
    marginLeft: theme.spacing.unit,
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
          <TableCell>Failed to receive score data</TableCell>
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
    const {classes, team} = this.props;

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

    const teamRow =
      team !== null ? (
        <TableRow>
          <TableCell>{team.ID}</TableCell>
          <TableCell>{team.Name}</TableCell>
          <TableCell>{team.Instance}</TableCell>
        </TableRow>
      ) : (
        <TableRow>
          <TableCell>Failed to receive team data</TableCell>
        </TableRow>
      );

    return (
      <div className={classes.root}>
        <Typography type="display1">Summary</Typography>
        <Paper className={classes.paper}>
          <Typography
            className={classes.secondary}
            color="secondary"
            type="subheading">
            Team detail
          </Typography>
          <Table className={classes.paperTable}>
            <TableHead>
              <TableRow>
                <TableCell>Team ID</TableCell>
                <TableCell>Name</TableCell>
                <TableCell>Instance</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>{teamRow}</TableBody>
          </Table>

          <Typography
            className={classes.secondary}
            color="secondary"
            type="subheading">
            Highlight scores
          </Typography>
          <Table className={classes.paperTable}>
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
        </Paper>
      </div>
    );
  }
}

Summary.propTypes = {
  classes: PropTypes.object.isRequired,
  data: PropTypes.array.isRequired,
  team: PropTypes.object,
};

export default withStyles(styles)(Summary);
