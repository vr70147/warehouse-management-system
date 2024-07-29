# Warehouse Management System

The Warehouse Management System is a microservice-based application designed to manage orders, inventory, shipping processes, customer management, user management, reporting, and account management. The system leverages Kafka for event-driven communication and uses various technologies to ensure scalability and maintainability.

## Table of Contents

- [Installation](#installation)
- [Environment Variables](#environment-variables)
- [Running the Application](#running-the-application)
- [API Documentation](#api-documentation)
- [Project Structure](#project-structure)
- [Services](#services)
  - [Order Service](#order-service)
  - [Inventory Service](#inventory-service)
  - [Shipping Service](#shipping-service)
  - [Customer Service](#customer-service)
  - [User Management Service](#user-management-service)
  - [Reporting and Analytics Service](#reporting-and-analytics-service)
  - [Account Management Service](#account-management-service)
- [Models](#models)
- [Handlers](#handlers)
- [Middleware](#middleware)
- [Kafka Integration](#kafka-integration)
- [Contributing](#contributing)
- [License](#license)

## Installation

### Clone the repository:

    git clone https://github.com/vr70147/warehouse-management-system.git
    cd warehouse-management-system

### Install dependencies for each service:

```bash
cd services/order-processing
go mod tidy

cd ../inventory-service
go mod tidy

cd ../shipping-service
go mod tidy

cd ../customer-service
go mod tidy

cd ../user-management-service
go mod tidy

cd ../reporting-analytics-service
go mod tidy

cd ../account-management-service
go mod tidy
```

## Ensure you have running instances of Kafka and Redis.

### Environment Variables

Create a .env file in the root directory of each service and add the following environment variables:

### Order Service

```bash
KAFKA_BROKERS=localhost:9092
ORDER_EVENT_TOPIC=order-events
INVENTORY_STATUS_TOPIC=inventory-status
LOW_STOCK_NOTIFICATION_TOPIC=low-stock-notification
SHIPPING_STATUS_TOPIC=shipping-status
REDIS_ADDR=localhost:6379
JWT_SECRET=your_jwt_secret
```

### Inventory Service

```bash
KAFKA_BROKERS=localhost:9092
INVENTORY_EVENT_TOPIC=inventory-events
ORDER_STATUS_TOPIC=order-status
REDIS_ADDR=localhost:6379
```

### Shipping Service

```bash
KAFKA_BROKERS=localhost:9092
SHIPPING_EVENT_TOPIC=shipping-events
ORDER_STATUS_TOPIC=order-status
REDIS_ADDR=localhost:6379
```

### Customer Service

```bash
KAFKA_BROKERS=localhost:9092
CUSTOMER_EVENT_TOPIC=customer-events
REDIS_ADDR=localhost:6379
```

### User Management Service

```bash
KAFKA_BROKERS=localhost:9092
USER_EVENT_TOPIC=user-events
REDIS_ADDR=localhost:6379
JWT_SECRET=your_jwt_secret
```

### Reporting and Analytics Service

```bash
KAFKA_BROKERS=localhost:9092
REPORTING_EVENT_TOPIC=reporting-events
REDIS_ADDR=localhost:6379
```

### Account Management Service

```bash
KAFKA_BROKERS=localhost:9092
ACCOUNT_EVENT_TOPIC=account-events
REDIS_ADDR=localhost:6379
```

## Running the Application

Initialize the environment variables and database connections for each service:

```bash
source .env
```

## Run each service:

### Order Service

```bash
cd services/order-processing
go run main.go
```

### Inventory Service

```bash
cd services/inventory-service
go run main.go
```

### Shipping Service

```bash
cd services/shipping-service
go run main.go
```

### Customer Service

```bash
cd services/customer-service
go run main.go
```

### User Management Service

```bash
cd services/user-management-service
go run main.go
```

### Reporting and Analytics Service

```bash
cd services/reporting-analytics-service
go run main.go
```

### Account Management Service

```bash
cd services/account-management-service
go run main.go
```

### API Documentation

Swagger is used to generate API documentation. Access the documentation for each service at:

- **Order Service:** http://localhost:8080/swagger/index.html
- **Inventory Service:** http://localhost:8081/swagger/index.html
- **Shipping Service:** http://localhost:8082/swagger/index.html
- **Customer Service:** http://localhost:8083/swagger/index.html
- **User Management Service:** http://localhost:8084/swagger/index.html
- **Reporting and Analytics Service:** http://localhost:8085/swagger/index.html
- **Account Management Service:** http://localhost:8086/swagger/index.html

## Project Structure

├── services
│ ├── order-processing
│ │ ├── internal
│ │ │ ├── api
│ │ │ │ └── handlers
│ │ │ │ └── OrdersHandler.go
│ │ │ ├── cache
│ │ │ │ └── redis.go
│ │ │ ├── initializers
│ │ │ │ ├── db.go
│ │ │ │ └── env.go
│ │ │ ├── kafka
│ │ │ │ ├── consumer.go
│ │ │ │ ├── producer.go
│ │ │ │ └── kafka.go
│ │ │ ├── model
│ │ │ │ ├── order.go
│ │ │ │ ├── customer.go
│ │ │ │ ├── user.go
│ │ │ │ ├── role.go
│ │ │ │ ├── department.go
│ │ │ │ └── event.go
│ │ │ ├── routes
│ │ │ │ └── routes.go
│ │ │ └── utils
│ │ │ └── notifications.go
│ │ ├── middleware
│ │ │ └── auth.go
│ │ ├── docs
│ │ │ └── swagger documentation files
│ │ ├── .env
│ │ ├── go.mod
│ │ ├── go.sum
│ │ └── main.go
│ ├── inventory-service
│ │ ├── internal
│ │ │ ├── api
│ │ │ │ └── handlers
│ │ │ ├── cache
│ │ │ ├── initializers
│ │ │ ├── kafka
│ │ │ ├── model
│ │ │ ├── routes
│ │ │ └── utils
│ │ ├── middleware
│ │ ├── docs
│ │ ├── .env
│ │ ├── go.mod
│ │ ├── go.sum
│ │ └── main.go
│ ├── shipping-service
│ │ ├── internal
│ │ │ ├── api
│ │ │ │ └── handlers
│ │ │ ├── cache
│ │ │ ├── initializers
│ │ │ ├── kafka
│ │ │ ├── model
│ │ │ ├── routes
│ │ │ └── utils
│ │ ├── middleware
│ │ ├── docs
│ │ ├── .env
│ │ ├── go.mod
│ │ ├── go.sum
│ │ └── main.go
│ ├── customer-service
│ │ ├── internal
│ │ │ ├── api
│ │ │ │ └── handlers
│ │ │ ├── cache
│ │ │ ├── initializers
│ │ │ ├── kafka
│ │ │ ├── model
│ │ │ ├── routes
│ │ │ └── utils
│ │ ├── middleware
│ │ ├── docs
│ │ ├── .env
│ │ ├── go.mod
│ │ ├── go.sum
│ │ └── main.go
│ ├── user-management-service
│ │ ├── internal
│ │ │ ├── api
│ │ │ │ └── handlers
│ │ │ ├── cache
│ │ │ ├── initializers
│ │ │ ├── kafka
│ │ │ ├── model
│ │ │ ├── routes
│ │ │ └── utils
│ │ ├── middleware
│ │ ├── docs
│ │ ├── .env
│ │ ├── go.mod
│ │ ├── go.sum
│ │ └── main.go
│ ├── reporting-analytics-service
│ │ ├── internal
│ │ │ ├── api
│ │ │ │ └── handlers
│ │ │ ├── cache
│ │ │ ├── initializers
│ │ │ ├── kafka
│ │ │ ├── model
│ │ │ ├── routes
│ │ │ └── utils
│ │ ├── middleware
│ │ ├── docs
│ │ ├── .env
│ │ ├── go.mod
│ │ ├── go.sum
│ │ └── main.go
│ ├── account-management-service
│ │ ├── internal
│ │ │ ├── api
│ │ │ │ └── handlers
│ │ │ ├── cache
│ │ │ ├── initializers
│ │ │ ├── kafka
│ │ │ ├── model
│ │ │ ├── routes
│ │ │ └── utils
│ │ ├── middleware
│ │ ├── docs
│ │ ├── .env
│ │ ├── go.mod
│ │ ├── go.sum
│ │ └── main.go
└── docker-compose.yml

## Services

### Order Service

Handles order processing, creation, updates, and status changes.

### Inventory Service

Manages inventory levels, updates, and low stock notifications.

### Shipping Service

Handles shipping status updates and notifications.

### Customer Service

Manages customer information and handles customer-related events.

### User Management Service

Handles user authentication, authorization, and management.

### Reporting and Analytics Service

Generates reports and analytics based on order, inventory, and shipping data.

### Account Management Service

Manages account information and handles account-related events.

## Models

Each service has its own set of models representing the data structures used within the service. Refer to the internal/model directory in each service for detailed definitions.

## Handlers

Handlers manage the API endpoints for each service. They include functionalities like:

- Creating, updating, and retrieving orders.
- Managing inventory levels and statuses.
- Handling shipping updates and notifications.
- Managing customer and user data.
- Generating reports and handling account management.

Refer to the internal/api/handlers directory in each service for detailed implementations.

## Middleware

Middleware functions include:

- CORS settings.
- JWT authentication and authorization.

Refer to the middleware directory in each service for detailed implementations.

## Kafka Integration

Kafka is used for event-driven communication between services. Each service has Kafka producers and consumers to handle relevant events.

### Order Service Kafka Activity

- **Producers:**
  - PublishOrderEvent: Publishes order events.
- **Consumers:**
  - ConsumerOrderEvent: Consumes order events.
  - ConsumerInventoryStatus: Consumes inventory status updates.
  - ConsumerShippingStatus: Consumes shipping status updates.

### Inventory Service Kafka Activity

- **Consumers:**
  - ConsumerOrderStatus: Consumes order status updates.

### Shipping Service Kafka Activity

- **Consumers:**
  - ConsumerOrderStatus: Consumes order status updates.

### Customer Service Kafka Activity

- **Producers:**
  - PublishCustomerEvent: Publishes customer events.
- **Consumers:**
  - ConsumerCustomerEvent: Consumes customer events.

### Reporting and Analytics Service Kafka Activity

- **Producers:**
  - PublishReportingEvent: Publishes reporting events.
- **Consumers:**
  - ConsumerReportingEvent: Consumes reporting events.

### Account Management Service Kafka Activity

- **Producers:**
  - PublishAccountEvent: Publishes account events.
- **Consumers:**
  - ConsumerAccountEvent: Consumes account events.

### Contributing

Contributions are welcome! Please submit a pull request or open an issue to discuss any changes.

### License

This project is licensed under the MIT License.
