build-all:build-proxy

build-proxy:
	env CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -i -o bin/saber-proxy-docker ./cmd/proxy