package service

import (
	"testing"
)

func TestParseAllDirs(t *testing.T) {
	rootDir := "/Users/silhouette/work-practice/repo_profile"
	if err := ParseAllDirs(rootDir); err != nil {
		t.Fatal(err)
	}
	t.Logf("ParseAllDirs success")
}
