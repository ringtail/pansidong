package types

import (
	"time"
	"github.com/ringtail/pansidong/backend"
)

const (
	defaultPoolSize = 20
	defaultInterval = 5 * 60 * time.Second
)

type GlobalConfig struct {
	Port int
	Host string
}

type Config struct {
	GlobalConfig *GlobalConfig
	Backend      *BackendConfigUnknown
	Memory       *MemoryConfig
}

type CacheConfig struct {
	Interval time.Duration
	Memory   *MemoryConfig
	Backend  BackendConfig
}

func (cc *CacheConfig) Valid() error {
	return nil
}

func (c *Config) Valid() error {
	return nil
}

func (c *Config) Config() *CacheConfig {
	cc := &CacheConfig{}
	if c.Backend != nil {
		backendType := c.Backend.Type
		cc.Backend = backend.CreateBackendConfig(backendType, c.Backend.Config)
	}
	if c.Memory != nil {
		cc.Memory = c.Memory
	} else {
		cc.Memory = &MemoryConfig{
			Size: defaultPoolSize,
		}
	}
	cc.Interval = defaultInterval
	return cc
}
