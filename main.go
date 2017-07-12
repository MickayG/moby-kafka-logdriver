package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/sdk"
	"strings"
)

var logLevels = map[string]logrus.Level{
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
}

func main() {
	levelVal := os.Getenv("LOG_LEVEL")
	setLogLevel(levelVal)

	addrs := getKafkaBrokersEnv()
	outputTopic := getKafkaTopicEnv()
	keyStrat := getKeyStrategyEnv()
	partitionStrat := getPartitionStrategyEnv()
	tag := getKeyTagEnv()

	client, err := CreateClient(addrs, partitionStrat)
	if err != nil {
		panic(err)
	}

	h := sdk.NewHandler(`{"Implements": ["LoggingDriver"]}`)
	setupDockerHandlers(&h, newDriver(&client, outputTopic, keyStrat, tag))
	if err := h.ServeUnix("kafka-logdriver", 0); err != nil {
		panic(err)
	}

	client.Close()
}

func getPartitionStrategyEnv() PartitionStrategy {
	partitionStrategyStr := os.Getenv("PARTITION_STRATEGY")
	if partitionStrategyStr == "" {
		partitionStrategyStr = "round_robin"
	}
	partitionStrat, err := getPartitionStrategyFromString(partitionStrategyStr)
	if err != nil {
		logrus.Error("unknown partition strategy", err)
		os.Exit(1)
	}
	return partitionStrat
}

func getKafkaBrokersEnv() []string {
	addrList := os.Getenv("KAFKA_BROKER_ADDR")
	if addrList == "" {
		logrus.Error("Missing environment var KAFKA_BROKER_ADDR")
		os.Exit(1)
	}
	addrs := strings.Split(addrList, ",")
	return addrs
}

func getKafkaTopicEnv() string {
	outputTopic := os.Getenv("LOG_TOPIC")
	if outputTopic == "" {
		outputTopic = "dockerlogs"
	}
	return outputTopic
}

func getKeyStrategyEnv() KeyStrategy {
	keyStrategyVar := os.Getenv("KEY_STRATEGY")
	if keyStrategyVar == "" {
		keyStrategyVar = "key_by_timestamp"
		logrus.Debug("Defaulting KEY_STRATEGY to key_by_timestamp")
	}

	keyStrat, err := getKeyStrategyFromString(keyStrategyVar)
	if err != nil {
		logrus.Error("unknown key strategy", err)
		os.Exit(1)
	}
	return keyStrat
}

func getKeyTagEnv() string {
	tag := os.Getenv("LOG_TAG")
	if tag == "" {
		tag = "common"
	}
	return tag
}

func setLogLevel(levelVal string) {
	if levelVal == "" {
		levelVal = "info"
	}
	if level, exists := logLevels[levelVal]; exists {
		logrus.SetLevel(level)
	} else {
		logrus.Error("invalid log level: ", levelVal)
		os.Exit(1)
	}
}
