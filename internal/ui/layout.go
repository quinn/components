package ui

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

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
				h.Link(h.Rel("stylesheet"), h.Href("/static/css/variables.css")),
				h.Link(h.Rel("stylesheet"), h.Href("/static/css/base.css")),
				h.Link(h.Rel("stylesheet"), h.Href("/static/css/layout.css")),
				h.Link(h.Rel("stylesheet"), h.Href("/static/css/ui/cards.css")),
				h.Link(h.Rel("stylesheet"), h.Href("/static/css/ui/content.css")),
				h.Link(h.Rel("stylesheet"), h.Href("/static/css/ui/forms.css")),
				h.Link(h.Rel("stylesheet"), h.Href("/static/css/ui/surfaces.css")),
				h.Link(h.Rel("stylesheet"), h.Href("/static/css/docs.css")),
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
						h.A(h.Class("site-brand"), h.Href("/"), g.Text("Components")),
					),
				),
				h.Main(h.Class(pageClass),
					g.Group(children),
				),
			),
		),
	)
}
