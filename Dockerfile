FROM golang:1.11-alpine as build

RUN apk add --no-cache git build-base

ADD . /src
WORKDIR /src

ENV GO111MODULE on
ENV CGO_ENABLED 0

RUN go build -o kube_event_exporter -ldflags '-extldflags "-static"'


FROM scratch

COPY --from=build /src/kube_event_exporter /kube_event_exporter

EXPOSE 9111

ENTRYPOINT  [ "/kube_event_exporter" ]