package kafkaproxy

import (
	"os"
	"time"

	"github.com/severuykhin/logfmt/v4"
)

type kafkaProxyConfig struct {
	Port         string
	Logger       logger
	BatchTimeout time.Duration
	BatchSize    int
	WriteTimeout time.Duration
}

var defaultPort = "80"
var defaultLogger = logfmt.New(os.Stdout, logfmt.L_INFO)

var defaultWriterBatchSize = 10
var defaultWriterBatchTimeout = time.Millisecond * 100
var defaultWriterWriteTimeout = time.Second

type optFunc func(*kafkaProxyConfig) *kafkaProxyConfig

func WithLogger(l logger) optFunc {
	return func(c *kafkaProxyConfig) *kafkaProxyConfig {
		c.Logger = l
		return c
	}
}

func WithPort(port string) optFunc {
	return func(c *kafkaProxyConfig) *kafkaProxyConfig {
		c.Port = port
		return c
	}
}

func WithBatchSize(bashSize int) optFunc {
	return func(c *kafkaProxyConfig) *kafkaProxyConfig {
		c.BatchSize = bashSize
		return c
	}
}

func WithBatchTimeout(batchTimeout time.Duration) optFunc {
	return func(c *kafkaProxyConfig) *kafkaProxyConfig {
		c.BatchTimeout = batchTimeout
		return c
	}
}

func WithWriteTimeout(writeTimeout time.Duration) optFunc {
	return func(c *kafkaProxyConfig) *kafkaProxyConfig {
		c.WriteTimeout = writeTimeout
		return c
	}
}
