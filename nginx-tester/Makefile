all: push

TAG ?= 1.1
PREFIX ?= bprashanth/nginx-tester

container:
	docker build -t $(PREFIX):$(TAG) .

push: container
	docker push $(PREFIX):$(TAG)

clean:
	docker rmi -f $(PREFIX):$(TAG) || true

.PHONY: container push
