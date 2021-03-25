VERSION = "0.0.1"

default:
	@echo "Helpers. Check Makefile for more instructions."
	echo "$(VERSION)"
	
watch:
	find -type f | egrep -i "*.go|*.ini" | entr -r go run *.go

build: clean-build amd64 arm
	@echo "Building amd64 and arm version"

clean-build:
	- rm release -Rf

amd64:
	mkdir -p release/linux-amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.ldVersion=$(VERSION)'" -o release/linux-amd64/esp32-iot -v -a *.go

arm:
	mkdir -p release/linux-arm
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -ldflags="-X 'main.ldVersion=$(VERSION)'"  -o release/linux-arm/esp32-iot -v -a *.go