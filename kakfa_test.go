package main

import (
	"testing"
	"github.com/Shopify/sarama/mocks"
	"github.com/Shopify/sarama"
	"github.com/docker/docker/daemon/logger"
	"time"
	"encoding/json"

	"github.com/stretchr/testify/assert"
)



func TestWriteMessage(t *testing.T) {
	config  := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer := mocks.NewAsyncProducer(t, config)


	expectedTime :=  time.Now()
	expectedSource := "containerABC"
	expectedLine := "I am a log message"
	expectedAttributes := make(map[string]string, 0)

	msg := logger.NewMessage()
	msg.Timestamp = expectedTime
	msg.Source = expectedSource
	msg.Line = []byte(expectedLine)
	msg.Attrs = expectedAttributes
	msg.Partial = false

	producer.ExpectInputAndSucceed()
	WriteMessage(*msg, expectedSource, producer)

	writtenMsg := <-producer.Successes()
	msgContentBytes, err := writtenMsg.Value.Encode()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	print("Recieved message: "+ string(msgContentBytes))

	var outputJson logger.Message
	json.Unmarshal(msgContentBytes, &outputJson)

	assert.Equal(t, expectedTime.String(), outputJson.Timestamp)
	assert.Equal(t, expectedLine, outputJson.Line)
}
