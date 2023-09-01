
IMAGE_NAME = averagemarcus/worst-webhook:latest

.PHONY: docker-build
docker-build:
	@docker build --platform linux/amd64 -t $(IMAGE_NAME) .

.PHONY: docker-publish
docker-publish:
	@docker push $(IMAGE_NAME)

.PHONY: deploy-pre-reqs
deploy-pre-reqs:
	@kubectl apply -f ./pre-reqs/

deploy-%: deploy-pre-reqs
	@kubectl apply -f ./scenarios/$*/manifests/
