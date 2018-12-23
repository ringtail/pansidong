package provider

import (
	"github.com/ringtail/pansidong/types"
	log "github.com/Sirupsen/logrus"
)

type Provider interface {
	Name() (name string)
	GetProxyList() (ips []*types.ProxyIP, err error)
}

type ProxyProviderManager struct {
	providers map[string]Provider
}

func (pm *ProxyProviderManager) RunOnce() ([]*types.ProxyIP) {
	ips := make([]*types.ProxyIP, 0)
	for _, p := range pm.providers {
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

// TODO add provider conf
func NewProxyProviderManager() *ProxyProviderManager {
	p := &ProxyProviderManager{}
	p.providers = make(map[string]Provider)
	p.providers["data5u"] = Data5uProviderSingleton
	p.providers["proxyList"] = GetProxyListProviderSingleton
	return p
}
