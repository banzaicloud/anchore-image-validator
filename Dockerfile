FROM golang:1-alpine

WORKDIR /go/src/github.com/banzaicloud/anchore-image-validator
COPY . .
RUN cd cmd && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /anchore-image-validator
CMD ["/anchore-image-validator"]
