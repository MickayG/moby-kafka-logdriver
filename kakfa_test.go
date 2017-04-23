package main

import (
	"testing"
	"github.com/Shopify/sarama/mocks"
	"github.com/Shopify/sarama"
	"time"
	"encoding/json"

	"github.com/stretchr/testify/assert"
)

func TestWriteSingleMessage(t *testing.T) {
	producer := newProducer(t)

	expectedTime := time.Now()
	expectedSource := "containerABC"
	expectedLine := "I am a log message"

	msg := newMessage(expectedTime, expectedSource, expectedLine)

	producer.ExpectInputAndSucceed()
	WriteMessage("topic1", msg, expectedSource, producer)

	writtenMsg := <-producer.Successes()

	assertLineMatch(t, expectedLine, writtenMsg)
}

func TestWriteMultipleMessagesToSameTopic(t *testing.T) {
	producer := newProducer(t)

	msg1 := newMessage(time.Now(), "1", "a")
	msg2 := newMessage(time.Now(), "2", "b")
	msg3 := newMessage(time.Now(), "3", "c")

	// Need to call this three times to expect three messages
	producer.ExpectInputAndSucceed()
	producer.ExpectInputAndSucceed()
	producer.ExpectInputAndSucceed()

	WriteMessage("topic1", msg1, msg1.Source, producer)
	WriteMessage("topic1", msg2, msg2.Source, producer)
	WriteMessage("topic1", msg3, msg3.Source, producer)

	out1 := <-producer.Successes()
	out2 := <-producer.Successes()
	out3 := <-producer.Successes()

	assert.NotNil(t, out1)
	assertLineMatch(t,"a", out1)
	assertTopic(t, "topic1", out1)

	assert.NotNil(t, out2)
	assertLineMatch(t,"b", out2)
	assertTopic(t, "topic1", out1)

	assert.NotNil(t, out3)
	assertLineMatch(t,"c", out3)
	assertTopic(t, "topic1", out1)
}

func TestWriteMultipleMessagesToDifferentTopics(t *testing.T) {
	producer := newProducer(t)

	msg1 := newMessage(time.Now(), "1", "a")
	msg2 := newMessage(time.Now(), "2", "b")
	msg3 := newMessage(time.Now(), "3", "c")

	// Need to call this three times to expect three messages
	producer.ExpectInputAndSucceed()
	producer.ExpectInputAndSucceed()
	producer.ExpectInputAndSucceed()

	WriteMessage("topic1", msg1, msg1.Source, producer)
	WriteMessage("topic2", msg2, msg2.Source, producer)
	WriteMessage("topic3", msg3, msg3.Source, producer)

	out1 := <-producer.Successes()
	out2 := <-producer.Successes()
	out3 := <-producer.Successes()

	assert.NotNil(t, out1)
	assertLineMatch(t,"a", out1)
	assertTopic(t, "topic1", out1)

	assert.NotNil(t, out2)
	assertLineMatch(t,"b", out2)
	assertTopic(t, "topic2", out1)

	assert.NotNil(t, out3)
	assertLineMatch(t,"c", out3)
	assertTopic(t, "topic3", out1)
}

func assertTopic(t *testing.T, expectedTopic string, message *sarama.ProducerMessage) {
	assert.Equal(t, expectedTopic, message.Topic)
}

func assertLineMatch(t *testing.T, expectedLine string, message *sarama.ProducerMessage) {
	msgContentBytes, err := message.Value.Encode()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var outputJson LogMessage
	json.Unmarshal(msgContentBytes, outputJson)
	assert.Equal(t, expectedLine, outputJson.Line)
}

func newMessage(expectedTime time.Time, expectedSource string, expectedLine string) (LogMessage) {
	var msg LogMessage
	msg.Timestamp = expectedTime
	msg.Source = expectedSource
	msg.Line = expectedLine
	msg.Partial = false
	return msg
}

func newProducer(t *testing.T) *mocks.AsyncProducer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer := mocks.NewAsyncProducer(t, config)
	return producer
}
