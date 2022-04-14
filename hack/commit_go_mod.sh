#!/bin/bash

set -e

git add go.mod go.sum deps.bzl

EMAIL=$(git show -s "--format=%ae" HEAD)

# If the last commit was made by dependabot
if [[ "$EMAIL" == *"dependabot[bot]@users.noreply.github.com"* ]]; then
    echo "Dependabot user detected"
    DEPENDABOT=true
fi

if [ -n "$(git diff-index --cached --name-only HEAD --)" ]; then
    # Amend dependabot commits, but not normal commits
    echo "Changes detected"
    if [ "$DEPENDABOT" = true ]; then
        git commit --amend --no-edit
    else
        git commit -m "chore: update go.mod bazel deps"
    fi

    echo "Pushing changes"
    if [ "$CI" = true ]; then
        git push --force-with-lease origin HEAD:$GITHUB_HEAD_REF
    fi
fi
