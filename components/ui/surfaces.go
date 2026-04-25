package ui

import (
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func SurfaceLink(href, className string, children ...g.Node) g.Node {
	return h.A(
		h.Class(joinClasses("surface", className)),
		h.Href(href),
		g.Group(children),
	)
}

func SurfaceBlock(className string, children ...g.Node) g.Node {
	return h.Article(
		h.Class(joinClasses("surface", className)),
		g.Group(children),
	)
}

func SurfaceHeader(title, meta string, actions ...g.Node) g.Node {
	return h.Div(h.Class("surface__header"),
		h.Div(h.Class("surface__main"),
			h.H3(h.Class("surface__title"), g.Text(title)),
			g.If(strings.TrimSpace(meta) != "",
				h.P(h.Class("surface__meta"),
					g.Text(meta),
				),
			),
		),
		g.If(len(actions) != 0,
			h.Div(h.Class("surface__actions"),
				g.Group(actions),
			),
		),
	)
}

func SurfaceDetail(text string) g.Node {
	return h.P(h.Class("surface__detail"),
		g.Text(text),
	)
}

func SurfaceStats(children ...g.Node) g.Node {
	return h.Div(h.Class("surface__stats"),
		g.Group(children),
	)
}
