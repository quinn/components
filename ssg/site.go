package ssg

import (
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strings"
	"sync"
)

// Site is the registry of routes and static mounts that make up a site.
//
// A Site is safe for concurrent use. The same value can register routes from
// multiple goroutines and be served while routes are still being added,
// although in practice all routes are usually registered up front.
type Site struct {
	mu       sync.RWMutex
	routes   map[string]Route
	statics  []staticMount
	notFound Renderer

	liveReload    bool
	reloadMu      sync.Mutex
	reloadClients map[chan struct{}]struct{}
}

// Route is a single addressable URL path together with the Renderer that
// produces its bytes.
type Route struct {
	Path     string
	Renderer Renderer
}

// staticMount mounts a filesystem subtree at a URL prefix. The Prefix is
// always normalised to start with "/" and never end with "/".
type staticMount struct {
	Prefix string
	FS     fs.FS
}

// New returns an empty Site ready for routes to be added.
func New() *Site {
	return &Site{routes: map[string]Route{}}
}

// Add registers a single route. It is the primary registration primitive.
//
// The path must begin with "/". Trailing slashes (other than on the root
// path) are stripped so "/about" and "/about/" refer to the same route.
// Re-registering an existing path replaces the previous renderer.
func (s *Site) Add(path string, r Renderer) {
	s.AddRoute(Route{Path: path, Renderer: r})
}

// AddRoute is like Add but accepts a Route value, which is convenient when
// generating routes in bulk.
func (s *Site) AddRoute(r Route) {
	if r.Renderer == nil {
		panic(fmt.Sprintf("ssg: nil renderer for route %q", r.Path))
	}
	r.Path = normalizePath(r.Path)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.routes[r.Path] = r
}

// AddRoutes registers many routes in one call.
func (s *Site) AddRoutes(routes []Route) {
	for _, r := range routes {
		s.AddRoute(r)
	}
}

// Static mounts a directory on disk at the given URL prefix. Files inside
// the directory are served verbatim by the dev server and copied as-is when
// the site is built.
func (s *Site) Static(urlPrefix, fsDir string) {
	s.StaticFS(urlPrefix, os.DirFS(fsDir))
}

// StaticFS mounts an fs.FS at the given URL prefix. This is the more
// flexible variant of Static and works with embedded filesystems.
func (s *Site) StaticFS(urlPrefix string, fsys fs.FS) {
	prefix := normalizePath(urlPrefix)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.statics = append(s.statics, staticMount{Prefix: prefix, FS: fsys})
}

// SetNotFound installs a custom renderer used when a request does not match
// any registered route or static mount. It only affects Site.Serve.
func (s *Site) SetNotFound(r Renderer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.notFound = r
}

// Routes returns a stable, alphabetically sorted snapshot of the registered
// routes.
func (s *Site) Routes() []Route {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Route, 0, len(s.routes))
	for _, r := range s.routes {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out
}

// lookup finds a registered route by path. The lookup normalises the path
// so trailing slashes do not matter.
func (s *Site) lookup(path string) (Route, bool) {
	p := normalizePath(path)
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.routes[p]
	return r, ok
}

// staticsSnapshot returns a copy of the static mounts so callers can iterate
// without holding the lock.
func (s *Site) staticsSnapshot() []staticMount {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]staticMount, len(s.statics))
	copy(out, s.statics)
	return out
}

// matchStatic returns the longest static mount whose prefix matches the
// given URL path along with the path remainder relative to the mount.
func (s *Site) matchStatic(urlPath string) (staticMount, string, bool) {
	mounts := s.staticsSnapshot()
	var best staticMount
	var bestSub string
	matched := false
	for _, m := range mounts {
		var sub string
		var ok bool
		switch {
		case m.Prefix == "/":
			sub = strings.TrimPrefix(urlPath, "/")
			ok = true
		case urlPath == m.Prefix:
			sub = ""
			ok = true
		case strings.HasPrefix(urlPath, m.Prefix+"/"):
			sub = strings.TrimPrefix(urlPath, m.Prefix+"/")
			ok = true
		}
		if !ok {
			continue
		}
		if !matched || len(m.Prefix) > len(best.Prefix) {
			best = m
			bestSub = sub
			matched = true
		}
	}
	return best, bestSub, matched
}

// normalizePath ensures a path starts with "/" and has no trailing slash
// (other than the root).
func normalizePath(p string) string {
	if p == "" {
		return "/"
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	if p == "/" {
		return p
	}
	return strings.TrimRight(p, "/")
}
