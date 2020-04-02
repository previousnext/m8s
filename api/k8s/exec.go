package k8s

import (
	"log"
	"net/http"

	apiutils "github.com/previousnext/m8s/api/utils"
	"github.com/previousnext/m8s/k8sclient"
	"golang.org/x/net/websocket"

	"github.com/previousnext/m8s/internal/podutils"
)

// Exec bash (shell) inside a container.
func (s Server) Exec(w http.ResponseWriter, r *http.Request) {
	pod, err := apiutils.Param(r, "pod")
	if err != nil {
		apiutils.Fatal(w, err)
		return
	}

	log.Println("Received request for pod:", pod)

	container, err := apiutils.Param(r, "container")
	if err != nil {
		apiutils.Fatal(w, err)
		return
	}

	log.Println("Received request for pod:", container)

	client, config, err := k8sclient.New(s.Master, s.Config)
	if err != nil {
		apiutils.Fatal(w, err)
		return
	}

	log.Println("Starting session")

	wws := websocket.Handler(func(ws *websocket.Conn) {
		opts := podutils.RunParams{
			Client:    client,
			Config:    config,
			Namespace: s.Namespace,
			Pod:       pod,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
			Writer:    ws,
			Reader:    ws,
			Container: container,
			Command: []string{
				"/bin/bash",
			},
		}

		err = podutils.Run(opts)
		if err != nil {
			apiutils.Fatal(w, err)
			return
		}
	})

	log.Println("Session finished")

	wws.ServeHTTP(w, r)
}
