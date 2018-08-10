import React from 'react';
import Breadcrumbs from './Breadcrumbs';

export default class Shell extends React.Component {

    ws(name, container) {
        let term, websocket;

        let proto = "ws"

        let loc = window.location;
        if (loc.protocol === "https:") {
            proto = "wss";
        }

        websocket = new WebSocket(`${proto}://${window.location.hostname}/api/v1/exec?pod=${name}&container=${container}`);

        websocket.onopen = function() {

            term = new Terminal({
                screenKeys: true,
                useStyle: true,
                cursorBlink: true,
            });

            term.on('data', function(data) {
                websocket.send(data);
            });

            term.on('title', function(title) {
                document.title = title;
            });

            term.open(document.getElementById('shell'));

            websocket.onmessage = function(evt) {
                term.write(evt.data);
            }

            websocket.onclose = function(evt) {
                term.destroy();

                // This means that when the user disconnects they are redirected back
                // to the main home page where they can schedule another shell session.
                document.location.href = "/";
            }

            websocket.onerror = function(evt) {
                if (typeof console.log !== "undefined") {
                    console.log(evt)
                }
            }
        }
    }

    render() {
        return (
            <main>
                <h1>Console</h1>
                <Breadcrumbs operation="console" name={this.props.match.params.name} container={this.props.match.params.container} />
                <pre id="shell" className="console" onLoad={this.ws(this.props.match.params.name, this.props.match.params.container)}></pre>
            </main>
        );
    }

}
