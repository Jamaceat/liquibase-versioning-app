package path

import (
	"log"
	"testing"
)

func Test(t *testing.T) {

	root := GetFindProjectRoot()

	if root == "" {
		t.Fatalf("Not path")
	}

	log.Println(root)

}
