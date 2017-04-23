package main

import (
	"testing"
	"bytes"
	"github.com/docker/docker/api/types/plugins/logdriver"
	"time"
	"github.com/docker/docker/pkg/ioutils"
	protoio "github.com/gogo/protobuf/io"
	"encoding/binary"
	"io"
	"github.com/Shopify/sarama/mocks"
	"github.com/docker/docker/daemon/logger"
)


func TestConsumesSingleLogMessagesFromDocker(t *testing.T) {
	producer := NewProducer(t)
	defer producer.Close()

	logMsg := newLogEntry("alpha")

	stream := createBufferForLogMessages([]logdriver.LogEntry{logMsg})

	lf := createLogPair(producer, stream)

	producer.ExpectInputAndSucceed()
	ConsumeLog(&lf, "topic")

	recvMsg := <-producer.Successes()
	assertLineMatch(t, "alpha", recvMsg)
}


func TestConsumesMultipleLogMessagesFromDocker(t *testing.T) {
	producer := NewProducer(t)
	defer producer.Close()

	stream := createBufferForLogMessages([]logdriver.LogEntry{
		newLogEntry("alpha"),
		newLogEntry("beta"),
		newLogEntry("charlie"),
		newLogEntry("delta"),
	})

	lf := createLogPair(producer, stream)

	producer.ExpectInputAndSucceed()
	producer.ExpectInputAndSucceed()
	producer.ExpectInputAndSucceed()
	producer.ExpectInputAndSucceed()
	ConsumeLog(&lf, "topic")

	assertLineMatch(t, "alpha", <-producer.Successes())
	assertLineMatch(t, "beta", <-producer.Successes())
	assertLineMatch(t, "charlie", <-producer.Successes())
	assertLineMatch(t, "delta", <-producer.Successes())
}




func createLogPair(producer *mocks.AsyncProducer, stream io.ReadCloser) logPair {
	var lf logPair
	lf.producer = producer
	lf.stream = stream
	lf.info = logger.Info{}
	return lf
}


func createBufferForLogMessages(logs []logdriver.LogEntry) io.ReadCloser {
	var buf bytes.Buffer

	protoWriter := protoio.NewUint32DelimitedWriter(&buf, binary.BigEndian)

	for _,log := range logs {
		protoWriter.WriteMsg(&log)
	}

	protoWriter.Close()

	closeFunc := func () error {
		return nil
	}

	readCloser := ioutils.NewReadCloserWrapper(&buf, closeFunc)
	return readCloser
}


func newLogEntry(line string) logdriver.LogEntry {
	var le logdriver.LogEntry
	le.Line = []byte(line)
	le.Source = "container"
	le.Partial = false
	le.TimeNano = time.Now().UnixNano()
	return le
}
