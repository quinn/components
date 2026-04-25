package site

import (
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/yosssi/gohtml"
	c "go.quinn.io/components/components/ui"
)

type componentGroup struct {
	Slug        string
	Name        string
	Description string
	Components  []componentDoc
}

type componentDoc struct {
	Name        string
	Signature   string
	Description string
	Examples    []componentExample
}

type componentExample struct {
	Title string
	Code  string
	Demo  g.Node
}

func toHTML(gomp g.Node) string {
	var b strings.Builder
	gomp.Render(&b)
	return gohtml.Format(b.String())
}

func allGroups() []componentGroup {
	return []componentGroup{
		layoutGroup(),
		dataGroup(),
		contentGroup(),
		buttonsGroup(),
		navigationGroup(),
		formsGroup(),
		surfacesGroup(),
	}
}

func layoutGroup() componentGroup {
	return componentGroup{
		Slug:        "layout",
		Name:        "Layout",
		Description: "Structural components for page and content organisation.",
		Components: []componentDoc{
			{
				Name:        "PageHeader",
				Signature:   "PageHeader(title string, actions ...g.Node) g.Node",
				Description: "Top-of-page heading with an optional row of action links or buttons.",
				Examples: []componentExample{
					{
						Title: "Basic",
						Code:  `ui.PageHeader("Repositories")`,
						Demo:  c.PageHeader("Repositories"),
					},
					{
						Title: "With actions",
						Code: `ui.PageHeader("Repositories",
    ui.PrimaryLink("New", "/new"),
    ui.SecondaryLink("Settings", "/settings"),
)`,
						Demo: c.PageHeader("Repositories",
							c.PrimaryLink("New", "#"),
							c.SecondaryLink("Settings", "#"),
						),
					},
				},
			},
			{
				Name:        "Card",
				Signature:   "Card(title string, children ...g.Node) g.Node",
				Description: "A bordered panel with an optional title bar. The primary container for grouping related content.",
				Examples: []componentExample{
					{
						Title: "With title",
						Code: `ui.Card("Details",
    ui.LabelledValue("Name", "my-backup"),
    ui.LabelledValue("Status", "Active"),
)`,
						Demo: c.Card("Details",
							c.LabelledValue("Name", "my-backup"),
							c.LabelledValue("Status", "Active"),
						),
					},
					{
						Title: "Without title",
						Code:  `ui.Card("", h.P(g.Text("Card body content.")))`,
						Demo:  c.Card("", h.P(g.Text("Card body content."))),
					},
				},
			},
		},
	}
}

func dataGroup() componentGroup {
	return componentGroup{
		Slug:        "data",
		Name:        "Data Display",
		Description: "Components for presenting key-value data and statistics.",
		Components: []componentDoc{
			{
				Name:        "StatGrid",
				Signature:   "StatGrid(children ...g.Node) g.Node",
				Description: "A responsive grid that arranges StatCards in a single row.",
				Examples: []componentExample{
					{
						Title: "With StatCards",
						Code: `ui.StatGrid(
    ui.StatCard("Snapshots", "142"),
    ui.StatCard("Size", "3.2 GB"),
    ui.StatCard("Last Run", "2 min ago"),
)`,
						Demo: c.StatGrid(
							c.StatCard("Snapshots", "142"),
							c.StatCard("Size", "3.2 GB"),
							c.StatCard("Last Run", "2 min ago"),
						),
					},
				},
			},
			{
				Name:        "StatCard",
				Signature:   "StatCard(label, value string) g.Node",
				Description: "A single statistic with a small label and a prominent value. Meant to be used inside a StatGrid.",
				Examples: []componentExample{
					{
						Title: "Standalone",
						Code:  `ui.StatCard("Uptime", "99.9%")`,
						Demo:  c.StatCard("Uptime", "99.9%"),
					},
				},
			},
			{
				Name:        "LabelledValue",
				Signature:   "LabelledValue(label, value string) g.Node",
				Description: "A two-column label–value row, typically stacked inside a Card.",
				Examples: []componentExample{
					{
						Title: "Stacked",
						Code: `ui.LabelledValue("Host", "server-01.example.com")
ui.LabelledValue("Port", "8080")
ui.LabelledValue("Status", "Running")`,
						Demo: g.Group([]g.Node{
							c.LabelledValue("Host", "server-01.example.com"),
							c.LabelledValue("Port", "8080"),
							c.LabelledValue("Status", "Running"),
						}),
					},
				},
			},
		},
	}
}

