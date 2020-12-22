FROM kubemq/gobuilder as builder
RUN  mkdir -p /go/src \
  && mkdir -p /go/bin \
  && mkdir -p /go/pkg
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH
ENV ADDR=0.0.0.0
ADD . $GOPATH/src/github.com/kubemq-io/receiver
WORKDIR $GOPATH/src/github.com/kubemq-io/receiver
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags='-w -s' -o receiver-run .
FROM alpine:latest
RUN apk add ca-certificates
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

ENV ADDR=0.0.0.0
RUN  mkdir -p /receiver

COPY --from=builder $GOPATH/src/github.com/kubemq-io/receiver/receiver-run /receiver
WORKDIR /receiver
RUN ls
RUN chmod u+x ./receiver-run
CMD ["./receiver-run"]

