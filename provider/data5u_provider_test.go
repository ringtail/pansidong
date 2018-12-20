package provider

import (
	"testing"
)

func TestData5uGetProxyList(t *testing.T) {
	ips, err := Data5uProviderSingleton.GetProxyList()
	if err != nil {
		t.Errorf("Failed to getProxyIp,because of %s", err.Error())
		return
	}
	for _, ip := range ips {
		t.Logf("successfully getProxyIp from %s : %s://%s:%s", ip.Refer, ip.Schema, ip.IP, ip.Port)
	}
}
