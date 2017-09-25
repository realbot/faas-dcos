TAG?=latest

build:
	docker build --build-arg http_proxy=$(http_proxy) --build-arg https_proxy=$(https_proxy) -f Dockerfile.alertmanager.dcos -t realbot/alertmanager:$(TAG) .
	docker build --build-arg http_proxy=$(http_proxy) --build-arg https_proxy=$(https_proxy) -f Dockerfile.prometheus.dcos -t realbot/prometheus:$(TAG) .

push:
	docker push realbot/alertmanager:$(TAG)
	docker push realbot/prometheus:$(TAG)