package mock

import (
	"encoding/json"
	"net/http"

	"github.com/previousnext/m8s/api/types"
	"github.com/previousnext/m8s/api/utils"
)

// Logs returns a stream of logs from a container.
func (s Server) Logs(w http.ResponseWriter, r *http.Request) {
	resp := types.LogsOutput{
		Logs: "This is a mock log entry!",
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		utils.Fatal(w, err)
	}
}
