package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func main() {
	//This section will start the HTTP server and expose
	//any metrics on the /metrics endpoint.
	eventClient, err := NewEventClient()
	if err != nil {
		log.Fatalf("Failed to create the event client: %s", err)
	}

	if err := prometheus.Register(NewExporter(eventClient)); err != nil {
		log.Fatalf("Failed to register exporter: %s", err)
	}

	http.Handle("/metrics", promhttp.Handler())
	log.Info("Beginning to serve on port :9111")
	log.Fatal(http.ListenAndServe(":9111", nil))
}
