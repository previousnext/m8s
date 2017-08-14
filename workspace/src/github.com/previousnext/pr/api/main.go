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
	cliCert      = kingpin.Flag("cert", "Certificate for TLS connection").Default("cert.pem").OverrideDefaultFromEnvar("TLS_CERT").String()
	cliKey       = kingpin.Flag("key", "Private key for TLS connection").Default("key.pem").OverrideDefaultFromEnvar("TLS_KEY").String()
	cliToken     = kingpin.Flag("token", "Token to authenticate against the API.").Default("").OverrideDefaultFromEnvar("AUTH_TOKEN").String()
	cliNamespace = kingpin.Flag("namespace", "Namespace to build environments.").Default("default").OverrideDefaultFromEnvar("NAMESPACE").String()

	// Black Death.
	cliBlackDeathImage   = kingpin.Flag("black-death-image", "Black Death image to deploy").Default("previousnext/k8s-black-death").OverrideDefaultFromEnvar("BLACK_DEATH_IMAGE").String()
	cliBlackDeathVersion = kingpin.Flag("black-death-version", "Version of Black Death to deploy").Default("0.0.1").OverrideDefaultFromEnvar("BLACK_DEATH_VERSION").String()

	// Traefik.
	cliTraefikImage   = kingpin.Flag("traefik-image", "Traefik image to deploy").Default("traefik").OverrideDefaultFromEnvar("TRAEFIK_IMAGE").String()
	cliTraefikVersion = kingpin.Flag("traefik-version", "Version of Traefik to deploy").Default("1.3").OverrideDefaultFromEnvar("TRAEFIK_VERSION").String()
	cliTraefikPort    = kingpin.Flag("traefik-port", "Assign this port to each node on the cluster for Traefik ingress").Default("80").OverrideDefaultFromEnvar("TRAEFIK_PORT").Int32()

	// SSH Server.
	cliSSHImage   = kingpin.Flag("ssh-image", "SSH server image to deploy").Default("previousnext/k8s-ssh-server").OverrideDefaultFromEnvar("SSH_IMAGE").String()
	cliSSHVersion = kingpin.Flag("ssh-version", "Version of SSH server to deploy").Default("0.0.5").OverrideDefaultFromEnvar("SSH_VERSION").String()
	cliSSHPort    = kingpin.Flag("ssh-port", "Assign this port to each node on the cluster for SSH ingress").Default("2222").OverrideDefaultFromEnvar("SSH_PORT").Int32()
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

	err = addons.CreateTraefik(client, *cliNamespace, *cliTraefikImage, *cliTraefikVersion, *cliTraefikPort)
	if err != nil {
		panic(err)
	}

	fmt.Println("Installing addon: ssh-server")

	err = addons.CreateSSHServer(client, *cliNamespace, *cliSSHImage, *cliSSHVersion, *cliSSHPort)
	if err != nil {
		panic(err)
	}

	fmt.Println("Installing addon: black-death")

	err = addons.CreateBlackDeath(client, *cliNamespace, *cliBlackDeathImage, *cliBlackDeathVersion)
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
