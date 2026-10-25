# Changelog

## [0.2.1](https://github.com/calionauta/daisyui-lint/compare/v0.2.0...v0.2.1) (2026-10-25)


### Fixed

* **cli:** help lists --css and sarif, exits 0, and points at SKILL.md ([#8](https://github.com/calionauta/daisyui-lint/issues/8)) ([0442c7e](https://github.com/calionauta/daisyui-lint/commit/0442c7e89fe068612ec1246d08f8e9f0626407eb))


### Documentation

* say what the linter is not, and fix the upstream URL ([#6](https://github.com/calionauta/daisyui-lint/issues/6)) ([4fdc0f1](https://github.com/calionauta/daisyui-lint/commit/4fdc0f1bfc18a63a7fc1e4fa2a092df8584f0b01))

## [0.2.0](https://github.com/calionauta/daisyui-lint/compare/v0.1.0...v0.2.0) (2026-10-25)


### Added

* add --check-update ([#2](https://github.com/calionauta/daisyui-lint/issues/2)) ([f65ed1d](https://github.com/calionauta/daisyui-lint/commit/f65ed1dec50859fbc2493f7f0edf040009b8d89c))

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