func contentGroup() componentGroup {
	return componentGroup{
		Slug:        "content",
		Name:        "Content",
		Description: "Text content, empty states, and code display.",
		Components: []componentDoc{
			{
				Name:        "EmptyState",
				Signature:   "EmptyState(message string) g.Node",
				Description: "A muted italic message shown when a list or section has no items.",
				Examples: []componentExample{
					{
						Title: "Default",
						Code:  `ui.EmptyState("No snapshots yet.")`,
						Demo:  c.EmptyState("No snapshots yet."),
					},
				},
			},
			{
				Name:        "ListOrEmpty",
				Signature:   "ListOrEmpty(className, emptyMessage string, items ...g.Node) g.Node",
				Description: "Renders items in a div, or falls back to EmptyState when the list is empty.",
				Examples: []componentExample{
					{
						Title: "With items",
						Code: `ui.ListOrEmpty("surface-list", "No items.",
    ui.SurfaceBlock("", ui.SurfaceHeader("Item One", "", nil)),
    ui.SurfaceBlock("", ui.SurfaceHeader("Item Two", "", nil)),
)`,
						Demo: c.ListOrEmpty("surface-list", "No items.",
							c.SurfaceBlock("", c.SurfaceHeader("Item One", "")),
							c.SurfaceBlock("", c.SurfaceHeader("Item Two", "")),
						),
					},
					{
						Title: "Empty",
						Code:  `ui.ListOrEmpty("surface-list", "Nothing to show.")`,
						Demo:  c.ListOrEmpty("surface-list", "Nothing to show."),
					},
				},
			},
			{
				Name:        "CodeBlock",
				Signature:   "CodeBlock(className, text string) g.Node",
				Description: "A preformatted code block with a subtle background.",
				Examples: []componentExample{
					{
						Title: "Default",
						Code:  "ui.CodeBlock(\"\", `restic backup /home/user`)",
						Demo:  c.CodeBlock("", "restic backup /home/user"),
					},
				},
			},
		},
	}
}

func buttonsGroup() componentGroup {
	return componentGroup{
		Slug:        "buttons",
		Name:        "Buttons & Links",
		Description: "Interactive elements for actions and navigation.",
		Components: []componentDoc{
			{
				Name:        "PrimaryLink",
				Signature:   "PrimaryLink(label, href string) g.Node",
				Description: "An anchor styled as a filled primary button.",
				Examples: []componentExample{
					{
						Title: "Default",
						Code:  `ui.PrimaryLink("Create Repository", "/new")`,
						Demo:  c.PrimaryLink("Create Repository", "#"),
					},
				},
			},
			{
				Name:        "PrimaryButton",
				Signature:   "PrimaryButton(label string) g.Node",
				Description: "A submit button styled as a filled primary button. Intended for use inside forms.",
				Examples: []componentExample{
					{
						Title: "Default",
						Code:  `ui.PrimaryButton("Save")`,
						Demo:  c.PrimaryButton("Save"),
					},
				},
			},
			{
				Name:        "SecondaryLink",
				Signature:   "SecondaryLink(label, href string) g.Node",
				Description: "A text-style navigation link with underline.",
				Examples: []componentExample{
					{
						Title: "Default",
						Code:  `ui.SecondaryLink("View all", "/list")`,
						Demo:  c.SecondaryLink("View all", "#"),
					},
				},
			},
			{
				Name:        "ButtonRow",
				Signature:   "ButtonRow(children ...g.Node) g.Node",
				Description: "A horizontal flex container for grouping buttons.",
				Examples: []componentExample{
					{
						Title: "Default",
						Code: `ui.ButtonRow(
    ui.PrimaryLink("Save", "/save"),
    ui.SecondaryLink("Cancel", "/"),
)`,
						Demo: c.ButtonRow(
							c.PrimaryLink("Save", "#"),
							c.SecondaryLink("Cancel", "#"),
						),
					},
				},
			},
		},
	}
}

