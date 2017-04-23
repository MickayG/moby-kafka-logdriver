package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types/plugins/logdriver"
	"github.com/docker/docker/daemon/logger"
	protoio "github.com/gogo/protobuf/io"
	"github.com/pkg/errors"
	"github.com/tonistiigi/fifo"

	"github.com/Shopify/sarama"

	"github.com/docker/docker/api/types/backend"
)

// An mapped version of logger.Message where Line is a String, not a byte array
type LogMessage struct {
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

type LogDriver interface {
	StartLogging(file string, logCtx logger.Info) error
	StopLogging(file string) error
	ReadLogs(info logger.Info, config logger.ReadConfig) (io.ReadCloser, error)
	GetCapability() logger.Capability
}

type KafkaDriver struct {
	mu     sync.Mutex
	logs   map[string]*logPair
	idx    map[string]*logPair
	logger logger.Logger
	client *sarama.Client
	outputTopic string
}

type logPair struct {
	stream io.ReadCloser
	info   logger.Info
	producer sarama.AsyncProducer
}

func newDriver(client *sarama.Client, outputTopic string) *KafkaDriver {
	return &KafkaDriver{
		logs: make(map[string]*logPair),
		idx:  make(map[string]*logPair),
		client: client,
		outputTopic: outputTopic,
	}
}

func (d *KafkaDriver) StartLogging(file string, logCtx logger.Info) error {
	d.mu.Lock()
	if _, exists := d.logs[file]; exists {
		d.mu.Unlock()
		return fmt.Errorf("logger for %q already exists", file)
	}
	d.mu.Unlock()

	if logCtx.LogPath == "" {
		logCtx.LogPath = filepath.Join("/var/log/docker", logCtx.ContainerID)
	}
	if err := os.MkdirAll(filepath.Dir(logCtx.LogPath), 0755); err != nil {
		return errors.Wrap(err, "error setting up logger dir")
	}

	logrus.WithField("id", logCtx.ContainerID).WithField("file", file).WithField("logpath", logCtx.LogPath).Debugf("Start logging")
	f, err := fifo.OpenFifo(context.Background(), file, syscall.O_RDONLY, 0700)
	if err != nil {
		return errors.Wrapf(err, "error opening logger fifo: %q", file)
	}

	d.mu.Lock()

	producer, err := CreateProducer(d.client)
	if err != nil {
		return errors.Wrapf(err,"unable to create kafka consumer")
	}

	lf := &logPair{ f, logCtx, producer}
	d.logs[file] = lf
	d.idx[logCtx.ContainerID] = lf

	d.mu.Unlock()

	go ConsumeLog(lf, d.outputTopic)

	return nil
}

func (d *KafkaDriver) StopLogging(file string) error {
	logrus.WithField("file", file).Debugf("Stop logging")
	d.mu.Lock()
	lf, ok := d.logs[file]
	if ok {
		lf.stream.Close()
		delete(d.logs, file)

		lf.producer.Close()
	}
	d.mu.Unlock()
	return nil
}

func (d* KafkaDriver) GetCapability() logger.Capability {
	return logger.Capability{ReadLogs: false}
}

func ConsumeLog(lf *logPair, topic string) {
	dec := protoio.NewUint32DelimitedReader(lf.stream, binary.BigEndian, 1e6)
	defer dec.Close()
	var buf logdriver.LogEntry
	for {
		// Check if there are any Kafka errors thus far
		select {
			case kafkaErr := <- lf.producer.Errors():
				// In the event of an error, continue to attempt to write messages
				logrus.Error("error recieved from Kafka", kafkaErr)
			default:
				//No errors, continue
		}

		if err := dec.ReadMsg(&buf); err != nil {
			if err == io.EOF {
				logrus.WithField("id", lf.info.ContainerID).WithError(err).Debug("shutting down log logger")
				lf.stream.Close()
				return
			}
			dec = protoio.NewUint32DelimitedReader(lf.stream, binary.BigEndian, 1e6)
		}

		var msg LogMessage
		msg.Line = string(buf.Line)
		msg.Source = buf.Source
		msg.Partial = buf.Partial
		msg.Timestamp = time.Unix(0, buf.TimeNano)

		err := WriteMessage(topic, msg, lf.info.ContainerID, lf.producer)
		if err != nil {
			logrus.WithField("id", lf.info.ContainerID).WithField("msg", msg).Error("Unable to write message to kafka", err)
		}

		buf.Reset()
	}
}

func (d *KafkaDriver) ReadLogs(info logger.Info, config logger.ReadConfig) (io.ReadCloser, error) {
	//TODO
	return nil, nil
}
