import React, { Component } from 'react'
import './App.css'
import Endpoint from './Endpoint'
import Header from './Header'
import NavMenu from './NavMenu'

class App extends Component {
  constructor() {
    super()
    this.state = {
      endpoints: window.apiData.endpoints
    }
  }
  onSearch(v) {
     this.setState({endpoints: window.apiData.endpoints.filter(e => e.path.toLowerCase().indexOf(v.toLowerCase()) > -1)})
  }
  renderEndpoints() {
    if (this.state.endpoints.length === 0) {
      return (
        <div className="nodata">No endpoints found</div>
      )
    }
    return (
      <div className="endpoints">
        {this.state.endpoints.map(ep => 
          <Endpoint endpoint={ep} key={ep.name}/>
        )}
      </div>
    )
  }
  render() {
    return (
      <div>
        <Header onSearch={this.onSearch.bind(this)} />
        <div className="app-container">
          <NavMenu />
          <div className="endpoint-container">
            {this.renderEndpoints()}
          </div>
        </div>
      </div>
    );
  }
}

export default App;
