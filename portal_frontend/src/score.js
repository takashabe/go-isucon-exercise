import React from 'react';
import axios from 'axios';
import Typography from 'material-ui/Typography';
import PropTypes from 'prop-types';

import HistoryTable from './table.js';

const styles = {
  width: '90%',
  margin: 'auto',
};

const Score = props => {
  const data = props.data;
  const histories = data.map(x => {
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
};

Score.prototype = {
  data: PropTypes.array.isRequired,
};

export default Score;
