package web

import (
	"fmt"
	"net"
	"net/http"
)

// CIDRProtectMiddleware is a middleware that limits the inner handlers to a list of allowed CIDR ranges.
type CIDRProtectMiddleware struct {
	handler http.Handler

	allowedNetworks []*net.IPNet
}

// NewCIDRProtectMiddleware creates a new middleware that only will allow access the provided CIDR ranges
// otherwise it will return 401 Unauthorized status.
// The source IP address is extracted from the requests headers or, if not found from the request itself.
// Refer to the function GetRemoteAddress(r *http.Request) for more information on the order of headers.
// The middleware will panic at initialization time if one of the allowedNetwork parameter cannot be parsed
// as a CIDR range.
// Also, this middleware always allows loopback address to reach inner handler.
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
	ip, err := GetRemoteAddr(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	isIPAllowed := ip.IsLoopback()
	for _, network := range cpm.allowedNetworks {
		if isIPAllowed {
			break
		}
		isIPAllowed = network.Contains(ip)
	}

	if isIPAllowed {
		cpm.handler.ServeHTTP(w, r)
	}
	http.Error(w, "endpoint is protected", http.StatusUnauthorized)
}
