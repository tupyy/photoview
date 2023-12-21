REACT_APP_API_ENDPOINT = 
UI_PUBLIC_URL=/
IMG ?= photoview
QUAY = quay.io/ctupangiu
GIT_COMMIT=$(shell git rev-list -1 HEAD --abbrev-commit)
TAG = v1.0-${GIT_COMMIT}
VERSION=v1.0
CONTAINER_TOOL=podman

build-frontend:
	cd ./ui && npm run build
.PHONY: build-frontend

build-backend: 
	cd ./api && go build -o ../photoview server.go

build: build-frontend build-backend
.PHONY: build

build-docker:
	$(CONTAINER_TOOL) build --ulimit=nofile=5000:5000 --build-arg=IMAGE_NAME="${QUAY}/${IMG}:${TAG}"  --build-arg=REACT_APP_API_ENDPOINT="${REACT_APP_API_ENDPOINT}" --build-arg=GIT_COMMIT="${GIT_COMMIT}" -t ${IMG}:latest .

docker-push: ## Push docker image with the manager.
	$(CONTAINER_TOOL) tag ${IMG}:latest $(QUAY)/${IMG}:${TAG}
	$(CONTAINER_TOOL) push $(QUAY)/${IMG}:${TAG}
.PHONY: docker-push

run-backend:
	cd ./api && go run server.go
.PHONY: run-backend

run-frontend:
	cd ./ui && npm start
.PHONY: run-frontend

run: run-frontend run-frontend
.PHONY: run

