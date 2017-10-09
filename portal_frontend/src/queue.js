import React from 'react';
import PropTypes from 'prop-types';
import {withStyles} from 'material-ui/styles';
import Avatar from 'material-ui/Avatar';
import deepOrange from 'material-ui/colors/deepOrange';
import deepPurple from 'material-ui/colors/deepPurple';
import axios from 'axios';
import Paper from 'material-ui/Paper';
import Typography from 'material-ui/Typography';
import Divider from 'material-ui/Divider';

const styles = theme => ({
  avatar: {
    margin: 5,
    width: 20,
    height: 20,
  },
  orangeAvatar: {
    margin: 5,
    width: 20,
    height: 20,
    color: '#fff',
    backgroundColor: deepOrange[500],
  },
  row: {
    display: 'flex',
    justifyContent: 'left',
  },
  paper: {
    width: '100%',
    marginTop: theme.spacing.unit,
    marginLeft: 'auto',
    marginRight: 'auto',
    overflowX: 'auto',
  },
  root: {
    width: '90%',
    marginTop: theme.spacing.unit * 3,
    marginLeft: 'auto',
    marginRight: 'auto',
  },
});

function QueueAvatar(props) {
  const {classes} = props;
  const data = props.data;
  const qa = data.map(x => {
    if (x.my_team === true) {
      return (
        <Avatar key={x.message_id} className={classes.orangeAvatar}>
          {' '}
        </Avatar>
      );
    } else {
      return (
        <Avatar key={x.message_id} className={classes.avatar}>
          {' '}
        </Avatar>
      );
    }
  });

  return (
    <div>
      <div className={classes.row}>{qa}</div>
      <Divider />
    </div>
  );
}

QueueAvatar.propTypes = {
  classes: PropTypes.object.isRequired,
  data: PropTypes.array.isRequired,
};

QueueAvatar = withStyles(styles)(QueueAvatar);

class Queues extends React.Component {
  constructor() {
    super();
    this.state = {
      activeQueues: [],
    };
  }

  componentWillMount() {
    this.handleQueues();
  }

  handleQueues() {
    axios.get('/api/queues', {withCredentials: true}).then(res => {
      this.setState({
        activeQueues: res.data,
      });
    });
  }

  render() {
    return (
      <div className={this.props.classes.root}>
        <Typography type="display1">Active queues</Typography>
        <QueueAvatar data={this.state.activeQueues} />
      </div>
    );
  }
}

export default withStyles(styles)(Queues);
