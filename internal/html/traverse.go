package html

import (
	"io"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

// Parse returns the parse tree for the HTML from the given Reader.
func Parse(r io.Reader) (*html.Node, error) {
	return html.Parse(r)
}

// First returns first node satisfying query inside given node.
func First(n *html.Node, query string) *html.Node {
	sel, err := cascadia.Parse(query)
	if err != nil {
		return &html.Node{}
	}
	return cascadia.Query(n, sel)
}

// FindAll returns all nodes satisfying query inside given node.
func FindAll(n *html.Node, query string) []*html.Node {
	sel, err := cascadia.Parse(query)
	if err != nil {
		return []*html.Node{}
	}
	return cascadia.QueryAll(n, sel)
}

// Attr returns given attribute value (or empty string if not found).
func Attr(n *html.Node, name string) string {
	for _, a := range n.Attr {
		if a.Key == name {
			return a.Val
		}
	}
	return ""
}

// Text returns text content of node.
func Text(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	return n.FirstChild.Data
}
