FROM golang:1-alpine

ENV NAME=anchore-image-validator

WORKDIR /go/src/github.com/banzaicloud/$NAME
COPY . .
RUN cd cmd && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /$NAME
CMD ["/$NAME"]
