package k8s

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/previousnext/m8s/api/types"
	apiutils "github.com/previousnext/m8s/api/utils"
	"github.com/previousnext/m8s/k8sclient"
	"k8s.io/client-go/pkg/api/v1"
)

// Logs returns a stream of logs from a container.
func (s Server) Logs(w http.ResponseWriter, r *http.Request) {
	pod, err := apiutils.Param(r, "pod")
	if err != nil {
		apiutils.Fatal(w, err)
		return
	}

	container, err := apiutils.Param(r, "container")
	if err != nil {
		apiutils.Fatal(w, err)
		return
	}

	client, _, err := k8sclient.New(s.Master, s.Config)
	if err != nil {
		apiutils.Fatal(w, err)
		return
	}

	opts := &v1.PodLogOptions{
		Container: container,
	}

	resp, err := client.CoreV1().Pods(s.Namespace).GetLogs(pod, opts).Stream()
	if err != nil {
		apiutils.Fatal(w, err)
		return
	}
	defer resp.Close()

	// Reading all the logs output into a buffer.
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp)

	// Structuring the data so it can be output as JSON.
	output := types.LogsOutput{
		Logs: buf.String(),
	}

	// Write the logs to the webservers output.
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		apiutils.Fatal(w, err)
	}
}
