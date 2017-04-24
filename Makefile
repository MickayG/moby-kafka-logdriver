NAME: mickyg/kafka-logdriver
TAG: 0.1

clean:
	@rm -rf kafka-logdriver.tar.gz

test:
	@go test

package:
	@docker build -t mickyg/kafka-logdriver:latest . -f Dockerfile.build
	@docker create --name kafka-logdriver-release mickyg/kafka-logdriver:latest
	@docker cp kafka-logdriver-release:kafka-logdriver.tar.gz .
	@docker rm -v kafka-logdriver-release

install: clean package
	@tar -xvf kafka-logdriver.tar.gz
	@docker plugin create ${NAME}:${TAG} kafka-logdriver

enable: install
	@docker plugin enable ${NAME}:${TAG}

push: enable
	@docker plugin push ${NAME}:${TAG}

