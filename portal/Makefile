SUBPACKAGES := $(shell go list ./... | grep -v /vendor/)
SHOW_ENV := $(shell env | grep GO_MESSAGE_QUEUE)

.PHONY: build test vet lint clean

build: cmd/portal/main.go
	cd cmd/portal && go build -a

start:
	cd cmd/portal && ./portal

test:
	$(SHOW_ENV)
	go test -v $(SUBPACKAGES)

vet:
	go vet $(SUBPACKAGES)

lint:
	golint $(SUBPACKAGES)

clean:
	rm cmd/portal/portal