func navigationGroup() componentGroup {
	return componentGroup{
		Slug:        "navigation",
		Name:        "Navigation",
		Description: "Filtering, chips, and navigation aids.",
		Components: []componentDoc{
			{
				Name:        "Chip",
				Signature:   "Chip(label, href string, active bool) g.Node",
				Description: "A small pill-shaped toggle link. The active variant inverts its colours.",
				Examples: []componentExample{
					{
						Title: "Inactive and active",
						Code: `ui.Chip("All", "/list", false)
ui.Chip("Running", "/list?status=running", true)
ui.Chip("Failed", "/list?status=failed", false)`,
						Demo: g.Group([]g.Node{
							c.Chip("All", "#", false),
							g.Text(" "),
							c.Chip("Running", "#", true),
							g.Text(" "),
							c.Chip("Failed", "#", false),
						}),
					},
				},
			},
			{
				Name:        "FilterRow",
				Signature:   "FilterRow(label string, children ...g.Node) g.Node",
				Description: "A labelled row of filter chips. Stack multiple FilterRows for multi-facet filtering.",
				Examples: []componentExample{
					{
						Title: "Default",
						Code: `ui.FilterRow("Status",
    ui.Chip("All", "/", false),
    ui.Chip("Running", "/?s=running", true),
    ui.Chip("Failed", "/?s=failed", false),
)`,
						Demo: c.FilterRow("Status",
							c.Chip("All", "#", false),
							c.Chip("Running", "#", true),
							c.Chip("Failed", "#", false),
						),
					},
				},
			},
		},
	}
}

func formsGroup() componentGroup {
	return componentGroup{
		Slug:        "forms",
		Name:        "Forms",
		Description: "Form fields, inputs, and submission layouts.",
		Components: []componentDoc{
			{
				Name:        "Field",
				Signature:   "Field(label, id string, control g.Node) g.Node",
				Description: "A generic label wrapper around any form control.",
				Examples: []componentExample{
					{
						Title: "With a text input",
						Code:  `ui.Field("Name", "name", h.Input(h.Type("text"), h.Name("name")))`,
						Demo:  c.Field("Name", "name", h.Input(h.Type("text"), h.Name("name"))),
					},
				},
			},
			{
				Name:        "TextInputField",
				Signature:   "TextInputField(label, name, placeholder, value string, required bool) g.Node",
				Description: "A labelled text input with optional placeholder and required attribute.",
				Examples: []componentExample{
					{
						Title: "Empty",
						Code:  `ui.TextInputField("Repository", "repo", "e.g. /srv/backup", "", true)`,
						Demo:  c.TextInputField("Repository", "repo", "e.g. /srv/backup", "", true),
					},
					{
						Title: "Pre-filled",
						Code:  `ui.TextInputField("Repository", "repo", "", "/srv/backup", false)`,
						Demo:  c.TextInputField("Repository", "repo-filled", "", "/srv/backup", false),
					},
				},
			},
			{
				Name:        "TextAreaField",
				Signature:   "TextAreaField(label, name, placeholder, value string) g.Node",
				Description: "A labelled multi-line text area.",
				Examples: []componentExample{
					{
						Title: "Default",
						Code:  `ui.TextAreaField("Notes", "notes", "Add notes…", "")`,
						Demo:  c.TextAreaField("Notes", "notes", "Add notes…", ""),
					},
				},
			},
			{
				Name:        "StackPostForm",
				Signature:   "StackPostForm(action string, children ...g.Node) g.Node",
				Description: "A vertically stacked POST form. Convenience wrapper around PostForm with the stack-form class.",
				Examples: []componentExample{
					{
						Title: "Default",
						Code: `ui.StackPostForm("/create",
    ui.TextInputField("Name", "name", "backup name", "", true),
    ui.TextAreaField("Description", "desc", "", ""),
    ui.ButtonRow(ui.PrimaryButton("Create")),
)`,
						Demo: c.StackPostForm("#",
							c.TextInputField("Name", "name-form", "backup name", "", true),
							c.TextAreaField("Description", "desc-form", "", ""),
							c.ButtonRow(c.PrimaryButton("Create")),
						),
					},
				},
			},
			{
				Name:        "InlinePostForm",
				Signature:   "InlinePostForm(action string, children ...g.Node) g.Node",
				Description: "An inline POST form for single-action controls like delete buttons.",
				Examples: []componentExample{
					{
						Title: "Default",
						Code:  `ui.InlinePostForm("/delete", ui.PrimaryButton("Delete"))`,
						Demo:  c.InlinePostForm("#", c.PrimaryButton("Delete")),
					},
				},
			},
		},
	}
}

