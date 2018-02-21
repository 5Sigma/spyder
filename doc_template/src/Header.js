import React, {Component} from 'react'


export default class Header extends Component {
  constructor() {
    super()
    this.state = {search: ''}
  }
  searchChange(v) {
    this.setState({search: v})
    this.props.onSearch(v)
  }
  render() {
    return (
      <div className="header-container">
        <div className="header-brand">
          {window.apiData.projectName}
          <input value={this.state.search} onChange={(e) => this.searchChange(e.target.value)} type="text" className="search" placeholder="search"/>
        </div>
        <div className="header-tools">
        </div>
      </div>
    )
  }
}

