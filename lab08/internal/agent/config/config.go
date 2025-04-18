package config

import "time"

type Config struct {
	TargetHost string
	TargetPort string

	RequestTimeout time.Duration

	Sleep  time.Duration
	Jitter float64

	AgentUUID string
	Endpoint  string
}

func NewConfig() *Config {
	return &Config{
		TargetHost:     "127.0.0.1",
		TargetPort:     "7777",
		RequestTimeout: 60 * time.Second,
		Sleep:          5 * time.Second,
		Jitter:         50.00,
		AgentUUID:      "",
		Endpoint:       "/",
	}
}
