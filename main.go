package main

import (
	"crypto/tls"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

type SdServer struct {
	Config *SdsConfig
}

var (
	app        = kingpin.New("SdServer", "Sales Discovery Service")
	configPath = app.Flag("config", "Configuration file for SD App").String()
)

// Initialize new SD object with passed local
func newSdInstance(config *SdsConfig) (*SdServer, error) {
	return &SdServer{
		Config: config,
	}, nil
}

func main() {
	logrus.Info("Starting up...")
	_, err := app.Parse(os.Args[1:])
	if err != nil {
		logrus.Fatal(err)
	}

	config, err := LoadConfig(*configPath)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load local [" + *configPath + "]")
	}

	sdi, err := newSdInstance(config)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize new SD instance")
	}

	// set up tls-enabled routing
	router := mux.NewRouter().StrictSlash(true)
	tlsConfig := &tls.Config{}
	tlsConfig.MinVersion = tls.VersionTLS12
	tlsConfig.InsecureSkipVerify = true

	server := &http.Server{Addr: config.ServiceConfig.ListenAddress, Handler: router, TLSConfig: tlsConfig}

	defer server.Close()

	router.HandleFunc("/salesData", sdi.handleGetData).Methods("GET")

	// todo: build out docker/kube/helm and allocate cert/key files for service
	//err = server.ListenAndServeTLS(local.ServiceConfig.CertFile, local.ServiceConfig.KeyFile)
	err = server.ListenAndServe()
	if err != nil {
		logrus.Fatal(err)
	}
}
