.PHONY:
.SILENT:

build:
	go build -o ./.bin/linkShortener cmd/main.go

run: build
	./.bin/linkShortener