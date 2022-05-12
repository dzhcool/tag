package service

import (
	"testing"
)

// go test -v ./service/ -test.run TestGenTag
func TestGenTag(t *testing.T) {
	gitSvc := NewGitSvc()

	var exists = []string{"v0.0.1"}
	tag, err := gitSvc.genTag(exists)
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	t.Logf("tag:%s", tag)
}
