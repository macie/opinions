package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/macie/opinions"
	"github.com/macie/opinions/internal/http"
	"github.com/macie/opinions/internal/security"
)

var AppVersion string // injected during build

func main() {
	log.SetFlags(0)
	log.SetPrefix("opinions: ")

	if err := security.Sandbox(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	config, err := NewAppConfig(os.Args[1:], AppVersion)
	if err != nil {
		log.Printf("invalid usage: %s\n", err)
		os.Exit(1)
	}
	if config.ShowVersion {
		fmt.Fprint(os.Stderr, config.Version())
		os.Exit(0)
	}

	ctx, cancel := NewAppContext(config)
	defer cancel()

	client := http.Client{AppVersion: AppVersion}
	services := []RemoteSearch{
		opinions.SearchHackerNews,
		opinions.SearchLemmy,
		opinions.SearchLobsters,
		opinions.SearchReddit,
	}

	wg := new(sync.WaitGroup)
	for _, s := range services {
		wg.Add(1)
		go func(searchFn RemoteSearch) {
			defer wg.Done()

			discussions, err := searchFn(ctx, client, config.Query)
			if err != nil {
				log.Println(err)
				return
			}

			PrintCommented(discussions)
		}(s)
	}
	wg.Wait()

	os.Exit(0)
}

// RemoteSearch represents function for searching on social news website.
type RemoteSearch func(context.Context, http.Client, string) ([]opinions.Discussion, error)

// PrintCommented prints to standard output searching results with non-zero
// comments for given collection.
func PrintCommented(discussions []opinions.Discussion) {
	for _, d := range discussions {
		if d.Comments == 0 {
			continue
		}
		fmt.Fprintf(os.Stdout, "%s\n", d)
	}
}
