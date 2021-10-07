//+build !fsres

package server

import "embed"

// NOTE: See resources-fs.go for a description of the design goal

//go:embed `partials/*.html`
//go:embed `*.html`
//go:embed `static/*`
var Resources embed.FS

const IsEmbedded = true

// for testing
func getEmbeddedResources() embed.FS {
	return Resources
}
