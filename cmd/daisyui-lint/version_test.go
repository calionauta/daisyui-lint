package main

import (
	"os"
	"regexp"
	"testing"
)

// The version constant must equal the newest CHANGELOG entry. A release that
// bumps the tag and the changelog but not the constant leaves every installed
// binary reporting the old version — invisible until someone compares strings.
// This is the release discipline `datastar-lint` established; the failure is a
// test, not an eyeball.
func TestVersionTracksChangelog(t *testing.T) {
	t.Parallel()
	raw, err := os.ReadFile("../../CHANGELOG.md")
	if err != nil {
		t.Fatalf("read CHANGELOG: %v", err)
	}
	m := regexp.MustCompile(`(?m)^## \[([0-9]+\.[0-9]+\.[0-9]+)\]`).FindStringSubmatch(string(raw))
	if m == nil {
		t.Fatal("no `## [x.y.z]` entry in CHANGELOG.md")
	}
	if Version != m[1] {
		t.Fatalf("Version = %q, newest CHANGELOG entry = %q: bump the constant", Version, m[1])
	}
}
