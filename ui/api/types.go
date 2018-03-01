package api

import "net/http"

// API interface used by the UI.
type API interface {
	List(http.ResponseWriter, *http.Request)
	Logs(http.ResponseWriter, *http.Request)
	Exec(http.ResponseWriter, *http.Request)
}
