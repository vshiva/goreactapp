import React, { Component } from 'react';
import logo from './logo.svg';
import './App.css';

class App extends Component {

  constructor(props) {
    super(props);
    this.state = {
      mode: null,
    };
  }

  componentDidMount() {
    fetch('/api/config')
      .then(response => response.json())
      .then(data => this.setState({ mode: data.mode }));
  }

  styleData = () => {
      if (this.state.mode == null) {
        return null
      }
      return {backgroundColor: this.state.mode}
  }

  render() {
    let style = this.styleData()
    return (
      <div className="App" >
        <header className="App-header" style={style}>
          <img src={logo} className="App-logo" alt="logo" />
          <p>
            Cool Crazy App on K8S
          </p>
        </header>
      </div>
    );
  }
}

export default App;
