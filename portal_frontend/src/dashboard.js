import React from 'react';
import axios from 'axios';
import 'babel-polyfill';

import Summary from './summary.js';
import Queues from './queue.js';
import Score from './score.js';
import Enqueue from './enqueue.js';
import BenchDetail from './bench_detail.js';

export default class Dashboard extends React.Component {
  constructor() {
    super();
    this.state = {
      history: [],
      team: null,
      detail: {
        data: null,
        message: '',
        open: false,
      },
    };

    this.handleOnClickDetail = this.handleOnClickDetail.bind(this);
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

  async handleOnClickDetail(id) {
    const state = await axios
      .get('/api/bench_detail/' + id, {withCredentials: true})
      .then(res => {
        return {
          open: true,
          message: '',
          data: JSON.parse(res.data.detail),
        };
      })
      .catch(e => {
        return {
          open: true,
          message: 'Failed to receive detail score',
          data: null,
        };
      });

    this.setState({detail: state});
  }

  render() {
    return (
      <div>
        <Summary
          data={this.state.history}
          team={this.state.team}
          detailOpen={this.handleOnClickDetail}
        />
        <Queues />
        <Enqueue />
        <Score
          data={this.state.history}
          detailOpen={this.handleOnClickDetail}
        />
        <BenchDetail detail={this.state.detail} />
      </div>
    );
  }
}
