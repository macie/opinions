package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/macie/opinions/security"
)

// AppConfig represets current app configuration.
type AppConfig struct {
	appVersion  string
	Query       string
	Timeout     time.Duration
	ShowVersion bool
}

// NewAppConfig combines command line arguments and app version into AppConfig.
func NewAppConfig(cliArgs []string, appVersion string) (AppConfig, error) {
	config := AppConfig{
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
		return AppConfig{}, fmt.Errorf("expected exactly 1 query but get %d: '%s'", len(f.Args()), strings.Join(f.Args(), "', '"))
	}
	config.Query = f.Args()[0]

	return config, nil
}

// Version returns string with full version description.
func (c *AppConfig) Version() string {
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

// NewAppContext creates cancellable app context with optional timeout.
func NewAppContext(config AppConfig) (context.Context, context.CancelFunc) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	if config.Timeout != 0 {
		return context.WithTimeout(ctx, config.Timeout)
	}

	return ctx, cancel
}
