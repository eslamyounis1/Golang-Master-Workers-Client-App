# Golang Master-Workers-Client Application

This repository demonstrates the Master-Worker pattern in Go, implementing a distributed system where a master node coordinates tasks among multiple worker nodes, and a client interacts with the master to submit tasks and retrieve results.

## Overview

The Master-Worker pattern is a concurrency design pattern that divides work among multiple worker goroutines or nodes to perform tasks concurrently, enhancing performance and scalability.

In this application:

- **Master**: Manages task distribution among workers and aggregates results.
- **Workers**: Perform assigned tasks and return results to the master.
- **Client**: Submits tasks to the master and retrieves processed results.

## Repository Structure

- **client/**: Contains the client application code.
- **master/**: Contains the master node implementation.
- **slave/**: Contains the worker node implementation.
- **mapreduce/**: Contains MapReduce-related code for task processing.
- **makefile**: Build automation script.
- **slave2.fasta**, **slave3.fasta**: Example data files for worker nodes.

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) installed (version 1.16 or higher recommended).

### Building the Application

Use the provided `makefile` to build the master, worker, and client applications:

```bash
make all
```

This command compiles the source code and generates executable binaries for each component.

### Running the Application

1. **Start the Master Node**:

   ```bash
   ./master/master
   ```

   The master node will initialize and listen for incoming connections from workers and clients.

2. **Start Worker Nodes**:

   ```bash
   ./slave/slave -id=1
   ./slave/slave -id=2
   ```

   Each worker node registers itself with the master and waits for tasks. Assign unique IDs to each worker using the `-id` flag.

3. **Run the Client**:

   ```bash
   ./client/client -task=input_data.txt
   ```

   The client submits a task (e.g., processing `input_data.txt`) to the master, which distributes the workload among available workers. Once processing is complete, the client retrieves the results.

### Configuration

- **Master**: Configure the master node's behavior and network settings in `master/config.go`.
- **Worker**: Set worker-specific parameters, such as processing capabilities and data paths, in `slave/config.go`.
- **Client**: Adjust client settings, including task parameters and result handling, in `client/config.go`.



