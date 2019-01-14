.PHONY: default
default: hdhr-signal-meter

SOURCES := $(shell find . -name '*.go')

hdhr-signal-meter: $(SOURCES)
	CGO_ENABLED=0 go build ${LDFLAGS}

.PHONY: restore
restore:
	dep ensure

.PHONY: test
test:
	go test -v ./...

.PHONY: testv
testv:
	go test -v ./... -args -v 6 -logtostderr true

.PHONY: all
all: restore hdhr-signal-meter test
