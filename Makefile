go-get:
	go get
windows: go-get
	GOOS=windows GOARC=arm64 go build
linux: go-get
	GOOS=linux GOARC=arm64 go build
all: windows linux
