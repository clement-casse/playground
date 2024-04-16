package plugintsauth

import (
	"context"
	"net"
	"net/http"
	"net/netip"

	"tailscale.com/client/tailscale"
)

// Config the plugin configuration.
type Config struct{}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// TailscaleAuth does some stuff.
type TailscaleAuth struct {
	name string
	next http.Handler
}

// New creates the Traefik plugin.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &TailscaleAuth{
		name: name,
		next: next,
	}, nil
}

func (ts *TailscaleAuth) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	remoteHost, remotePort := req.Header.Get("Remote-Addr"), req.Header.Get("Remote-Port")
	if remoteHost == "" || remotePort == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	remoteAddrStr := net.JoinHostPort(remoteHost, remotePort)
	remoteAddr, err := netip.ParseAddrPort(remoteAddrStr)
	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	tailscale.WhoIs(req.Context(), remoteAddr.String())

	ts.next.ServeHTTP(rw, req)
}
