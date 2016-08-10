IMAGE_NAME=block_landing
IMAGE_DIR=.


all:
	@echo "Targets: build, image, deploy"

$(IMAGE_NAME):
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -a -ldflags '-s -w' .

image: build
	docker build -t $(IMAGE_NAME) $(IMAGE_DIR)

deploy: image
	docker run -d --name=$(IMAGE_NAME) -p 9000:9000 --restart=always $(IMAGE_NAME)

clean:
	-docker rm -v -f $(IMAGE_NAME)
	-docker rmi $(IMAGE_NAME)
	-rm ./$(IMAGE_NAME)

build: $(IMAGE_NAME)

run: deploy
