import React from 'react';

class Header extends React.Component {

    render() {
        return (
            <header className="header">
                <a className="logo" href="/" alt="m8s User Interface">
                    <img src="/dist/logo.svg"></img>
                </a>
                <nav className="nav">
                    <ul className="nav__level-1">
                        <li className="nav__item"><a href="https://github.com/previousnext/m8s">Github</a></li>
                        <li className="nav__item"><a href="https://github.com/previousnext/m8s/tree/master/docs">Docs</a></li>
                    </ul>
                </nav>
            </header>
        );
    }

}

export default Header;