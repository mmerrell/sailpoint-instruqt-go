#!/bin/bash
# Run this from inside your local temporal-go-hands-on folder, in a normal
# Terminal (not through Cowork/Claude) — it needs real delete permissions,
# which the sandbox that wrote these files didn't have.
#
#   cd ~/Projects/temporal-go-hands-on
#   chmod +x finish-migration.sh
#   ./finish-migration.sh
#
# This repo already has `origin` pointing at github.com/mmerrell/temporal-go-hands-on,
# so there's no repo-creation step — just cleanup, commit, push.
#
# What it does:
#   1. Deletes the old VM-era *-workshop-host scripts (superseded by the new
#      *-workshop container-era scripts already sitting next to them)
#   2. Deletes .idea/ and a stray probe file from testing
#   3. Stages everything, commits
#   4. Pushes to origin/main
#
# It pauses before committing so you can eyeball git status first.

set -euo pipefail

echo "== Step 1: remove old VM-era (*-workshop-host) lifecycle scripts =="
rm -f track/track_scripts/setup-workshop-host
for d in 01-converting 02-child-workflows 03-parallel-activities 04-local-activities; do
  rm -f "track/$d/setup-workshop-host" \
        "track/$d/check-workshop-host" \
        "track/$d/solve-workshop-host" \
        "track/$d/cleanup-workshop-host"
done

echo "== Step 2: remove .idea/ and stray test file =="
rm -rf .idea
rm -f testfile_probe.txt

echo "== Step 3: make sure new scripts are executable =="
chmod +x track/track_scripts/setup-workshop
for d in 01-converting 02-child-workflows 03-parallel-activities 04-local-activities 05-saga; do
  chmod +x "track/$d/setup-workshop" "track/$d/check-workshop" \
           "track/$d/solve-workshop" "track/$d/cleanup-workshop"
done

echo
echo "== Review before committing =="
git status
echo
read -p "Looks right? Stage everything and commit? [y/N] " confirm
if [[ "$confirm" != "y" && "$confirm" != "Y" ]]; then
  echo "Stopped before committing. Nothing pushed."
  exit 0
fi

git add -A
git commit -m "Convert to containers, add Saga exercise, drop SailPoint branding

- track/config.yml: VM -> containers (ghcr.io/mmerrell/temporal-go-sandbox)
- Rename hostname workshop-host -> workshop everywhere
- Rename *-workshop-host lifecycle scripts -> *-workshop
- Switch VS Code/code-server tab -> native type:code editor tab
- Add bash,run fences + tab-switch buttons to all challenge instructions
- Rebuild docker/Dockerfile as a standalone sandbox image (Go 1.23 + Temporal CLI + baked exercises)
- Add Exercise 5: Saga pattern (compensating transactions), ported from temporal-java-hands-on
- Remove SailPoint references from README, exercise READMEs, assignment.md
- Fix stale lab1-/lab2- paths in README and build-image.yml"

echo
echo "== Step 4: push =="
git push origin main

echo
echo "Done. https://github.com/mmerrell/temporal-go-hands-on"
echo
echo "Still worth doing by hand:"
echo "  - go build ./... in each practice/ and solution/ dir under 5_saga/ — the"
echo "    sandbox that wrote this had no Go toolchain to verify it compiles."
echo "  - Retire the old repos once you're happy: mmerrell/sailpoint-instruqt-go"
echo "    (superseded by this one) — this script does not touch that repo."
echo "  - You can now safely ignore/delete ~/Projects/sailpoint-instruqt-go's copy"
echo "    of migrate-to-temporal-go-hands-on.sh — it targeted the wrong remote."
