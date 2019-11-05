package proxy

import ()

// GetProxyIPList get proxy ips from api
func GetProxyIPList() []string {
	return proxies
}

var proxies = []string{
	"http://124.158.161.58:8080",
}
