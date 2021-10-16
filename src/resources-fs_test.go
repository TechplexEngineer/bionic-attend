//+build fsres

package src

import (
	"log"
	"os"
	"testing"
)

func Test_setupResources(t *testing.T) {
	packageDir, _ := os.Getwd()

	// Since the test working directory for this test is different
	// then the working directory when main is run
	for {
		if fileExists("go.mod") {
			break
		}
		err := os.Chdir("..")
		if err != nil {
			t.Fatalf("unable to find go.mod - %s", err)
		}
		cwd, _ := os.Getwd()
		if cwd == "/" {
			t.Fatalf("unable to find main.go, checked all parent directories")
		}
		log.Printf("Changing dir to %s", cwd)
	}

	dirFS := setupResources()

	filename := "partials/layout.html"

	_, err := dirFS.Open(filename)
	if err != nil {
		getwd, _ := os.Getwd()
		log.Printf("Started in package dir %s, changed to %s", packageDir, getwd)
		t.Fatalf("unable to open '%s' - %s", filename, err)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
