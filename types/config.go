package types

import (
	"time"
	"encoding/json"
	"log"
)

const (
	DefaultPoolSize = 20
	defaultInterval = 5 * 60 * time.Second
)

type GlobalConfig struct {
	Port int
	Host string
}

type BackendConfig interface {
	Name() string
	Config() interface{}
}

type BoltDBConfig struct {
	Path string
}

func (bc *BoltDBConfig) Name() string {
	return "boltdb"
}

func (bc *BoltDBConfig) Config() interface{} {
	return bc
}

type BackendConfigUnknown struct {
	Type   string
	Config string
}

type MemoryConfig struct {
	Size int
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
		cc.Backend = CreateBackendConfig(backendType, c.Backend.Config)
	}
	if c.Memory != nil {
		cc.Memory = c.Memory
	} else {
		cc.Memory = &MemoryConfig{
			Size: DefaultPoolSize,
		}
	}
	cc.Interval = defaultInterval
	return cc
}

func CreateBackendConfig(backendType string, j_string string) BackendConfig {
	switch backendType {
	case "boltdb":
		bc := &BoltDBConfig{}
		if j_string == "" {
			return bc
		}
		err := json.Unmarshal([]byte(j_string), bc)
		if err != nil {
			log.Fatalf("Failed to create %s backend config,because of %s", backendType, err.Error())
		}
		return bc
	}
	return nil
}
