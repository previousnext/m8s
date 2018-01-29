package k8s

// Server for interacting with the K8s API Server.
type Server struct {
	Namespace string
	Master    string
	Config    string
}
