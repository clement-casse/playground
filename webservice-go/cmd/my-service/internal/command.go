package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

var (
	version = "develop"
	date    = "unknown"

	defaultLogLevel = slog.LevelInfo.String()
)

const (
	jsonLoggerFlag = "json"
	listenAddrFlag = "listen"
	logLevelFlag   = "log-level"

	defaultListeningAddress = ""
	defaultListeningPort    = "8080"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		SilenceUsage:  true,
		SilenceErrors: true,
		Use:           "my-service",
		Args:          cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) (err error) {
			return runMain(cmd.Flags())
		},
	}

	cmd.Flags().Bool(jsonLoggerFlag, false, `Whether or not log in the JSON format`)
	cmd.Flags().String(listenAddrFlag, fmt.Sprintf("%s:%s", defaultListeningAddress, defaultListeningPort), `Address and port used by the application to serve HTTP requests`)
	cmd.Flags().String(logLevelFlag, defaultLogLevel, fmt.Sprintf(`How verbose is my-service (either: %s, %s, %s, %s)`, slog.LevelDebug.String(), slog.LevelInfo.String(), slog.LevelWarn.String(), slog.LevelError.String()))

	cmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Version of my-service",
		Long:  "Prints the version of the my-service tool and exits",
		RunE: func(cmd *cobra.Command, _ []string) (err error) {
			if version != "develop" {
				cmd.Println(fmt.Sprintf("%s version %s (%s)", cmd.Parent().Name(), version, date))
			}
			if info, ok := debug.ReadBuildInfo(); ok {
				cmd.Println(fmt.Sprintf("%s version %s (%s)", cmd.Parent().Name(), info.Main.Version, date))
			}
			err = fmt.Errorf("cannot read binary version")
			return
		},
	})

	return cmd
}

func runMain(flags *flag.FlagSet) error {
	logLevelStr, err := flags.GetString(logLevelFlag)
	if err != nil {
		return err
	}
	var logLevel slog.Level
	switch logLevelStr {
	case slog.LevelDebug.String():
		logLevel = slog.LevelDebug
	case slog.LevelInfo.String():
		logLevel = slog.LevelInfo
	case slog.LevelWarn.String():
		logLevel = slog.LevelWarn
	case slog.LevelError.String():
		logLevel = slog.LevelError
	default:
		return fmt.Errorf("Unknown Log level %s", logLevelStr)
	}

	jsonFlag, err := flags.GetBool(jsonLoggerFlag)
	if err != nil {
		return err
	}
	var logHandler slog.Handler
	if jsonFlag {
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	}
	logger := slog.New(logHandler)

	logger.Info("Starting My Service", "version", version)

	listenAddr, err := flags.GetString(listenAddrFlag)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	server := &http.Server{Addr: listenAddr} //nolint:gosec // Ignoring G114: Use of net/http serve function that has no support for setting timeouts.

	srvErr := make(chan error, 1)
	go func() {
		srvErr <- server.ListenAndServe()
	}()

	select {
	case err := <-srvErr:
		return err
	case <-ctx.Done():
		logger.WarnContext(ctx, "caught ^C signal, terminating my-service")
		stop()
	}

	return server.Shutdown(ctx)
}
