

all:
	@echo "Targets: build, deploy"


run: deploy

build:
	docker build -t block_landing .


deploy:
	docker run -d -p 0.0.0.0:9000:9000 --restart=always block_landing
