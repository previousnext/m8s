package mock

import (
	"encoding/json"
	"net/http"

	"github.com/previousnext/m8s/ui/api/types"
	"github.com/previousnext/m8s/ui/api/utils"
)

// List returns a list of environments.
func (s Server) List(w http.ResponseWriter, r *http.Request) {
	resp := []types.Environment{
		{
			Name:   "project-1234-topic-1",
			Domain: "project-1234-topic-1.example.io",
			Containers: []string{
				"app",
				"mail",
				"search",
				"db",
			},
		},
		{
			Name:   "project-1234-topic-2",
			Domain: "project-1234-topic-2.example.io",
			Containers: []string{
				"app",
				"mail",
				"search",
				"db",
			},
		},
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		utils.Fatal(w, err)
	}
}
