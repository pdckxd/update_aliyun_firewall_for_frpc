all:
	go build -ldflags "-s -w"

amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w"

