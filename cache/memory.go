package cache

import (
	"github.com/ringtail/pansidong/types"
	"github.com/virtual-kubelet/virtual-kubelet/providers/aliyun/ingress/errors"
)

const (
	NotFound = "NotFound"
)

type Cache struct {
	cursor int
	ips    []*types.ProxyIP
}

func (c *Cache) Next(options *types.ListOptions) ([]*types.ProxyIP, error) {
	size := len(c.ips)
	if size == 0 {
		return nil, errors.New(NotFound)
	}
	ips := make([]*types.ProxyIP, 0)
	if options != nil && options.Limit != 0 {
		for i := 0; i < options.Limit; i ++ {
			r_i := (i + c.cursor) % size
			ips = append(ips, c.ips[r_i])
			c.cursor = c.cursor + 1
		}
	}
	return ips, nil
}

func (c *Cache) List(options *types.ListOptions) ([]*types.ProxyIP, error) {
	return c.ips, nil
}

func (c *Cache) Get(key string) (*types.ProxyIP, error) {
	for _, ip := range c.ips {
		if key == ip.IP {
			return ip, nil
		}
	}
	return nil, errors.New(NotFound)
}

func (c *Cache) Expire(key string) error {
	for i, ip := range c.ips {
		if key == ip.IP {
			c.ips = append(c.ips[:i], c.ips[i+1:]...)
		}
	}
	return nil
}

func (c *Cache) Refresh(ips []*types.ProxyIP, options *types.RefreshOptions) error {
	if options != nil && options.Force == true {
		c.ips = make([]*types.ProxyIP, len(ips))
		c.ips = ips
	} else {
		for _, ip := range ips {
			exists := false
			for _, oIp := range c.ips {
				if oIp.IP == ip.IP {
					exists = true
				}
			}
			if exists == true {
				continue
			}
			c.ips = append(c.ips, ip)
		}
	}
	return nil
}

func NewCache(config *types.MemoryConfig) types.MemoryStore {
	return &Cache{}
}
