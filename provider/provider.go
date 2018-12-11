package provider

import (
	"github.com/ringtail/pansidong/types"
	log "github.com/Sirupsen/logrus"
	"time"
)

type Provider interface {
	Name() (name string)
	CheckHealth() (healthy bool, nextTickTime time.Time)
	GetProxyList() (ips []*types.ProxyIP, err error)
}

type ProxyProviderManager struct {
	providers map[string]Provider
}

func (pm *ProxyProviderManager) RunOnce() ([]*types.ProxyIP) {
	now := time.Now()
	ips := make([]*types.ProxyIP, 0)
	for _, p := range pm.providers {
		healthy, nextTickTime := p.CheckHealth()
		if healthy == false && now.Before(nextTickTime) {
			continue
		}
		// run once provider
		p_ips, err := p.GetProxyList()
		if err != nil {
			log.Warningf("Failed to get GetProxyList from %s,because of %s", p.Name(), err.Error())
			continue
		}
		ips = append(ips, p_ips...)
	}
	return ips
}
