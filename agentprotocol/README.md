# The ContainerSSH Agent protocol

The ContainerSSH Agent protocol allows for forwarding several types of connections within the container: X11, TCP, etc. This description describes how the protocol works and how to integrate it.

## Concepts

The server and the client communicate over the standard input/output using the container APIs. (Docker and Kubernetes) The agent is running just like any other program would in the container.

### Server

The server in this context is the ContainerSSH agent. It waits for connection requests from the client (ContainerSSH) and opens the corresponding sockkets.

### Client

The client in this context is ContainerSSH, opening a connection by sending requests to the agent via the standard input/output using the container API (Docker or Kubernetes).

## Protocol

The protocol consists of messages sent in [CBOR-encoding](https://cbor.io/) in a back-to-back fashion.