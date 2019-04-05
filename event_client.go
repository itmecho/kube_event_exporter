package main

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type EventClient struct {
	kube                *kubernetes.Clientset
	IncludeNormalEvents bool
	TimeWindowMinutes   time.Duration
}

func NewEventClient() (client EventClient, err error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return
	}

	client = EventClient{
		kube: kubeClient,

		// TODO make this configurable
		TimeWindowMinutes: 5,
	}

	return
}

func (e EventClient) Scrape(ch chan<- prometheus.Metric) error {

	opts := metav1.ListOptions{FieldSelector: "type==Warning"}

	if e.IncludeNormalEvents {
		opts = metav1.ListOptions{}
	}

	list, err := e.kube.CoreV1().Events("").List(opts)
	if err != nil {
		return err
	}

	timeWindow := &metav1.Time{
		Time: time.Now().Add(-e.TimeWindowMinutes * time.Minute),
	}

	for _, event := range list.Items {
		if event.LastTimestamp.Before(timeWindow) {
			continue
		}
		metric, err := prometheus.NewConstMetric(
			prometheus.NewDesc(
				"kubernetes_event_count",
				"Kubernetes events",
				[]string{"kubernetes_namespace", "event_type", "event_name", "event_message", "object_kind", "object_name", "event_reason"},
				nil,
			),
			prometheus.GaugeValue,
			float64(event.Count),
			event.InvolvedObject.Namespace,
			event.Type,
			event.Name,
			event.Message,
			event.InvolvedObject.Kind,
			event.InvolvedObject.Name,
			event.Reason,
		)
		if err != nil {
			return err
		}
		log.Println("adding metric to channel")
		ch <- metric
	}

	return nil
}