func surfacesGroup() componentGroup {
	return componentGroup{
		Slug:        "surfaces",
		Name:        "Surfaces",
		Description: "Container components for list items and detail blocks.",
		Components: []componentDoc{
			{
				Name:        "SurfaceLink",
				Signature:   "SurfaceLink(href, className string, children ...g.Node) g.Node",
				Description: "A clickable surface that acts as a link. Highlights on hover.",
				Examples: []componentExample{
					{
						Title: "Default",
						Code: `h.Div(h.Class("surface-list"),
    ui.SurfaceLink("/repo/1", "",
        ui.SurfaceHeader("my-backup", "3 snapshots"),
        ui.SurfaceDetail("Last run: 2 minutes ago"),
    ),
    ui.SurfaceLink("/repo/2", "",
        ui.SurfaceHeader("my-backup", "3 snapshots"),
        ui.SurfaceDetail("Last run: 2 minutes ago"),
    ),
)`,
						Demo: h.Div(h.Class("surface-list"),
							c.SurfaceLink("#", "",
								c.SurfaceHeader("my-backup", "3 snapshots"),
								c.SurfaceDetail("Last run: 2 minutes ago"),
							),
							c.SurfaceLink("#", "",
								c.SurfaceHeader("my-backup", "3 snapshots"),
								c.SurfaceDetail("Last run: 2 minutes ago"),
							),
						),
					},
				},
			},
			{
				Name:        "SurfaceBlock",
				Signature:   "SurfaceBlock(className string, children ...g.Node) g.Node",
				Description: "A non-interactive surface container for detail views.",
				Examples: []componentExample{
					{
						Title: "Default",
						Code: `ui.SurfaceBlock("",
    ui.SurfaceHeader("Snapshot #42", "2024-01-15 09:30"),
    ui.SurfaceStats(
        ui.StatCard("Files", "1,204"),
        ui.StatCard("Size", "842 MB"),
    ),
)`,
						Demo: c.SurfaceBlock("",
							c.SurfaceHeader("Snapshot #42", "2024-01-15 09:30"),
							c.SurfaceStats(
								c.StatCard("Files", "1,204"),
								c.StatCard("Size", "842 MB"),
							),
						),
					},
				},
			},
			{
				Name:        "SurfaceHeader",
				Signature:   "SurfaceHeader(title, meta string, actions ...g.Node) g.Node",
				Description: "Title row for a surface, with optional meta text and action buttons.",
				Examples: []componentExample{
					{
						Title: "With meta and action",
						Code: `ui.SurfaceHeader("my-backup", "running",
    ui.PrimaryLink("View", "/view"),
)`,
						Demo: c.SurfaceHeader("my-backup", "running",
							c.PrimaryLink("View", "#"),
						),
					},
				},
			},
			{
				Name:        "SurfaceDetail",
				Signature:   "SurfaceDetail(text string) g.Node",
				Description: "A muted detail line within a surface.",
				Examples: []componentExample{
					{
						Title: "Default",
						Code:  `ui.SurfaceDetail("Includes /home, /etc, /var/log")`,
						Demo:  c.SurfaceDetail("Includes /home, /etc, /var/log"),
					},
				},
			},
			{
				Name:        "SurfaceStats",
				Signature:   "SurfaceStats(children ...g.Node) g.Node",
				Description: "A grid of small stats inside a surface. Typically wraps StatCard children.",
				Examples: []componentExample{
					{
						Title: "Default",
						Code: `ui.SurfaceStats(
    ui.StatCard("Files", "1,204"),
    ui.StatCard("Added", "12"),
    ui.StatCard("Removed", "3"),
    ui.StatCard("Modified", "5"),
    ui.StatCard("Size", "842 MB"),
)`,
						Demo: c.SurfaceStats(
							c.StatCard("Files", "1,204"),
							c.StatCard("Added", "12"),
							c.StatCard("Removed", "3"),
							c.StatCard("Modified", "5"),
							c.StatCard("Size", "842 MB"),
						),
					},
				},
			},
		},
	}
}
