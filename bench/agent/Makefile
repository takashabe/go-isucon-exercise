.PHONY: build test vet lint clean

build: cmd/agent/main.go
	cd cmd/agent && go build -a

test:
	go test -v

vet:
	go vet

lint:
	golint

clean:
	rm cmd/agent/agent
