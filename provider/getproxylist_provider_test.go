package provider

import (
	"testing"
)

func TestGetProxyListProviderGetSingleIp(t *testing.T) {
	ip, err := GetProxyListProviderSingleton.getProxyIp()
	if err != nil {
		t.Errorf("Failed to getProxyIp,because of %s", err.Error())
		return
	}
	t.Logf("successfully getProxyIp from %s : %s://%s:%s", ip.Refer, ip.Schema, ip.IP, ip.Port)
}

func TestGetProxyListProviderGetSingleIps(t *testing.T) {
	ips, err := GetProxyListProviderSingleton.GetProxyList()
	if err != nil {
		t.Errorf("Failed to getProxyIp,because of %s", err.Error())
		return
	}
	for _, ip := range ips {
		t.Logf("successfully getProxyIp from %s : %s://%s:%s", ip.Refer, ip.Schema, ip.IP, ip.Port)
	}

}
