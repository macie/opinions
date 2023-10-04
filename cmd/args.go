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

// AppConfig represets CLI parameters given by user.
type AppConfig struct {
	Timeout     time.Duration
	Query       string
	ShowVersion bool
}

// parse creates AppConfig from given command line arguments.
func parse(args []string) (error, AppConfig) {
	config := AppConfig{}
	f := flag.NewFlagSet("opinions", flag.ContinueOnError)

	f.DurationVar(&config.Timeout, "timeout", 0, "max running time. Valid time units: ns, us, ms, s, m, h")
	f.BoolVar(&config.ShowVersion, "version", false, "version")
	if err := f.Parse(args); err != nil {
		return err, config
	}

	if config.ShowVersion {
		return nil, config
	}

	if len(f.Args()) != 1 {
		return fmt.Errorf("expected exactly 1 query but get %d: '%s'", len(f.Args()), strings.Join(f.Args(), "', '")), AppConfig{}
	}
	config.Query = f.Args()[0]

	return nil, config
}

// appContext creates context using AppConfig.
func appContext(config AppConfig) (context.Context, context.CancelFunc) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	if config.Timeout != 0 {
		return context.WithTimeout(ctx, config.Timeout)
	}

	return ctx, cancel
}
