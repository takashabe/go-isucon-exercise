import React from 'react';
import axios from 'axios';

export default class Enqueue extends React.Component {
  componentWillMount() {
    axios
      .post('/api/enqueue', {withCredentials: true})
      .then(res => {
        console.log(res);
      })
      .catch(e => console.log(JSON.stringify(e.response.data)));
  }

  render() {
    return <div>Enqueue</div>;
  }
}
