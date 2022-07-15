
build:
	go build -ldflags="-s -w" -o="dist/"

all: linux windows darwin

linux:
	GOARCH="amd64" GOOS="linux" go build -ldflags="-s -w" -o="dist/linux/swl"

windows:
	GOARCH="amd64" GOOS="windows" go build -ldflags="-s -w" -o="dist/windows/"

darwin:
	GOARCH="amd64" GOOS="darwin" go build -ldflags="-s -w" -o="dist/darwin/"
