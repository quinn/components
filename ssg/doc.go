// Package ssg provides primitives for building statically generated sites
// as a Go library rather than via templates and content-folder conventions.
//
// The model is intentionally small. A Site is a registry of routes. Each
// route maps a URL path to a Renderer that produces its bytes. Once routes
// are registered the same Site can either be served by a live HTTP server
// (Site.Serve) or written out to a directory tree (Site.Build).
//
// The library is unopinionated about how content is produced. Renderers can
// wrap gomponents nodes, raw bytes, or arbitrary writer functions. Pages can
// be added one at a time or generated in bulk from any data source the
// caller has in hand.
//
// A minimal example:
//
//	site := ssg.New()
//	site.Add("/", ssg.HTML(homePage()))
//	site.Add("/about", ssg.HTML(aboutPage()))
//	site.Add("/feed.xml", ssg.Bytes("application/rss+xml", feedBytes))
//	site.Static("/assets", "./public")
//
//	if *build {
//	    return site.Build("./dist")
//	}
//	site.EnableLiveReload()
//	return site.Serve(":8080")
package ssg
