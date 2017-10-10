import React from 'react';
import axios from 'axios';

import Summary from './summary.js';
import Queues from './queue.js';
import Score from './score.js';
import Enqueue from './enqueue.js';

export default class Dashboard extends React.Component {
  constructor() {
    super();
    this.state = {
      history: [],
      team: null,
    };
  }

  componentWillMount() {
    this.updateTeam();
    this.updateHistory();
  }

  updateTeam() {
    axios.get('/api/team', {withCredentials: true}).then(res => {
      this.setState({
        team: res.data,
      });
    });
  }

  updateHistory() {
    axios.get('/api/history', {withCredentials: true}).then(res => {
      this.setState({
        history: res.data,
      });
    });
  }

  render() {
    return (
      <div>
        <Summary data={this.state.history} team={this.state.team} />
        <Queues />
        <Enqueue />
        <Score data={this.state.history} />
      </div>
    );
  }
}
