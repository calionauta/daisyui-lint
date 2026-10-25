# Changelog

## [0.1.0] - 2026-10-25

First release. Extracted from the `gogogo` template, where the rules ran as a
repo-specific gate; the reusable half now ships as a binary so any project can
run it.

### Added

- `daisyui-lint`, a linter for DaisyUI class usage: unknown classes (with a
  nearest-match suggestion), parts without their container, components on the
  wrong form element, classes retired in DaisyUI 5 (with the replacement in the
  message), duplicate size/colour modifiers on one element, and off-scale
  arbitrary values.
- Vocabulary **derived from the vendored DaisyUI skill** rather than maintained
  by hand: a skill reshape is an error, not a silent narrowing, and a version
  bump widens the checks.
- A **bundled snapshot** of the vocabulary so the binary works in a project that
  has no vendored skill; `TestBundledVocabularyMatchesSkill` fails when the
  snapshot is stale, and `make vocabulary` regenerates it.
- `--format sarif` for GitHub code scanning, `--format json` for agents,
  `--strict` to promote warnings, `--allow`/`.daisyui-allow` for accepted
  values, `--css` to teach it the project's own class names, `--version` that
  names the vocabulary source and the DaisyUI release it maps to.
- A composite GitHub Action wrapping the binary.
