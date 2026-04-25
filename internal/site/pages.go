package site

import (
	"bytes"
	"fmt"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	comp "go.quinn.io/components/components/ui"
	layout "go.quinn.io/components/internal/ui"
)

// chromaFormatter and chromaStyle are configured once at package init so we
// pay the chroma setup cost a single time per build.
var (
	chromaFormatter = chromahtml.New(
		chromahtml.WithClasses(false),
		chromahtml.PreventSurroundingPre(true),
	)
	chromaStyle = func() *chroma.Style {
		if s := styles.Get("github-dark"); s != nil {
			return s
		}
		return styles.Fallback
	}()
)

// highlightedCode renders code with chroma using inline styles, wrapped in a
// <pre class={class}> tag. On any failure it falls back to plain text inside
// the same <pre> so docs never fail to render.
func highlightedCode(class, lang, code string) g.Node {
	lexer := lexers.Get(lang)
	if lexer == nil {
		return h.Pre(h.Class(class), g.Text(code))
	}
	iter, err := lexer.Tokenise(nil, code)
	if err != nil {
		return h.Pre(h.Class(class), g.Text(code))
	}
	var buf bytes.Buffer
	if err := chromaFormatter.Format(&buf, chromaStyle, iter); err != nil {
		return h.Pre(h.Class(class), g.Text(code))
	}
	return h.Pre(h.Class(class), g.Raw(buf.String()))
}

func docsPage(title, activeSlug string, groups []componentGroup, content ...g.Node) g.Node {
	return layout.Layout(title, "docs",
		docsNav(groups, activeSlug),
		h.Div(h.Class("docs-content"),
			g.Group(content),
		),
	)
}

func homePage(groups []componentGroup) g.Node {
	total := 0
	for _, grp := range groups {
		total += len(grp.Components)
	}

	cards := make([]g.Node, 0, len(groups))
	for _, grp := range groups {
		cards = append(cards, groupCard(grp))
	}

	return docsPage("Components", "", groups,
		comp.PageHeader("Components"),
		h.P(h.Class("docs-intro"),
			g.Text(fmt.Sprintf(
				"A library of %d Go UI components built with gomponents, "+
					"organised into %d groups. Each component is a function "+
					"that returns a gomponents Node.",
				total, len(groups),
			)),
		),
		h.Div(h.Class("component-grid"),
			g.Group(cards),
		),
	)
}

func groupCard(grp componentGroup) g.Node {
	return h.A(h.Class("group-card"), h.Href(layout.Path("/"+grp.Slug)),
		h.Span(h.Class("name"), g.Text(grp.Name)),
		h.Span(h.Class("desc"), g.Text(grp.Description)),
		h.Span(h.Class("count"),
			g.Text(fmt.Sprintf("%d components", len(grp.Components))),
		),
	)
}

func componentGroupPage(grp componentGroup, groups []componentGroup) g.Node {
	sections := make([]g.Node, 0, len(grp.Components))
	for _, c := range grp.Components {
		sections = append(sections, componentSection(c))
	}
	return docsPage(grp.Name, grp.Slug, groups,
		comp.PageHeader(grp.Name, comp.SecondaryLink("All Components", layout.Path("/"))),
		h.P(h.Class("docs-intro"), g.Text(grp.Description)),
		g.Group(sections),
	)
}

func componentSection(c componentDoc) g.Node {
	examples := make([]g.Node, 0, len(c.Examples))
	for _, ex := range c.Examples {
		examples = append(examples, exampleBlock(ex))
	}
	return h.Section(h.Class("component-section"),
		h.Div(h.Class("component-meta"),
			h.H2(g.Text(c.Name)),
			g.If(c.Description != "",
				h.P(h.Class("component-description"), g.Text(c.Description)),
			),
			highlightedCode("component-signature", "go", c.Signature),
		),
		g.Group(examples),
	)
}

func exampleBlock(ex componentExample) g.Node {
	return h.Div(h.Class("example"),
		g.If(ex.Title != "",
			h.Div(h.Class("title"), g.Text(ex.Title)),
		),
		h.Div(h.Class("demo"), ex.Demo),
		h.Div(h.Class("code"),
			highlightedCode("code-block", "go", ex.Code),
		),
	)
}

func docsNav(groups []componentGroup, activeSlug string) g.Node {
	items := []g.Node{
		navLink(layout.Path("/"), "Overview", activeSlug == ""),
	}
	for _, grp := range groups {
		items = append(items,
			h.Span(h.Class("group"), g.Text(grp.Name)),
			navLink(layout.Path("/"+grp.Slug), grp.Name, activeSlug == grp.Slug),
		)
	}
	return h.Nav(h.Class("docs-nav"), g.Group(items))
}

func navLink(href, label string, active bool) g.Node {
	attrs := []g.Node{h.Href(href), g.Text(label)}
	if active {
		attrs = append(attrs, g.Attr("aria-current", "page"))
	}
	return h.A(attrs...)
}
