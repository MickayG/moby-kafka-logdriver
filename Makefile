.DEFAULT_GOAL := test

NAME=mickyg/kafka-logdriver
TAG=latest

clean:
	-rm -rf kafka-logdriver.tar.gz
	-rm -rf kafka-logdriver

test:
	@go get -t
	@go test

package:
	@docker build -t tempo . -f Dockerfile.build
	@docker create --name kafka-logdriver-release tempo
	@docker cp kafka-logdriver-release:kafka-logdriver.tar.gz .
	@docker rm -v kafka-logdriver-release
	@docker image rm tempo

install: clean package
	@tar -xvf kafka-logdriver.tar.gz
	@docker plugin create ${NAME}:${TAG} kafka-logdriver
	@echo "Now configure the KAFKA_BROKER with 'docker plugin set ${NAME}:${TAG} KAFKA_BROKER_ADDR=< Broker list here >'"
	@echo "Once configured, run 'make enable' to enable the plugin"

enable:
	@docker plugin enable ${NAME}:${TAG}

uninstall:
	-docker plugin disable ${NAME}:${TAG}
	@docker plugin rm ${NAME}:${TAG}

push:
	@docker plugin push ${NAME}:${TAG}
