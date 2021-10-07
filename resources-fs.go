// +build fsres

package server

import (
	"io/fs"
	"os"
)

// The design goal is to allow the system to run in two modes, one where template files are embedded in the binary, another
// where template files are read from disk on each request. Reading the files on each request should not be used in a
// production system but allows for rapid template development.
//
// This file implements the version that allows direct access to the filesystem.
var Resources fs.FS

const IsEmbedded = false

func init() {
	Resources = setupResources()
}

func setupResources() fs.FS {
	return os.DirFS(".")
}
