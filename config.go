package kafkaproxy

import (
	"os"
	"time"

	"github.com/severuykhin/logfmt/v4"
)

type kafkaWriterConfig struct {
	Brokers      []string
	BatchTimeout time.Duration
	BatchSize    int
	WriteTimeout time.Duration
}

var defaultPort = "80"
var defaultLogger = logfmt.New(os.Stdout, logfmt.L_INFO)

var defaultWriterBatchSize = 10
var defaultWriterBatchTimeout = time.Millisecond * 100
var defaultWriterWriteTimeout = time.Second

type optFunc func(*kafkaproxy) *kafkaproxy

func WithLogger(l logger) optFunc {
	return func(k *kafkaproxy) *kafkaproxy {
		k.Logger = l
		return k
	}
}

func WithPort(port string) optFunc {
	return func(k *kafkaproxy) *kafkaproxy {
		k.Port = port
		return k
	}
}

func WithBatchSize(bashSize int) optFunc {
	return func(k *kafkaproxy) *kafkaproxy {
		k.KafkaWriterConfig.BatchSize = bashSize
		return k
	}
}

func WithBatchTimeout(bashTimeout time.Duration) optFunc {
	return func(k *kafkaproxy) *kafkaproxy {
		k.KafkaWriterConfig.BatchTimeout = bashTimeout
		return k
	}
}

func WithWriteTimeout(writeTimeout time.Duration) optFunc {
	return func(k *kafkaproxy) *kafkaproxy {
		k.KafkaWriterConfig.WriteTimeout = writeTimeout
		return k
	}
}
