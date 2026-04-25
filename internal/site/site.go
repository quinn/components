// Package site defines the static site shared between the dev server and
// the build command. Keep the route registrations here so both cmd/dev and
// cmd/build stay in sync.
package site

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"go.quinn.io/components/ssg"
)

// New returns the registered site. It currently contains a single blank
// page plus the shared CSS mount.
func New() *ssg.Site {
	s := ssg.New()
	s.Add("/", ssg.HTML(home()))
	s.Static("/static/css", "./css")
	return s
}

func home() g.Node {
	return h.Doctype(
		h.HTML(
			h.Lang("en"),
			h.Head(
				h.Meta(h.Charset("utf-8")),
				h.Meta(
					h.Name("viewport"),
					h.Content("width=device-width, initial-scale=1"),
				),
				h.TitleEl(g.Text("Components")),
			),
			h.Body(
				h.H1(g.Text("Components")),
			),
		),
	)
}
