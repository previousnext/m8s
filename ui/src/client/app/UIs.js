import React from 'react';

export default class UIs extends React.Component {

    change(event) {
        document.location.href = event.target.value;
    }

    render() {
        return (
            <select key={this.props.name} onChange={this.change}>
                <option value="none">---</option>
                <option key="solr" value={`${this.props.base_url}/solr`}>Solr</option>
                <option key="mailhog" value={`${this.props.base_url}/mailhog`}>MailHog</option>
            </select>
        );
    }

}
