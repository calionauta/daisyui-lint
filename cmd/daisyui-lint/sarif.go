// SCOPE:core - DO NOT REMOVE - SARIF output for the standalone linter.
//
// SARIF (Static Analysis Results Interchange Format) is what GitHub code
// scanning consumes: with `--format sarif` the findings appear as PR annotations
// and in the Security tab, which is how a linter reaches a team that will not
// read CI logs. The rule ids are the diagnostic codes — they are the contract,
// so the descriptor list is stable and always complete, not only the codes that
// happened to fire.
package main

import (
	"encoding/json"
	"io"
	"sort"

	"github.com/calionauta/daisyui-lint/daisyui"
)

const sarifSchema = "https://json.schemastore.org/sarif-2.1.0.json"

// ruleDescriptions is the human half of each stable code.
var ruleDescriptions = map[daisyui.Code]string{
	daisyui.CodeUnknownClass:      "class is not declared by the vendored DaisyUI skill",
	daisyui.CodePartWithoutParent: "DaisyUI part used without its container in the same file",
	daisyui.CodeHostMismatch:      "DaisyUI component used on the wrong element",
	daisyui.CodeRetiredClass:      "class retired in DaisyUI 5 (use the replacement in the message)",
	daisyui.CodeDuplicateModifier: "two sizes or two colors on one element",
	daisyui.CodeArbitraryValue:    "arbitrary value steps off the theme scale",
}

type sarifLog struct {
	Version string     `json:"version"`
	Schema  string     `json:"$schema"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool    sarifTool     `json:"tool"`
	Results []sarifResult `json:"results"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name           string      `json:"name"`
	Version        string      `json:"version"`
	InformationURI string      `json:"informationUri,omitempty"`
	Rules          []sarifRule `json:"rules"`
}

type sarifRule struct {
	ID               string    `json:"id"`
	ShortDescription sarifText `json:"shortDescription"`
	HelpURI          string    `json:"helpUri,omitempty"`
}

type sarifText struct {
	Text string `json:"text"`
}

type sarifResult struct {
	RuleID    string       `json:"ruleId"`
	Level     string       `json:"level"`
	Message   sarifText    `json:"message"`
	Locations []sarifPlace `json:"locations"`
}

type sarifPlace struct {
	PhysicalLocation sarifPhysical `json:"physicalLocation"`
}

type sarifPhysical struct {
	ArtifactLocation sarifArtifact `json:"artifactLocation"`
	Region           sarifRegion   `json:"region"`
}

type sarifArtifact struct {
	URI string `json:"uri"`
}

type sarifRegion struct {
	StartLine int `json:"startLine"`
}

// writeSARIF renders the findings as a one-run SARIF log. results is never null:
// a consumer that expects an array must not have to special-case "clean".
func writeSARIF(w io.Writer, findings []daisyui.Finding, docsURL string) error {
	results := make([]sarifResult, 0, len(findings))
	for _, f := range findings {
		level := "warning"
		if f.Severity == daisyui.Error {
			level = "error"
		}
		results = append(results, sarifResult{
			RuleID:  string(f.Code),
			Level:   level,
			Message: sarifText{Text: f.Message},
			Locations: []sarifPlace{{PhysicalLocation: sarifPhysical{
				ArtifactLocation: sarifArtifact{URI: f.File},
				Region:           sarifRegion{StartLine: f.Line},
			}}},
		})
	}

	rules := make([]sarifRule, 0, len(ruleDescriptions))
	for _, code := range sortedCodes() {
		rules = append(rules, sarifRule{
			ID:               string(code),
			ShortDescription: sarifText{Text: ruleDescriptions[code]},
			HelpURI:          docsURL + "#" + string(code),
		})
	}

	log := sarifLog{
		Version: "2.1.0",
		Schema:  sarifSchema,
		Runs: []sarifRun{{
			Tool: sarifTool{Driver: sarifDriver{
				Name:    "daisyui-lint",
				Version: Version,
				Rules:   rules,
			}},
			Results: results,
		}},
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(log)
}

// sortedCodes keeps the rule list deterministic (map order is not).
func sortedCodes() []daisyui.Code {
	codes := make([]daisyui.Code, 0, len(ruleDescriptions))
	for code := range ruleDescriptions {
		codes = append(codes, code)
	}
	sort.Slice(codes, func(i, j int) bool { return codes[i] < codes[j] })
	return codes
}
