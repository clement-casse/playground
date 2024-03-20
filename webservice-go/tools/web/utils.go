package web

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

// GetRemoteAddr guesses the sender IP address based on http header X-Forwarded-For & X-Real-Ip.
// If the headers are note found the function falls back to the IP address of the sender.
// Header X-Forwarded-For has a higher precedence over X-Real-Ip. If multiple IP addresses are
// in the X-Forwarded-For header, the function will return the last one which is not private.
func GetRemoteAddr(r *http.Request) (net.IP, error) {
	// ref: https://husobee.github.io/golang/ip-address/2015/12/17/remote-ip-go.html
	for _, headerKey := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		addresses := strings.Split(r.Header.Get(headerKey), ",")
		// go from right to left until we get a public address to find the address right before our proxy.
		for i := len(addresses) - 1; i >= 0; i-- {
			thisIP := net.ParseIP(strings.TrimSpace(addresses[i])) // header can contain spaces
			if thisIP.IsGlobalUnicast() && !thisIP.IsPrivate() {
				return thisIP, nil
			}
		}
	}
	var ip net.IP
	strIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil, err
	}
	ip = net.ParseIP(strIP)
	if ip == nil {
		return nil, fmt.Errorf("IP address format of '%s' is invalid. error: %w", strIP, err)
	}
	return ip, nil
}
