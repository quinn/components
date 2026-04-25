package ui

import (
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func Field(label, id string, control g.Node) g.Node {
	return h.Label(
		h.Class("field"),
		h.For(id),
		h.Span(h.Class("field__label"), g.Text(label)),
		control,
	)
}

func TextInputField(label, name, placeholder, value string, required bool) g.Node {
	nodes := []g.Node{
		h.ID(name),
		h.Name(name),
		h.Type("text"),
	}
	if required {
		nodes = append(nodes, h.Required())
	}
	if strings.TrimSpace(placeholder) != "" {
		nodes = append(nodes, h.Placeholder(placeholder))
	}
	if value != "" {
		nodes = append(nodes, h.Value(value))
	}

	return Field(label, name, h.Input(nodes...))
}

func TextAreaField(label, name, placeholder, value string) g.Node {
	nodes := []g.Node{
		h.ID(name),
		h.Name(name),
	}
	if strings.TrimSpace(placeholder) != "" {
		nodes = append(nodes, h.Placeholder(placeholder))
	}

	children := []g.Node{}
	if value != "" {
		children = append(children, g.Text(value))
	}

	return Field(label, name, h.Textarea(
		append(nodes, children...)...,
	))
}

func StackPostForm(action string, children ...g.Node) g.Node {
	return PostForm("stack-form", action, children...)
}

func InlinePostForm(action string, children ...g.Node) g.Node {
	return PostForm("inline-form", action, children...)
}

func PostForm(className, action string, children ...g.Node) g.Node {
	nodes := []g.Node{
		h.Action(action),
		h.Method("post"),
	}
	if strings.TrimSpace(className) != "" {
		nodes = append([]g.Node{h.Class(className)}, nodes...)
	}
	nodes = append(nodes, g.Group(children))

	return h.Form(nodes...)
}

func ButtonRow(children ...g.Node) g.Node {
	return h.Div(h.Class("button-row"),
		g.Group(children),
	)
}
