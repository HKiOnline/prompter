# Release preparation

Homebrew on macOS (Apple Silicon and Intel) is the primary distribution target. GitHub Releases also hosts Linux and Windows binaries. Releases use GoReleaser OSS and the tag-triggered `.github/workflows/release.yaml` workflow.

## Local validation (does not publish)

Install Go as specified in `go.mod`, Make, and GoReleaser v2.18.3 or newer v2. Run from the repository root:

```sh
make test
goreleaser check
go run ./tools/release-notes -tag v0.8.0 -output scratch/release-notes.md
export PROMPTER_VERSION="$(go run ./tools/release-notes -version)"
goreleaser release --snapshot --clean --skip=publish --release-notes=scratch/release-notes.md
```

Replace the example tag with the intended version. The release-notes command validates the committed version.json and latest changelog heading, rejects mismatched tags, and extracts only that release's curated changes. Omitting `-tag` generates a local preview; CI always supplies the actual Git tag. The snapshot version uses PROMPTER_VERSION when provided, otherwise the latest Git tag version. Snapshots do not update GitHub or the tap, and their generated download URLs are not published URLs.

`--clean` replaces generated `dist/` artifacts. Inspect `dist/homebrew/Casks/prompter.rb`: only macOS ARM64 and AMD64 downloads should be present. Confirm five tar.gz archives (macOS/Linux), three Windows ZIP archives, and the SHA-256 checksum file. The binaries include SQLite and use `-s -w` stripping. Inspect `scratch/release-notes.md` for Homebrew-first instructions and the intended change section.

Run the native snapshot executable with `--version`. To smoke-test its MCP transport, set `PROMPTER_TEST_BINARY` to its absolute path and run `tests/test_server.sh`. The test checks initialize, prompts/list, and prompts/get responses with a bounded deadline and isolated storage. Missing, malformed, or error responses fail the command. Cross-built targets also need runtime testing on their actual operating systems.

## Release metadata and dates

Before cutting a release, confirm the latest CHANGELOG.md heading has the intended release version and actual release date. Update version.json to the matching version and a UTC baseline timestamp (`YYYY-MM-DDTHH:mm:ssZ`), then run the checks again. Do not advance the changelog date merely because a preparation check ran on a later day.

The JSON date is a build-time baseline, not the changelog's release date. Makefile and GoReleaser stamp the real build time and build identifier into each executable. Those generated `+ci...`/`+lb...` suffixes do not belong in the release tag: normally committed `build` is empty, giving a tag such as `v0.8.0`. MCP reports the version without build metadata.

For a prerelease, set preRelease (for example `rc.1`), use a matching changelog heading (`0.8.0-rc.1`) and tag (`v0.8.0-rc.1`). GoReleaser marks it as a prerelease and skips the stable Homebrew tap update. Its release notes direct users to archives rather than the stable cask.

## Publishing (only when explicitly authorized)

1. Review and merge the release changes to main; confirm the worktree is clean and metadata/dates are final.
2. Create and push the matching annotated version tag.
3. The workflow validates the tag/metadata/changelog, runs `make test`, validates GoReleaser configuration, then publishes the GitHub release and updates the tap.
4. Verify the release notes, checksums, and assets, and confirm `Casks/prompter.rb` was committed to `hkionline/homebrew-tap`.
5. Install with `brew install --cask hkionline/tap/prompter`, verify `prompter --version`, and exercise an MCP client on Intel and Apple Silicon Macs. Test `brew upgrade --cask hkionline/tap/prompter` when an upgrade is available.

The GitHub release uses the workflow's GITHUB_TOKEN with contents write permission. Updating the separate public tap requires HOMEBREW_TAP_TOKEN with contents write access to that repository. A non-publishing snapshot cannot verify that token's permissions.

The cask's postflight hook removes quarantine from the unsigned executable. Signing/notarization can be introduced later with the same cask name. Linux deb/rpm packages and Windows installers/package-manager entries are separate future work.
