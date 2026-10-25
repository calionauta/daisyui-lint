# daisyui-lint

A linter for [DaisyUI](https://daisyui.com) class usage. It checks markup against
the DaisyUI vocabulary — derived from the vendored skill, not hand-maintained —
and catches the mistakes that compile, render, and still look almost right.

```bash
# Pin the version: a gate that resolves @latest changes under you.
go install github.com/calionauta/daisyui-lint/cmd/daisyui-lint@v0.2.0
daisyui-lint --version      # tool version, vocabulary source, DaisyUI release
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
| `DUPLICATE_MODIFIER` | two sizes or two colours on one element — a component takes one of each |
| `ARBITRARY_VALUE` | an off-scale value (`text-[13px]`) on a style axis; warning by default, error under `--strict` |

These are the failures a compiler cannot see: the classes are valid strings, the
markup renders, and the defect shows up as a control that looks right and behaves
wrong, or as design drift that accumulates one plausible value at a time.

## When it pays

You do not need this to build a DaisyUI app. You need it when markup changes
faster than anyone reads the class list — agents writing views, a v4 to v5
migration, contributors who have not memorised the vocabulary. It costs about
0.3s on a mid-size app and prints one line when clean.

## What this is not

- **Not a CSS build check.** Whether a used class ended up in your compiled
  bundle is a build-time question (Tailwind emits what it scans); use your build
  gate for that. This lints the vocabulary, not the bundle.
- **Not your design system's gate.** Classes your own stylesheets define are
  declared with `--css` and are never judged as DaisyUI's.
- **Not a visual or runtime check.** It cannot tell you a control behaves wrong;
  it tells you the class cannot mean what the markup implies.

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
| `--check-update` | ask GitHub for the newest release and print it if it is newer; never fails the run, never runs implicitly (a linter in CI must not pay a network round-trip per invocation) |
| `--version` | tool version, where the vocabulary came from, and the DaisyUI release it maps to |

Exit codes: `0` clean, `1` findings at error severity, `2` usage or setup error.

`.templ`, `.html` and directories are linted; `node_modules/`, `vendor/`,
`.git/`, `testdata/`, `site/` and `dist/` are skipped.

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

The skill is the upstream artifact (`daisyui.com/SKILL.md`, vendored under
`skills/daisyui/`, pinned by `UPSTREAM_SHA`). The binary embeds a **generated
snapshot** of it so it works in a project with no vendored copy;
`make vocabulary` regenerates the snapshot and `TestBundledVocabularyMatchesSkill`
fails when it is stale. A reshape of a component section is an error rather than
a silent narrowing. `--version` always says which vocabulary a run used.

## GitHub Action

```yaml
- uses: calionauta/daisyui-lint@v0.2.0
  with:
    paths: web/
    strict: "true"
    css: "src/*.css"
    version: v0.2.0   # keep in step with the ref above
```

The action is composite: it resolves the linter with
`go run github.com/calionauta/daisyui-lint/cmd/daisyui-lint@<version>`. `version`
defaults to `latest`, independently of the action ref — pin both, or the gate
moves when you push.

## Releases

Conventional commits on `master` drive [release-please](https://github.com/googleapis/release-please):
it keeps one open release PR, and merging it bumps the annotated `Version`
constant, writes `CHANGELOG.md`, tags, and creates the release. The tag then runs
goreleaser, which builds the four binaries and uploads them with a checksum file.
Nothing is tagged by hand.

## License

MIT.
