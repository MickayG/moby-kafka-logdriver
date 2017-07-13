package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/sdk"
)

var logLevels = map[string]logrus.Level{
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
}

func main() {
	levelVal := os.Getenv(ENV_LOG_LEVEL)
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


func setLogLevel(levelVal string) {
	if levelVal == "" {
		levelVal = DEFAULT_LOG_LEVEL
	}
	if level, exists := logLevels[levelVal]; exists {
		logrus.SetLevel(level)
	} else {
		logrus.Error("invalid log level: ", levelVal)
		os.Exit(1)
	}
}
