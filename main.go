package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func main() {
	portEnv := os.Getenv("KEE_PORT")
	if portEnv == "" {
		portEnv = "9111"
	}
	portString := fmt.Sprintf(":%s", portEnv)

	eventClient, err := NewEventClient()
	if err != nil {
		log.Fatalf("Failed to create the event client: %s", err)
	}

	eventClient.IncludeNormalEvents = os.Getenv("KEE_INCLUDE_NORMAL_EVENTS") == "true"

	if err := prometheus.Register(NewExporter(eventClient)); err != nil {
		log.Fatalf("Failed to register exporter: %s", err)
	}

	http.Handle("/metrics", promhttp.Handler())
	log.Infof("Beginning to serve on port %s", portString)
	log.Fatal(http.ListenAndServe(portString, nil))
}
