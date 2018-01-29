import React from 'react';

class Options extends React.Component {

    change(event) {
        document.location.href = event.target.value;
    }

    render() {
        var items = this.props.containers.map(function(container) {
            return (
                <option key={container} value={"/e/"+this.props.name+"/"+container+"/"+this.props.operation}>{container}</option>
            )
        }, this);

        return (
            <select key={this.props.name} onChange={this.change}>
                <option value="none">---</option>
                {items}
            </select>
        );
    }

}

export default Options;