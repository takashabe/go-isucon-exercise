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
      detail: {
        data: null,
        message: '',
        open: false,
      },
    };

    this.handleOnClickDetail = this.handleOnClickDetail.bind(this);
    this.handleOnRequestDetailClose = this.handleOnRequestDetailClose.bind(
      this,
    );
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

  handleOnClickDetail(id) {
    axios
      .get('/api/bench_detail/' + id, {withCredentials: true})
      .then(res => {
        this.setState({
          detail: {
            open: true,
            message: '',
            data: JSON.parse(res.data.detail),
          },
        });
      })
      .catch(e => {
        this.setState({
          detail: {
            open: true,
            message: 'Failed to receive detail score',
            data: null,
          },
        });
      });
  }

  handleOnRequestDetailClose() {
    this.setState({
      detail: {
        open: false,
        message: '',
        data: null,
      },
    });
  }

  render() {
    return (
      <div>
        <Summary
          data={this.state.history}
          team={this.state.team}
          detail={this.state.detail}
          detailOpen={this.handleOnClickDetail}
          detailClose={this.handleOnRequestDetailClose}
        />
        <Queues />
        <Enqueue />
        <Score
          data={this.state.history}
          detail={this.state.detail}
          detailOpen={this.handleOnClickDetail}
          detailClose={this.handleOnRequestDetailClose}
        />
      </div>
    );
  }
}
