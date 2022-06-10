all: docker-build docker-run

docker-build:
	DOCKER_BUILDKIT=1 docker build -t park .

docker-run:
	docker rm -f park
	docker run --memory 2G --log-opt max-size=5M --log-opt max-file=3 -p 5000:5000 -p 5432:5432 --name park -t park

mod:
	go mod tidy && go mod vendor && go install ./...
