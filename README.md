# Kubernetes Events Exporter

A prometheus exporter that exposes kubernetes events as prometheus metrics. By default the exporter only exposes Warning events but it can be configured via an environment variable to expose normal events as well.

# Configuration
To configure the exporter, the following environment variables are available:

* `KEE_PORT` - The port for the exporter to listen on
* `KEE_INCLUDE_NORMAL_EVENTS` - A boolean to control the inclusion of Normal typed events

# Example Kubernetes spec
```yaml
apiVersion: v1
kind: Service
metadata:
  name: kube-event-exporter
  namespace: default
  annotations:
    prometheus.io/scrape: "true"
spec:
  ports:
  - name: http
    port: 9111
    targetPort: 9111
  selector:
    app: kube-event-exporter
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kube-event-exporter
  namespace: default
spec:
  replicas: 3
  selector:
    matchLabels:
      app: kube-event-exporter
  template:
    metadata:
      labels:
        app: kube-event-exporter
    spec:
      containers:
      - name: kube-event-exporter
        image: itmecho/kube_event_exporter
        ports:
        - name: http
          containerPort: 9111
        resources:
          requests:
            cpu: 100m
            memory: 50Mi
          limits:
            cpu: 200m
            memory: 100Mi
```