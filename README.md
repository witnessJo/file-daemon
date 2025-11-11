---
author: witnessjo
---

# file-sentinel

Kubernetes DaemonSet application that monitors filesystem directories
and stores file metadata in PostgreSQL.

## Prerequisites

-   Go 1.21.5+
-   Docker
-   Minikube
-   Helm 3.x

## Quick Start

### 1. Create Minikube Cluster

``` bash
make minikube
```

Creates a 2-node cluster (harman, harman-m02) with containerd and
cilium.

### 2. Build Application Image

``` bash
make docker
```

Builds the Docker image and loads it into Minikube.

### 3. Deploy PostgreSQL Database

``` bash
make helm-postgres
```

Deploys PostgreSQL using Helm with Bitnami chart.

Database credentials (helm-postgres/values.yaml):

-   User: postgres
-   Password: postgres
-   Database: postgres
-   Port: 5432 (internal), 30432 (NodePort)

### 4. Deploy file-sentinel DaemonSet

``` bash
make helm-sentinel
```

This command:

-   Creates local ~/harman directory
-   Host-mounts path can be changed in Makefile `$(HOST_DIR_PATH)`
-   Starts minikube mount in background
-   Deploys DaemonSet to all nodes

## Verify Deployment

``` bash
# Check pods
kubectl get pods -l app=file-sentinel-file-sentinel

# View logs
kubectl logs -l app=file-sentinel-file-sentinel --tail=20

# Query database
kubectl exec helm-postgres-postgresql-0 -- \
  env PGPASSWORD=postgres psql -U postgres -d postgres \
  -c "SELECT * FROM file_infos ORDER BY created_at DESC LIMIT 5;"
```

## Configuration

Environment variables (helm-sentinel/values.yaml):

-   `TARGET_DIR_PATH`: Directory to monitor (default: /mnt/harman)
-   `MINUTE_CYCLE`: Scan interval in minutes (default: 1)
-   `LOG_LEVEL`: DEBUG, INFO, or ERROR (default: DEBUG)

Database credentials are stored in Kubernetes Secret (db-secret).

## Testing with Files

Add files to your local directory (visible on control plane node):

``` bash
echo "test" > ~/harman/test.txt
```

Add files to worker node:

``` bash
minikube ssh -p harman -n harman-m02 \
  "sudo mkdir -p /home/witnessjo/harman && \
   echo 'worker' | sudo tee /home/witnessjo/harman/worker.txt"
```

## Database Schema

``` sql
CREATE TABLE file_infos (
    id SERIAL PRIMARY KEY,
    node_name VARCHAR NOT NULL,
    mount_path VARCHAR NOT NULL,
    file_list JSONB,
    created_at TIMESTAMP NOT NULL
);
```

## Makefile Targets

| Target | Description |
|--------|-------------|
| `make minikube` | Create 2-node Minikube cluster |
| `make docker` | Build and load Docker image |
| `make helm-postgres` | Deploy PostgreSQL database |
| `make helm-sentinel` | Deploy file-sentinel DaemonSet |
| `make unmount` | Stop minikube mount process |
| `make build` | Build Go binary |
| `make test` | Run tests |

## Architecture

-   DaemonSet: One pod per node
-   Volume: hostPath (node-local storage)
-   Database: PostgreSQL with Ent ORM
-   Periodic Scan: time.Ticker (configurable interval)
-   Authentication: Kubernetes Secret

## Project Structure

``` example
file-daemon/
|-- cmd/file-sentinel/          # Application code
|-- ent/schema/                 # Ent ORM schema
|-- helm-postgres/              # PostgreSQL Helm chart
|-- helm-sentinel/              # DaemonSet Helm chart
|-- Dockerfile
|-- Makefile
+-- README.org
```
