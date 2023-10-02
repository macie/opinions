package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

var AppVersion = time.Now().Format("2006.01.02-dev150405")

func main() {
	log.SetFlags(0)
	log.SetPrefix("opinions: ")

	err, config := parse(os.Args[1:])
	if err != nil {
		log.Println("invalid usage")
		os.Exit(1)
	}
	if config.ShowVersion {
		fmt.Fprintf(os.Stderr, "opinions %s\n", AppVersion)
		os.Exit(0)
	}

	_, cancel := appContext(config)
	defer cancel()

	fmt.Fprintln(os.Stderr, "ERROR: not implemented")
	os.Exit(1)
}
