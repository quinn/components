package ssg_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"go.quinn.io/components/ssg"
)

// post is the kind of value a caller might already have in hand: a slice of
// records loaded from disk, a database, or somewhere else.
type post struct {
	Slug, Title string
}

func page(title string, body ...g.Node) g.Node {
	return h.HTML(
		h.Head(h.TitleEl(g.Text(title))),
		h.Body(body...),
	)
}

// Example shows the whole workflow: register a handful of routes, generate
// more from data, and serve the result over HTTP. The same Site could be
// passed to Build to write the tree to disk instead.
func Example() {
	site := ssg.New()

	site.Add("/", ssg.HTML(page("home", h.H1(g.Text("Welcome")))))
	site.Add("/about", ssg.HTML(page("about", h.P(g.Text("About me.")))))

	posts := []post{
		{Slug: "hello-world", Title: "Hello, World"},
		{Slug: "second-post", Title: "Second Post"},
	}
	for _, p := range posts {
		site.Add("/blog/"+p.Slug, ssg.HTML(page(p.Title, h.H1(g.Text(p.Title)))))
	}

	site.Add("/feed.xml", ssg.Bytes("application/rss+xml", []byte("<rss/>")))

	// Hand the handler to httptest just to exercise the API in this example.
	srv := httptest.NewServer(site.Handler())
	defer srv.Close()

	for _, p := range []string{"/", "/about", "/blog/hello-world", "/feed.xml"} {
		fmt.Println(p, "->", status(srv.Client().Get(srv.URL+p)))
	}
	// Output:
	// / -> 200
	// /about -> 200
	// /blog/hello-world -> 200
	// /feed.xml -> 200
}

func status(resp *http.Response, err error) int {
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	return resp.StatusCode
}
