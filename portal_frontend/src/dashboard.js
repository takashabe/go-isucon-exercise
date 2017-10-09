import React from 'react';
import axios from 'axios';

import Score from './score.js';
import Enqueue from './enqueue.js';

export default class Dashboard extends React.Component {
  render() {
    return (
      <div>
        Hello Dashboard
        <Queues />
        <Enqueue />
        <Score />
      </div>
    );
  }
}

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
    axios
      .get('/api/queues', {withCredentials: true})
      .then(res => {
        this.setState({
          activeQueues: res.data,
        });
      })
      .catch(e => {
        console.log(JSON.stringify(e.response.data));
      });
  }

  render() {
    const style = {
      myTeam: {
        textDecoration: 'underline',
      },
    };
    const queues = this.state.activeQueues.map(x => {
      if (x.my_team) {
        return (
          <li key={x.message_id} style={style.myTeam}>
            {x.message_id}
          </li>
        );
      }
      return <li key={x.message_id}>{x.message_id}</li>;
    });
    return <ul>{queues}</ul>;
  }
}
