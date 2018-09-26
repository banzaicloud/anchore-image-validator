FROM golang:1.11-alpine

ADD . /go/src/github.com/banzaicloud/anchore-image-validator
WORKDIR /go/src/github.com/banzaicloud/anchore-image-validator
RUN apk update && apk add ca-certificates make git curl mercurial

RUN make vendor
RUN go build -o /tmp/anchore-image-validator ./cmd

FROM alpine:3.7

COPY --from=0 /tmp/anchore-image-validator /usr/local/bin/anchore-image-validator
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

RUN adduser -D anchore-image-validator

USER anchore-image-validator

ENTRYPOINT ["/usr/local/bin/anchore-image-validator"]