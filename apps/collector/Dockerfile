#FROM kubemq/gobuilder as builder
FROM kubemq/gobuilder-ubuntu as builder
ARG VERSION
ARG GIT_COMMIT
ARG BUILD_TIME
ENV GOPATH=/go
ENV PATH=$GOPATH:$PATH
ENV ADDR=0.0.0.0
ADD . $GOPATH/github.com/kubemq-hub/kubemq-collector
WORKDIR $GOPATH/github.com/kubemq-hub/kubemq-collector
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -mod=vendor -installsuffix cgo -ldflags="-w -s -X main.version=$VERSION" -o kubemq-collector-run .
FROM registry.access.redhat.com/ubi8/ubi-minimal
MAINTAINER KubeMQ info@kubemq.io
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH
RUN mkdir /kubemq-collector
COPY --from=builder $GOPATH/github.com/kubemq-hub/kubemq-collector/kubemq-collector-run ./kubemq-collector
COPY --from=builder $GOPATH/github.com/kubemq-hub/kubemq-collector/dist ./kubemq-collector/dist
RUN chown -R 1001:root  /kubemq-collector && chmod g+rwX  /kubemq-collector
WORKDIR kubemq-collector
USER 1001
EXPOSE 8085
CMD ["./kubemq-collector-run"]
