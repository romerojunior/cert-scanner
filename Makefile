BINARY=scanner

.PHONY: pre
pre:
	go get github.com/golang/dep/cmd/dep
	dep ensure

.PHONY: build
build:
	env GOOS=linux GOARCH=386 go build -o $(BINARY)-linux-i386

.PHONY: test
test: pre
	go test -v
