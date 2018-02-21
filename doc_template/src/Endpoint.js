import React, { Component } from 'react'
import Example from './Example'
// import htmlEncode from 'htmlencode'


export default class Endpoint extends Component {
  constructor() {
    super()
    this.state = {
    }
  }
  renderQueryParams() {
    if (this.props.endpoint.method !== 'GET' || this.props.endpoint.params.length === 0) {
      return null;
    }
    return (
      <div className="params">
        <h3>Query Parameters</h3>
        {this.props.endpoint.params.map(param =>
          <div className="param" key={param.name}>
            <span className="name">{param.name}</span>
            <span className="description">{param.description}</span>
          </div>
        )}
      </div>
    )
  }
  renderHeaders() {
    if (!this.props.endpoint.headers) {
      return null;
    }
    if (this.props.endpoint.headers.length === 0) {
      return null;
    }
    return (
      <div className="headers">
        <h3>Request Headers</h3>
        {this.props.endpoint.headers.map(header => 
          <div>
            <div className="header-name">{header.key}</div>
            <div className="header-description">{header.description}</div>
            <div className="header-example">Example: <em>{header.example}</em></div>
          </div>
        )}
      </div>
    )
  }
  render() {
    return (
      <div className="endpoint">
        <div className="info">
          <a name={this.props.endpoint.path}>
            <h2 className={this.props.endpoint.method}>
              {this.props.endpoint.name}
              <span className="method-label">{this.props.endpoint.method}</span>
            </h2>
          </a>
          <div className="url">
            {window.apiData.baseUrl &&  
              <span className="baseurl">
                { window.apiData.baseUrl }/
              </span>
            }
            {this.props.endpoint.url}
          </div>
          <div className="description"
            dangerouslySetInnerHTML={{ __html: this.props.endpoint.documentation}} />
            {this.renderHeaders()}
            {this.renderQueryParams()}
        </div>
        <Example endpoint={this.props.endpoint} />
      </div>
    )
  }
}
