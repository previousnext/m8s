import React from 'react';

import EnvironmentTable from './EnvironmentTable'
import queryString from 'query-string'

export default class Home extends React.Component {

    constructor(props) {
        super(props);
        const params = queryString.parse(props.location.search);
        this.state = {
          envs : [],
          search: params.search ? params.search : '',
        };
        this.updateInput = this.updateInput.bind(this)
    }

    componentDidMount() {
        fetch(`/api/v1/list`, {
            credentials: "same-origin"
        }).then(result=>result.json())
          .then(envs=>this.setState({envs}))
    }

    componentWillUnmount() {
        this.serverRequest.abort();
    }

    updateInput(e) {
      this.props.history.push(`${this.props.match.url}?search=${e.target.value}`);
      this.setState({
        search: e.target.value,
      })
    }

    render() {
        return (
            <div className="table__stretch-wrapper">
                <div className="form__item">
                  <label htmlFor="search">Search</label>
                  <input
                    type="text"
                    id="search"
                    onChange={this.updateInput}
                    value={this.state.search}
                  />
                </div>
                <EnvironmentTable envs={this.state.envs} search={this.state.search}/>
            </div>
        );
    }

}
