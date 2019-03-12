# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=synflood
CMD_PATH=./cmd/$(BINARY_NAME)

build:
	$(GOBUILD) -o dist/$(BINARY_NAME) $(CMD_PATH)
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o dist/$(BINARY_NAME) $(CMD_PATH)
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)