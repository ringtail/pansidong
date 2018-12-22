package cache

import (
	"github.com/ringtail/pansidong/types"
	"time"
	"github.com/ringtail/pansidong/provider"
	log "github.com/Sirupsen/logrus"
	"sync"
)

const (
	NotFound = "NotFound"
)

type ListOptions struct {
	Limit int
}

type RefreshOptions struct {
	Force bool
}

type Store interface {
	Next(options *ListOptions) ([]*types.ProxyIP, error)
	List(options *ListOptions) ([]*types.ProxyIP, error)
	Get(key string) (*types.ProxyIP, error)
	Expire(key string) error
	Refresh(ips []*types.ProxyIP, options *RefreshOptions) error
}

type CacheConfig struct {
	Interval time.Duration
}

type CacheManager struct {
	sync.Mutex
	busy            bool
	Config          CacheConfig
	Store           Store
	ProviderManager *provider.ProxyProviderManager
}

func (cm *CacheManager) Loop(stopChan chan struct{}) {
	go cm.loop()
	<-stopChan
}

func (cm *CacheManager) loop() {
	ticker := time.NewTicker(cm.Config.Interval)
	for {
		select {
		case <-ticker.C:
			cm.runOnce()
		}
	}
}

func (cm *CacheManager) runOnce() {
	if cm.busy == true {
		return
	}
	defer func() {
		cm.busy = false
	}()
	cm.Lock()
	cm.busy = true
	ips := cm.ProviderManager.RunOnce()
	err := cm.Store.Refresh(ips, &RefreshOptions{
		Force: false,
	})
	if err != nil {
		log.Errorf("Failed to fetch ips from provider,because of %s", err.Error())
		return
	}
	cm.Unlock()
}

func (cm *CacheManager) ExpireIp(ip string) error {
	return nil
}

func (cm *CacheManager) FetchNextIps(options ListOptions) ([]*types.ProxyIP, error) {
	return nil, nil
}
