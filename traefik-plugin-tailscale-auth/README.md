# Traefik plugin: Tailscale Authentication Middleware

This plugin is intended to be run on Traefik instances that serve HTTP requests over a *Tailscale Network* (tailnet).
It is inspired by the [`nginx-auth` binary distributed by Tailscale][1]: 

## References

1. [tailscale binary `nginx-auth`][1]

[1]: https://pkg.go.dev/tailscale.com/cmd/nginx-auth