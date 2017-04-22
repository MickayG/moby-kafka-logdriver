package main

import (
	"github.com/Shopify/sarama"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/daemon/logger"
	"time"
	"github.com/docker/docker/api/types/backend"
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

// An mapped version of logger.Message where Line is a String, not a byte array
type stringMessage struct {
	Line      string
	Source    string
	Timestamp time.Time
	Attrs     backend.LogAttributes
	Partial   bool

	// Err is an error associated with a message. Completeness of a message
	// with Err is not expected, tho it may be partially complete (fields may
	// be missing, gibberish, or nil)
	Err error

}

// Converts a log message into JSON and writes it to the producer, ready to be written to the broker
//  Returns an error if any occurred.
func WriteMessage(msg logger.Message, containerId string, producer sarama.AsyncProducer) error {

	strMsg := asStringMessage(msg)

	asJson, err := json.Marshal(strMsg)
	if err != nil {
		logrus.WithField("id", containerId).WithError(err).WithField("message", strMsg).Error("error converting log message to json")
		return err
	}

	producer.Input() <- &sarama.ProducerMessage{Topic: "logs", Key: sarama.StringEncoder(strMsg.Source), Value: sarama.StringEncoder(asJson), Timestamp: strMsg.Timestamp}

	return nil
}
func asStringMessage(message logger.Message) stringMessage {
	var strMsg stringMessage
	strMsg.Line = string(message.Line)
	strMsg.Timestamp = message.Timestamp
	strMsg.Partial = message.Partial
	strMsg.Attrs = message.Attrs
	strMsg.Source = message.Source
	strMsg.Err = message.Err
	return strMsg
}