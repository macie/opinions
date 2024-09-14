package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/macie/opinions/internal/security"
)

// appConfig represents current app configuration.
type appConfig struct {
	appVersion  string
	Query       string
	Timeout     time.Duration
	ShowVersion bool
}

// newAppConfig combines command line arguments and app version into AppConfig.
func newAppConfig(cliArgs []string, appVersion string) (appConfig, error) {
	config := appConfig{
		appVersion: appVersion,
	}
	f := flag.NewFlagSet("opinions", flag.ContinueOnError)

	f.DurationVar(&config.Timeout, "timeout", 0, "max running time. Valid time units: ns, us, ms, s, m, h")
	f.BoolVar(&config.ShowVersion, "version", false, "version")
	if err := f.Parse(cliArgs); err != nil {
		return config, err
	}

	if config.ShowVersion {
		return config, nil
	}

	if len(f.Args()) != 1 {
		return appConfig{}, fmt.Errorf("expected exactly 1 query but get %d: '%s'", len(f.Args()), strings.Join(f.Args(), "', '"))
	}
	config.Query = f.Args()[0]

	return config, nil
}

// version returns string with full version description.
func (c *appConfig) version() string {
	ver := c.appVersion
	if ver == "" {
		ver = time.Now().Format("2006.01.02-dev150405")
	}
	build := ""
	if security.IsHardened {
		build = " (hardened)"
	}
	return fmt.Sprintf("opinions %s%s\n", ver, build)
}

// newAppContext creates cancellable app context with optional timeout.
func newAppContext(config appConfig) (context.Context, context.CancelFunc) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	if config.Timeout != 0 {
		return context.WithTimeout(ctx, config.Timeout)
	}

	return ctx, cancel
}
