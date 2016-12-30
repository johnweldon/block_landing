IMAGE_BIN=block_landing
IMAGE_NAME=$(IMAGE_BIN)
IMAGE_DIR=.
BUILD_DIR="/go/src/github.com/johnweldon/block_landing"

REV=$(shell git rev-parse --short HEAD)
PWD=$(shell pwd)

all:
	@echo "Targets: build, image, deploy"

$(IMAGE_NAME):
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -a -ldflags '-s -w' .

image: build
	docker build -t $(IMAGE_NAME):$(REV) -t $(IMAGE_NAME):latest $(IMAGE_DIR)

push:
	docker push $(IMAGE_NAME):latest
	docker push $(IMAGE_NAME):$(REV)

build: $(IMAGE_NAME)
	docker run --rm -v "$(PWD)":"$(BUILD_DIR)" -w "$(BUILD_DIR)" golang:alpine go build -v 

.PHONY: all image build clean push
