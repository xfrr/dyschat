# DYCHAT

This repository contains the source code for the DyChat project. The project is a distributed chat application that allows users to communicate with each other in real time. The application is built using NATS JetStream as the messaging system and persistence layer. The application is built using the following technologies:

### Messaging System & Persistence Layer
- [NATS JetStream](https://docs.nats.io/jetstream/jetstream)
- [NATS Go Client](https://github.com/nats-io/nats.go)

### Monitoring
- [Prometheus](https://prometheus.io/)
- [Grafana](https://grafana.com/)
- [Jaeger](https://www.jaegertracing.io/)
- [OpenTelemetry](https://opentelemetry.io/)
- [OpenTelemetry Go Client](https://github.com/open-telemetry/opentelemetry-go)

## Getting Started
### Prerequisites
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Go](https://golang.org/)

### Running the Application
1. Clone the repository
2. Run `make run_<mode>` where `<mode>` is one of the following:
    - `dev`: Runs the application in development mode
    - `prod`: Runs the application in production mode

### Running the Tests
1. Clone the repository
2. Run `make test`
3. Run `make test_coverage` to generate a coverage report

## Architecture
The application is composed of the following services:

### Common
These environment variables are used by all services.

| Name                     | Description                                         | Default Value                       |
| ------------------------ | --------------------------------------------------- | ----------------------------------- |
| `DYCHAT_ENV_MODE `       | The environment mode for the application            | `dev`                               |
| `DYCHAT_NATS_URL`        | The URL of the NATS server to consume messages from | `nats://localhost:4222`             |
| `DYCHAT_JAEGER_URL`      | The URL of the Jaeger server                        | `http://localhost:14268/api/traces` |
| `DYCHAT_PROMETHEUS_ADDR` | The address of the Prometheus server                | `localhost:9090`                    |
| `DYCHAT_GRAFANA_ADDR`    | The address of the Grafana server                   | `localhost:3000`                    |
| `DYCHAT_REDIS_ADDR`      | The address of the Redis server                     | `localhost:6379`                    |
| `DYCHAT_REDIS_PASS`      | The password for the Redis server                   | ``                                  |
| `DYCHAT_REDIS_DATABASE`  | The database for the Redis server                   | `0`                                 |


### Authentication Service
This service is responsible for authenticating users and generating JWT tokens.

| Name                          | Description                                           | Default Value     |
| ----------------------------- | ----------------------------------------------------- | ----------------- |
| `DYCHAT_AUTH_GRPC_PORT`       | The address of the authentication service gRPC server | `localhost:50051` |
| `DYCHAT_AUTH_SECRET`          | The secret used to sign JWT tokens                    | `secret`          |
| `DYCHAT_AUTH_AUTH_LOG_LEVEL ` | The log level for the authentication service          | `info`            |

### Message Writer Service
This service is responsible for writing messages to the database. It listens for messages sent to a Room and writes them to the database.

| Name                          | Description                                         | Default Value |
| ----------------------------- | --------------------------------------------------- | ------------- |
| `DYCHAT_MSG_WRITER_ENV_MODE ` | The environment mode for the message writer service | `dev`         |
| `DYCHAT_MSG_WRITER_LOG_LEVEL` | The log level for the message writer service        | `info`        |

### Rooms Service
This service is responsible for managing the rooms.

| Name                     | Description                                  | Default Value     |
| ------------------------ | -------------------------------------------- | ----------------- |
| `DYCHAT_ROOMS_GRPC_PORT` | The address of the rooms service gRPC server | `localhost:50052` |
| `DYCHAT_ROOMS_LOG_LEVEL` | The log level for the rooms service          | `info`            |

### Websocket Agent Service
This service exposes a websocket endpoint that allows to receive and send messages in real time to a room.

| Name                        | Description                                   | Default Value     |
| --------------------------- | --------------------------------------------- | ----------------- |
| `DYCHAT_WS_AGENT_PORT`      | The address of the websocket agent service    | `localhost:50053` |
| `DYCHAT_WS_AGENT_LOG_LEVEL` | The log level for the websocket agent service | `info`            |



