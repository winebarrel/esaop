.PHONY: all
all: vet build

.PHONY: build
build:
	go build ./cmd/open-esa

.PHONY: vet
vet:
	go vet ./...

.PHONY: clean
clean:
	rm -rf open-esa open-esa.exe

