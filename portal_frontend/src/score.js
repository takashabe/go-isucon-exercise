import React from 'react';
import axios from 'axios';

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
    axios
      .get('/api/history', {withCredentials: true})
      .then(res => {
        console.log(JSON.stringify(res.data));
        this.setState({
          history: res.data,
        });
      })
      .catch(e => {
        console.log(e);
      });
  }

  handleDetail(id) {
    axios
      .get('/api/bench_detail/' + id, {withCredentials: true})
      .then(res => {
        console.log(JSON.stringify(res.data));
        this.setState({
          detail: res.data,
        });
      })
      .catch(e => {
        console.log(e);
      });
  }

  render() {
    const history = this.state.history.map(x => {
      const timestamp = new Date(x.submitted_at * 1000);
      return (
        <li key={x.id} onClick={id => this.handleDetail(x.id)}>
          Detail: {x.summary}-{x.score} at {timestamp.toLocaleString()}
        </li>
      );
    });
    // TODO: table layout
    return (
      <div>
        <ul>{history}</ul>
        <Detail source={this.state.detail} />
      </div>
    );
  }
}

class Detail extends React.Component {
  render() {
    console.log(this.props.detail);
    return null;
  }
}
