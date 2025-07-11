#!/usr/bin/env bash
set -euo pipefail

chart_yaml=charts/latestver/Chart.yaml

# Patch latest App-Version into Chart
yq -iP ".appVersion = \"v${TAG_VERSION}\"" ${chart_yaml}
yq -iP ".version = \"${TAG_VERSION}\"" ${chart_yaml}

# Validate there has been a change before adding
git diff --exit-code -- ${chart_yaml} >/dev/null && exit 0 || true

git add ${chart_yaml}
