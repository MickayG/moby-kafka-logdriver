package main

import (
	"github.com/Shopify/sarama"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	_ "strings"
	"errors"
)

// Create a Sarama Kafka client with the broker list
// Will pass back the client and any errors
func CreateClient(brokerList []string, partitionStrategy PartitionStrategy) (sarama.Client, error){
	conf := sarama.NewConfig()

	switch partitionStrategy {
	case PARTITION_KEY_HASH:
		conf.Producer.Partitioner = sarama.NewHashPartitioner
	case PARTITION_ROUND_ROBIN:
		conf.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	default:
		err := errors.New("Unknown partition strategy: " + string(partitionStrategy))
		return nil, err
	}


	return sarama.NewClient(brokerList, conf)
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