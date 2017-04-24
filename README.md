# Kafka Log Driver  [![Build Status](https://travis-ci.org/MickayG/moby-kafka-logdriver.svg?branch=master)](https://travis-ci.org/MickayG/moby-kafka-logdriver)


An upcoming build of Moby/Docker (expected to be 17.05) will support log driver plugins which allows people to add extra functionality to handle logs coming out of docker containers.

This plugin allows users to route all Moby/Docker logs to Kafka.

## Installation

***This plugin requires at least Docker version 17.05***

There are two ways to install the plugin, via Dockerhub or building
from this repository.

### Install the plugin from Dockerhub:

```
docker plugin install mickyg/kafka-logdriver:latest
```
Now configure it as per the Configuration section below, **you must set the KAFKA_BROKER_ADDR option**, then enable the plugin:
```
docker plugin enable mickyg/kafka-logdriver:latest
```
Now it's installed! To use the log driver for a given container, use the `--logdriver` flag. For example, to start the hello-world
container with all of it's logs being sent to Kafka, run the command:
```
docker run --log-driver mickyg/kafka-logdriver:latest hello-world
```

### Install from source

Clone the project
```
git clone https://github.com/MickayG/moby-kafka-logdriver.git
cd moby-kafka-logdriver
```
Build the plugin and install it
```
make install
```
Set the KAFKA_BROKER_ARG variable. In the example below the host 192.168.0.1 is a Kafka broker.
```
docker plugin set mickyg/kafka-logdriver:latest KAFKA_BROKER_ADDR="192.168.0.1:9092"
```
Enable the plugin:
```
make enable
```
Now test it! Connect to your broker and consume from the "dockerlogs" topic
(Topic can be changed via environment variable, see below). Then launch a container:
```
docker run --log-driver mickyg/kafka-logdriver:latest hello-world
```

## Configuration


Once the plugin has been installed, you can modify the below arguments with the command

`docker plugin set kafka-logdriver <OPTION>=<VALUE>`

For example, to change the topic to "logs"

`docker plugin set kafka-logdriver LOG_TOPIC=logs`


| Option | Description | Default |
| -------|------------| --------|
|KAFKA_BROKER_ADDR|**(Required)** Comma delimited list of Kafka brokers. | |
|LOG_TOPIC| Topic to which logs will be written to | dockerlogs |
|KEY_STRATEGY| Method in which Kafka methods should be keyed. Options are: <br>*key_by_timestamp* - Key each message by the timestamp of the log message <br>*key_by_container_id* - Key each message with the container id. | key_by_timestamp
|PARTITION_STRATEGY| Kafka partitioner type. Options are:<br>*round_robin* - Write to each partition one after another, i.e equally distributed<br>*key_hash* - Partition based on the hash of the message key|round_robin|
|LOG_LEVEL| Log level of the internal logger. Options: debug, info, warn, error|info|

