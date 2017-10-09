import React from 'react';
import axios from 'axios';
import Typography from 'material-ui/Typography';

import HistoryTable from './table.js';

const styles = {
  width: '90%',
  margin: 'auto',
};

export default class Score extends React.Component {
  constructor() {
    super();
    this.state = {
      history: [],
      detail: null,
    };
  }

  componentWillMount() {
    this.updateHistory();
  }

  updateHistory() {
    axios.get('/api/history', {withCredentials: true}).then(res => {
      this.setState({
        history: res.data,
      });
    });
  }

  render() {
    const data = this.state.history.map(x => {
      const timestamp = new Date(x.submitted_at * 1000);
      return {
        id: x.id,
        summary: x.summary,
        score: x.score,
        timestamp: timestamp.toLocaleString(),
      };
    });

    return (
      <div style={styles}>
        <Typography type="display1">Scores</Typography>
        <HistoryTable data={data} />
      </div>
    );
  }
}
