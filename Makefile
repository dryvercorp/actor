all: go-get go-generate

go-get:
	@go get golang.org/x/tools/cmd/stringer

go-generate:
	@go generate
