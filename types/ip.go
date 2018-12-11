package types

// ProxyIP
type ProxyIP struct {
	IP            string `json:"ip"`
	Port          string `json:"port"`
	Schema        string `json:"schema"`
	CreateTime    int    `json:"create_time"`
	ExpiredTime   int    `json:"_"`
	LastCheckTime int    `json:"_"`
	Priority      int    `json:"priority"`
	Refer         string `json:"refer"`
}
