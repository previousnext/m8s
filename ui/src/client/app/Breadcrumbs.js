import React from 'react';

export default class Breadcrumbs extends React.Component {


    render() {
        return (
            <nav className="breadcrumbs" role="navigation" aria-label="You are here">
                <ol className="breadcrumbs__menu">
                    <li class="breadcrumbs__item"><a href="/">Home</a> /</li>
                    <li class="breadcrumbs__item"><a href={`/?search=${this.props.name}`}>{this.props.name}</a> /</li>
                    <li class="breadcrumbs__item">{this.props.container} /</li>
                    <li class="breadcrumbs__item">{this.props.operation}</li>
                </ol>
            </nav>
        );
    }

}
