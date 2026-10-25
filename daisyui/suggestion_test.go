package daisyui

import (
	"strings"
	"testing"
)

// When no declared class is close enough for a suggestion, the message tells the
// reader where the vocabulary comes from. That hint named daisyui.com/llms.txt,
// which is not the skill: the canonical artifact is daisyui.com/SKILL.md (the
// URL the vendored copy, the drift job and the README all use).
func TestStaleSkillHintNamesTheRealURL(t *testing.T) {
	spec := &Component{
		Base:    "badge",
		Classes: map[string]bool{"badge-success": true},
		Parts:   map[string]bool{},
	}
	msg := unknownClassMessage(spec, "badge-zebra")
	if !strings.Contains(msg, "daisyui.com/SKILL.md") {
		t.Errorf("message does not point at the skill: %q", msg)
	}
	if strings.Contains(msg, "llms.txt") {
		t.Errorf("message still names llms.txt: %q", msg)
	}
}
