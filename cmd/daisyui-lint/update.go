// SCOPE:core - DO NOT REMOVE - explicit update check.
//
// `--check-update` asks the GitHub releases API what the newest tag is and says
// so when it is newer than this binary. It is EXPLICIT on purpose: a linter runs
// in pre-commit and in CI, where a network round-trip on every invocation is a
// latency tax and a failure mode (offline, rate limit, proxy) that has nothing to
// do with the code being linted. `datastar-lint` checks silently on each run with
// a 2s timeout; we took the useful half and left the tax.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// releasesAPI is the endpoint `--check-update` reads. It is a variable so a test
// can point it at an httptest server.
var releasesAPI = "https://api.github.com/repos/calionauta/daisyui-lint/releases/latest"

// latestRelease is the only field of the GitHub payload this tool needs.
type latestRelease struct {
	TagName string `json:"tag_name"`
}

// checkForUpdate returns the latest released version (without the leading "v")
// when it is newer than current, or "" when the binary is up to date.
func checkForUpdate(current, apiURL string, client *http.Client) (string, error) {
	resp, err := client.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("releases API answered %s", resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<16))
	if err != nil {
		return "", err
	}
	var rel latestRelease
	if err := json.Unmarshal(body, &rel); err != nil {
		return "", fmt.Errorf("releases API payload: %w", err)
	}
	latest := strings.TrimPrefix(strings.TrimSpace(rel.TagName), "v")
	if latest == "" {
		return "", fmt.Errorf("releases API returned no tag")
	}
	if newerVersion(current, latest) {
		return latest, nil
	}
	return "", nil
}

// newerVersion reports whether latest is a strictly higher major.minor.patch than
// current. A leading "v" and any pre-release suffix are ignored: the question
// here is only "should this binary be replaced by something newer".
func newerVersion(current, latest string) bool {
	cur, curPre, okCur := semverTriple(current)
	next, nextPre, okNext := semverTriple(latest)
	if !okCur || !okNext {
		return false
	}
	for i := range 3 {
		if next[i] != cur[i] {
			return next[i] > cur[i]
		}
	}
	// Same numbers: a final release is newer than its own pre-release.
	return curPre && !nextPre
}

// semverTriple parses "v1.2.3" or "1.2.3-rc1" into three numbers plus whether it
// is a pre-release (a "-" or "+" suffix).
func semverTriple(v string) ([3]int, bool, bool) {
	v = strings.TrimPrefix(strings.TrimSpace(v), "v")
	pre := false
	if i := strings.IndexAny(v, "-+"); i >= 0 {
		pre = true
		v = v[:i]
	}
	parts := strings.Split(v, ".")
	if len(parts) != 3 {
		return [3]int{}, false, false
	}
	var out [3]int
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return [3]int{}, false, false
		}
		out[i] = n
	}
	return out, pre, true
}

// runUpdateCheck prints the result and returns the process exit code. An
// unreachable API is a note, not a failure: the flag asks a question, it does not
// gate anything.
func runUpdateCheck(stdout, stderr io.Writer, current string) int {
	latest, err := checkForUpdate(current, releasesAPI, &http.Client{Timeout: 3 * time.Second})
	if err != nil {
		fmt.Fprintf(stderr, "daisyui-lint: update check failed: %v\n", err)
		return 0
	}
	if latest == "" {
		fmt.Fprintf(stdout, "daisyui-lint %s is up to date\n", current)
		return 0
	}
	fmt.Fprintf(stdout, "daisyui-lint %s is available (running %s)\n", latest, current)
	fmt.Fprintf(stdout, "  go install github.com/calionauta/daisyui-lint/cmd/daisyui-lint@v%s\n", latest)
	return 0
}
