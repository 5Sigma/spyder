import React, { Component} from 'react'
import _ from 'lodash'

export default class NavMenu extends Component {
  getPaths() {
    var paths = {}
    window.apiData.endpoints.forEach(ep => {
      let vals = _.get(paths, ep.path.replace(/\//g,'.'))
      if (vals) {
        _.set(paths, ep.path.replace(/\//g,'.'), vals.append(ep))
      } else {
        _.set(paths, ep.path.replace(/\//g,'.'), [ep])
      }
    })
    return paths
  }
  buildLinks(def, path) {
    return (
      <div>
        {Object.keys(def).map(key => {
          if (_.isArrayLike(def[key])) {
            return (
              <div key={`container-${key}`}>
                {def[key].map(ep => 
                  <div className="nav-item" key={ep["path"]}>
                    <a href={`#${ep.path}`}>{ep.name.split('/')[ep.name.split('/').length - 1]}</a>
                  </div>
                )}
              </div>
            )
          } else if (_.isObject(def[key])) {
            return (
              <div className="section" key={`${path}/${key}`}>
                <div className="section-header">{key}</div>
                {this.buildLinks(def[key])}
              </div>
            )
          }
          return null;
        })}
      </div>
    )
  }
  navLinks() {
    return this.buildLinks(this.getPaths(), "/")
  }
  render() {
    return (
      <div className="nav-menu">
        {this.navLinks()}
      </div>
    )
  }
}
