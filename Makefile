DOCKER_REPO := public.ecr.aws/winebarrel/openesa

.PHONY: all
all: vet build

.PHONY: build
build:
	go build ./cmd/openesa

.PHONY: vet
vet:
	go vet ./...

.PHONY: clean
clean:
	rm -rf openesa openesa.exe

.PHONY: image
image:
	docker build . -t $(DOCKER_REPO)

.PHONY: push
push: image
	docker push $(DOCKER_REPO)
