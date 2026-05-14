---
description: "Commit all changes, push to remote, and create a pull request"
argument-hint: "<optional: PR title or description>"
---

Do the following steps in order. Stop and report if any step fails.

1. **Commit**: Look at all staged and unstaged changes. Create a well-crafted commit (or multiple atomic commits if the changes are logically separate). Every commit message MUST use [Conventional Commits v1.0.0](https://www.conventionalcommits.org/en/v1.0.0/): `<type>[optional scope][!]: <description>`, with optional body/footer as needed. Prefer the most accurate type such as `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, or `chore`. Use `!` and/or a `BREAKING CHANGE:` footer for breaking changes. Follow repository-local style only when it is compatible with Conventional Commits; Conventional Commits takes precedence.
2. **Push**: Push the current branch to the remote. If no upstream is set, push with `-u origin HEAD`.
3. **Create PR**: Create a pull request using `gh pr create`. Target the repository's default branch. The PR title MUST also use Conventional Commits v1.0.0 header format: `<type>[optional scope][!]: <description>`. Write a clear body summarizing the changes.

If the user provided arguments, use them to guide the PR title/description while preserving the Conventional Commits title format: $ARGUMENTS

IMPORTANT:
- Do NOT force push.
- Do NOT push to main/master directly. If on main/master, create a new branch first with a descriptive name before committing.
- If there are no changes to commit, skip the commit step and just push + create PR.
- If a PR already exists for this branch, report the existing PR URL instead of creating a duplicate.
