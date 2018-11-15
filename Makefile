.PHONY: build

build:
	go build

docker: build
	docker build -t itmecho/kube_events_exporter .