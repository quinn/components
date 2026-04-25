// Command dev runs the site with the live-reload dev server. Combined with
// watchexec (see the justfile) it rebuilds and reloads the browser on every
// source change.
package main

import (
	"flag"
	"log"
	"os"

	"go.quinn.io/components/internal/site"
)

func main() {
	addr := flag.String("addr", defaultAddr(), "address the dev server listens on")
	flag.Parse()

	s := site.New()
	s.EnableLiveReload()

	log.Printf("dev: http://localhost%s", *addr)
	if err := s.Serve(*addr); err != nil {
		log.Fatal(err)
	}
}

// defaultAddr resolves the listen address from ADDR or PORT, falling back
// to :8080 so the binary is useful without any environment set.
func defaultAddr() string {
	if addr := os.Getenv("ADDR"); addr != "" {
		return addr
	}
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	return ":8080"
}
