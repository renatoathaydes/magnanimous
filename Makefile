.PHONY: all
all: test install

ifeq ($(GOPATH),)
  GOPATH=~/go
endif

# runs all tests and benchmarks
.PHONY: bench
bench: test
	cd tests && go test -bench=.

# runs all tests
.PHONY: test
test:
	go test ./...

# installs go-hash
.PHONY: install
install:
	go install

# builds the website for local deployment
.PHONY: website
website: install
	magnanimous -globalctx=_local_global_context website

# builds the website for GitHub deployment
.PHONY: website-github
website-github: install
	magnanimous website

# serve the website locally
.PHONY: serve
serve: website
	go get -u github.com/vwochnik/gost
	gost website/target

# deploy to GitHub
.PHONY: deploy
deploy: website-github
	./deploy-to-github.sh

# build a smaller executable without symbols and debug info for all supported OSs and ARCHs
.PHONY: release release-linux release-windows release-darwin

release-linux:
	env GOOS=linux env GOARCH=amd64 go build -ldflags "-s -w" -o releases/magnanimous-linux-amd64
	env GOOS=linux env GOARCH=386 go build -ldflags "-s -w" -o releases/magnanimous-linux-386

release-windows:
	env GOOS=windows env GOARCH=amd64 go build -ldflags "-s -w" -o releases/magnanimous-windows-amd64
	env GOOS=windows env GOARCH=386 go build -ldflags "-s -w" -o releases/magnanimous-windows-386

release-darwin:
	env GOOS=darwin env GOARCH=amd64 go build -ldflags "-s -w" -o releases/magnanimous-darwin-amd64	

release: test release-linux release-windows release-darwin

# clean build artifacts, i.e. everything that is not source code.
# Does not remove the installed binary.
.PHONY: clean
clean:
	rm -f magnanimous
	rm -rf releases
	rm -rf website/target

# uninstall the binary
.PHONY: uninstall
uninstall:
	rm -f $(GOPATH)/bin/magnanimous
