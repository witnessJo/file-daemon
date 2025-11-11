include .env

HOST_DIR_PATH ?= $(HOME)/harman

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
	minikube delete -p harman || true
	minikube start -p harman --nodes 2 \
 --container-runtime=containerd --cni=cilium \
 --cpus=2 --memory=2048mb --disk-size=20gb

.PHONY: helm-postgres
helm-postgres:
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo update
	helm dependency update helm-postgres
	helm upgrade --install helm-postgres helm-postgres

# download ent-go mandatory binaries
.PHONY: ent-install
ent-install:
	go install entgo.io/ent/cmd/ent@latest

# build ent-go ORM
.PHONY: ent
ent:
	ent generate --feature sql/upsert ./ent/schema

.PHONY: docker
docker:
	docker build -t file-sentinel:latest .
	minikube image load file-sentinel:latest -p harman

.PHONY: mount
mount:
	@echo "Mounting local directory to all Minikube nodes..."
	@echo "Keep this running in a separate terminal!"
	@mkdir -p $(HOST_DIR_PATH)
	minikube mount $(HOST_DIR_PATH):$(HOST_DIR_PATH) -p harman

.PHONY: helm-sentinel
helm-sentinel:
	@mkdir -p $(HOST_DIR_PATH)
	@echo "Stopping any existing minikube mount..."
	@-pgrep -f "minikube mount" | xargs -r kill 2>/dev/null || true
	@sleep 1
	@echo "Starting minikube mount in background..."
	@nohup minikube mount $(HOST_DIR_PATH):$(HOST_DIR_PATH) -p harman > /tmp/minikube-mount.log 2>&1 &
	@sleep 3
	@echo "Deploying helm chart..."
	@helm dependency update helm-sentinel
	@helm upgrade --install file-sentinel helm-sentinel \
	--set app.hostDirPath=$(HOST_DIR_PATH)
	@echo ""
	@echo "✓ Deployment complete!"
	@echo "Mount running in background (see /tmp/minikube-mount.log for logs)"
	@echo "To stop mount: make unmount"
