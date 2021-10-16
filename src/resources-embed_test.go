//+build !fsres

package src

import (
	"log"
	"testing"
)

func Test_getResources(t *testing.T) {
	dirFs := getEmbeddedResources()

	_, err := dirFs.Open("index.html")
	if err != nil {
		t.Fatalf("unable to find index.html")
	}

	res, err := dirFs.ReadDir("static")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	log.Printf("%#v", res)
}
