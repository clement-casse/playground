//go:build linux

// Command that bootstraps the github.com/0xERR0R/blocky dns server by also starting
// tailscale and authenticating to it with the auth_key provided in the environment
// variable TAILSCALE_AUTHKEY.
package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"path"
	"time"

	blockycmd "github.com/0xERR0R/blocky/cmd"
	"golang.org/x/sync/errgroup"
	tailscalecli "tailscale.com/cmd/tailscale/cli"
)

var (
	tailscaleSocketDir = "/var/run/tailscale/"
	tailscaleStateDir  = "/var/lib/tailscale/"
	tailscaleCacheDir  = "/var/cache/tailscale/"

	tailscaleStateFile  = path.Join(tailscaleStateDir, "tailscaled.state")
	tailscaleSocketFile = path.Join(tailscaleSocketDir, "tailscaled.sock")

	tsAuthKey = os.Getenv("TAILSCALE_AUTHKEY")
)

func main() {
	for _, dir := range []string{tailscaleSocketDir, tailscaleStateDir, tailscaleCacheDir} {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			log.Fatalf("cannot create tailscale directories, error: %s", err)
		}
	}

	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		tailscaled := exec.CommandContext(
			ctx,
			"/app/tailscaled",
			"--state="+tailscaleStateFile,
			"--socket="+tailscaleSocketFile,
		)
		tailscaled.Stdout = os.Stdout
		tailscaled.Stderr = os.Stderr
		return tailscaled.Run()
	})

	g.Go(func() error {
		time.Sleep(5 * time.Second)
		// include tailscale binary as a part of the go program
		return tailscalecli.Run([]string{"up",
			"--accept-routes",
			"--accept-dns",
			"--authkey=" + tsAuthKey,
			"--advertise-tags=tag:server",
		})
	})

	g.Go(func() error {
		// include blocky as a part of the program
		blocky := blockycmd.NewRootCommand()
		return blocky.Execute()
	})

	if err := g.Wait(); err != nil {
		log.Fatalf("One of the programs returned an error: %s", err)
	}
}
