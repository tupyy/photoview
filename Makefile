build-frontend:
	cd ./ui && npm install
.PHONY: build-frontend

run-backend:
	cd ./api && go run server.go
.PHONY: run-backend

run-frontend:
	cd ./ui && npm start
.PHONY: run-frontend

run: run-frontend run-frontend
.PHONY: run

