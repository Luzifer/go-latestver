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
