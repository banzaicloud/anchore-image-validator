NAME = "anchore-image-validator"

.PHONY: build

build:
	@cd cmd && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${NAME}
