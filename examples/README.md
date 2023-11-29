# go-docker-builder Examples

Explore practical examples that demonstrate the capabilities and usage of go-docker-builder. These examples provide a hands-on learning experience, guiding you through specific scenarios to enhance your understanding of building and managing Docker images with go-docker-builder.

## List of Examples

| Example | Description |
|---|---|
| [build-and-push](https://github.com/apenella/go-docker-builder/tree/master/examples/build-and-push) | Demonstrates the basic build and push workflow using go-docker-builder. Follows a simple scenario to build and push a Docker image. |
| [build-and-push-join-context](https://github.com/apenella/go-docker-builder/tree/master/examples/build-and-push-join-context) | Illustrates how to join multiple build contexts when using go-docker-builder to build and push Docker images. |
| [build-git-context](https://github.com/apenella/go-docker-builder/tree/master/examples/build-git-context) | Shows how to use a Git repository as the build context in go-docker-builder. |
| [build-git-context-auth](https://github.com/apenella/go-docker-builder/tree/master/examples/build-git-context-auth) | Extends the previous example by demonstrating authentication for accessing a private Git repository as the build context. |
| [build-path-context](https://github.com/apenella/go-docker-builder/tree/master/examples/build-path-context) | Guides you through building a Docker image using a specific path as the build context in go-docker-builder. |
| [copy-remote](https://github.com/apenella/go-docker-builder/tree/master/examples/copy-remote) | Illustrates how to copy files remotely from a source location to the build context in go-docker-builder. |
| [push](https://github.com/apenella/go-docker-builder/tree/master/examples/push) | Demonstrates the process of pushing a locally built Docker image to a Docker registry using go-docker-builder. |

## Running the Examples

To run each example, use the provided Makefile to prepare the environment and execute the example. The following commands are available:

- To run the example: `make example`
- To run the example as a functional test: `make test`
