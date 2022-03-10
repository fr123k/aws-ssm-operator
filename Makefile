.PHONY: build

IMAGE?=fr123k/aws-ssm-operator
VERSION?=v0.3.2

build:
	docker build -t ${IMAGE}:${VERSION} --platform linux/amd64 ./

push: build
	docker push ${IMAGE}:${VERSION}
