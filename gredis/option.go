package gredis

import (
	redis "github.com/redis/go-redis/v9"
)

// 分片
type RingOptions struct {
	redis.RingOptions
}

// 集群
type ClusterOptions struct {
	redis.ClusterOptions
}

// 哨兵
type FailoverOptions struct {
	redis.FailoverOptions
}

type SingleOptions struct {
	redis.Options
}

type Options struct {
	configPath string
	Mode       string          `json:"mode" yaml:"mode"` // single, cluster, ring, sentinel
	Cluster    ClusterOptions  `json:"cluster" yaml:"cluster"`
	Ring       RingOptions     `json:"ring" yaml:"ring"`
	Sentinel   FailoverOptions `json:"sentinel" yaml:"sentinel"`
	Single     SingleOptions   `json:"Single" yaml:"Single"`
}

const (
	DefaultMaxRetries   = 5
	DefaultPoolSize     = 100
	DefaultDialTimeout  = 500
	DefaultReadTimeout  = 100
	DefaultWriteTimeout = 100
	DefaultPoolTimeout  = 1000
	DefaultIdleTimeout  = 60000
)

type Option func(*Options)

// WithConfigPath sets application path.
func WithConfigPath(configPath string) Option {
	return func(opts *Options) {
		opts.configPath = configPath
	}
}
