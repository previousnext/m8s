package mock

import (
	"io"
	"net/http"
	"os/exec"

	"github.com/kr/pty"
	"golang.org/x/net/websocket"
)

// Exec bash (shell) inside a container.
func (s Server) Exec(w http.ResponseWriter, r *http.Request) {
	wws := websocket.Handler(func(ws *websocket.Conn) {
		cmd := exec.Command("/bin/sh")

		tty, err := pty.Start(cmd)
		if err != nil {
			panic(err)
		}
		defer tty.Close()

		go func() {
			io.Copy(ws, tty)
		}()

		go func() {
			io.Copy(ws, tty)
		}()

		go func() {
			io.Copy(tty, ws)
		}()

		cmd.Wait()
	})

	wws.ServeHTTP(w, r)
}
