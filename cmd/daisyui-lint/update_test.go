package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Semver comparison is the whole point of the check: "0.10.0" must beat "0.9.0"
// (string comparison would say otherwise), a "v" prefix is noise, and a
// pre-release suffix does not make a release newer.
func TestNewerVersion(t *testing.T) {
	t.Parallel()
	cases := []struct {
		current, latest string
		want            bool
	}{
		{"0.1.0", "0.2.0", true},
		{"0.1.0", "v0.2.0", true},
		{"v0.1.0", "0.1.0", false},
		{"0.9.0", "0.10.0", true},
		{"0.10.0", "0.9.0", false},
		{"1.0.0", "0.9.9", false},
		{"0.2.0", "0.2.0-rc1", false},
		{"0.2.0-rc1", "0.2.0", true},
		{"not-a-version", "0.2.0", false},
	}
	for _, tc := range cases {
		if got := newerVersion(tc.current, tc.latest); got != tc.want {
			t.Errorf("newerVersion(%q, %q) = %v, want %v", tc.current, tc.latest, got, tc.want)
		}
	}
}

// The API shape is one field; a newer tag is reported, the same tag is not, and
// an unreachable API is an error the caller turns into a note.
func TestCheckForUpdate(t *testing.T) {
	t.Parallel()
	respond := func(tag string, status int) *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			if status == http.StatusOK {
				_, _ = w.Write([]byte(`{"tag_name":"` + tag + `","name":"x"}`))
			}
		}))
	}

	newer := respond("v0.2.0", http.StatusOK)
	defer newer.Close()
	got, err := checkForUpdate("0.1.0", newer.URL, newer.Client())
	if err != nil || got != "0.2.0" {
		t.Fatalf("checkForUpdate = %q, %v; want 0.2.0", got, err)
	}

	same := respond("v0.1.0", http.StatusOK)
	defer same.Close()
	if got, err := checkForUpdate("0.1.0", same.URL, same.Client()); err != nil || got != "" {
		t.Fatalf("up-to-date must report nothing, got %q, %v", got, err)
	}

	broken := respond("", http.StatusInternalServerError)
	defer broken.Close()
	if _, err := checkForUpdate("0.1.0", broken.URL, broken.Client()); err == nil ||
		!strings.Contains(err.Error(), "500") {
		t.Fatalf("expected an API error, got %v", err)
	}
}

// An unreachable or failing API is a note, never a failing exit code: the flag
// asks a question, it does not gate anything.
func TestRunUpdateCheckNeverFailsTheRun(t *testing.T) {
	t.Parallel()
	broken := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer broken.Close()
	old := releasesAPI
	releasesAPI = broken.URL
	defer func() { releasesAPI = old }()

	var out, errOut strings.Builder
	if code := runUpdateCheck(&out, &errOut, "0.1.0"); code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if !strings.Contains(errOut.String(), "update check failed") {
		t.Fatalf("expected a note on stderr, got %q", errOut.String())
	}
}
