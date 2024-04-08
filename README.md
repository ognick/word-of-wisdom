# Word of Wisdom TCP Server with Proof of Work
## Overview
The "Word of Wisdom" is a TCP server that provides inspirational quotes to clients. To protect the server from Distributed Denial of Service (DDoS) attacks, a Proof of Work (PoW) mechanism is implemented.

## Proof of Work Protocol
The server and client interact using the following "challenge-response" protocol:

1. The client sends a "request service" message to the server.
2. The server generates a "challenge" for the client to solve.
3. The server sends the "challenge" to the client.
4. The client solves the "challenge".
5. The client sends the response solution to the server.
6. The server verifies the "solution" provided by the client.
7. If the solution is correct, the server grants access to the client and provides the requested service.

The Proof of Work mechanism requires the client to perform a computationally expensive task (solving the "challenge") before the server grants access to the service. This helps to mitigate DDoS attacks by making it difficult for bots to overwhelm the server with requests.
Learn More [here](https://en.wikipedia.org/wiki/Proof_of_work)

## Implementation

The server is implemented using the Go programming language. It utilizes the `crypto/argon2` package to generate and verify the "challenges". 

**Argon2** was chosen for its strong resistance to hardware-based attacks, such as those using ASICs or GPUs, and its ability to adapt the computational complexity to the server's needs.

The client implementation is also provided in Go, demonstrating the interaction with the server using the "challenge-response" protocol.

## Usage

Run test and benchmarks:
```
make test
```

Run demonstration:
```
make run
```