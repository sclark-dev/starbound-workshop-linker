
build:
	go build -ldflags="-s -w" -o="dist/"

all: linux windows darwin zip

linux:
	GOARCH="amd64" GOOS="linux" go build -ldflags="-s -w" -o="dist/linux/swl"

windows:
	GOARCH="amd64" GOOS="windows" go build -ldflags="-s -w" -o="dist/windows/swl.exe"

darwin:
	GOARCH="amd64" GOOS="darwin" go build -ldflags="-s -w" -o="dist/darwin/swl"

zip:
	rm -f dist/starbound-workshop-linker.zip
	zip dist/starbound-workshop-linker.zip README.md LICENSE
	cd dist && zip -r starbound-workshop-linker.zip darwin linux windows