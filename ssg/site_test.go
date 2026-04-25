package ssg

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func TestNormalizePath(t *testing.T) {
	cases := map[string]string{
		"":          "/",
		"/":         "/",
		"about":     "/about",
		"/about":    "/about",
		"/about/":   "/about",
		"/a/b/":     "/a/b",
		"/feed.xml": "/feed.xml",
	}
	for in, want := range cases {
		if got := normalizePath(in); got != want {
			t.Errorf("normalizePath(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestAddAndLookup(t *testing.T) {
	s := New()
	s.Add("/", Text("home"))
	s.Add("about", Text("about"))    // no leading slash
	s.Add("/blog/post/", Text("p1")) // trailing slash

	if _, ok := s.lookup("/"); !ok {
		t.Fatal("missing root")
	}
	if _, ok := s.lookup("/about"); !ok {
		t.Fatal("missing /about")
	}
	if _, ok := s.lookup("/about/"); !ok {
		t.Fatal("trailing slash should match")
	}
	if _, ok := s.lookup("/blog/post"); !ok {
		t.Fatal("missing /blog/post")
	}
	if _, ok := s.lookup("/missing"); ok {
		t.Fatal("/missing should not match")
	}

	got := s.Routes()
	if len(got) != 3 {
		t.Fatalf("Routes len = %d, want 3", len(got))
	}
	wantOrder := []string{"/", "/about", "/blog/post"}
	for i, r := range got {
		if r.Path != wantOrder[i] {
			t.Errorf("Routes[%d] = %s, want %s", i, r.Path, wantOrder[i])
		}
	}
}

func TestAddNilRendererPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil renderer")
		}
	}()
	New().Add("/x", nil)
}

func TestRouteOutputPath(t *testing.T) {
	cases := map[string]string{
		"/":          "index.html",
		"/about":     "about/index.html",
		"/blog/post": "blog/post/index.html",
		"/feed.xml":  "feed.xml",
		"/sitemap.xml":          "sitemap.xml",
		"/static/styles.css":    "static/styles.css",
	}
	for in, want := range cases {
		if got := routeOutputPath(in); got != want {
			t.Errorf("routeOutputPath(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestBuildWritesRoutesAndStatics(t *testing.T) {
	s := New()
	s.Add("/", Text("home"))
	s.Add("/about", Text("about"))
	s.Add("/feed.xml", Bytes("application/rss+xml", []byte("<rss/>")))

	staticFS := fstest.MapFS{
		"styles.css":      &fstest.MapFile{Data: []byte("body{}")},
		"img/logo.svg":    &fstest.MapFile{Data: []byte("<svg/>")},
	}
	s.StaticFS("/assets", staticFS)

	dir := t.TempDir()
	if err := s.Build(dir); err != nil {
		t.Fatalf("Build: %v", err)
	}

	wantFiles := map[string]string{
		"index.html":           "home",
		"about/index.html":     "about",
		"feed.xml":             "<rss/>",
		"assets/styles.css":    "body{}",
		"assets/img/logo.svg":  "<svg/>",
	}
	for rel, want := range wantFiles {
		got, err := os.ReadFile(filepath.Join(dir, rel))
		if err != nil {
			t.Errorf("read %s: %v", rel, err)
			continue
		}
		if string(got) != want {
			t.Errorf("%s = %q, want %q", rel, got, want)
		}
	}
}

func TestServeRoute(t *testing.T) {
	s := New()
	s.Add("/", Text("hello"))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	s.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	if got := rr.Body.String(); got != "hello" {
		t.Errorf("body = %q, want %q", got, "hello")
	}
	if got := rr.Header().Get("Content-Type"); got != ContentTypeText {
		t.Errorf("content-type = %q, want %q", got, ContentTypeText)
	}
}

func TestServeStaticMount(t *testing.T) {
	s := New()
	staticFS := fstest.MapFS{
		"styles.css":       &fstest.MapFile{Data: []byte("body{}")},
		"index.html":       &fstest.MapFile{Data: []byte("<h1>root</h1>")},
		"sub/page.html":    &fstest.MapFile{Data: []byte("<h1>sub</h1>")},
	}
	s.StaticFS("/assets", staticFS)

	cases := map[string]string{
		"/assets/styles.css":    "body{}",
		"/assets":               "<h1>root</h1>",
		"/assets/":              "<h1>root</h1>",
		"/assets/sub/page.html": "<h1>sub</h1>",
	}
	for url, want := range cases {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, url, nil)
		s.Handler().ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("%s: status = %d, want 200", url, rr.Code)
			continue
		}
		if got := rr.Body.String(); got != want {
			t.Errorf("%s: body = %q, want %q", url, got, want)
		}
	}
}

func TestServeNotFoundDefault(t *testing.T) {
	s := New()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	s.Handler().ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rr.Code)
	}
}

func TestServeCustomNotFound(t *testing.T) {
	s := New()
	s.SetNotFound(Bytes(ContentTypeHTML, []byte("<p>nope</p>")))
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	s.Handler().ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "nope") {
		t.Errorf("body = %q, want to contain %q", rr.Body.String(), "nope")
	}
}

func TestLiveReloadInjection(t *testing.T) {
	s := New()
	s.Add("/", Bytes(ContentTypeHTML, []byte("<html><body><p>hi</p></body></html>")))

	// Without live reload no script is injected.
	rr := httptest.NewRecorder()
	s.Handler().ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
	if strings.Contains(rr.Body.String(), "EventSource") {
		t.Fatal("script injected when live reload disabled")
	}

	// With live reload the script should appear before </body>.
	s.EnableLiveReload()
	rr = httptest.NewRecorder()
	s.Handler().ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
	body := rr.Body.String()
	if !strings.Contains(body, "EventSource") {
		t.Fatal("script not injected")
	}
	if !strings.Contains(body, "<script>") {
		t.Fatal("script tag missing")
	}
	if idx := strings.Index(body, "<script>"); idx == -1 || idx > strings.Index(body, "</body>") {
		t.Fatalf("script not before </body>: %s", body)
	}
}

func TestRendererTypes(t *testing.T) {
	if got := Text("ok").ContentType(); got != ContentTypeText {
		t.Errorf("Text content type = %q", got)
	}
	if got := JSON([]byte(`{}`)).ContentType(); got != ContentTypeJSON {
		t.Errorf("JSON content type = %q", got)
	}
	r := Func("text/csv", func(w io.Writer) error {
		_, err := io.WriteString(w, "a,b\n")
		return err
	})
	if got := r.ContentType(); got != "text/csv" {
		t.Errorf("Func content type = %q", got)
	}
	var buf bytes.Buffer
	if err := r.Render(&buf); err != nil {
		t.Fatalf("Func render: %v", err)
	}
	if buf.String() != "a,b\n" {
		t.Errorf("Func body = %q", buf.String())
	}
}
