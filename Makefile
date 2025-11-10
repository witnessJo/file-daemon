include .env

.PHONY: all
all: build run

.PHONY: build
build:
	go build -o ./build/file-sentinel ./cmd/file-sentinel

.PHONY: run
run:
	go run ./cmd/file-sentinel

.PHONY: test
test:
	go test ./...

.PHONY: minikube
minikube:
	minikube start -p harman --nodes 2 \
 --container-runtime=containerd --cni=cilium \
 --cpus=2 --memory=2048mb --disk-size=20gb

.PHONY: helm-postgres
helm-postgres:
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo update
	helm install postgres bitnami/postgresql -f values.yaml

.PHONY: ent-install
ent-install:
	go install entgo.io/ent/cmd/ent@latest

.PHONY: ent
ent:
	ent generate --feature sql/upsert ./ent/schema

.PHONY: docker
docker:
	docker build -t file-sentinel:latest .
