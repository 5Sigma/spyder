import React, { Component } from 'react'
// import Highlight from 'react-highlight'
import SyntaxHighlighter from 'react-syntax-highlighter';
import 'highlight.js/styles/atelier-savanna-dark.css';
import { atelierSavannaDark } from 'react-syntax-highlighter/dist/styles';


export default class Example extends Component {
  constructor(props) {
    super()
    this.state = {
      show: 'Request'
    }
    console.log(props.endpoint.requestExample)
    this.requestExample = JSON.stringify(props.endpoint.request, null, 2)
    this.responseExample = JSON.stringify(props.endpoint.response, null, 2)
  }
  renderLink(key) {
    let cls = ""
    if (this.state.show === key) {
      cls = 'active'
    }
    return (
      <button className={cls} onClick={() => this.setState({show: key}) }>{key}</button>
    )
  }
  buildJsExample() {
    return `
let payload = ${this.requestExample};
fetch('${this.props.endpoint.url}', payload)
.then(r => r.json()).then(data => {
  // process
})
`;
  }
    buildCurlExample() {
      return `
      curl ${this.props.endpoint.url} -D '${JSON.stringify(JSON.parse(this.requestExample))}'
`;
  }
  renderExample() {
    switch (this.state.show) {
      case 'Response':
        return (
          <SyntaxHighlighter language='json' style={atelierSavannaDark}>
            {this.responseExample || 'None provided'} 
          </SyntaxHighlighter>
        )
      case 'Request':
        return (
          <SyntaxHighlighter language='json' style={atelierSavannaDark}>
            {this.requestExample || 'None provided'} 
          </SyntaxHighlighter>
        )
      case 'JS':
        return (
          <SyntaxHighlighter language='javascript' style={atelierSavannaDark}>
            {this.buildJsExample()} 
          </SyntaxHighlighter>
        )
      case 'cURL':
        return (
          <SyntaxHighlighter language='sh' style={atelierSavannaDark}>
            {this.buildCurlExample()} 
          </SyntaxHighlighter>
        )
      default:
        return null; 
    }
  }
  render() {
    return (
      <div className="example">
        <div className="example-menu">
          {this.renderLink('Request')}
          {this.renderLink('Response')}
          {this.renderLink('JS')}
        </div>
        {this.renderExample()}
      </div>
    )
  }
}
