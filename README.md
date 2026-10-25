# daisyui-lint

A linter for [DaisyUI](https://daisyui.com) class usage. It checks markup against
the DaisyUI vocabulary — derived from the vendored skill, not hand-maintained —
and catches the mistakes that compile, render, and still look almost right.

```bash
go install github.com/calionauta/daisyui-lint/cmd/daisyui-lint@latest
daisyui-lint --version
daisyui-lint --strict web/
```

Prebuilt binaries (linux/darwin, amd64/arm64) are attached to each
[release](https://github.com/calionauta/daisyui-lint/releases).

## Codes

| Code | What it means |
|---|---|
| `UNKNOWN_CLASS` | a class that starts with a component already in use but is not one the skill declares — with the nearest declared class in the message |
| `PART_WITHOUT_PARENT` | a part (`card-body`, `dropdown-content`, …) without its component in the same file |
| `HOST_MISMATCH` | a component the skill documents on exactly one form element, used on another (`toggle` on a `<button>` holds no state) |
| `RETIRED_CLASS` | a DaisyUI 4 name, with the v5 replacement in the message |
| `DUPLICATE_MODIFIER` | two sizes or two colours in one class list — `{MODIFIER}` is one of each |
| `ARBITRARY_VALUE` | an off-scale value (`text-[13px]`) on a style axis; warning by default, error under `--strict` |

These are the failures a compiler cannot see: the classes are valid strings, the
markup renders, and the defect shows up as a control that looks right and behaves
wrong, or as design drift that accumulates one plausible value at a time.

## Usage

```bash
daisyui-lint [flags] [paths...]
```

| Flag | Meaning |
|---|---|
| `--skill PATH` | the vendored DaisyUI skill (`SKILL.md`); default: search upward for `.agents/skills/daisyui/SKILL.md` or `skills/daisyui/SKILL.md`, then fall back to the bundled snapshot |
| `--css GLOBS` | comma-separated globs of YOUR stylesheets: classes they define are not judged as DaisyUI's (a local `chat-layout` is not a typo of `chat`) |
| `--strict` | promote warnings to errors (recommended in an app's CI) |
| `--allow PATH` | accepted arbitrary-value tokens, one per line; default: `.daisyui-allow`, searched upward |
| `--format text\|json\|sarif` | text (default), machine-readable JSON, or SARIF for GitHub code scanning |
| `--emit-vocabulary` | print the parsed skill as JSON (this is how the bundled snapshot is generated) |
| `--version` | tool version, where the vocabulary came from, and the DaisyUI release it maps to |

Exit codes: `0` clean, `1` findings at error severity, `2` usage or setup error.

`.templ`, `.html` and directories are linted; `node_modules/`, `vendor/`,
`testdata/`, `site/` and `dist/` are skipped.

## Severity is the policy

Off-scale values are a judgement, not a bug: `ARBITRARY_VALUE` warns by default,
fails under `--strict`, and a project accepts specific tokens in
`.daisyui-allow`:

```
# one token per line; # starts a comment
text-[10px]
min-h-[7rem]
```

No suppression comments and no ignore-ranges: the allowlist is a reviewed file,
and the codes are the contract.

## Where the vocabulary comes from

The skill is the upstream artifact (`daisyui.com/llms.txt`, vendored under
`skills/daisyui/`, pinned by `UPSTREAM_SHA`). The binary embeds a **generated
snapshot** of it so it works in a project with no vendored copy;
`make vocabulary` regenerates the snapshot and `TestBundledVocabularyMatchesSkill`
fails when it is stale. A reshape of a component section is an error rather than
a silent narrowing. `--version` always says which vocabulary a run used.

## GitHub Action

```yaml
- uses: calionauta/daisyui-lint@v0.1.0
  with:
    paths: web/
    strict: "true"
    css: "src/*.css"
```

## Releases

Conventional commits on `master` drive [release-please](https://github.com/googleapis/release-please):
it keeps one open release PR, and merging it bumps the annotated `Version`
constant, writes `CHANGELOG.md`, tags, and creates the release. The tag then runs
goreleaser, which builds the four binaries and uploads them with a checksum file.
Nothing is tagged by hand.

## License

MIT.
