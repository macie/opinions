package main

import (
	"context"
	"errors"
	"fmt"
	"io"
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
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					// context errors are handled near the end of the program
					return
				}
				log.Printf("cannot search: %s\n", err)
				return
			}

			if _, err := FprintlnCommented(os.Stdout, discussions...); err != nil {
				log.Printf("cannot print results to stdout: %s\n", err)
			}
		}(s)
	}
	wg.Wait()

	switch ctx.Err() {
	case nil:
		// no error
	case context.Canceled:
		log.Println("searching was cancelled by user")
		os.Exit(1)
	case context.DeadlineExceeded:
		log.Println("searching needs more time than expected")
		os.Exit(1)
	default:
		log.Printf("searching was interrupted: %s\n", ctx.Err())
		os.Exit(1)
	}

	os.Exit(0)
}

// RemoteSearch represents function for searching on social news website.
type RemoteSearch func(context.Context, opinions.GetRequester, string) ([]opinions.Discussion, error)

// FprintlnCommented writes to w discussions with non-zero comments, each in
// a new line.
func FprintlnCommented(w io.Writer, discussions ...opinions.Discussion) (int, error) {
	buf := make([]byte, 0, 1024)
	for _, d := range discussions {
		if d.Comments == 0 {
			continue
		}
		buf = append(buf, d.String()...)
		buf = append(buf, '\n')
	}

	return w.Write(buf)
}
