FROM kubemq/gobuilder as builder
RUN  mkdir -p /go/src \
  && mkdir -p /go/bin \
  && mkdir -p /go/pkg
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH
ENV ADDR=0.0.0.0
ADD . $GOPATH/src/github.com/kubemq-io/sender
WORKDIR $GOPATH/src/github.com/kubemq-io/sender
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags='-w -s' -o sender-run .
FROM alpine:latest
RUN apk add ca-certificates
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

ENV ADDR=0.0.0.0
RUN  mkdir -p /sender

COPY --from=builder $GOPATH/src/github.com/kubemq-io/sender/sender-run /sender
WORKDIR /sender
RUN ls
RUN chmod u+x ./sender-run
CMD ["./sender-run"]

