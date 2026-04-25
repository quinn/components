package ui

import (
	"fmt"
	"os"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// BasePath is prepended to every internal URL emitted by Layout and other
// page renderers. It is read once from the BASE_PATH env var so the same
// build binary can produce a site rooted at "/" for local dev or under a
// subpath like "/components" for GitHub Pages.
var BasePath = func() string {
	p := strings.TrimRight(os.Getenv("BASE_PATH"), "/")
	if p != "" && !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return p
}()

// Path prefixes p with BasePath. p must start with "/".
func Path(p string) string {
	if p == "" {
		return BasePath + "/"
	}
	return BasePath + p
}

func Layout(title, pageClass string, children ...g.Node) g.Node {
	return h.Doctype(
		h.HTML(
			h.Lang("en"),
			h.Head(
				h.Meta(h.Charset("utf-8")),
				h.Meta(
					h.Name("viewport"),
					h.Content("width=device-width, initial-scale=1"),
				),
				h.TitleEl(g.Text(fmt.Sprintf("%s · Components", title))),
				h.Link(h.Rel("stylesheet"), h.Href(Path("/static/css/variables.css"))),
				h.Link(h.Rel("stylesheet"), h.Href(Path("/static/css/base.css"))),
				h.Link(h.Rel("stylesheet"), h.Href(Path("/static/css/layout.css"))),
				h.Link(h.Rel("stylesheet"), h.Href(Path("/static/css/ui/cards.css"))),
				h.Link(h.Rel("stylesheet"), h.Href(Path("/static/css/ui/content.css"))),
				h.Link(h.Rel("stylesheet"), h.Href(Path("/static/css/ui/forms.css"))),
				h.Link(h.Rel("stylesheet"), h.Href(Path("/static/css/ui/surfaces.css"))),
				h.Link(h.Rel("stylesheet"), h.Href(Path("/static/css/docs.css"))),
				h.Script(
					h.Src("https://cdn.jsdelivr.net/npm/htmx.org@2.0.9/dist/htmx.min.js"),
					g.Attr("crossorigin", "anonymous"),
					h.Defer(),
				),
				h.Script(
					h.Src("https://cdn.jsdelivr.net/npm/htmx-ext-sse@2.2.4"),
					g.Attr("crossorigin", "anonymous"),
					h.Defer(),
				),
			),
			h.Body(
				h.Header(h.Class("site-header"),
					h.Div(h.Class("inner"),
						h.A(h.Class("site-brand"), h.Href(Path("/")), g.Text("Components")),
					),
				),
				h.Main(h.Class(pageClass),
					g.Group(children),
				),
				h.Footer(h.Class("site-footer"),
					h.Div(h.Class("inner"),
						h.Div(
							h.A(h.Href("https://github.com/quinn/components"), g.Text("github")),
						),
						h.Div(
							h.A(h.Href("https://quinn.io"), g.Text("quinn.io")),
							g.Text(" • "),
							h.A(h.Href("https://loosecollective.dev"), g.Text("loosecollective.dev")),
						),
					),
				),
				h.Script(g.Raw(`
function switchTab(btn) {
	var example = btn.closest('.example');
	var tab = btn.getAttribute('data-tab');

	// Update tab buttons
	example.querySelectorAll('.code-tab').forEach(function(t) {
		t.classList.toggle('active', t.getAttribute('data-tab') === tab);
	});

	// Update code blocks
	example.querySelectorAll('.code').forEach(function(c) {
		c.classList.toggle('active', c.classList.contains('code-' + tab));
	});
}
`)),
			),
		),
	)
}
