// Command build renders the site to a directory on disk.
package main

import (
	"flag"
	"log"

	"go.quinn.io/components/internal/site"
)

func main() {
	out := flag.String("o", "./dist", "output directory")
	flag.Parse()

	s := site.New()
	if err := s.Build(*out); err != nil {
		log.Fatal(err)
	}
	log.Printf("built site to %s", *out)
}
