// Package site defines the static site shared between the dev server and
// the build command. Keep the route registrations here so both cmd/dev and
// cmd/build stay in sync.
package site

import (
	"go.quinn.io/components/ssg"
)

// New returns the registered site with all component documentation pages.
func New() *ssg.Site {
	s := ssg.New()
	groups := allGroups()

	s.Add("/", ssg.HTML(homePage(groups)))
	for _, grp := range groups {
		s.Add("/components/"+grp.Slug, ssg.HTML(componentGroupPage(grp, groups)))
	}

	s.Static("/static/css", "./css")
	return s
}
