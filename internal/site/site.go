// Package site defines the static site shared between the dev server and
// the build command. Keep the route registrations here so both cmd/dev and
// cmd/build stay in sync.
package site

import (
	"os"

	"go.quinn.io/components/ssg"
)

// New returns the registered site with all component documentation pages.
func New() *ssg.Site {
	s := ssg.New()
	groups := allGroups()

	s.Add("/", ssg.HTML(homePage(groups)))
	for _, grp := range groups {
		s.Add("/"+grp.Slug, ssg.HTML(componentGroupPage(grp, groups)))
	}

	s.Static("/static/css", "./css")

	cssVariables := concatenateCSSFiles([]string{
		"./css/variables.css",
	})

	cssContent := concatenateCSSFiles([]string{
		"./css/base.css",
		"./css/layout.css",
		"./css/docs.css",
		"./css/ui/cards.css",
		"./css/ui/content.css",
		"./css/ui/forms.css",
		"./css/ui/surfaces.css",
	})

	s.Add("/variables.css", ssg.Bytes("text/css; charset=utf-8", cssVariables))
	s.Add("/components.css", ssg.Bytes("text/css; charset=utf-8", cssContent))

	return s
}

// concatenateCSSFiles reads and concatenates the given CSS files.
// Each file is separated by a newline. If a file cannot be read,
// the error is logged and that file is skipped.
func concatenateCSSFiles(paths []string) []byte {
	var result []byte
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			// Log error but continue with other files
			continue
		}
		result = append(result, data...)
		result = append(result, '\n')
	}
	return result
}
