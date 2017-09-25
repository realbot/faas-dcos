TAG?=0.1

build:
	docker build --build-arg http_proxy=$(http_proxy) --build-arg https_proxy=$(https_proxy) -t realbot/faas-dcos:$(TAG) .

push:
	docker push realbot/faas-dcos:$(TAG)