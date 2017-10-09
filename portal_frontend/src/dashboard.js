import React from 'react';
import axios from 'axios';

import Queues from './queue.js';
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
