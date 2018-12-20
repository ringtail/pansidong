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
	"github.com/virtual-kubelet/virtual-kubelet/providers/aliyun/ingress/errors"
	log "github.com/Sirupsen/logrus"
)

const (
	GetProxyListProviderName     = "GetProxyListProvider"
	GetProxyListProviderEndpoint = "https://api.getproxylist.com/proxy?protocol[]=http"
	DefaultMaxIps                = 5
)

var (
	GetProxyListProviderSingleton *GetProxyListProvider
)

type GetProxyListResponse struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Error    string `json:"error"`
}

type GetProxyListProvider struct {
	Endpoint string
}

func (gp *GetProxyListProvider) Name() (name string) {
	return GetProxyListProviderName
}

func (gp *GetProxyListProvider) GetProxyList() (ips []*types.ProxyIP, err error) {
	ips = make([]*types.ProxyIP, 0)
	for i := 0; i < DefaultMaxIps; i ++ {
		ip, err := gp.getProxyIp()
		if err != nil {
			log.Warnf("Failed to GetProxyList,because of %s", err.Error())
			break;
		}
		ips = append(ips, ip)
	}
	return
}

func (gp *GetProxyListProvider) getProxyIp() (proxyIp *types.ProxyIP, err error) {
	gr := &GetProxyListResponse{}
	resp, err := http.Get(gp.Endpoint)
	if err != nil {
		log.Warnf("Failed to getProxyIp from %s,because of %s", gp.Endpoint, err.Error())
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warnf("Failed to read body from resp,because of %s", err.Error())
		return nil, err
	}
	err = json.Unmarshal(body, gr)
	if err != nil {
		log.Warnf("Failed to unmarshal body from resp body,because of %s", err.Error())
		return nil, err
	}

	if gr.IP == "" {
		if gr.Error != "" {
			return nil, errors.New(gr.Error)
		}
		return nil, errors.New("EmptyProxyIp")
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
