package cmd

import (
	"fmt"
	"net/http"

	"gopkg.in/alecthomas/kingpin.v2"

	dashboardapi "github.com/previousnext/m8s/api"
)

type cmdDashboardAPI struct {
	Port       int32
	Namespace  string
	KubeMaster string
	KubeConfig string
	Mock       bool
}

func (cmd *cmdDashboardAPI) run(c *kingpin.ParseContext) error {
	var (
		mux = http.NewServeMux()
		srv = dashboardapi.New(cmd.KubeMaster, cmd.KubeConfig, cmd.Namespace, cmd.Mock)
	)

	mux.HandleFunc("/api/v1/list", srv.List)
	mux.HandleFunc("/api/v1/exec", srv.Exec)
	mux.HandleFunc("/api/v1/logs", srv.Logs)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", cmd.Port),
		Handler: mux,
	}

	return server.ListenAndServe()
}

// API declares the "api" sub command.
func API(app *kingpin.Application) {
	c := new(cmdDashboardAPI)

	cmd := app.Command("api", "API server for the UI component").Action(c.run)
	cmd.Flag("port", "Port to access requests.").Default("8080").Envar("M8S_UI_PORT").Int32Var(&c.Port)
	cmd.Flag("namespace", "Namespace to query resources.").Default("m8s").Envar("M8S_UI_NAMESPACE").StringVar(&c.Namespace)
	cmd.Flag("kube-master", "Address of the Kubernetes master.").Envar("M8S_UI_KUBE_MASTER").StringVar(&c.KubeMaster)
	cmd.Flag("kube-config", "Path to the Kubernetes config file.").Envar("M8S_UI_KUBE_CONFIG").StringVar(&c.KubeConfig)
	cmd.Flag("mock", "Run this API server with the mock data backend.").Envar("M8S_UI_MOCK").BoolVar(&c.Mock)
}
