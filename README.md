# FragmentDB

FragmentDB is a distributed sharded database implemented in Golang.

## Features

- Distributed nodes with HTTP REST API.
- Data sharding with encryption for security.
- Recovery mechanism that periodically synchronizes data between nodes.
- Graceful shutdown using context and OS signal handling.

## Getting Started

### Prerequisites

- [Go](https://golang.org/) 1.16+
- Git

### Build and Run

1. Clone the repository:
   ```bash
   git clone https://github.com/itsfuad/FragmentDB.git
   cd FragmentDB
   ```

2. Build the project:
   ```bash
   go build -o fragmentdb
   ```

3. Prepare configuration files (e.g., `config1.json`, `config2.json`, `config3.json`) in the repository root. An example (`config1.json`):
   ```json
   {
       "node_id": "node1",
       "port": 8081,
       "peer_nodes": ["localhost:8082", "localhost:8083"],
       "data_path": "./data1",
       "secret_key": "your-32-byte-secret-key-here-12345",
       "shard_count": 3
   }
   ```

4. Run your nodes (for example, using the provided script):
   ```bash
   chmod +x scripts/run.sh
   ./scripts/run.sh
   ```

### Usage

- **Store Data**: Use a POST request to `/put` with JSON payload:
  ```json
  {
      "key": "your-key",
      "value": "your-data"
  }
  ```
- **Retrieve Data**: Use a GET request to `/get/your-key`.

- **Sync Data**: Peer nodes synchronize automatically via the `/sync` endpoint.

## API Endpoints

- **/put**: POST JSON payload to store data (e.g., `{"key": "your-key", "value": "your-value"}`).
- **/get/{key}**: GET to retrieve and reconstruct stored data.
- **/sync**: GET to expose current node data for recovery synchronization.

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository.
2. Create a feature branch (`git checkout -b feature/your-feature`).
3. Commit your changes.
4. Push and create a Pull Request.
5. Ensure all tests pass:
   ```bash
   go test ./...
   ```

## Code Verification Requirements

This project enforces code formatting and verified signatures:

- **Code Formatting:** All Go source files must be formatted using `go fmt`. To check and format your code, run:
  ```bash
  go fmt ./...
  ```
  The CI workflow will fail if any unformatted code is detected.

- **Verified Signature:** The repository requires signed commits.

Make sure both requirements are met before pushing your changes.

## Testing

Run tests using:
   ```bash
   go test ./...
   ```

## Additional Notes

- Recovery sync runs periodically based on a 5-minute interval.
- Nodes shutdown gracefully when receiving SIGINT or SIGTERM.

## Hosting

To host FragmentDB in production, consider the following:

- Deploy nodes across different servers to ensure redundancy.
- Use a load balancer or DNS round-robin for client requests.
- Secure the communication channel (e.g., HTTPS, VPN) between nodes.
- Monitor logs and use a centralized logging solution.
- Regularly back up configuration files and data directories.

## License

This project is licensed under the MOZILLA PUBLIC LICENSE 2.0 - see the [LICENSE](LICENSE) file for details.
