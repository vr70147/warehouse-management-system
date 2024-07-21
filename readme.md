# Warehouse Management System

## Introduction

Welcome to the Warehouse Management System project. This system is designed to efficiently manage warehouse operations including accounts, orders, inventory, and shipping. The system is built using Go with a microservices architecture and integrates Kafka for messaging.

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Architecture](#architecture)
- [Installation](#installation)
- [Usage](#usage)
- [API Documentation](#api-documentation)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Accounts Management**: Manage user accounts, roles, and permissions.
- **Orders Management**: Process and track orders.
- **Inventory Management**: Track and manage inventory levels.
- **Shipping Management**: Handle shipping and logistics.
- **Kafka Integration**: Use Kafka for event-driven architecture and messaging between services.
- **Docker Integration**: Use Docker for containerization and deployment.
- **Kubernetes Integration**: Use Kubernetes for orchestration and management.
- **PostgreSQL Integration**: Use PostgreSQL for database management.
- **Redis Integration**: Use Redis for caching and message queuing.

## Architecture

The system follows a microservices architecture with the following services:

- **Accounts Management Service**: Handles user accounts, roles, and permissions.
- **Order Processing Service**: Manages order processing and status updates.
- **Inventory Management Service**: Keeps track of inventory levels and updates.
- **Shipping Management Service**: Manages shipping and delivery of orders.

## System Structure

```sh
warehouse-management-system/
├── docker-compose.yml
├── get_helm.sh
├── k8s/
│ ├── accounts-management/
│ │ ├── configmap.yaml
│ │ ├── deployment.yaml
│ │ ├── ingress.yaml
│ │ └── service.yaml
│ ├── ingress-nginx/
│ │ └── ingress-nginx.yaml
│ ├── inventory-management/
│ │ ├── configmap.yaml
│ │ ├── deployment.yaml
│ │ ├── ingress.yaml
│ │ └── service.yaml
│ ├── order-processing/
│ │ ├── configmap.yaml
│ │ ├── deployment.yaml
│ │ ├── ingress.yaml
│ │ └── service.yaml
│ ├── postgres/
│ │ ├── secret.yaml
│ │ ├── user-management.yaml
│ │ └── values.yaml
│ ├── reporting-analytics/
│ │ ├── configmap.yaml
│ │ ├── deployment.yaml
│ │ ├── ingress.yaml
│ │ └── service.yaml
│ ├── shipping-receiving/
│ │ ├── configmap.yaml
│ │ ├── deployment.yaml
│ │ ├── ingress.yaml
│ │ └── service.yaml
│ └── user-management/
│ ├── configmap.yaml
│ ├── deployment.yaml
│ ├── ingress.yaml
│ ├── secret.yaml
│ └── service.yaml
├── readme.md
├── services/
│ ├── accounts-management/
│ │ ├── .env
│ │ ├── Dockerfile
│ │ ├── accounts-management
│ │ ├── go.mod
│ │ ├── go.sum
│ │ ├── internal/
│ │ │ ├── api/
│ │ │ │ ├── handlers/
│ │ │ │ │ └── accountHandler.go
│ │ │ │ └── routes/
│ │ │ │ └── routes.go
│ │ │ ├── cache/
│ │ │ │ └── redis_client.go
│ │ │ ├── initializers/
│ │ │ │ ├── db.go
│ │ │ │ └── loadEnvVariables.go
│ │ │ ├── middleware/
│ │ │ │ └── middleware.go
│ │ │ ├── model/
│ │ │ │ └── model.go
│ │ │ └── tests/
│ │ │ └── routes_test.go
│ │ └── main.go
│ ├── inventory-management/
│ │ ├── Dockerfile
│ │ ├── docs/
│ │ │ ├── docs.go
│ │ │ ├── swagger.json
│ │ │ └── swagger.yaml
│ │ ├── go.mod
│ │ ├── go.sum
│ │ ├── handlers.go
│ │ ├── internal/
│ │ │ ├── api/
│ │ │ │ ├── handlers/
│ │ │ │ │ ├── categoriesHandler.go
│ │ │ │ │ ├── productsHandler.go
│ │ │ │ │ ├── stocksHandler.go
│ │ │ │ │ └── suppliersHandler.go
│ │ │ │ └── routes/
│ │ │ │ └── routes.go
│ │ │ ├── cache/
│ │ │ │ └── redis_client.go
│ │ │ ├── initializers/
│ │ │ │ ├── db.go
│ │ │ │ └── loadEnvVariables.go
│ │ │ ├── middleware/
│ │ │ │ └── middleware.go
│ │ │ ├── model/
│ │ │ │ └── model.go
│ │ │ └── tests/
│ │ │ └── routes_test.go
│ │ ├── kafka/
│ │ │ └── kafka_new_order_consumer.go
│ │ └── main.go
│ ├── order-processing/
│ │ ├── Dockerfile
│ │ ├── docs/
│ │ │ ├── docs.go
│ │ │ ├── swagger.json
│ │ │ └── swagger.yaml
│ │ ├── go.mod
│ │ ├── go.sum
│ │ ├── internal/
│ │ │ ├── api/
│ │ │ │ ├── handlers/
│ │ │ │ │ └── OrdersHandler.go
│ │ │ │ ├── middleware/
│ │ │ │ │ └── middleware.go
│ │ │ │ └── routes/
│ │ │ │ └── routes.go
│ │ │ ├── cache/
│ │ │ │ └── redis_client.go
│ │ │ ├── initializers/
│ │ │ │ ├── db.go
│ │ │ │ └── loadEnvVariables.go
│ │ │ ├── kafka/
│ │ │ │ └── kafka.go
│ │ │ ├── model/
│ │ │ │ └── model.go
│ │ │ └── tests/
│ │ │ └── orders_test.go
│ │ └── main.go
│ ├── reporting-analytics/
│ │ ├── Dockerfile
│ │ ├── go.mod
│ │ ├── go.sum
│ │ ├── internal/
│ │ │ ├── api/
│ │ │ │ ├── handlers/
│ │ │ │ │ └── analyticsHandler.go
│ │ │ │ ├── middleware/
│ │ │ │ │ └── middleware.go
│ │ │ │ └── routes/
│ │ │ │ └── routes.go
│ │ │ ├── cache/
│ │ │ │ └── redis_client.go
│ │ │ ├── initializers/
│ │ │ │ ├── db.go
│ │ │ │ └── loadEnvVariables.go
│ │ │ ├── model/
│ │ │ │ └── model.go
│ │ │ └── tests/
│ │ │ └── analytics_test.go
│ │ ├── kafka/
│ │ │ ├── consumer.go
│ │ │ └── producer.go
│ │ └── main.go
│ ├── shipping-receiving/
│ │ ├── Dockerfile
│ │ ├── docs/
│ │ │ ├── docs.go
│ │ │ ├── swagger.json
│ │ │ └── swagger.yaml
│ │ ├── go.mod
│ │ ├── go.sum
│ │ ├── internal/
│ │ │ ├── api/
│ │ │ │ ├── handlers/
│ │ │ │ │ ├── receivingHandler.go
│ │ │ │ │ ├── shippingHandler.go
│ │ │ │ │ └── test.db
│ │ │ │ ├── middleware/
│ │ │ │ │ └── middleware.go
│ │ │ │ └── routes/
│ │ │ │ └── routes.go
│ │ │ ├── cache/
│ │ │ │ └── redis_client.go
│ │ │ ├── initializers/
│ │ │ │ ├── db.go
│ │ │ │ └── loadEnvVariables.go
│ │ │ ├── model/
│ │ │ │ └── model.go
│ │ │ └── tests/
│ │ │ └── receiving_test.go
│ │ └── main.go
│ └── user-management/
│ ├── .env
│ ├── Dockerfile
│ ├── docs/
│ │ ├── docs.go
│ │ ├── swagger.json
│ │ └── swagger.yaml
│ ├── go.mod
│ ├── go.sum
│ ├── internal/
│ │ ├── api/
│ │ │ ├── handlers/
│ │ │ │ ├── rolesHandler.go
│ │ │ │ ├── userHandlers.go
│ │ │ │ ├── userHandlers_test.go
│ │ │ │ └── user_handlers_test.go
│ │ │ ├── middleware/
│ │ │ │ └── middleware.go
│ │ │ └── routes/
│ │ │ └── routes.go
│ │ ├── cache/
│ │ │ └── redis_client.go
│ │ ├── initializers/
│ │ │ ├── db.go
│ │ │ └── loadEnvVariables.go
│ │ ├── middleware/
│ │ │ └── middleware.go
│ │ ├── model/
│ │ │ └── model.go
│ │ ├── services/
│ │ │ └── userService.go
│ │ └── tests/
│ │ └── routes_test.go
│ ├── kafka/
│ │ └── kafka.go
│ └── main.go
└── readme.md
```

## Installation

To set up the project locally, follow these steps:

1. **Clone the repository**:

   ```sh
   git clone https://github.com/yourusername/warehouse-management-system.git
   cd warehouse-management-system
   ```

2. **Set up environment variables**:

   - Create a `.env` file in the root directory and add the necessary environment variables:

     ```env
     JWT_SECRET=your-jwt-secret
     KAFKA_BROKERS=your-kafka-brokers

     ```

3. **Install dependencies**:

   ```sh
   go mod tidy

   ```

4. **Run the services**:

   - You can run each service individually. For example, to run the Accounts Management Service:

     ```sh
     cd services/accounts-management
     go run main.go

     ```

## Usage

Each service exposes a set of RESTful API endpoints. You can interact with these endpoints using tools like `curl` or Postman.

## API Documentation

### Accounts Management Service

- **Create Account:** POST /signup
- **Login: POST /login**
- **Logout: POST /logout**
- **Get Users: GET /users**
- **Update User: PUT /users/:id**
- **Soft Delete User: DELETE /users/:id**
- **Hard Delete User: DELETE /users/hard/:id**
- **Recover User: POST /users/:id/recover**

### Order Processing Service

- **Create Order: POST /orders**
- **Get Orders: GET /orders**
- **Update Order: PUT /orders/:id**
- **Soft Delete Order: DELETE /orders/:id**
- **Hard Delete Order: DELETE /orders/:id/hard**
- **Recover Order: POST /orders/:id/recover**

### Inventory Management Service

    Get Inventory: GET /inventory
    Update Inventory: PUT /inventory/:id

### Shipping Management Service

    Create Shipping: POST /shipping
    Get Shippings: GET /shipping
    Update Shipping: PUT /shipping/:id

## Testing

### To run the tests for the project:

    Run Unit Tests:

```sh

go test ./...
```

## Run End-to-End Tests:

    Ensure Kafka and other dependencies are running.
    Run the tests:

```sh

        go test -tags=e2e ./tests
```

## Contributing

We welcome contributions! Please follow these steps to contribute:

    Fork the repository.
    Create a new branch for your feature or bugfix.

```sh

git checkout -b feature/your-feature-name

```

Commit your changes.

```sh

git commit -m "Description of your changes"

```

Push to your branch.

```sh

git push origin feature/your-feature-name

```

Create a Pull Request.
