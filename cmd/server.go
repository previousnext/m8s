package cmd

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/pkg/errors"
	"github.com/previousnext/m8s/k8sclient"
	pb "github.com/previousnext/m8s/pb"
	"github.com/previousnext/m8s/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	promlog "github.com/prometheus/common/log"
	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/alecthomas/kingpin.v2"
)

type cmdServer struct {
	Port    int32
	TLSCert string
	TLSKey  string

	Token     string
	Namespace string

	CacheDirectories string
	CacheSize        string
	CacheType        string

	LetsEncryptEmail  string
	LetsEncryptDomain string
	LetsEncryptCache  string

	PrometheusPort string
	PrometheusPath string

	KubeMaster string
	KubeConfig string

	DockerCfg string
}

func (cmd *cmdServer) run(c *kingpin.ParseContext) error {
	promlog.Info("Starting Prometheus Endpoint")

	go metrics(cmd.PrometheusPort, cmd.PrometheusPath)

	promlog.Info("Starting Listener")

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", cmd.Port))
	if err != nil {
		return errors.Wrap(err, "failed to start listener")
	}

	client, config, err := k8sclient.New(cmd.KubeMaster, cmd.KubeConfig)
	if err != nil {
		return errors.Wrap(err, "failed to get Kubernetes client")
	}

	promlog.Info("Configuring Server")

	// Create a new server which adheres to the GRPC interface.
	srv, err := server.New(server.Input{
		Client:    client,
		Config:    config,
		Token:     cmd.Token,
		Namespace: cmd.Namespace,
		DockerCfg: cmd.DockerCfg,
		Cache: server.InputCache{
			Directories: cmd.CacheDirectories,
			Type:        cmd.CacheType,
			Size:        cmd.CacheSize,
		},
	})
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

	cmd.Flag("cache-dirs", "Directories which will be cached between builds").Default("composer:/root/.composer,yarn:/usr/local/share/.cache/yarn").OverrideDefaultFromEnvar("M8S_CACHE_DIRS").StringVar(&c.CacheDirectories)
	cmd.Flag("cache-size", "Size of the filesystem for persistent cache storage").Default("100Gi").OverrideDefaultFromEnvar("M8S_CACHE_SIZE").StringVar(&c.CacheSize)
	cmd.Flag("cache-type", "StorageClass which you wish to use to provision the cache storage").Default("standard").OverrideDefaultFromEnvar("M8S_CACHE_TYPE").StringVar(&c.CacheType)

	// Lets Encrypt.
	cmd.Flag("lets-encrypt-email", "Email address to register with Lets Encrypt certificate").Default("admin@previousnext.com.au").OverrideDefaultFromEnvar("M8S_LETS_ENCRYPT_EMAIL").StringVar(&c.LetsEncryptEmail)
	cmd.Flag("lets-encrypt-domain", "Domain to use for Lets Encrypt certificate").Default("").OverrideDefaultFromEnvar("M8S_LETS_ENCRYPT_DOMAIN").StringVar(&c.LetsEncryptDomain)
	cmd.Flag("lets-encrypt-cache", "Cache directory to use for Lets Encrypt").Default("/tmp").OverrideDefaultFromEnvar("M8S_LETS_ENCRYPT_CACHE").StringVar(&c.LetsEncryptCache)

	// Promtheus.
	cmd.Flag("prometheus-port", "Prometheus metrics port").Default(":9000").OverrideDefaultFromEnvar("M8S_METRICS_PORT").StringVar(&c.PrometheusPort)
	cmd.Flag("prometheus-path", "Prometheus metrics path").Default("/metrics").OverrideDefaultFromEnvar("M8S_METRICS_PATH").StringVar(&c.PrometheusPath)

	// Kubernetes.
	cmd.Flag("kube-master", "Address of the Kubernetes master.").Envar("M8S_UI_KUBE_MASTER").StringVar(&c.KubeMaster)
	cmd.Flag("kube-config", "Path to the Kubernetes config file.").Envar("M8S_UI_KUBE_CONFIG").StringVar(&c.KubeConfig)

	// Docker Registry.
	cmd.Flag("dockercfg", "https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry").Default("").Envar("M8S_DOCKERCFG").StringVar(&c.DockerCfg)
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
