package ui

import (
	"fmt"

	"github.com/labstack/echo/v4"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func Layout(c echo.Context, title, pageClass string, children ...g.Node) g.Node {
	return h.Doctype(
		h.HTML(
			h.Lang("en"),
			h.Head(
				h.Meta(h.Charset("utf-8")),
				h.Meta(
					h.Name("viewport"),
					h.Content("width=device-width, initial-scale=1"),
				),
				h.TitleEl(g.Text(fmt.Sprintf("%s · restic-gui", title))),
				h.Link(h.Rel("stylesheet"), h.Href("/static/css/variables.css")),
				h.Link(h.Rel("stylesheet"), h.Href("/static/css/base.css")),
				h.Link(h.Rel("stylesheet"), h.Href("/static/css/layout.css")),
				h.Link(h.Rel("stylesheet"), h.Href("/static/css/ui/cards.css")),
				h.Link(h.Rel("stylesheet"), h.Href("/static/css/ui/content.css")),
				h.Link(h.Rel("stylesheet"), h.Href("/static/css/ui/forms.css")),
				h.Link(h.Rel("stylesheet"), h.Href("/static/css/ui/surfaces.css")),
				h.Link(h.Rel("stylesheet"), h.Href("/static/css/pages/dashboard.css")),
				h.Link(h.Rel("stylesheet"), h.Href("/static/css/pages/instance.css")),
				h.Link(h.Rel("stylesheet"), h.Href("/static/css/pages/run.css")),
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
						h.A(h.Class("site-brand"), h.Href("/"), g.Text("restic-gui")),
					),
				),
				h.Main(h.Class(pageClass),
					g.Group(children),
				),
			),
		),
	)
}

func PageHeader(title string, actions ...g.Node) g.Node {
	return h.Section(h.Class("page-header"),
		h.H1(g.Text(title)),
		g.If(len(actions) != 0,
			h.Nav(h.Class("page-header__nav"),
				g.Group(actions),
			),
		),
	)
}

func Card(title string, children ...g.Node) g.Node {
	return h.Section(h.Class("card"),
		g.If(title != "",
			h.H2(h.Class("card__title"), g.Text(title)),
		),
		h.Div(h.Class("card__body"),
			g.Group(children),
		),
	)
}

func StatGrid(children ...g.Node) g.Node {
	return h.Div(h.Class("stat-grid"),
		g.Group(children),
	)
}

func StatCard(label, value string) g.Node {
	return h.Div(h.Class("stat-card"),
		h.P(h.Class("stat-card__label"), g.Text(label)),
		h.P(h.Class("stat-card__value"), g.Text(value)),
	)
}

func LabelledValue(label, value string) g.Node {
	return h.Div(h.Class("labelled-value"),
		h.P(h.Class("labelled-value__label"), g.Text(label)),
		h.P(h.Class("labelled-value__value"), g.Text(value)),
	)
}

func PrimaryLink(label, href string) g.Node {
	return h.A(
		h.Class("button button--primary"),
		h.Href(href),
		g.Text(label),
	)
}

func PrimaryButton(label string) g.Node {
	return h.Button(
		h.Class("button button--primary"),
		h.Type("submit"),
		g.Text(label),
	)
}

func SecondaryLink(label, href string) g.Node {
	return h.A(
		h.Class("nav-link"),
		h.Href(href),
		g.Text(label),
	)
}

func Chip(label, href string, active bool) g.Node {
	cls := "chip"
	if active {
		cls = "chip chip--active"
	}
	return h.A(
		h.Class(cls),
		h.Href(href),
		g.Text(label),
	)
}

func FilterRow(label string, children ...g.Node) g.Node {
	return h.Div(h.Class("filter-row"),
		h.Span(h.Class("filter-row__label"), g.Text(label)),
		h.Div(h.Class("filter-row__options"),
			g.Group(children),
		),
	)
}
