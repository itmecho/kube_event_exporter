package main

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type EventClient struct {
	kube                *kubernetes.Clientset
	includeNormalEvents bool
}

func NewEventClient(includeNormalEvents bool) (client EventClient, err error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return
	}

	client = EventClient{
		kube:                kubeClient,
		includeNormalEvents: includeNormalEvents,
	}

	return
}

func (e EventClient) Scrape(ch chan<- prometheus.Metric) error {

	opts := metav1.ListOptions{FieldSelector: "type==Warning"}

	if e.includeNormalEvents {
		opts = metav1.ListOptions{}
	}

	list, err := e.kube.CoreV1().Events("").List(opts)
	if err != nil {
		return err
	}

	for _, event := range list.Items {
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
