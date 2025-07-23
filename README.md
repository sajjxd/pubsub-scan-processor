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
## Testing

After running docker compose up, you can verify that everything is working correctly by checking the container logs and querying the database directly.

1. Check the Logs

You should see logs from both services streaming in your terminal.

The scanner will log every message it publishes.
The processor will log every message it successfully processes and stores.

Example logs

```
scanner-1    | Published message for IP: 1.1.1.123
processor-1  | Processed record for 1.1.1.123:45678 (HTTP)
```

2. Query the Database

The processor service will create a scans.db file inside a data/ directory in your project root (I've added this to make it easier to test).

You can query this database from a separate terminal window to inspect the stored data

```
sqlite3 scans.db
SELECT * FROM scans ORDER BY last_scanned DESC LIMIT 5;
```

Sample Output

```
1.1.1.154|22|SSH|service response: 88|2025-07-23 19:32:20
1.1.1.92|80|HTTP|service response: 41|2025-07-23 19:32:21
1.1.1.217|53|DNS|service response: 12|2025-07-23 19:32:22
```
