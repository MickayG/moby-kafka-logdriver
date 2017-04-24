package main

import (
	"fmt"
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

	client, err := CreateClient(addrs, partitionStrat)
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to connect to kafka:", err)
		os.Exit(1)
	}

	h := sdk.NewHandler(`{"Implements": ["LoggingDriver"]}`)
	setupDockerHandlers(&h, newDriver(&client, outputTopic, keyStrat))
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
		fmt.Fprintln(os.Stderr, "unknown partition strategy", err)
		os.Exit(1)
	}
	return partitionStrat
}

func getKafkaBrokersEnv() []string {
	addrList := os.Getenv("KAFKA_BROKER_ADDR")
	if addrList == "" {
		fmt.Fprintln(os.Stderr, "Missing environment var KAFKA_BROKER_ADDR")
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
	}
	keyStrat, err := getKeyStrategyFromString(keyStrategyVar)
	if err != nil {
		fmt.Fprintln(os.Stderr, "unknown key strategy", err)
		os.Exit(1)
	}
	return keyStrat
}

func setLogLevel(levelVal string) {
	if levelVal == "" {
		levelVal = "info"
	}
	if level, exists := logLevels[levelVal]; exists {
		logrus.SetLevel(level)
	} else {
		fmt.Fprintln(os.Stderr, "invalid log level: ", levelVal)
		os.Exit(1)
	}
}
