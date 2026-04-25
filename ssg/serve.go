package ssg

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// liveReloadPath is the SSE endpoint browsers connect to when live reload
// is enabled. Hosted under an underscored prefix to avoid colliding with
// user routes.
const liveReloadPath = "/_ssg/reload"

// liveReloadScript is injected into HTML responses when live reload is
// enabled. It opens a long-lived SSE stream and reloads the page on either
// an explicit reload event or when the connection re-establishes after a
// disconnect (e.g. after the server restarts via watchexec).
const liveReloadScript = `<script>(function(){
  if (window.__ssg_live) return; window.__ssg_live = true;
  var hadConnection = false;
  function connect(){
    var es = new EventSource(` + "\"" + liveReloadPath + "\"" + `);
    es.addEventListener("connected", function(){
      if (hadConnection) { location.reload(); return; }
      hadConnection = true;
    });
    es.addEventListener("reload", function(){ location.reload(); });
    es.onerror = function(){ es.close(); setTimeout(connect, 500); };
  }
  connect();
})();</script>`

// EnableLiveReload turns on the live reload SSE endpoint and HTML script
// injection used by Serve. It has no effect on Build.
func (s *Site) EnableLiveReload() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.liveReload = true
}

// Reload notifies every connected live-reload client to reload the page.
// Useful for callers that watch their own data sources and want to push
// updates without restarting the server.
func (s *Site) Reload() {
	s.reloadMu.Lock()
	defer s.reloadMu.Unlock()
	for ch := range s.reloadClients {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

// Serve runs an HTTP server on addr that serves the registered routes.
// It blocks until the server returns an error.
func (s *Site) Serve(addr string) error {
	server := &http.Server{Addr: addr, Handler: s.Handler()}
	log.Printf("ssg: serving on http://%s", addr)
	return server.ListenAndServe()
}

// Handler returns the http.Handler used by Serve. It is exported so the
// site can be embedded into an existing HTTP stack or tested directly.
func (s *Site) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == liveReloadPath {
			s.handleReload(w, r)
			return
		}

		if route, ok := s.lookup(r.URL.Path); ok {
			s.serveRoute(w, r, route)
			return
		}

		if mount, sub, ok := s.matchStatic(r.URL.Path); ok {
			s.serveStatic(w, r, mount, sub)
			return
		}

		s.serveNotFound(w, r)
	})
}

func (s *Site) serveRoute(w http.ResponseWriter, _ *http.Request, route Route) {
	body, err := renderToBytes(route.Renderer)
	if err != nil {
		log.Printf("ssg: render %s: %v", route.Path, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	contentType := route.Renderer.ContentType()
	if s.liveReload && strings.HasPrefix(contentType, "text/html") {
		body = injectReloadScript(body)
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	_, _ = w.Write(body)
}

func (s *Site) serveStatic(w http.ResponseWriter, r *http.Request, mount staticMount, sub string) {
	name := sub
	if name == "" {
		name = "index.html"
	}

	// Resolve directory requests to their index.html if present.
	if info, err := fs.Stat(mount.FS, name); err == nil && info.IsDir() {
		name = strings.TrimSuffix(name, "/") + "/index.html"
	}

	f, err := mount.FS.Open(name)
	if err != nil {
		s.serveNotFound(w, r)
		return
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil || info.IsDir() {
		s.serveNotFound(w, r)
		return
	}

	if ct := guessContentType(name); ct != "" {
		w.Header().Set("Content-Type", ct)
	}
	w.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10))
	if seeker, ok := f.(io.ReadSeeker); ok {
		http.ServeContent(w, r, name, info.ModTime(), seeker)
		return
	}
	_, _ = io.Copy(w, f)
}

func (s *Site) serveNotFound(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	notFound := s.notFound
	s.mu.RUnlock()
	if notFound == nil {
		http.NotFound(w, r)
		return
	}
	body, err := renderToBytes(notFound)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	contentType := notFound.ContentType()
	if s.liveReload && strings.HasPrefix(contentType, "text/html") {
		body = injectReloadScript(body)
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write(body)
}

func (s *Site) handleReload(w http.ResponseWriter, r *http.Request) {
	if !s.liveReload {
		http.NotFound(w, r)
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := make(chan struct{}, 4)
	s.reloadMu.Lock()
	if s.reloadClients == nil {
		s.reloadClients = map[chan struct{}]struct{}{}
	}
	s.reloadClients[ch] = struct{}{}
	s.reloadMu.Unlock()
	defer func() {
		s.reloadMu.Lock()
		delete(s.reloadClients, ch)
		s.reloadMu.Unlock()
	}()

	fmt.Fprint(w, "event: connected\ndata: ok\n\n")
	flusher.Flush()

	ping := time.NewTicker(30 * time.Second)
	defer ping.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ch:
			fmt.Fprint(w, "event: reload\ndata: 1\n\n")
			flusher.Flush()
		case <-ping.C:
			fmt.Fprint(w, ": ping\n\n")
			flusher.Flush()
		}
	}
}

// injectReloadScript inserts the live reload script just before </body>, or
// at the end of the body if no closing tag is present.
func injectReloadScript(body []byte) []byte {
	script := []byte(liveReloadScript)
	idx := bytes.LastIndex(body, []byte("</body>"))
	if idx < 0 {
		out := make([]byte, 0, len(body)+len(script))
		out = append(out, body...)
		out = append(out, script...)
		return out
	}
	out := make([]byte, 0, len(body)+len(script))
	out = append(out, body[:idx]...)
	out = append(out, script...)
	out = append(out, body[idx:]...)
	return out
}
