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

	addrList := os.Getenv("KAFKA_BROKER_ADDR")
	if addrList == "" {
		fmt.Fprintln(os.Stderr, "Missing environment var KAFKA_BROKER_ADDR")
		os.Exit(1)
	}
	addrs := strings.Split(addrList, ",")

	outputTopic := os.Getenv("LOG_TOPIC")
	if outputTopic == "" {
		outputTopic = "dockerlogs"
	}

	client, err := CreateClient(addrs)
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to connect to kafka", err)
	}

	h := sdk.NewHandler(`{"Implements": ["LoggingDriver"]}`)
	handlers(&h, newDriver(&client, outputTopic))
	if err := h.ServeUnix("kafka-logdriver", 0); err != nil {
		panic(err)
	}

	client.Close()
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
