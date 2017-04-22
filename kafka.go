package main

import (
	"github.com/Shopify/sarama"
	"encoding/json"
	"github.com/Sirupsen/logrus"
)

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
func WriteMessage(msg LogMessage, containerId string, producer sarama.AsyncProducer) error {

	asJson, err := json.Marshal(msg)
	if err != nil {
		logrus.WithField("id", containerId).WithError(err).WithField("message", msg).Error("error converting log message to json")
		return err
	}

	producer.Input() <- &sarama.ProducerMessage{Topic: "logs", Key: sarama.StringEncoder(msg.Source), Value: sarama.StringEncoder(asJson), Timestamp: msg.Timestamp}

	return nil
}

