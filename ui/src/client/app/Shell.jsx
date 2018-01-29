import React from 'react';

class Shell extends React.Component {

    ws(name, container) {
        var term, websocket;

        var proto = "ws"

        var loc = window.location;
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
                if (typeof console.log == "function") {
                    console.log(evt)
                }
            }
        }
    }

    render() {
        return (
            <pre id="shell" className="console" onLoad={this.ws(this.props.match.params.name,this.props.match.params.container)}></pre>
        );
    }

}

export default Shell;