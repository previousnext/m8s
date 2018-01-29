import React from 'react';

import Options from './Options.jsx';
import UIs from './UIs.jsx';

class Home extends React.Component {

    constructor(props) {
        super(props);
        this.state = {envs : []};
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

    render() {
        var tbody = this.state.envs.map(function(env) {
            return (
                <tr key={env.name}>
                    <td className="table--enlarged"><a href={"//"+env.domain}>{env.name}</a></td>
                    <td>
                        <Options operation="logs" name={env.name} containers={env.containers} />
                    </td>
                    <td>
                        <Options operation="shell" name={env.name} containers={env.containers} />
                    </td>
                    <td>
                        { /* @todo, Consolidate with the Options component. */ }
                        <UIs name={env.name} base_url={"//"+env.domain} />
                    </td>
                </tr>
            )
        });

        return (
            <div className="table__stretch-wrapper">
                <table className="js-table--responsive">
                    <thead>
                    <tr>
                        <th>Domain</th>
                        <th>Logs</th>
                        <th>Console</th>
                        <th>UIs</th>
                    </tr>
                    </thead>
                    <tbody>
                        {tbody}
                    </tbody>
                </table>
            </div>
        );
    }

}

export default Home;