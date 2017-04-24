.DEFAULT_GOAL := test

NAME=mickyg/kafka-logdriver
TAG=latest

clean:
	@rm -rf kafka-logdriver.tar.gz
	@rm -rf kafka-logdriver

test:
	@go get -t
	@go test

package:
	@docker build -t mickyg/kafka-logdriver:latest . -f Dockerfile.build
	@docker create --name kafka-logdriver-release mickyg/kafka-logdriver:latest
	@docker cp kafka-logdriver-release:kafka-logdriver.tar.gz .
	@docker rm -v kafka-logdriver-release

install: clean package
	@tar -xvf kafka-logdriver.tar.gz
	@docker plugin create ${NAME}:${TAG} kafka-logdriver

enable:
	@docker plugin enable ${NAME}:${TAG}

push:
	@docker plugin push ${NAME}:${TAG}
