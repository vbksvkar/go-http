package clients

import (
	"net/http"
	"time"
)

type clientConfig struct {
	timeout            time.Duration
	maxIdleConnections int
}

type ClientOptionFunc func(c *clientConfig)

func WithTimeout(timeout time.Duration) ClientOptionFunc {
	return func(c *clientConfig) {
		c.timeout = timeout
	}
}

func WithMaxIdleConnections(maxIdleConnections int) ClientOptionFunc {
	return func(c *clientConfig) {
		c.maxIdleConnections = maxIdleConnections
	}
}

func New(opts ...ClientOptionFunc) *http.Client {
	config := &clientConfig{
		timeout:            10 * time.Second,
		maxIdleConnections: 10,
	}

	defaultTransport := http.DefaultTransport.(*http.Transport).Clone()
	defaultTransport.MaxIdleConnsPerHost = config.maxIdleConnections
	defaultTransport.MaxIdleConns = config.maxIdleConnections
	defaultTransport.MaxConnsPerHost = config.maxIdleConnections

	for _, opt := range opts {
		opt(config)
	}

	return &http.Client{
		Timeout:   config.timeout,
		Transport: defaultTransport,
	}
}
