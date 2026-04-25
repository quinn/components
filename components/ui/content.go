package ui

import (
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func EmptyState(message string) g.Node {
	return h.P(h.Class("empty-state"),
		g.Text(message),
	)
}

func ListOrEmpty(className, emptyMessage string, items ...g.Node) g.Node {
	if len(items) == 0 {
		return EmptyState(emptyMessage)
	}

	return h.Div(h.Class(className),
		g.Group(items),
	)
}

func CodeBlock(className, text string) g.Node {
	return h.Pre(h.Class(joinClasses("code-block", className)),
		g.Text(text),
	)
}

func joinClasses(classes ...string) string {
	values := make([]string, 0, len(classes))
	for _, className := range classes {
		className = strings.TrimSpace(className)
		if className == "" {
			continue
		}
		values = append(values, className)
	}

	return strings.Join(values, " ")
}
