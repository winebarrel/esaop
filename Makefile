DOCKER_REPO := public.ecr.aws/winebarrel/esaop

.PHONY: all
all: vet build

.PHONY: build
build:
	go build ./cmd/esaop

.PHONY: vet
vet:
	go vet ./...

.PHONY: clean
clean:
	rm -rf esaop esaop.exe

.PHONY: image
image:
	docker build . --platform linux/amd64 -t $(DOCKER_REPO)

.PHONY: push
push: image
	docker push $(DOCKER_REPO)
