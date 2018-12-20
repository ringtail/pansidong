package provider

import (
	"github.com/ringtail/pansidong/types"
	"time"
	"github.com/gocolly/colly"
	log "github.com/Sirupsen/logrus"
	"strings"
)

const (
	Data5uProviderName = "Data5uProvider"
	Data5uEndpoint     = "http://www.data5u.com/free/index.html"
)

var Data5uProviderSingleton *Data5uProvider

type Data5uProvider struct {
	Endpoint string
}

func (d *Data5uProvider) Name() (name string) {
	return Data5uProviderName
}
func (d *Data5uProvider) CheckHealth() (healthy bool, nextTickTime time.Time) {
	return
}

func (d *Data5uProvider) GetProxyList() (ips []*types.ProxyIP, err error) {
	ips = make([]*types.ProxyIP, 0)
	c := colly.NewCollector()

	defer func() {
		if err := recover(); err != nil {
			log.Warningf("%s panic because of some unknown panic : %v", Data5uProviderName, err)
		}
	}()

	// Find and visit all links
	c.OnHTML("body > div.wlist > ul > li:nth-child(2) > ul.l2", func(e *colly.HTMLElement) {
		nodes := e.DOM.Children()

		ip := &types.ProxyIP{
			IP:         strings.Trim(nodes.Eq(0).Text(), "\b]"),
			Port:       nodes.Eq(1).Text(),
			Schema:     nodes.Eq(3).Text(),
			CreateTime: time.Now().Second(),
		}
		ips = append(ips, ip)
	})

	c.Visit(d.Endpoint)
	return ips, nil
}

func init() {
	Data5uProviderSingleton = &Data5uProvider{
		Endpoint: Data5uEndpoint,
	}
}
