package provider

/*
	https://getproxylist.com/ is a very popular proxy ip provider
*/

import (
	"github.com/ringtail/pansidong/types"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"time"
	"strconv"
)

const (
	GetProxyListProviderName     = "GetProxyListProvider"
	GetProxyListProviderEndpoint = "https://api.getproxylist.com/proxy?protocol[]=http"
)

var (
	GetProxyListProviderSingleton *GetProxyListProvider
)

type GetProxyListResponse struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

type GetProxyListProvider struct {
	Endpoint string
}

func (gp *GetProxyListProvider) Name() (name string) {
	return GetProxyListProviderName
}

func (gp *GetProxyListProvider) CheckHealth() (healthy bool, nextTickTime int) {
	return false, 0
}

func (gp *GetProxyListProvider) GetProxyList() (ips []*types.ProxyIP, err error) {
	return
}

func (gp *GetProxyListProvider) getProxyIp() (proxyIp *types.ProxyIP, err error) {
	gr := &GetProxyListResponse{}
	resp, err := http.Get(gp.Endpoint)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, gr)
	if err != nil {
		return nil, err
	}

	proxyIp = &types.ProxyIP{}
	proxyIp.CreateTime = time.Now().Second()
	proxyIp.IP = gr.IP
	proxyIp.Port = strconv.Itoa(gr.Port)
	proxyIp.Refer = GetProxyListProviderName
	proxyIp.Schema = gr.Protocol
	return
}

func init() {
	GetProxyListProviderSingleton = &GetProxyListProvider{
		Endpoint: GetProxyListProviderEndpoint,
	}
}
