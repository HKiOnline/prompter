# Pull-request checks

`.github/workflows/pull-request.yaml` runs on pull requests to any branch. Checks run independently; newer commits cancel outdated runs for the same PR. The workflow has read-only repository permissions and does not use release or Homebrew credentials.

| Check | Coverage |
|---|---|
| Actionlint | All GitHub Actions workflows, including embedded shell checks when ShellCheck is available |
| Go formatting | All tracked Go files; reports differences without rewriting files |
| Go vet | Filesystem-only and SQLite-enabled builds |
| Tests | `make test` on Linux and macOS: both Go test variants plus the built-binary MCP smoke test |
| Go race detection | Both build variants on Linux |
| Go vulnerability scan | Reachable known Go vulnerabilities in both build variants, using govulncheck |
| GoReleaser configuration | `goreleaser check`, without publishing |
| Version and changelog | Committed metadata, latest changelog section, and release-note validation tests; no Git tag required |
| Dependency review | Newly introduced vulnerable dependencies, failing at low severity or above |

Go jobs select their toolchain from `go.mod`. Actionlint and govulncheck are pinned in the workflow; update those tool versions explicitly when needed. Vulnerability data is fetched when the scan runs. Dependabot's Go-module updates do not update tools referenced only in workflow shell commands.

## Dependency automation

`.github/dependabot.yml` requests weekly updates for Go modules and GitHub Actions, with up to five open version-update PRs per ecosystem. Dependabot must see this configuration on the default branch before it starts scheduled updates. Update PRs run the same validation workflow.

Dependency review requires the repository's dependency graph to be enabled (and an eligible repository/plan; available for public repositories). Repository administrators can confirm this under GitHub Settings → Advanced Security. Dependabot alerts and automated security fixes are separate repository settings from the version-update configuration.

## Merge requirements

Once the checks have run on GitHub, select their status names in the branch protection rule or ruleset to make them required before merging. Adding a workflow alone does not make its checks mandatory. The release workflow separately checks the actual version tag and reruns tests before publication.
