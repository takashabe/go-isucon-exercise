build_linux:
	GOOS=linux GOARCH=amd64 go build -o bin/app_linux

build_darwin:
	GOOS=darwin GOARCH=amd64 go build -o bin/app_darwin

test:
	go test -v

debug_run:
	GO_ROUTER_ENABLE_LOGGING=1 go run app.go
