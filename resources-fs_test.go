//+build fsres

package server

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func Test_setupResources(t *testing.T) {
	packageDir, _ := os.Getwd()

	// Since the test working directory is different than the working directory when main is run
	for {
		if fileExists("main.go") {
			break
		}
		err := os.Chdir("..")
		if err != nil {
			t.Fatalf("unable to find main.go - %s", err)
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
		fmt.Printf("Started in package dir %s, changed to %s", packageDir, getwd)
		t.Fatalf("unable to open '%s' - %s", filename, err)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
