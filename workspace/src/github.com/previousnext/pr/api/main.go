package main

import (
	"fmt"
	"net"

	"github.com/alecthomas/kingpin"
	"github.com/previousnext/pr/api/k8s/addons"
	pb "github.com/previousnext/pr/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"k8s.io/client-go/rest"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
)

var (
	cliPort      = kingpin.Flag("port", "Port to run this service on").Default("80").OverrideDefaultFromEnvar("PORT").Int32()
	cliCert      = kingpin.Flag("cert", "Certificate for TLS connection").OverrideDefaultFromEnvar("TLS_CERT").Default("cert.pem").String()
	cliKey       = kingpin.Flag("key", "Private key for TLS connection").OverrideDefaultFromEnvar("TLS_KEY").Default("key.pem").String()
	cliToken     = kingpin.Flag("token", "Token to authenticate against the API.").Default("").String()
	cliNamespace = kingpin.Flag("namespace", "Namespace to build environments.").Default("default").String()

	// Black Death.
	cliBlackDeathVersion = kingpin.Flag("black-death-version", "Version of Black Death to deploy").Default("0.0.1").String()

	// Traefik.
	cliTraefikVersion = kingpin.Flag("traefik-version", "Version of Traefik to deploy").Default("1.3").String()
	cliTraefikPort    = kingpin.Flag("traefik-port", "Assign this port to each node on the cluster for Traefik ingress").Default("80").Int32()

	// SSH Server.
	cliSSHVersion = kingpin.Flag("ssh-version", "Version of SSH server to deploy").Default("0.0.5").String()
	cliSSHPort    = kingpin.Flag("ssh-port", "Assign this port to each node on the cluster for SSH ingress").Default("2222").Int32()
)

func main() {
	kingpin.Parse()

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", *cliPort))
	if err != nil {
		panic(err)
	}

	// Attempt to connect to the K8s API Server.
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	client, err := client.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	fmt.Println("Installing addon: traefik")

	err = addons.CreateTraefik(client, *cliNamespace, *cliTraefikVersion, *cliTraefikPort)
	if err != nil {
		panic(err)
	}

	fmt.Println("Installing addon: ssh-server")

	err = addons.CreateSSHServer(client, *cliNamespace, *cliSSHVersion, *cliSSHPort)
	if err != nil {
		panic(err)
	}

	fmt.Println("Installing addon: black-death")

	err = addons.CreateBlackDeath(client, *cliNamespace, *cliBlackDeathVersion)
	if err != nil {
		panic(err)
	}

	fmt.Println("Booting API")

	// Create a new server which adheres to the GRPC interface.
	srv := server{
		client: client,
		config: config,
	}

	creds, err := credentials.NewServerTLSFromFile(*cliCert, *cliKey)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterPRServer(grpcServer, srv)
	grpcServer.Serve(listen)
}
