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

    handleSubmit(e) {
        e.preventDefault();
    }

    render() {
        return (
            <main className="table__stretch-wrapper">
                <h1 class="visually-hidden">Dashboard home</h1>
                <form className="search" onSubmit={this.handleSubmit}>
                    <div className="form__item">
                        <label htmlFor="search" className="visually-hidden">Search</label>
                        <input
                            type="search"
                            id="search"
                            placeholder="Search environments"
                            onChange={this.updateInput}
                            value={this.state.search}
                        />
                    </div>
                    <input type="submit" className="button--search" />
                </form>
                <EnvironmentTable envs={this.state.envs} search={this.state.search}/>
            </main>
        );
    }

}
