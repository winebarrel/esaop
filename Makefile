DOCKER_REPO := public.ecr.aws/winebarrel/open-esa

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

.PHONY: image
image:
  docker build . -t $(DOCKER_REPO)

.PHONY: push
push: image
  docker push $(DOCKER_REPO)
