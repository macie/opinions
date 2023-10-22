package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"
)

// AppConfig represets current app configuration.
type AppConfig struct {
	Query       string
	Timeout     time.Duration
	Version     string
	ShowVersion bool
}

// NewAppConfig combines command line arguments and app version into AppConfig.
func NewAppConfig(cliArgs []string, appVersion string) (AppConfig, error) {
	if appVersion == "" {
		appVersion = time.Now().Format("2006.01.02-dev150405")
	}
	config := AppConfig{
		Version: appVersion,
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

// appContext creates context using AppConfig.
func appContext(config AppConfig) (context.Context, context.CancelFunc) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	if config.Timeout != 0 {
		return context.WithTimeout(ctx, config.Timeout)
	}

	return ctx, cancel
}
