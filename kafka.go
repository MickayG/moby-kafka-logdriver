package main

import (
	"github.com/Shopify/sarama"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"strings"
	"errors"
)

type KeyStrategy int

const (
	KEY_BY_CONTAINER_ID KeyStrategy = iota
	KEY_BY_TIMESTAMP    KeyStrategy = iota
)

func getKeyStrategyFromString(keyStrategyString string) (KeyStrategy, error) {
	// Trim and whitespace and lowercase the string so it matches
	// no matter what someone has put in
	switch strings.TrimSpace(strings.ToLower(keyStrategyString)) {
	case "key_by_container_id":
		return KEY_BY_CONTAINER_ID, nil
	case "key_by_timestamp":
		return KEY_BY_TIMESTAMP, nil
	default:
		return 0, errors.New("Unknown keying strategy " + keyStrategyString)
	}
}

// Create a Sarama Kafka client with the broker list
// Will pass back the client and any errors
func CreateClient(brokerList []string) (sarama.Client, error){
	return sarama.NewClient(brokerList, nil)
}

// Create a producer from a Sarama Kafka client
func CreateProducer(client *sarama.Client) (sarama.AsyncProducer, error){
	return sarama.NewAsyncProducerFromClient(*client)
}

// Converts a log message into JSON and writes it to the producer, ready to be written to the broker
//  Returns an error if any occurred.
func WriteMessage(topic string, msg LogMessage, containerId string, keyStrategy KeyStrategy,producer sarama.AsyncProducer) error {

	asJson, err := json.Marshal(msg)
	if err != nil {
		logrus.WithField("id", containerId).WithError(err).WithField("message", msg).Error("error converting log message to json")
		return err
	}

	key := keyBy(msg, containerId, keyStrategy)
	producer.Input() <- &sarama.ProducerMessage{Topic: topic, Key: sarama.StringEncoder(key), Value: sarama.StringEncoder(asJson), Timestamp: msg.Timestamp}

	return nil
}

func keyBy(msg LogMessage, containerId string,  strategy KeyStrategy) string {
	switch strategy {
		case KEY_BY_CONTAINER_ID:
			return containerId
		case KEY_BY_TIMESTAMP:
			return string(msg.Timestamp.Unix())
	default:
		logrus.WithField("keyStrategy", strategy).Error("Unknown key strategy. Defaulting to KEY_BY_CONTAINER_ID")
		return keyBy(msg, containerId, KEY_BY_CONTAINER_ID)
	}
}