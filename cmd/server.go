package cmd

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"

	pb "github.com/previousnext/m8s/pb"
	"github.com/previousnext/m8s/server"
	"github.com/previousnext/m8s/server/k8s/addons/ssh-server"
	"github.com/previousnext/m8s/server/k8s/addons/traefik"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/alecthomas/kingpin.v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type cmdServer struct {
	Port    int32
	TLSCert string
	TLSKey  string

	Token     string
	Namespace string

	FilesystemSize string

	LetsEncryptEmail  string
	LetsEncryptDomain string
	LetsEncryptCache  string

	TraefikImage   string
	TraefikVersion string
	TraefikPort    int32

	SSHImage   string
	SSHVersion string

	DockerCfgRegistry string
	DockerCfgUsername string
	DockerCfgPassword string
	DockerCfgEmail    string
	DockerCfgAuth     string

	PrometheusPort   string
	PrometheusPath   string
	PrometheusApache int32
}

func (cmd *cmdServer) run(c *kingpin.ParseContext) error {
	log.Println("Starting Prometheus Endpoint")

	go metrics(cmd.PrometheusPort, cmd.PrometheusPath)

	log.Println("Starting Server")

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", cmd.Port))
	if err != nil {
		panic(err)
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	log.Println("Installing addon: traefik")

	err = traefik.Create(client, cmd.Namespace, cmd.TraefikImage, cmd.TraefikVersion, cmd.TraefikPort)
	if err != nil {
		panic(err)
	}

	log.Println("Installing addon: ssh-server")

	err = ssh_server.Create(client, cmd.Namespace, cmd.SSHImage, cmd.SSHVersion)
	if err != nil {
		panic(err)
	}

	log.Println("Booting API")

	// Create a new server which adheres to the GRPC interface.
	srv, err := server.New(client, config, cmd.Token, cmd.Namespace, cmd.FilesystemSize, cmd.PrometheusApache, cmd.DockerCfgRegistry, cmd.DockerCfgUsername, cmd.DockerCfgPassword, cmd.DockerCfgEmail, cmd.DockerCfgAuth)
	if err != nil {
		panic(err)
	}

	var creds credentials.TransportCredentials

	// Attempt to load user provided certificates.
	// If no certificates are provided, fallback to Lets Encrypt.
	if cmd.TLSCert != "" && cmd.TLSKey != "" {
		creds, err = credentials.NewServerTLSFromFile(cmd.TLSCert, cmd.TLSKey)
		if err != nil {
			panic(err)
		}
	} else {
		creds, err = getLetsEncrypt(cmd.LetsEncryptDomain, cmd.LetsEncryptEmail, cmd.LetsEncryptCache)
		if err != nil {
			panic(err)
		}
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterM8SServer(grpcServer, srv)
	return grpcServer.Serve(listen)
}

// Server declares the "server" sub command.
func Server(app *kingpin.Application) {
	c := new(cmdServer)

	cmd := app.Command("server", "Run the M8s server").Action(c.run)
	cmd.Flag("port", "Port to run this service on").Default("443").OverrideDefaultFromEnvar("M8S_PORT").Int32Var(&c.Port)
	cmd.Flag("cert", "Certificate for TLS connection").Default("").OverrideDefaultFromEnvar("M8S_TLS_CERT").StringVar(&c.TLSCert)
	cmd.Flag("key", "Private key for TLS connection").Default("").OverrideDefaultFromEnvar("M8S_TLS_KEY").StringVar(&c.TLSKey)

	cmd.Flag("token", "Token to authenticate against the API.").Default("").OverrideDefaultFromEnvar("M8S_AUTH_TOKEN").StringVar(&c.Token)
	cmd.Flag("namespace", "Namespace to build environments.").Default("default").OverrideDefaultFromEnvar("M8S_NAMESPACE").StringVar(&c.Namespace)

	cmd.Flag("fs-size", "Size of the filesystem for persistent storage").Default("100Gi").OverrideDefaultFromEnvar("M8S_FS_SIZE").StringVar(&c.FilesystemSize)

	// Lets Encrypt.
	cmd.Flag("lets-encrypt-email", "Email address to register with Lets Encrypt certificate").Default("admin@previousnext.com.au").OverrideDefaultFromEnvar("M8S_LETS_ENCRYPT_EMAIL").StringVar(&c.LetsEncryptEmail)
	cmd.Flag("lets-encrypt-domain", "Domain to use for Lets Encrypt certificate").Default("").OverrideDefaultFromEnvar("M8S_LETS_ENCRYPT_DOMAIN").StringVar(&c.LetsEncryptDomain)
	cmd.Flag("lets-encrypt-cache", "Cache directory to use for Lets Encrypt").Default("/tmp").OverrideDefaultFromEnvar("M8S_LETS_ENCRYPT_CACHE").StringVar(&c.LetsEncryptCache)

	// Traefik.
	cmd.Flag("traefik-image", "Traefik image to deploy").Default("traefik").OverrideDefaultFromEnvar("M8S_TRAEFIK_IMAGE").StringVar(&c.TraefikImage)
	cmd.Flag("traefik-version", "Version of Traefik to deploy").Default("1.3").OverrideDefaultFromEnvar("M8S_TRAEFIK_VERSION").StringVar(&c.TraefikVersion)
	cmd.Flag("traefik-port", "Assign this port to each node on the cluster for Traefik ingress").Default("80").OverrideDefaultFromEnvar("M8S_TRAEFIK_PORT").Int32Var(&c.TraefikPort)

	// SSH Server.
	cmd.Flag("ssh-image", "SSH server image to deploy").Default("previousnext/k8s-ssh-server").OverrideDefaultFromEnvar("M8S_SSH_IMAGE").StringVar(&c.SSHImage)
	cmd.Flag("ssh-version", "Version of SSH server to deploy").Default("2.1.0").OverrideDefaultFromEnvar("M8S_SSH_VERSION").StringVar(&c.SSHVersion)

	// DockerCfg.
	cmd.Flag("dockercfg-registry", "Registry for Docker Hub credentials").Default("").OverrideDefaultFromEnvar("M8S_DOCKERCFG_REGISTRY").StringVar(&c.DockerCfgRegistry)
	cmd.Flag("dockercfg-username", "Username for Docker Hub credentials").Default("").OverrideDefaultFromEnvar("M8S_DOCKERCFG_USERNAME").StringVar(&c.DockerCfgUsername)
	cmd.Flag("dockercfg-password", "Password for Docker Hub credentials").Default("").OverrideDefaultFromEnvar("M8S_DOCKERCFG_PASSWORD").StringVar(&c.DockerCfgPassword)
	cmd.Flag("dockercfg-email", "Email for Docker Hub credentials").Default("").OverrideDefaultFromEnvar("M8S_DOCKERCFG_EMAIL").StringVar(&c.DockerCfgEmail)
	cmd.Flag("dockercfg-auth", "Auth token for Docker Hub credentials").Default("").OverrideDefaultFromEnvar("M8S_DOCKERCFG_AUTH").StringVar(&c.DockerCfgAuth)

	// Promtheus.
	cmd.Flag("prometheus-port", "Prometheus metrics port").Default(":9000").OverrideDefaultFromEnvar("M8S_METRICS_PORT").StringVar(&c.PrometheusPort)
	cmd.Flag("prometheus-path", "Prometheus metrics path").Default("/metrics").OverrideDefaultFromEnvar("M8S_METRICS_PATH").StringVar(&c.PrometheusPath)
	cmd.Flag("prometheus-apache-exporter", "Prometheus metrics port for Apache on built environments").Default("9117").OverrideDefaultFromEnvar("M8S_METRICS_APACHE_PORT").Int32Var(&c.PrometheusApache)
}

// Helper function for serving Prometheus metrics.
func metrics(port, path string) {
	http.Handle(path, promhttp.Handler())
	log.Fatal(http.ListenAndServe(port, nil))
}

// Helper function for adding Lets Encrypt certificates.
func getLetsEncrypt(domain, email, cache string) (credentials.TransportCredentials, error) {
	manager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache(cache),
		HostPolicy: autocert.HostWhitelist(domain),
		Email:      email,
	}
	return credentials.NewTLS(&tls.Config{GetCertificate: manager.GetCertificate}), nil
}
