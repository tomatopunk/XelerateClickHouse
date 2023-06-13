# Project Title: clickhouse-benchmark

## Description

clickhouse-benchmark is a command-line tool built with Go for benchmarking and testing the performance of ClickHouse databases using the clickhouse-go driver. This tool provides commands to describe the database schema, write data, read data, and initialize the database for benchmarking purposes.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Commands](#commands)
- [Contributing](#contributing)
- [License](#license)

## Installation

To install clickhouse-benchmark, you can download the binary for your operating system from the GitHub releases page. Alternatively, you can build it from source by following these steps:

1. Ensure you have Go installed and configured on your system.
2. Clone the repository: `git clone https://github.com/your-username/clickhouse-benchmark.git`.
3. Navigate to the project directory: `cd clickhouse-benchmark`.
4. Build the binary: `go build -o clickhouse-benchmark`.

## Usage

To use clickhouse-benchmark, open your terminal and execute the binary with the desired command and options. The general syntax is as follows:

```bash
./clickhouse-benchmark [command] [options]
```

## Commands

clickhouse-benchmark supports the following commands:

### desc

The `desc` command allows you to describe the schema of a ClickHouse database table. It provides information about the columns, their types, and any indexes or constraints.

```bash
./clickhouse-benchmark desc
```

### init

The `init` command initializes the ClickHouse database for benchmarking by creating the necessary tables and performing any required setup.

```bash
./clickhouse-benchmark init
```

### read

The `read` command benchmarks the read performance of the ClickHouse database by executing read queries. You can specify the start time, end time, time step, and SQL query for the benchmark.

```bash
./clickhouse-benchmark read --start [start-time] --end [end-time] --step [time-step] --sql [query]
```

### write

The `write` command benchmarks the write performance of the ClickHouse database by writing data. You can specify the bucket count, bucket size, and concurrency limit for the benchmark.

```bash
./clickhouse-benchmark write --b [bucket-count] --n [bucket-size] --c [concurrency-limit]
```

Note: Replace `[start-time]`, `[end-time]`, `[time-step]`, `[query]`, `[bucket-count]`, `[bucket-size]`, and `[concurrency-limit]` with the actual values for your benchmark.

## Make Usage

The Makefile in your project provides several useful commands for building and pushing Docker images. Here is an example of how you can use it:

```bash
make build-push
```

This command will log in to the Docker registry using the provided username and password, create a buildx builder instance, and build the Docker image using the specified Dockerfile for the target platforms (linux/amd64 and linux/arm64). The image will be tagged with the registry, image name, and version. Finally, the image will be pushed to the Docker registry.

Make sure you have Docker installed and configured on your system before running the `make build-push` command.

## Adding License

To add a license header to your Go files, you can use the following command in your project:

```bash
make add-license
```

This command will iterate over all the Go files specified in the `GO_FILES` variable and add the license header to each file. The license header is defined in the `LICENSE_HEADER` variable.

Please ensure that the `GO_FILES` variable is properly configured in your Makefile to include all the relevant Go files in your project. Additionally, make sure to set the `LICENSE_HEADER` variable to the desired license text.

## Usage of build/k8s.yaml

The `build/k8s.yaml` file in your project is used for configuring and deploying your application in a Kubernetes cluster. You can use the following steps to utilize this file:

1. Make sure you have a Kubernetes cluster set up and configured.
2. Apply the configuration from the `build/k8s.yaml` file using the following command:

   ```bash
   kubectl apply -f build/k8s.yaml
   ```

   This command will create the necessary resources (such as deployments, services, etc.) in your Kubernetes cluster based on the specifications in the `build/k8s.yaml` file.

   Note: Ensure that you have the `kubectl` command-line tool installed and properly configured to connect to your Kubernetes cluster.

3. Monitor the deployment and check the status of your application using the appropriate Kubernetes commands, such as `kubectl get deployments`, `kubectl get pods`, or `kubectl get services`.

## Contributing

Contributions to clickhouse-benchmark are welcome! If you encounter any issues or have suggestions for improvement, please open an issue on the GitHub repository.

## License

This project is licensed under the Apache License. See the [LICENSE](LICENSE) file for more information.