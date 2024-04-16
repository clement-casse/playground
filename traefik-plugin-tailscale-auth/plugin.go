package plugintsauth

import (
	"context"
	"net"
	"net/http"
	"net/netip"
	"os"
	"strings"

	"tailscale.com/client/tailscale"
)

// Config the plugin configuration.
type Config struct {
	Socket string `yaml:"socket"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// TailscaleAuth does some stuff.
type TailscaleAuth struct {
	name string
	next http.Handler

	localClient *tailscale.LocalClient
}

// New creates the Traefik plugin.
func New(_ context.Context, next http.Handler, _ *Config, name string) (http.Handler, error) {
	return &TailscaleAuth{
		name: name,
		next: next,

		localClient: &tailscale.LocalClient{UseSocketOnly: true},
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

	info, err := ts.localClient.WhoIs(req.Context(), remoteAddr.String())
	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		os.Stderr.WriteString("can't look up " + remoteAddr.String() + " : " + err.Error())
		return
	}

	h := rw.Header()
	h.Set("Tailscale-Login", strings.Split(info.UserProfile.LoginName, "@")[0])
	h.Set("Tailscale-User", info.UserProfile.LoginName)
	h.Set("Tailscale-Name", info.UserProfile.DisplayName)
	h.Set("Tailscale-Profile-Picture", info.UserProfile.ProfilePicURL)

	// h.Set("Tailscale-Tailnet", tailnet)

	ts.next.ServeHTTP(rw, req)
}
