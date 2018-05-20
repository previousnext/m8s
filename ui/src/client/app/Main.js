import React from 'react';
import { Switch, Route } from 'react-router-dom'

import Home from './Home';
import Logs from './Logs';
import Shell from './Shell';

class Main extends React.Component {

    render() {
        return (
            <main className="page__middle">
                <Switch>
                    <Route exact path='/' component={Home}/>
                    <Route exact path='/e/:name/:container/logs' component={Logs}/>
                    <Route exact path='/e/:name/:container/shell' component={Shell}/>
                </Switch>
            </main>
        );
    }

}

export default Main;