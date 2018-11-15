.PHONY: build

build:
	go build

docker:
	docker build -t itmecho/kube_event_exporter .