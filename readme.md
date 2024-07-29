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
- [Database](#database)
  - [Postgres](#postgres)
  - [Redis](#redis)
- [Models](#models)
- [Handlers](#handlers)
- [Middleware](#middleware)
- [Kafka Integration](#kafka-integration)
- [Contributing](#contributing)
- [License](#license)

## Installation

1. Clone the repository:

   git clone https://github.com/vr70147/warehouse-management-system.git
   cd warehouse-management-system

2. Install dependencies for each service:

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

3. Ensure you have running instances of Kafka, Redis, and PostgreSQL.

4. Set up PostgreSQL databases:

Install PostgreSQL (if not already installed):

For Debian-based systems (Ubuntu):

```bash
sudo apt-get update
sudo apt-get install postgresql postgresql-contrib
```

For Red Hat-based systems (Fedora, CentOS):

```bash
sudo yum install postgresql-server postgresql-contrib
sudo postgresql-setup initdb
sudo systemctl start postgresql
```

### Create databases and users for each service:

Switch to the postgres user:

```bash
sudo -i -u postgres
```

### Create databases:

```bash
createdb order_db
createdb inventory_db
createdb shipping_db
createdb customer_db
createdb user_db
createdb reporting_db
createdb account_db
```

### Create users and grant privileges:

```bash
psql

CREATE USER order_user WITH ENCRYPTED PASSWORD 'your_postgres_password';
CREATE USER inventory_user WITH ENCRYPTED PASSWORD 'your_postgres_password';
CREATE USER shipping_user WITH ENCRYPTED PASSWORD 'your_postgres_password';
CREATE USER customer_user WITH ENCRYPTED PASSWORD 'your_postgres_password';
CREATE USER user_user WITH ENCRYPTED PASSWORD 'your_postgres_password';
CREATE USER reporting_user WITH ENCRYPTED PASSWORD 'your_postgres_password';
CREATE USER account_user WITH ENCRYPTED PASSWORD 'your_postgres_password';

GRANT ALL PRIVILEGES ON DATABASE order_db TO order_user;
GRANT ALL PRIVILEGES ON DATABASE inventory_db TO inventory_user;
GRANT ALL PRIVILEGES ON DATABASE shipping_db TO shipping_user;
GRANT ALL PRIVILEGES ON DATABASE customer_db TO customer_user;
GRANT ALL PRIVILEGES ON DATABASE user_db TO user_user;
GRANT ALL PRIVILEGES ON DATABASE reporting_db TO reporting_user;
GRANT ALL PRIVILEGES ON DATABASE account_db TO account_user;

\q

```

Exit the postgres user session:

```bash
exit
```

### Environment Variables

Create a .env file in the root directory of each service and add the following environment variables:

### Order Service

```bash
PORT=8083
KAFKA_BROKERS=localhost:9092
ORDER_EVENT_TOPIC=order-events
INVENTORY_STATUS_TOPIC=inventory-status
SHIPPING_STATUS_TOPIC=shipping-status
LOW_STOCK_NOTIFICATION_TOPIC=low-stock-notifications
USER_SERVICE_URL=http://localhost:8080
CUSTOMER_SERVICE_URL=http://localhost:8087
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=<your_redis_password>
POSTGRES_USER=<your_postgres_user>
POSTGRES_PASSWORD=<your_postgres_password>
POSTGRES_DB=order_processing
POSTGRES_HOST=127.0.0.1
POSTGRES_PORT=5432
TOKEN_EXPIRED_IN=60m
TOKEN_MAXAGE=60
TOKEN_SECRET=<your_token_secret>
EMAIL_ADDRESS=<your_email_address>
EMAIL_HOST=<your_email_host>
EMAIL_PASSWORD=<your_email_password>

```

### Inventory Management Service

```bash
PORT=8081
KAFKA_BROKERS=localhost:9092
ORDER_EVENT_TOPIC=order-events
INVENTORY_STATUS_TOPIC=inventory-status
SHIPPING_STATUS_TOPIC=shipping-status
LOW_STOCK_NOTIFICATION_TOPIC=low-stock-notifications
USER_SERVICE_URL=http://localhost:8080
ORDER_SERVICE_URL=http://localhost:8082
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=<your_redis_password>
POSTGRES_USER=<your_postgres_user>
POSTGRES_PASSWORD=<your_postgres_password>
POSTGRES_DB=inventory_management
POSTGRES_HOST=127.0.0.1
POSTGRES_PORT=5432
TOKEN_EXPIRED_IN=60m
TOKEN_MAXAGE=60
TOKEN_SECRET=<your_token_secret>

```

### Shipping Receiving Service

```bash
PORT=8082
KAFKA_BROKERS=localhost:9092
SHIPPING_EVENT_TOPIC=shipping-events
ORDER_EVENT_TOPIC=order-events
INVENTORY_STATUS_TOPIC=inventory-status
SHIPPING_STATUS_TOPIC=shipping-status
LOW_STOCK_NOTIFICATION_TOPIC=low-stock-notifications
USER_SERVICE_URL=http://localhost:8080
ORDER_SERVICE_URL=http://localhost:8083
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=<your_redis_password>
POSTGRES_USER=<your_postgres_user>
POSTGRES_PASSWORD=<your_postgres_password>
POSTGRES_DB=shipping_receiving
POSTGRES_HOST=127.0.0.1
POSTGRES_PORT=5432
TOKEN_EXPIRED_IN=60m
TOKEN_MAXAGE=60
TOKEN_SECRET=<your_token_secret>

```

### Customer Service

```bash
PORT=8087
KAFKA_BROKERS=localhost:9092
USER_SERVICE_URL=http://localhost:8080
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=<your_redis_password>
POSTGRES_USER=<your_postgres_user>
POSTGRES_PASSWORD=<your_postgres_password>
POSTGRES_DB=customer_service
POSTGRES_HOST=127.0.0.1
POSTGRES_PORT=5432
TOKEN_EXPIRED_IN=60m
TOKEN_MAXAGE=60
TOKEN_SECRET=<your_token_secret>
EMAIL_ADDRESS=<your_email_address>
EMAIL_HOST=<your_email_host>
EMAIL_PASSWORD=<your_email_password>
```

### User Management Service

```bash
PORT=8080
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=<your_redis_password>
POSTGRES_USER=<your_postgres_user>
POSTGRES_PASSWORD=<your_postgres_password>
POSTGRES_DB=user_management
POSTGRES_HOST=127.0.0.1
POSTGRES_PORT=5432
TOKEN_EXPIRED_IN=60m
TOKEN_MAXAGE=60
TOKEN_SECRET=<your_token_secret>
```

### Reporting and Analytics Service

```bash
PORT=8084
KAFKA_BROKERS=localhost:9092
USER_SERVICE_URL=http://localhost:8080
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=<your_redis_password>
POSTGRES_USER=<your_postgres_user>
POSTGRES_PASSWORD=<your_postgres_password>
POSTGRES_DB=report_analytics
POSTGRES_HOST=127.0.0.1
POSTGRES_PORT=5432
TOKEN_EXPIRED_IN=60m
TOKEN_MAXAGE=60
TOKEN_SECRET=<your_token_secret>
```

### Account Management Service

```bash
PORT=8086
KAFKA_BROKERS=localhost:9092
USER_SERVICE_URL=http://localhost:8080
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=<your_redis_password>
POSTGRES_USER=<your_postgres_user>
POSTGRES_PASSWORD=<your_postgres_password>
POSTGRES_DB=accounts_management
POSTGRES_HOST=127.0.0.1
POSTGRES_PORT=5432
TOKEN_EXPIRED_IN=60m
TOKEN_MAXAGE=60
TOKEN_SECRET=<your_token_secret>
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

## API Documentation

Swagger is used to generate API documentation. Access the documentation for each service at:

- **Order Service:** http://localhost:8080/swagger/index.html
- **Inventory Service:** http://localhost:8081/swagger/index.html
- **Shipping Service:** http://localhost:8082/swagger/index.html
- **Customer Service:** http://localhost:8083/swagger/index.html
- **User Management Service:** http://localhost:8084/swagger/index.html
- **Reporting and Analytics Service:** http://localhost:8085/swagger/index.html
- **Account Management Service:** http://localhost:8086/swagger/index.html

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

## Databases

Each service has its own database. Refer to the internal/initializers directory in each service for detailed implementations.

### PostgreSQL

PostgreSQL is used as the primary database for all services. Each service connects to its own PostgreSQL database instance to store and retrieve data.

### Redis

Redis is used for caching data and storing session information across all services.

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
