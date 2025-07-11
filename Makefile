default: frontend_lint build

build:
	go build -ldflags "-X main.version=$(git describe --tags --always || echo dev)"

frontend_prod: export NODE_ENV=production
frontend_prod: frontend

frontend: node_modules
	corepack yarn@1 node ci/build.mjs

frontend_lint: node_modules
	corepack yarn@1 eslint --fix src

node_modules:
	corepack yarn@1 install --production=false --frozen-lockfile

go_test:
	go test -cover -v ./...
	golangci-lint run

helm_lint:
	helm lint charts/latestver

.PHONY: frontend

# --- Documentation

gendoc: .venv
	.venv/bin/python3 ci/gendoc.py $(shell grep -l '@module ' internal/fetcher/*.go) >docs/config.md
	git add docs/config.md

.venv:
	python -m venv .venv
	.venv/bin/pip install -r ci/requirements.txt
