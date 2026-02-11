# CI/CD Pipeline Improvement Plan

## A. Eliminate Duplicate PR Runs
- [x] Remove `pull_request` trigger from `ci.yml` so PRs are owned entirely by `pr.yml`
- [x] Reconcile the format check inconsistency: `ci.yml` uses `goimports`, `pr.yml` now also uses `goimports`

## B. Centralize Go Version and Setup with a Composite Action
- [x] Create `.github/actions/setup-go/action.yml` composite action encapsulating checkout, setup-go, and dependency download
- [x] Define the Go version once (composite action default `'1.24'`) instead of hardcoding in every workflow
- [x] Replace the repeated checkout + setup-go + deps boilerplate in all jobs with the composite action

## C. Add Concurrency Groups
- [x] Add `concurrency` block to `ci.yml`, `pr.yml`, `codeql.yml`, and `release.yml` to cancel superseded runs on the same branch/PR
  - `ci.yml`, `codeql.yml`: `group: workflow-ref`, `cancel-in-progress: true`
  - `pr.yml`: `group: workflow-pr_number`, `cancel-in-progress: true`
  - `release.yml`: `group: workflow-ref`, `cancel-in-progress: false` (safety)

## D. Pin All Actions to Tagged Versions
- [x] Replace `securego/gosec@master` with `securego/gosec@v2.23.0`
- [x] Replace `golangci-lint version: latest` with `version: v2.5.0` (matches Makefile)
- [x] Update `codecov/codecov-action@v4` to `@v5`
- [x] Audited all other action references — all current (`checkout@v4`, `setup-go@v5`, `github-script@v7`, `codeql-action@v3`, `action-gh-release@v2`)

## E. Update Go Matrix to Supported Versions
- [x] Update the test matrix in `ci.yml` from `['1.21', '1.22', '1.23']` to `['1.24', '1.25']`
- [x] Update `go.mod` minimum version from `1.21` to `1.24` (matches oldest matrix version)
- [x] Uncomment the `intrange` and `copyloopvar` linters in `.golangci.yml` (minimum Go version is now 1.24)

## F. Consolidate Test Steps
- [x] Remove the separate Unix/Windows test steps in `ci.yml` and use a single step for all platforms

## G. Extract PR Comment Helper
- [x] Factor the find-or-update bot comment JS logic into `.github/scripts/upsert-comment.js` shared helper, used by both coverage and benchmark comment steps in `pr.yml`

## H. Extract Coverage Threshold to a Variable
- [x] Define the `80%` coverage threshold once as a workflow-level `env.COVERAGE_THRESHOLD` variable
- [x] Reference it in both the check step and the comment template in `pr.yml` to prevent drift

## I. Adopt Semver Release Tagging (replaces current CalVer)

### Migration Strategy
- [x] Initial semver release will be `v1.0.0` (CalVer tags filtered out via `awk -F'[v.]' '$2 < 100'`)
- [x] Replace CalVer logic in `release.yml` with semver derivation using `--sort=-version:refname`
- [x] Determine bump type from conventional commit prefixes (`feat!:`/`BREAKING CHANGE` → major, `feat:` → minor, else → patch)
- [x] Support manual override of bump type via `workflow_dispatch` input (auto/major/minor/patch)
- [x] Tag creation produces annotated tags with the semver version
- [x] GitHub Release title/body references semver version

### Changelog Generation
- [x] Replace `sed`-based placeholder substitution with direct file writes using a category loop
- [x] Keep conventional commit categorization (Features, Bug Fixes, etc.) in streamlined ~50-line script
- [x] Include "Full Changelog" diff link in release notes

### Go Module Compatibility
- [x] Module path remains `github.com/mstephenholl/go-solar` (no `/v2` suffix needed until v2.0.0)
- [ ] Verify that `go install github.com/mstephenholl/go-solar@v1.x.x` resolves correctly once first semver tag is published

## J. Implementation Order

| Phase | Items | Status |
|-------|-------|--------|
| 1 | **B** (composite action), **E** (Go matrix) | Done |
| 2 | **A** (dedupe PR runs), **C** (concurrency), **D** (pin actions) | Done |
| 3 | **F** (consolidate tests), **G** (PR comment helper), **H** (coverage threshold) | Done |
| 4 | **I** (semver + changelog) | Done |
