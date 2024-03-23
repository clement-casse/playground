package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"path"
	"time"

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
		blocky := exec.CommandContext(
			ctx,
			"/app/blocky",
			"--config=./config.yml",
		)
		blocky.Stdout = os.Stdout
		blocky.Stderr = os.Stderr
		return blocky.Run()
	})

	g.Go(func() error {
		ticker := time.NewTicker(10 * time.Minute)
		for {
			<-ticker.C
			blockyRefresh := exec.CommandContext(ctx, "/app/blocky", "refresh")
			if err := blockyRefresh.Run(); err != nil {
				return err
			}
		}
	})

	if err := g.Wait(); err != nil {
		log.Fatalf("One of the programs returned an error: %s", err)
	}
}
