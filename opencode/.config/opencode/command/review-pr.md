---
description: "Address all reviewer comments on a PR: fix valid issues, reply to false positives, resolve threads, then push"
argument-hint: "<PR URL>"
---

Check the comments made by a reviewer in this PR: $ARGUMENTS

For each comment:
1. Determine whether the comment raises a valid issue or is just a false positive.
2. If the issue is valid, implement a fix. Otherwise, reply in the thread with a clear explanation of why it's a false positive.
3. Resolve the conversation / thread.

Once ALL comments are addressed and ALL conversations / threads are resolved, push the changes.
DO NOT push before resolving all the conversations / threads.
