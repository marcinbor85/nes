VERSION = $(shell git describe --always --dirty)

BINARY_NAME = nes

VERSION_VAR = github.com/marcinbor85/nes/config.Version

define compile
	CGO_ENABLED=0 GOARCH=amd64 GOOS=$(1) \
		go build -a \
		-ldflags="-X ${VERSION_VAR}=${VERSION}" \
		-installsuffix cgo -o $(2) $(3)
endef

build: main.go
	$(call compile,linux,${BINARY_NAME},$<)
	$(call compile,windows,${BINARY_NAME}.exe,$<)

clean:
	go clean
	rm -f ${BINARY_NAME}
	rm -f ${BINARY_NAME}.exe
