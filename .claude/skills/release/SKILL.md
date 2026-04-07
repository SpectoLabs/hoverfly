---
name: release
description: Automate the full hoverfly release process — version bump, GitHub release, wait for CI, trigger Homebrew update.
disable-model-invocation: true
---

# Hoverfly Release

Release version: **$ARGUMENTS**

## Current state

- Current version: !`grep 'hoverfly.version' core/hoverfly.go | grep -o 'v[0-9]*\.[0-9]*\.[0-9]*'`
- Branch: !`git branch --show-current`
- Working tree clean: !`git status --porcelain | wc -l | xargs`

## Instructions

Execute the following phases in order. Confirm with the user before moving to each phase. Stop immediately if any step fails.

### Phase 1 — Version bump & push

1. **Validate** that `$ARGUMENTS` matches the pattern `vX.Y.Z` (e.g. `v1.13.0`). If not, stop and ask the user for a valid version.
2. **Verify** the working tree is clean (`git status --porcelain` is empty) and the current branch is `master`. If not, stop and tell the user.
3. Run: `make update-version VERSION=$ARGUMENTS`
4. Push: `git push origin master`

Tell the user the version commit has been pushed and that a CircleCI build has been triggered (no need to wait for it).

### Phase 2 — Create GitHub release

1. Create the release with auto-generated notes:
   ```
   gh release create $ARGUMENTS --generate-notes --target master
   ```
2. Show the user the release URL.
3. Tell the user this has triggered the CircleCI `deploy-release` job, which will build cross-platform binaries and Docker images. This typically takes 20+ minutes.

### Phase 3 — Wait for release assets

Poll once per minute until **all 7** expected zip bundles appear in the release assets. The expected files are:

- `hoverfly_bundle_OSX_amd64.zip`
- `hoverfly_bundle_OSX_arm64.zip`
- `hoverfly_bundle_windows_amd64.zip`
- `hoverfly_bundle_windows_386.zip`
- `hoverfly_bundle_linux_amd64.zip`
- `hoverfly_bundle_linux_386.zip`
- `hoverfly_bundle_linux_arm64.zip`

To check, run:
```
gh release view $ARGUMENTS --json assets --jq '.assets[].name'
```

Each poll iteration:
- Count how many of the 7 expected files are present
- Report progress to the user: "X/7 assets uploaded..."
- Sleep 60 seconds between checks
- After 45 minutes with no completion, warn the user and ask whether to keep waiting

### Phase 4 — Trigger Homebrew update

Once all 7 assets are confirmed:

1. Trigger the Homebrew formula bump workflow:
   ```
   gh workflow run homebrew-bump-formula.yml -f version=$ARGUMENTS
   ```
2. Tell the user the workflow has been triggered and they need to manually merge the resulting PR in `SpectoLabs/homebrew-tap`.

### Done

Summarize what was completed:
- Version bumped to `$ARGUMENTS`
- GitHub release created with auto-generated notes
- All 7 platform bundles uploaded by CircleCI
- Homebrew formula update triggered
- Remaining manual step: merge the PR in `SpectoLabs/homebrew-tap`
