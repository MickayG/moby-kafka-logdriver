FROM  golang:1.8

COPY . /go/src/github.com/MickayG/moby-kafka-logdriver
WORKDIR /go/src/github.com/MickayG/moby-kafka-logdriver
RUN go build --ldflags '-extldflags "-static"' -o kafka-logdriver

# Prepare the directory for packaging
RUN mkdir -p /kafka-logdriver/rootfs/usr/bin/
RUN cp kafka-logdriver /kafka-logdriver/rootfs/usr/bin/
RUN cp config.json /kafka-logdriver

WORKDIR /
RUN tar -zcvf kafka-logdriver.tar.gz kafka-logdriver