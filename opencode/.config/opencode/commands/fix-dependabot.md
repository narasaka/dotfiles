---
description: Resolve Dependabot security alerts — 1 issue = 1 PR
---

Resolve open Dependabot security alerts for this repository. $ARGUMENTS

Base branch: `$1` if provided, otherwise `dev`.

## Current open alerts

<alerts>
!`gh api "repos/$(gh repo view --json nameWithOwner -q .nameWithOwner)/dependabot/alerts" --jq '.[] | select(.state == "open") | {number, dependency: .dependency.package.name, ecosystem: .dependency.package.ecosystem, severity: .security_advisory.severity, summary: .security_advisory.summary, manifest: .dependency.manifest_path, vulnerable_range: .security_vulnerability.vulnerable_version_range, first_patched: .security_vulnerability.first_patched_version.identifier, ghsa_id: .security_advisory.ghsa_id}'`
</alerts>

## Instructions

If there are no open alerts, report that and stop.

Group alerts that share the same GHSA ID — they are one logical vulnerability.

For **each unique GHSA**, on a fresh branch off the latest base branch:

1. **Branch**: `fix/<manifest-scope>-<short-kebab-description>` off the latest base branch.
2. **Fix**: Apply the minimal dependency upgrade that satisfies the patched version. Use the ecosystem's native tooling (`go get`, `uv lock`, `npm install`, `cargo update`, etc.). Run the tidy/lock step (`go mod tidy`, `uv lock`, etc.).
3. **Build**: Verify the build passes for the affected package.
4. **Commit**: Atomic conventional commit — `fix(<scope>): <short description> (<GHSA-ID>)`.
5. **PR**: Push and open a PR against the base branch with:
   - Which alert numbers are resolved (link them)
   - GHSA ID and vulnerability summary
   - Before/after version table
   - Verification steps taken

**1 unique GHSA = 1 branch = 1 PR.**

After all PRs are created, return to the original branch and list all created PRs.
