package web

import (
	"fmt"
	"net"
	"net/http"
)

// CIDRProtectMiddleware
type CIDRProtectMiddleware struct {
	handler http.Handler

	allowedNetworks []*net.IPNet
}

// NewCIDRProtectMiddleware
func NewCIDRProtectMiddleware(allowedNetworks ...string) *CIDRProtectMiddleware {
	cpm := &CIDRProtectMiddleware{
		handler:         nil,
		allowedNetworks: make([]*net.IPNet, 0, len(allowedNetworks)),
	}
	for _, allowedNetwork := range allowedNetworks {
		_, network, err := net.ParseCIDR(allowedNetwork)
		if err != nil {
			panic(fmt.Sprintf("cannot parse network CIDR %s", allowedNetwork))
		}
		cpm.allowedNetworks = append(cpm.allowedNetworks, network)
	}
	return cpm
}

func (cpm *CIDRProtectMiddleware) Chain(handler http.Handler) http.Handler {
	cpm.handler = handler
	return cpm
}

func (cpm *CIDRProtectMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	strIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot parse host:port from '%s', error: %s", r.RemoteAddr, err), http.StatusInternalServerError)
	}
	ip := net.ParseIP(strIP)
	if ip == nil {
		http.Error(w, fmt.Sprintf("cannot parse IP address %s, error: %s", strIP, err), http.StatusInternalServerError)
	}
	isIPAllowed := ip.IsLoopback()
	for _, network := range cpm.allowedNetworks {
		if isIPAllowed {
			break
		}
		if network == nil {
			continue
		}
		isIPAllowed = network.Contains(ip)
	}

	if isIPAllowed {
		cpm.handler.ServeHTTP(w, r)
	}
	http.Error(w, "endpoint is protected", http.StatusUnauthorized)
}
