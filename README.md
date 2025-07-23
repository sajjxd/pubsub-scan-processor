# Pub/Sub Scan Processor

This project contains two Go services that simulate a security scan processing pipeline using a message queue.

## Services
* **Scanner**: A mock service that generates random scan results and publishes them to a Google Cloud Pub/Sub topic.
* **Processor**: A service that consumes scan results from the Pub/Sub subscription, normalizes the data, and stores the latest result for each unique `(ip, port, service)` in a SQLite database.

---
## Local Development

### Prerequisites
* Docker
* Docker Compose

### Running the Stack
To build and run all services locally, execute the following command from the project root:
```bash
docker compose up --build
```

This will start the Pub/Sub emulator, create the necessary topic and subscription, and run the `scanner` and `processor` services.

To run multiple `processor` instances for higher throughput, you can use the `--scale` flag:
```bash
docker compose up --build --scale processor=3
```
