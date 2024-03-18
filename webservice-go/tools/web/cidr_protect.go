package web

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

// CIDRProtectMiddleware is a middleware that limits the inner handlers to a list of allowed CIDR ranges.
type CIDRProtectMiddleware struct {
	handler http.Handler

	allowedNetworks []*net.IPNet
}

// NewCIDRProtectMiddleware creates a new middleware that only will allow access the provided CIDR ranges
// otherwise it will return 401 Unauthorized status.
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
	var ip net.IP
	// ref: https://husobee.github.io/golang/ip-address/2015/12/17/remote-ip-go.html
	for _, headerKey := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		addresses := strings.Split(r.Header.Get(headerKey), ",")
		// go from right to left until we get a public address to finf the address right before our proxy.
		for i := len(addresses) - 1; i >= 0; i-- {
			thisIP := net.ParseIP(strings.TrimSpace(addresses[i])) // header can contain spaces
			if thisIP.IsGlobalUnicast() && !thisIP.IsPrivate() {
				ip = thisIP
				break
			}
		}
	}
	if len(ip) == 0 { // means that it has not been initialized by headers (covers v4 & v6)
		strIP, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot parse host:port from '%s', error: %s", r.RemoteAddr, err), http.StatusInternalServerError)
		}
		ip = net.ParseIP(strIP)
		if ip == nil {
			http.Error(w, fmt.Sprintf("cannot parse IP address %s, error: %s", strIP, err), http.StatusInternalServerError)
		}
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
