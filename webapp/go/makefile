build_linux:
	GOOS=linux GOARCH=amd64 go build -o app_linux

build_darwin:
	GOOS=darwin GOARCH=amd64 go build -o app_darwin

test:
	GO_ROUTER_ENABLE_LOGGING=1 go test -v $(glide novendor)

debug_run:
	GO_ROUTER_ENABLE_LOGGING=1 go run app.go
