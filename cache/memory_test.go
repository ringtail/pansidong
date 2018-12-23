package cache

import (
	"testing"
	"github.com/ringtail/pansidong/types"
	"github.com/magiconair/properties/assert"
)

var (
	m     types.MemoryStore
	ips   []*types.ProxyIP
	n_ips []*types.ProxyIP
)

func init() {
	m = NewCache()
	ips = []*types.ProxyIP{
		&types.ProxyIP{
			IP:   "0.0.0.0",
			Port: "80",
		},
		&types.ProxyIP{
			IP:   "127.0.0.1",
			Port: "81",
		},
	}
	n_ips = []*types.ProxyIP{
		&types.ProxyIP{
			IP:   "192.168.0.1",
			Port: "80",
		},
	}
}

func TestRefreshMemory(t *testing.T) {
	err := m.Refresh(ips, nil)
	if err != nil {
		t.Error("Failed to refresh cache")
		return
	}
	l_ips, err := m.List(nil)
	if err != nil {
		t.Error("Failed to list ips")
		return
	}
	assert.Equal(t, len(l_ips), len(ips), "Failed to pass refresh keys")
	t.Log("pass refresh memory")
}

func TestNextMemory(t *testing.T) {
	for i := 0; i < len(ips)*2; i ++ {
		ips, err := m.Next(&types.ListOptions{
			Limit: 1,
		})
		if err != nil || len(ips) == 0 {
			t.Errorf("Failed to pass next ips,because of %s", err.Error())
			return
		}
		t.Log(ips[0].IP, ips[0].Port)
		assert.Equal(t, len(ips), 1, "Next ips size is invalid")
	}
	t.Log("pass test next memory")
}

func TestExpire(t *testing.T) {
	err := m.Expire("0.0.0.0")
	if err != nil {
		t.Errorf("Failed to expire ip because of %s", err.Error())
		return
	}
	ip, err := m.Get("0.0.0.0")
	if ip == nil || err != nil {
		t.Log("pass expire key")
		return
	}
	t.Error("Failed to expire key")
}

func TestRefreshAppend(t *testing.T) {
	err := m.Refresh(n_ips, &types.RefreshOptions{
		Force: false,
	})
	if err != nil {
		t.Errorf("Failed to append and refresh memory,because of %s", err.Error())
	}
	ips, err := m.List(nil)
	assert.Equal(t, len(ips), len(ips)+len(n_ips)-1, "Failed to pass refresh append")
	t.Log("pass refresh append")
}
