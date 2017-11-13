package cmd

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/pkg/errors"
	pb "github.com/previousnext/m8s/pb"
	"github.com/previousnext/m8s/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	promlog "github.com/prometheus/common/log"
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

	CacheSize string
	CacheType string

	LetsEncryptEmail  string
	LetsEncryptDomain string
	LetsEncryptCache  string

	PrometheusPort   string
	PrometheusPath   string
	PrometheusApache int32
}

func (cmd *cmdServer) run(c *kingpin.ParseContext) error {
	promlog.Info("Starting Prometheus Endpoint")

	go metrics(cmd.PrometheusPort, cmd.PrometheusPath)

	promlog.Info("Starting Listener")

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", cmd.Port))
	if err != nil {
		return errors.Wrap(err, "failed to start listener")
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		return errors.Wrap(err, "failed to get cluster config")
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "failed to get kubernetes client")
	}

	promlog.Info("Configuring Server")

	// Create a new server which adheres to the GRPC interface.
	srv, err := server.New(client, config, cmd.Token, cmd.Namespace, cmd.CacheType, cmd.CacheSize, cmd.PrometheusApache)
	if err != nil {
		return errors.Wrap(err, "failed to start server")
	}

	promlog.Info("Configuring TLS")

	var creds credentials.TransportCredentials

	// Attempt to load user provided certificates.
	// If no certificates are provided, fallback to Lets Encrypt.
	if cmd.TLSCert != "" && cmd.TLSKey != "" {
		promlog.Info("Loading TLS certificates from the filesystem")

		creds, err = credentials.NewServerTLSFromFile(cmd.TLSCert, cmd.TLSKey)
		if err != nil {
			return errors.Wrap(err, "failed to load tls from the filesystem")
		}
	} else {
		promlog.Info("Generating TLS certificates from LetsEncrypt")

		creds, err = getLetsEncrypt(cmd.LetsEncryptDomain, cmd.LetsEncryptEmail, cmd.LetsEncryptCache)
		if err != nil {
			return errors.Wrap(err, "failed to load tls from lets encrypt")
		}
	}

	promlog.Info("Booting Server")

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

	cmd.Flag("token", "Token to authenticate against the API.").Default("").OverrideDefaultFromEnvar("M8S_TOKEN").StringVar(&c.Token)
	cmd.Flag("namespace", "Namespace to build environments.").Default("default").OverrideDefaultFromEnvar("M8S_NAMESPACE").StringVar(&c.Namespace)

	cmd.Flag("cache-size", "Size of the filesystem for persistent cache storage").Default("100Gi").OverrideDefaultFromEnvar("M8S_CACHE_SIZE").StringVar(&c.CacheSize)
	cmd.Flag("cache-type", "StorageClass which you wish to use to provision the cache storage").Default("standard").OverrideDefaultFromEnvar("M8S_CACHE_TYPE").StringVar(&c.CacheType)

	// Lets Encrypt.
	cmd.Flag("lets-encrypt-email", "Email address to register with Lets Encrypt certificate").Default("admin@previousnext.com.au").OverrideDefaultFromEnvar("M8S_LETS_ENCRYPT_EMAIL").StringVar(&c.LetsEncryptEmail)
	cmd.Flag("lets-encrypt-domain", "Domain to use for Lets Encrypt certificate").Default("").OverrideDefaultFromEnvar("M8S_LETS_ENCRYPT_DOMAIN").StringVar(&c.LetsEncryptDomain)
	cmd.Flag("lets-encrypt-cache", "Cache directory to use for Lets Encrypt").Default("/tmp").OverrideDefaultFromEnvar("M8S_LETS_ENCRYPT_CACHE").StringVar(&c.LetsEncryptCache)

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
