package main

import (
	"context"
	"fmt"

	"net/http"
	"os"
	"os/signal"
	"syscall"

	server "github.com/vshiva/goreactapp/internal/http"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var serverCommand = cli.Command{
	Name:   "server",
	Usage:  "Start Console Server",
	Action: serverAction,
	Flags:  serverFlags,
}

var serverFlags = []cli.Flag{
	cli.IntFlag{
		Name:   "health-port",
		Value:  12273,
		EnvVar: "HEALTH_PORT",
	},
	cli.IntFlag{
		Name:   "metrics-port",
		Value:  12274,
		EnvVar: "METRICS_PORT",
	},
	cli.IntFlag{
		Name:   "port",
		Value:  4443,
		EnvVar: "PORT",
	},
	cli.StringFlag{
		Name:   "mode",
		EnvVar: "MODE",
		Usage:  "Application mode. Either Blue or Green.",
	},
	cli.StringFlag{
		Name:  "bind-address",
		Value: "0.0.0.0",
		Usage: "The IP address on which to listen for the -port port. If blank, all interfaces will be used (0.0.0.0). (default 0.0.0.0)",
	},
	cli.StringFlag{
		Name:   "tls-cert-file",
		Value:  "/etc/tls/goreactapp/cert.pem",
		Usage:  "File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated after server cert).",
		EnvVar: "CERT_FILE",
	},
	cli.StringFlag{
		Name:   "tls-private-key-file",
		Value:  "/etc/tls/goreactapp/private-key.pem",
		Usage:  "File containing the default x509 private key matching --tls-cert-file.",
		EnvVar: "PRIVATE_KEY_FILE",
	}}

var serverAction = func(c *cli.Context) error {
	log.Info("Starting app")

	log.Debug("Parsing server options")
	o, err := parseServerOptions(c)
	if err != nil {
		log.WithError(err).Error("Unable to validate arguments")
		return errorExitCode
	}

	errc := make(chan error, 3)

	// Shutdown on SIGINT, SIGTERM
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// Start health server
	healthService := NewHealthService()
	go func() {
		log.WithField("port", o.HealthPort).Info("Starting health service")
		err := healthService.ListenAndServe(fmt.Sprintf(":%d", o.HealthPort))
		errc <- errors.Wrap(err, "health service returned an error")
	}()

	// Start metrics server
	go func() {
		log.WithField("port", o.MetricsPort).Info("Starting metrics server")
		http.Handle("/metrics", prometheus.Handler())
		errc <- http.ListenAndServe(fmt.Sprintf(":%d", o.MetricsPort), nil)
	}()

	webServer, err := server.New(o.BindAddress, o.Port, c.String("mode"))
	if err != nil {
		log.WithError(err).Error("Unable to start fnops server")
		return errorExitCode
	}

	// Start  server
	go func() {
		log.WithField("port", o.Port).Info("Starting console server")
		useSSL := true
		if _, err := os.Stat(o.CertFile); os.IsNotExist(err) {
			log.WithField("certFile", o.CertFile).Info("cert file doesn't exist, not using SSL")
			useSSL = false
		}
		if _, err := os.Stat(o.KeyFile); os.IsNotExist(err) {
			log.WithField("keyFile", o.KeyFile).Info("Key file doesn't exist, not using SSL")
			useSSL = false
		}
		var err error
		if useSSL {
			err = webServer.ListenAndServeTLS(o.CertFile, o.KeyFile)
		} else {
			err = webServer.ListenAndServe()
		}
		errc <- errors.Wrap(err, "web server returned an error")
	}()

	err = <-errc
	log.WithError(err).Info("Shutting down")

	// Gracefully shutdown the health server
	healthService.Shutdown(context.Background())
	webServer.Shutdown(context.Background())

	return nil
}

type serverOptions struct {
	Port        int
	HealthPort  int
	MetricsPort int
	BindAddress string
	CertFile    string
	KeyFile     string
}

func parseServerOptions(c *cli.Context) (options *serverOptions, err error) {

	port := c.Int("port")
	healthPort := c.Int("health-port")

	if healthPort == port {
		return nil, errors.New("health-port and port cannot be the same")
	}

	metricsPort := c.Int("metrics-port")
	if metricsPort == port {
		return nil, errors.New("metrics-port and port cannot be the same")
	}

	if metricsPort == healthPort {
		return nil, errors.New("metrics-port and health-port cannot be the same")
	}

	options = &serverOptions{
		Port:        port,
		HealthPort:  healthPort,
		MetricsPort: metricsPort,
		BindAddress: c.String("bind-address"),
		CertFile:    c.String("tls-cert-file"),
		KeyFile:     c.String("tls-private-key-file"),
	}

	return
}
