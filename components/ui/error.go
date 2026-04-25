package ui

import (
	"net/http"

	"github.com/labstack/echo/v4"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func ErrorPage(c echo.Context, code int, message string) g.Node {
	title := http.StatusText(code)
	if title == "" {
		title = "Application Error"
	}

	return Layout(c, title, "dashboard-page",
		PageHeader(title),
		Card("",
			h.P(g.Text(message)),
			ButtonRow(
				PrimaryLink("Back to dashboard", "/"),
			),
		),
	)
}
