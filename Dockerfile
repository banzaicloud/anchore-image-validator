FROM golang:1.16.4-alpine AS builder

RUN apk add --update --no-cache ca-certificates git

RUN mkdir -p /build
WORKDIR /build

COPY go.* /build/
RUN go mod download
COPY . /build
RUN go install ./cmd

FROM alpine:3.13.5

COPY --from=builder /go/bin/cmd /usr/local/bin/anchore-image-validator
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/usr/local/bin/anchore-image-validator"]
