package cache

import (
	"github.com/ringtail/pansidong/types"
	"time"
	log "github.com/Sirupsen/logrus"
	"sync"
	"github.com/ringtail/pansidong/provider"
	"fmt"
	"github.com/ringtail/pansidong/backend"
)

type CacheManager struct {
	sync.Mutex
	busy            bool
	Config          *types.CacheConfig
	Cache           types.MemoryStore
	Backend         types.BackendStore
	ProviderManager *provider.ProxyProviderManager
}

func (cm *CacheManager) Loop(stopChan chan struct{}) {
	go cm.loop()
	<-stopChan
}

func (cm *CacheManager) loop() {
	cm.runOnce()
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
	// unlock and remove busy status finally
	defer func() {
		cm.busy = false
		cm.Unlock()
	}()
	cm.Lock()
	cm.busy = true
	ips := cm.ProviderManager.RunOnce()
	err := cm.Cache.Refresh(ips, &types.RefreshOptions{
		Force: false,
	})
	if err != nil {
		log.Errorf("Failed to fetch ips from provider,because of %s", err.Error())
		return
	}

	if cm.Backend != nil {
		err := cm.Backend.Insert(ips, &types.InsertOptions{
			Force: true,
		})
		if err != nil {
			log.Errorf("Failed to insert ips to backend,because of %s", err.Error())
			return
		}
	}
}

func (cm *CacheManager) ExpireIp(ip string) error {
	err := cm.Cache.Expire(ip)
	if err != nil {
		log.Warningf("Failed to expire key from cache,because of %s", err.Error())
		return err
	}

	if cm.Backend != nil {
		err := cm.Backend.Expire(ip)
		if err != nil {
			log.Warnf("Failed to expire key from backend,because of %s", err.Error())
		}
	}
	return nil
}

func (cm *CacheManager) FetchNextIps(options *types.ListOptions) ([]*types.ProxyIP, error) {
	return cm.Cache.Next(options)
}

func NewCacheManager(config *types.CacheConfig) (*CacheManager, error) {
	if err := config.Valid(); err != nil {
		return nil, fmt.Errorf("Failed to create cache manager because of %s", err.Error())
	}

	cm := &CacheManager{
		busy:            false,
		Config:          config,
		ProviderManager: provider.NewProxyProviderManager(),
	}

	cm.Cache = NewCache(config.Memory)

	if config.Backend != nil {
		log.Infof("backend is enabled,current backend is %s", config.Backend.Name())
		cm.Backend = backend.BackendFactory(config.Backend)
		size := config.Memory.Size
		ips, err := cm.Backend.List(&types.ListOptions{
			Limit: size,
		})
		if err != nil {
			return nil, fmt.Errorf("Failed to load cache data from backend,because of %s", err.Error())
		}

		err = cm.Cache.Refresh(ips, &types.RefreshOptions{
			Force: true,
		})

		if err != nil {
			return nil, fmt.Errorf("Failed to refresh ip from backend,because of %s", err.Error())
		}
		log.Infof("Success refresh from backend %s with %d ips", config.Backend.Name(), len(ips))
	}
	return cm, nil
}
