package main

import (
	"bytes"
	"strings"
	"testing"
)

// The help text is a literal block, so it drifts silently from the flags it is
// supposed to describe. It already had: --css was missing (a user reading
// --help never learned their own stylesheets can be exempted) and --format
// advertised "text or json" long after SARIF shipped. The stale-skill hint also
// named daisyui.com/llms.txt instead of SKILL.md.
func TestHelpListsEveryFlagAndExitsZero(t *testing.T) {
	var out, errb bytes.Buffer
	code := run([]string{"--help"}, &out, &errb)
	help := errb.String()

	for _, flag := range []string{
		"--skill", "--allow", "--strict", "--css", "--format",
		"--version", "--emit-vocabulary", "--check-update",
	} {
		if !strings.Contains(help, flag) {
			t.Errorf("help does not mention %s", flag)
		}
	}
	for _, format := range []string{"text", "json", "sarif"} {
		if !strings.Contains(help, format) {
			t.Errorf("help does not mention the %s output format", format)
		}
	}
	if strings.Contains(help, "llms.txt") {
		t.Error("help still names llms.txt; the skill lives at daisyui.com/SKILL.md")
	}
	if code != 0 {
		t.Errorf("--help exit code = %d, want 0 (asking for help is not a usage error)", code)
	}
}
