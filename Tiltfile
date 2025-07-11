# Install Node deps on change of package.json
local_resource(
  'yarn',
  cmd='corepack yarn@1 install', # Not using the make target to edit the lockfile
  deps=['package.json'],
)

# Rebuild frontend if source files change
local_resource(
  'frontend',
  cmd='make frontend',
  deps=['src'],
  resource_deps=['yarn'],
)

# Update go.sum file and ensure modules are available
local_resource(
  'go-sum',
  cmd='go mod tidy',
  deps=['go.mod'],
)

# Rebuild and run Go webserver on code changes
local_resource(
  'server',
  cmd='make build',
  deps=[
    '*.go',
    'go.mod',
    'go.sum',
    'frontend',
    'internal',
    'pkg',
  ],
  ignore=['src'],
  serve_cmd='./go-latestver --listen=:64824',
  readiness_probe=probe(
    http_get=http_get_action(64824, path='/v1/healthz'),
    initial_delay_secs=1,
  ),
  resource_deps=[
    'frontend',
    'go-sum',
  ],
)
