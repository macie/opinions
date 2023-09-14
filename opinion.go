package main

import (
	"fmt"
)

// Discussion is a representation of discussion inside social media service.
type Discussion struct {
	Service  string
	URL      string
	Title    string
	Source   string
	Comments int
}

// String returns string representation of discussion metadata.
func (d *Discussion) String() string {
	return fmt.Sprintf("%s\t%s\t%s\t%s", d.Service, d.URL, d.Title, d.Source)
}
