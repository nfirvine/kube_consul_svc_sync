# Build Stage
FROM golang:1.9-alpine3.7 AS build-stage

LABEL app="build-kube_consul_svc_sync"
LABEL REPO="https://github.com/nfirvine/kube_consul_svc_sync"

ENV GOROOT=/usr/lib/go \
    GOPATH=/gopath \
    GOBIN=/gopath/bin \
    PROJPATH=/gopath/src/github.com/nfirvine/kube_consul_svc_sync

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /gopath/src/github.com/nfirvine/kube_consul_svc_sync
WORKDIR /gopath/src/github.com/nfirvine/kube_consul_svc_sync

RUN make build-alpine

# Final Stage
FROM alpine:3.7

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/nfirvine/kube_consul_svc_sync"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/kube_consul_svc_sync/bin

WORKDIR /opt/kube_consul_svc_sync/bin

COPY --from=build-stage /gopath/src/github.com/nfirvine/kube_consul_svc_sync/bin/kube_consul_svc_sync /opt/kube_consul_svc_sync/bin/
RUN chmod +x /opt/kube_consul_svc_sync/bin/kube_consul_svc_sync

CMD /opt/kube_consul_svc_sync/bin/kube_consul_svc_sync