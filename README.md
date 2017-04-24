Kafka Log Driver  [![Build Status](https://travis-ci.org/MickayG/moby-kafka-logdriver.svg?branch=master)](https://travis-ci.org/MickayG/moby-kafka-logdriver)
=======================

An upcoming build of Moby/Docker (expected to be 17.05) will
support log driver plugins which allows people to add extra
functionality to handle logs coming out of docker containers.

This plugin allows users to route all Moby/Docker logs to Kafka.

Installation
-----------------
**This plugin requires at least Docker version 17.05**



Configuration
----------------

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

