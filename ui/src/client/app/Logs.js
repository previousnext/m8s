import React from 'react';

export default class Logs extends React.Component {

    constructor(props) {
        super(props);
        this.state = {logs : "Loading..."}
    }

    componentDidMount() {
        fetch(`/api/v1/logs?pod=${this.props.match.params.name}&container=${this.props.match.params.container}`, {
            credentials: "same-origin"
        }).then(response => response.json())
          .then(data => this.setState({ logs: data.logs }))
    }

    componentWillUnmount() {
        this.serverRequest.abort();
    }

    render() {
        return (
            <pre id="logs" className="console">{this.state.logs}</pre>
        );
    }

}
