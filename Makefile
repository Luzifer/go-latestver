default: lint build

build: node_modules
	node ci/build.mjs

lint: node_modules
	./node_modules/.bin/eslint \
		--ext .js,.vue \
		--fix \
		src

node_modules:
	npm ci

go_test:
	go test -cover -v ./...
	golangci-lint run

# --- Documentation

gendoc: .venv
	.venv/bin/python3 ci/gendoc.py $(shell grep -l '@module ' internal/fetcher/*.go) >docs/config.md
	git add docs/config.md

.venv:
	python -m venv .venv
	.venv/bin/pip install -r ci/requirements.txt
