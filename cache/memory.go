package cache

import (
	"github.com/ringtail/pansidong/types"
	"github.com/virtual-kubelet/virtual-kubelet/providers/aliyun/ingress/errors"
)



type MemoryStore struct {
	cursor int
	ips    []*types.ProxyIP
}

func (ms *MemoryStore) Next(options *ListOptions) ([]*types.ProxyIP, error) {
	size := len(ms.ips)
	if size == 0 {
		return nil, errors.New(NotFound)
	}
	ips := make([]*types.ProxyIP, 0)
	if options != nil && options.Limit != 0 {
		for i := 0; i < options.Limit; i ++ {
			r_i := (i + ms.cursor) % size
			ips = append(ips, ms.ips[r_i])
			ms.cursor = ms.cursor + 1
		}
	}
	return ips, nil
}

func (ms *MemoryStore) List(options *ListOptions) ([]*types.ProxyIP, error) {
	return ms.ips, nil
}

func (ms *MemoryStore) Get(key string) (*types.ProxyIP, error) {
	for _, ip := range ms.ips {
		if key == ip.IP {
			return ip, nil
		}
	}
	return nil, errors.New(NotFound)
}

func (ms *MemoryStore) Expire(key string) error {
	for i, ip := range ms.ips {
		if key == ip.IP {
			ms.ips = append(ms.ips[:i], ms.ips[i+1:]...)
		}
	}
	return nil
}

func (ms *MemoryStore) Refresh(ips []*types.ProxyIP, options *RefreshOptions) error {
	if options != nil && options.Force == true {
		ms.ips = make([]*types.ProxyIP, len(ips))
		ms.ips = ips
	} else {
		for _, ip := range ips {
			exists := false
			for _, oIp := range ms.ips {
				if oIp.IP == ip.IP {
					exists = true
				}
			}
			if exists == true {
				continue
			}
			ms.ips = append(ms.ips, ip)
		}
	}
	return nil
}

func NewMemoryStore() Store {
	return &MemoryStore{}
}
