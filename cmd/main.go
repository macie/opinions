package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/macie/opinions"
	"github.com/macie/opinions/security"
)

var (
	AppVersion        string
	DefaultAppVersion = time.Now().Format("2006.01.02-dev150405")
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("opinions: ")

	if err := security.Sandbox(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	err, config := parse(os.Args[1:])
	if err != nil {
		log.Printf("invalid usage: %s\n", err)
		os.Exit(1)
	}
	if config.ShowVersion {
		if AppVersion == "" {
			AppVersion = DefaultAppVersion
		}
		fmt.Fprintf(os.Stderr, "opinions %s\n", AppVersion)
		os.Exit(0)
	}

	ctx, cancel := appContext(config)
	defer cancel()

	services := []RemoteSearch{
		opinions.SearchHackerNews,
		opinions.SearchLobsters,
	}

	wg := new(sync.WaitGroup)
	for _, s := range services {
		wg.Add(1)
		go func(searchFn RemoteSearch) {
			defer wg.Done()
			PrintCommented(ctx, config.Query, searchFn)
		}(s)
	}
	wg.Wait()

	os.Exit(0)
}

// RemoteSearch represents function for searching on social news website.
type RemoteSearch func(context.Context, string) ([]opinions.Discussion, error)

// PrintCommented prints to standard output searching results with non-zero
// comments for given query.
func PrintCommented(ctx context.Context, query string, search RemoteSearch) {
	results, err := search(ctx, query)
	if err != nil {
		log.Println(err)
		return
	}

	for _, discussion := range results {
		if discussion.Comments == 0 {
			continue
		}
		fmt.Fprintf(os.Stdout, "%s\n", discussion)
	}
}