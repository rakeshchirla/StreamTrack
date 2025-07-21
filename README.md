# StreamTrack 🚀
StreamTrack is a cloud-native, real-time event tracking system designed to ingest, process, and analyze high-volume data streams. It's built on a modern, decoupled microservices architecture and utilizes a professional-grade technology stack, making it a powerful foundation for any data analytics platform.

This project demonstrates a complete end-to-end workflow, from local development with Docker Compose to a fully automated CI/CD pipeline deploying the application to a cloud-based Kubernetes cluster.

## Core Architecture 🏗️
The system is designed around an event-driven architecture to ensure high throughput, low latency, and resilience. Each component is decoupled and serves a specific purpose.

    graph TD
        subgraph Client
            User[👩‍💻 User/Client]
        end

        subgraph "Kubernetes Cluster (Cloud)"
            subgraph "Application"
                API[/"API Service (Go)"/]
                Worker[⚙️ Worker Service (Go)]
            end
    
            subgraph "Data Infrastructure"
                Kafka[<--> Kafka Topic]
                ClickHouse[(ClickHouse DB)]
            end
    
            User -- "1. POST /track (JSON Event)" --> API
            API -- "2. Publish Event" --> Kafka
            Kafka -- "3. Consume Event" --> Worker
            Worker -- "4. Persist Data" --> ClickHouse
            User -- "5. GET /activities" --> API
            API -- "6. Query Data" --> ClickHouse
        end

### How It Works
Ingestion (API Service): A lightweight service written in Go exposes an HTTP endpoint (/track). It receives activity events as JSON payloads. Its only job is to validate the request and immediately publish it to a Kafka topic. This makes the API incredibly fast and prevents it from getting blocked by database operations.

Decoupling (Kafka): Apache Kafka acts as the central nervous system of the application. It's a distributed message bus that buffers all incoming events. This decoupling is crucial: if the worker service or database goes down, events are safely stored in Kafka, ready to be processed when the services come back online.

Processing (Worker Service): A separate Go microservice acts as a dedicated consumer. It continuously listens to the Kafka topic, reads events in batches, and is responsible for persisting them into the database. This isolates the data-writing logic from the client-facing API.

Storage (ClickHouse): We use ClickHouse, a high-performance, column-oriented OLAP (Online Analytical Processing) database. It is specifically designed to ingest and query massive amounts of event or log data with sub-second latency, making it the perfect choice for an analytics platform.

Querying (API Service): The API service also exposes a GET /activities endpoint. When a user requests data, the API directly queries the ClickHouse database to fetch and return the results.

## Technology Stack 🛠️
Category

Technology

Purpose

Backend

Go

For building high-performance, concurrent microservices.

Messaging

Apache Kafka

Resilient, high-throughput message bus for decoupling services.

Database

ClickHouse

High-performance OLAP database for real-time analytics.

Containerization

Docker

Containerizing all applications and services.

Orchestration

Kubernetes (K8s)

Production-grade deployment, scaling, and management.

Infrastructure

Terraform

Infrastructure as Code (IaC) for provisioning the K8s cluster.

CI/CD

Jenkins

Automating the entire build, push, and deploy pipeline.

## Getting Started
There are two ways to run this project: locally for development and in the cloud for a production-like setup.

### Local Development (with Docker Compose)
This is the quickest way to get the application running on your machine.

Prerequisites:

Docker

Docker Compose

Steps:

Clone the repository:

git clone https://github.com/your-username/your-repo-name.git
cd your-repo-name

Build and start all services:

docker-compose up --build

Open a new terminal and test the endpoints:

Track a new event:

curl -X POST -H "Content-Type: application/json" \
-d '{"user_id": "local_user", "action": "app_started"}' \
http://localhost:8080/track

View all events:

curl http://localhost:8080/activities

### Cloud Deployment (Terraform & Kubernetes)
This deploys the application to a cloud-native environment managed by Kubernetes.

Prerequisites:

A cloud provider account (e.g., Google Cloud, AWS)

Terraform CLI

kubectl CLI

Helm CLI

Steps:

Provision the Infrastructure:
Navigate to the terraform directory and use Terraform to create the Kubernetes cluster.

cd terraform
terraform init
terraform apply

This will build the cluster and output a command to configure kubectl.

Deploy Dependencies with Helm:
Helm is the package manager for Kubernetes. Use it to deploy production-ready instances of Kafka and ClickHouse.

# Example for Kafka
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install kafka bitnami/kafka --namespace kafka --create-namespace

# Example for ClickHouse
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install clickhouse bitnami/clickhouse --namespace clickhouse --create-namespace

Note: You will need to update the service names in the Kubernetes deployment files to match the names created by these Helm charts.

Deploy the Application:
Navigate to the k8s directory and apply the application manifests.

cd ../k8s
kubectl apply -f .

This will create the deployments and services for the api and worker. Find the public IP address of the api-service to interact with your deployed application.

## CI/CD Pipeline with Jenkins автоматизация
This project includes a Jenkinsfile that defines a complete, automated pipeline.

Pipeline Stages:

Checkout: Clones the source code from the Git repository.

Build API Image: Compiles the api Go application, builds its Docker image, and tags it with the build number.

Build Worker Image: Does the same for the worker application.

Push to Registry: Pushes both newly built images to a container registry (like Docker Hub or Google Container Registry).

Deploy to Kubernetes: Uses kubectl to trigger a rolling update for the api and worker deployments in the Kubernetes cluster, deploying the new images with zero downtime.

## Future Scope & Improvements 🔮
This project provides a solid foundation. Here are some exciting features that could be added:

📈 Monitoring & Alerting: Integrate Prometheus to scrape custom metrics from the Go services and Grafana to build real-time dashboards for system health and performance.

🌐 Web UI Dashboard: Build a frontend application (e.g., using React or Vue.js) to provide a user-friendly interface for viewing and querying activities.

🔍 Enhanced API: Add more powerful querying capabilities to the API, such as filtering by user_id, date ranges, and performing aggregations.

🔐 Authentication & Authorization: Secure the API endpoints using a method like JWT or OAuth2.

📜 Schema Enforcement: Integrate a Schema Registry (e.g., Confluent Schema Registry) to enforce a strict data structure for events published to Kafka, preventing data quality issues.

⚡ Performance Tuning: Conduct load testing to identify bottlenecks and tune the Go services, Kafka, and ClickHouse for even higher throughput.
