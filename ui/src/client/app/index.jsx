import React from 'react';
import {render} from 'react-dom';
import { BrowserRouter } from 'react-router-dom'

import Header from './Header.jsx';
import Main from './Main.jsx';

class App extends React.Component {
    render () {
        return (
            <div className="page page__container">
                <Header />
                <Main />
            </div>
        );
    }
}

// https://medium.com/@pshrmn/a-simple-react-router-v4-tutorial-7f23ff27adf
render(<BrowserRouter><App /></BrowserRouter>, document.getElementById('app'));