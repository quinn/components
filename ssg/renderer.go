package ssg

import (
	"bytes"
	"io"
	"mime"
	"path/filepath"

	g "maragu.dev/gomponents"
)

// Common MIME types used by the helper renderers.
const (
	ContentTypeHTML = "text/html; charset=utf-8"
	ContentTypeText = "text/plain; charset=utf-8"
	ContentTypeJSON = "application/json"
)

// Renderer produces the bytes for a single route. Implementations must be
// safe to call multiple times. They should also be deterministic for
// Site.Build to produce stable output.
type Renderer interface {
	// ContentType returns the MIME type, including any charset, that should
	// be used when serving or storing the rendered bytes.
	ContentType() string
	// Render writes the body to w.
	Render(w io.Writer) error
}

// HTML wraps a gomponents Node so it can be registered as a Renderer. The
// resulting renderer always reports text/html as its content type.
func HTML(node g.Node) Renderer {
	return htmlRenderer{node: node}
}

type htmlRenderer struct{ node g.Node }

func (h htmlRenderer) ContentType() string       { return ContentTypeHTML }
func (h htmlRenderer) Render(w io.Writer) error { return h.node.Render(w) }

// Bytes returns a Renderer that emits the supplied bytes with the given
// content type. The returned renderer holds a reference to the slice; do
// not mutate it after calling Bytes.
func Bytes(contentType string, data []byte) Renderer {
	return bytesRenderer{typ: contentType, data: data}
}

type bytesRenderer struct {
	typ  string
	data []byte
}

func (b bytesRenderer) ContentType() string { return b.typ }
func (b bytesRenderer) Render(w io.Writer) error {
	_, err := w.Write(b.data)
	return err
}

// Text is shorthand for Bytes with a text/plain content type.
func Text(s string) Renderer {
	return Bytes(ContentTypeText, []byte(s))
}

// JSON is shorthand for Bytes with an application/json content type.
func JSON(data []byte) Renderer {
	return Bytes(ContentTypeJSON, data)
}

// Func adapts an arbitrary writer function as a Renderer with the given
// content type. Useful when the body is streamed or generated lazily.
func Func(contentType string, fn func(io.Writer) error) Renderer {
	return funcRenderer{typ: contentType, fn: fn}
}

type funcRenderer struct {
	typ string
	fn  func(io.Writer) error
}

func (f funcRenderer) ContentType() string         { return f.typ }
func (f funcRenderer) Render(w io.Writer) error    { return f.fn(w) }

// renderToBytes is a small helper used by serve.go and tests.
func renderToBytes(r Renderer) ([]byte, error) {
	var buf bytes.Buffer
	if err := r.Render(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// guessContentType returns a MIME type for a path based on its extension,
// falling back to application/octet-stream.
func guessContentType(name string) string {
	ext := filepath.Ext(name)
	if ext == "" {
		return "application/octet-stream"
	}
	if t := mime.TypeByExtension(ext); t != "" {
		return t
	}
	return "application/octet-stream"
}
