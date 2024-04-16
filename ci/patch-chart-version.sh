#!/usr/bin/env bash
set -euo pipefail

chart_yaml=charts/latestver/Chart.yaml

# Patch latest App-Version into Chart
yq -iP ".appVersion = \"$(git describe --tags --abbrev=0)\"" ${chart_yaml}

# Validate there has been a change before patching the chart version
git diff --exit-code -- ${chart_yaml} >/dev/null && exit 0 || true

# There were changes, we need to patch the chart version
chart_ver=$(yq '.version' ${chart_yaml})
yq -iP ".version = \"$(semver -i minor ${chart_ver})\"" ${chart_yaml}

git add ${chart_yaml}
