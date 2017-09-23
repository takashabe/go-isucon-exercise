import React from 'react';
import ReactDOM from 'react-dom';

class Root extends React.Component {
  render() {
    return <p>Hello React!</p>;
  }
}

ReactDOM.render(<Root />, document.getElementById('root'));
