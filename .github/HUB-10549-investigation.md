# HUB-10549 Investigation Result

## Finding: No `actions/missing-workflow-permissions` CodeQL alerts exist

The CodeQL code scanning alerts API was queried for all pages of
`tool:CodeQL` + `rule:actions/missing-workflow-permissions` + `state:open`
alerts on the default branch (`master`):

```
gh api '/repos/happyreturns/fedex/code-scanning/alerts?tool_name=CodeQL&state=open&per_page=100&page=1' \
  --jq '.[] | select(.rule.id == "actions/missing-workflow-permissions")'
```

**Result:** HTTP 404 — `"no analysis found"`

## Root Cause

This repository has never had CodeQL run against it — there are no code
scanning analyses on record. The repo has no `.github/workflows/` directory;
only `.github/dependabot.yml` exists.

Because there are no workflow files, CodeQL has never been configured to run,
and there are therefore zero `actions/missing-workflow-permissions` alerts
to fix.

## Recommendation

This ticket can be **closed as not applicable** for `happyreturns/fedex`.
There are no GitHub Actions workflows in this repository that require
permissions hardening.
