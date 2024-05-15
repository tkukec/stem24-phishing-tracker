docker-build-dev:
	DOCKER_BUILDKIT=1 docker build --ssh default -t hrzg1do20.hr.asseco-see.local/voice/stem24-backend:dev .

docker-push-dev:
	docker login --username voice --password 4YUeP799cQB886Q hrepo.asseco-see.hr
	docker push hrzg1do20.hr.asseco-see.local/voice/stem24-backend:dev

docker-build-version:
	ssh-add
	DOCKER_BUILDKIT=1 docker build --ssh default -t hrzg1do20.hr.asseco-see.local/voice/stem24-backend:${v} .

docker-push-version:
	docker login --username voice --password 4YUeP799cQB886Q hrepo.asseco-see.hr
	docker push hrzg1do20.hr.asseco-see.local/voice/stem24-backend:${v}

delete-logs:
	rm -rf logs/*

docker-sip-restart:
	docker-compose up -d stem24-backend

dev-update: delete-logs docker-build-dev docker-sip-restart

generate-docs:
	swag init -g cmd/stem24-backend/main.go -o api/openapi --parseVendor --pd --parseDepth 1 -ot yaml
