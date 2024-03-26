# go-docker-builder

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) ![Test](https://github.com/apenella/go-docker-builder/actions/workflows/ci.yaml/badge.svg) [![Go Reference](https://pkg.go.dev/badge/github.com/apenella/go-docker-builder.svg)](https://pkg.go.dev/github.com/apenella/go-docker-builder)

The `go-docker-builder` library serves as a wrapper over the Docker Client SDK, offering a collection of packages to simplify the development workflow for common Docker use cases, including image building, pushing or authentication to Docker registries.

This library not only handles Docker registry authentication but also prepares the Docker build context for use with the Docker Client SDK. It supports build context sourced either from the local path or a Git repository.

- [go-docker-builder](#go-docker-builder)
  - [Install](#install)
  - [Use cases](#use-cases)
    - [Build](#build)
      - [Context](#context)
        - [Path](#path)
        - [Git](#git)
      - [Context filesystem](#context-filesystem)
    - [Push](#push)
    - [Copy](#copy)
  - [Authentication](#authentication)
    - [Docker registry](#docker-registry)
    - [Git server](#git-server)
      - [Basic auth](#basic-auth)
      - [SSH key](#ssh-key)
      - [SSH agent](#ssh-agent)
  - [Response](#response)
  - [Examples](#examples)
  - [References](#references)
  - [License](#license)

## Install

To install the latest stable version of go-docker-builder, run the following command:

```sh
go get -u github.com/apenella/go-docker-builder@v0.8.4
```

## Use cases

`go-docker-builder` library provides you with the packages to resolve the following use cases:

- **Build**: Build a docker images
- **Push**: Push a docker image to a docker registry
- **Copy**: Copy a docker image from one Docker registry to another one

### Build

The `build` package's purpose is to build docker images.
To perform a build action you must create a `DockerBuildCmd` instance.

```go
// DockerBuilderCmd
type DockerBuildCmd struct {
  // Cli is the docker api client
  Cli types.DockerClienter
  // ImageName is the name of the image
  ImageName string
  // ImageBuildOptions from docker sdk
  ImageBuildOptions *dockertypes.ImageBuildOptions
  // ImagePushOptions from docker sdk
  ImagePushOptions *dockertypes.ImagePushOptions
  // PushAfterBuild when is true images are automatically pushed to registry after build
  PushAfterBuild bool
  // Response manages responses from docker client
  Response types.Responser
  // UseNormalizedNamed when is true tags are transformed to a fully qualified reference
  UseNormalizedNamed bool
  // RemoveAfterPush when is true images are removed from local after push
  RemoveAfterPush bool
}
```

Below there is a recipe to build docker images using `go-docker-builder`:

1. Create a docker client from the Docker client SDK

```go
dockerCli, err = client.NewClientWithOpts(client.FromEnv)
if err != nil {
  return err
}
```

2. Give a name to the image

```go
registry := "registry.go-docker-builder.test"
imageName := strings.Join([]string{registry, "alpine"}, "/")
```

3. In case you need a custom response, create a response object

```go
res := response.NewDefaultResponse(
  response.WithTransformers(
    transformer.Prepend("my-custom-response"),
  ),
)
```

4. Create `DockerBuildCmd` instance

```go
dockerBuilder := &build.DockerBuildCmd{
  Cli:       dockerCli,
  ImageName: imageName,
  Response:  res,
}
```

5. Create a docker build context

```go
imageDefinitionPath := filepath.Join(".", "files")
dockerBuildContext := &contextpath.PathBuildContext{
  Path: imageDefinitionPath,
}
```

6. Add the docker build Context to `DockerBuildCmd`

```go
err = dockerBuilder.AddBuildContext(dockerBuildContext)
if err != nil {
  return err
}
```

7. Include extra docker image tags, when you need them

```go
dockerBuilder.AddTags(strings.Join([]string{imageName, "tag1"}, ":"))
```

8. Include authorization either for pull, push or both

```go
err = dockerBuilder.AddAuth(username, password, registry)
if err != nil {
  return err
}
```

9. Include build arguments, when you need them

```go
err = dockerBuilder.AddBuildArgs("key", "value")
if err != nil {
  return err
}
```

10. Start the Docker image build

```go
err = dockerBuilder.Run(context.TODO())
if err != nil {
  return err
}
```

#### Context

The Docker build context is a set of files required to create a docker image. `go-docker-builder` library supports two kinds of sources to create the Docker build context: `path` and `git`

##### Path

When files are located on your localhost, the Docker build context must be created as `path`, and it is only required to define the local folder where files are located.

A `PathBuildContext` instance represents a Docker build context as the path.

```go
// PathBuildContext creates a build context from path
type PathBuildContext struct {
  // Path is context location on the local system
  Path string
}
```

##### Git

When files are located on a git repository, the Docker build context must be created as `git`.

A `GitBuildContext` instance represents a Docker build context as git.

```go
// GitBuildContext defines a build context from a git repository
type GitBuildContext struct {
  // Path must be set when docker build context is located in a subpath inside the repository
  Path string
  // Repository which will be used as docker build context
  Repository string
  // Reference is the name of the branch to clone. By default is used 'master'
  Reference string
  // Auth
  Auth auth.GitAuther
}
```

To define a `git` Docker build context, the only required attribute is `Repository`, although it accepts other configurations such as the `Reference` name (branch, commit or tag), `Auth` which is used to be authenticated over the git server or `Path` which let you define as Docker build context base a subfolder inside the repository.

`go-docker-builder` uses [go-git](https://github.com/go-git/go-git) library to manage `git` Docker build context.

#### Context filesystem

Context filesystem has been created as an intermediate filesystem between the source files and Docker build context.

Context filesystem is built on top of an [afero](https://github.com/spf13/afero) filesystem. It supports **taring** the entire filesystem and also **joining** multiple filesystems.

### Push

Package `push` purpose is to push Docker images to a Docker registry.
To perform a push action you must create a `DockerPushBuildCmd` instance.

```go
// DockerPushCmd is used to push images to docker registry
type DockerPushCmd struct {
  // Cli is the docker client to use
  Cli types.DockerClienter
  // ImagePushOptions from docker sdk
  ImagePushOptions *dockertypes.ImagePushOptions
  // ImageName is the name of the image
  ImageName string
  // Tags is a list of the images to push
  Tags []string
  // Response manages the docker client output
  Response types.Responser
  // UseNormalizedNamed when is true tags are transformed to a fully qualified reference
  UseNormalizedNamed bool
  // RemoveAfterPush when is true the image from local is removed after push
  RemoveAfterPush bool
}
```

Below there is a recipe to build docker images using `go-docker-builder`:

1. Create a docker client from Docker client SDK

```go
dockerCli, err = client.NewClientWithOpts(client.FromEnv)
if err != nil {
  return err
}
```

2. Give a name to the image

```go
registry := "registry.go-docker-builder.test"
imageName := strings.Join([]string{registry, "alpine"}, "/")
```

3. Create `DockerPushCmd` instance

```go
dockerPush := &push.DockerPushCmd{
  Cli:       dockerCli,
  ImageName: imageName,
}
```

4. Add authorization, when is required

```go
user := "myregistryuser"
pass := "myregistrypass"
dockerPush.AddAuth(user, pass)
```

5. Push the image to the Docker registry

```go
err = dockerPush.Run(context.TODO())
if err != nil {
  return err
}
```

### Copy

The `copy` package can be understood as such a `push` use case variation. Its purpose is to push images either from your localhost or from a Docker registry to another Docker registry. It can be also used to copy images from one Docker registry namespace to another namespace. 
To perform a copy action you must create a `DockerImageCopyCmd` instance.

```go
// DockerCopyImageCmd is used to copy images to docker registry. Copy image is understood as tag an existing image and push it to a docker registry
type DockerImageCopyCmd struct {
  // Cli is the docker client to use
  Cli types.DockerClienter
  // ImagePushOptions from docker sdk
  ImagePullOptions *dockertypes.ImagePullOptions
  // ImagePushOptions from docker sdk
  ImagePushOptions *dockertypes.ImagePushOptions
  // SourceImage is the name of the image to be copied
  SourceImage string
  // TargetImage is the name of the copied image
  TargetImage string
  // Tags is a copied images tags list
  Tags []string
  // UseNormalizedNamed when is true tags are transformed to a fully qualified reference
  UseNormalizedNamed bool
  // RemoteSource when is true the source image is pulled from registry before push it to its destination
  RemoteSource bool
  // RemoveAfterPush when is true the image from local is removed after push
  RemoveAfterPush bool
  // Response manages the docker client output
  Response types.Responser
}
```

You can use the `copy` package when:

- You need to copy one image from [dockerhub](https://hub.docker.com/) to a private Docker registry
- The Docker images you use on your staging environments need to be promoted somewhere else before using them in the production environment.

## Authentication

`go-docker-builder` library contains a bunch of packages to manage authentication either to the Docker registry or the git server.

### Docker registry

Package `github.com/apenella/go-docker-builder/pkg/auth/docker` provides a set of functions to create the authentication items required by the Docker registry. The Docker registry may require authentication either for `push` or `pull` operations.

Packages [build](#build), [push](#push) or [copy](#copy) already use that package on its authentication methods, and is not necessary to use it directly.

### Git server

Git server authentication is needed when the Docker build context is located on a git repository and the git server requires it to authorize you for cloning any repository.

`go-docker-builder` supports several git authentication methods such as `basic auth`, `ssh-key` or `ssh-agent`.

All these authorization methods generate an `AuthMethod` (github.com/go-git/go-git/v5/plumbing/transport).

#### Basic auth

`Basic auth` requires a `username` and `password`. It must be used when communication is done through `http/s`.

```go
type BasicAuth struct {
  Username string
  Password string
}
```

#### SSH key

You can use an `ssh key` when want to connect to the git server over SSH. It requires your private key location (`PkFile`). In case your key is being protected by a password, you have to set it on `PkPassword` attribute. Finally, when the git user is not `git`, you can define it on `GitSSHUser` attribute.

```go
type KeyAuth struct {
  GitSSHUser string
  PkFile     string
  PkPassword string
}
```

#### SSH agent

To authenticate to the git server, `ssh-agent` method uses the SSH agent running on your host. When your SSH user is not `git`, you can define it on `GitSSHUser` attribute.

```go
type SSHAgentAuth struct {
  GitSSHUser string
}
```

## Response

To control how Docker output is sent to the user, `go-docker-build` provides `response` package.
By default, [DockerBuildCmd](#build), [DockerPushCmd](#push) and [DockerCopyCmd](#copy) instances use a basic configuration of `response` that prints Docker output to `os.Stdout`. But you could customize the Docker client SDK response, define your `response` instance, and pass it to them.

`response` receives `ImageBuildResponse` or `ImagePushResponse` items, unmarshals into a struct and finally prepares a user-friendly output.
To customize the Docker output coming from response items, it could receive a list of transformer functions. That function must complain the type `TransformerFunc` defined on `github.com/apenella/go-common-utils/transformer/string`.

```go
// TransformerFunc is used to enrich or update messages before to be printed out
type TransformerFunc func(string) string
```

Below there is a [build-and-push](https://github.com/apenella/go-docker-builder/tree/master/examples/build-and-push) snipped that shows how to define your own `response`:

```go
res := response.NewDefaultResponse(
  response.WithTransformers(
    transformer.Prepend("buildAndPush"),
  ),
  response.WithWriter(w),
)

dockerBuilder := &build.DockerBuildCmd{
  Cli:               dockerCli,
  ImageBuildOptions: &dockertypes.ImageBuildOptions{},
  ImagePushOptions:  &dockertypes.ImagePushOptions{},
  PushAfterBuild:    true,
  RemoveAfterPush:   true,
  ImageName:         imageName,
  Response:          res,
}
```

## Examples

On folder [examples](https://github.com/apenella/go-docker-builder/tree/master/examples), you could find some `go-docker-build` examples. Among those examples, you could find how to build images using distinct Docker build context, how to authenticate to Docker registry or git server, etc.

To run any example, the repository is provided with some resources that let you start an ephemeral environment where examples can run. Each environment runs on `docker compose` and starts a Docker registry, a git server and a client container where the example runs. Those environments are also used to run the functional test.
Each example is also provided by a `Makefile` which helps you to start the examples or tests.

```shell
❯ make help

 Executing example build-and-push

 help                 list allowed targets
 start                start docker registry
 cleanup              cleanup example environment
 generate-certs       generate certificate for go-docker-builder.test
 cleanup-certs        cleanup certificates
 prepare              prepare docker images required to run the example or test
 example              executes the examples
 test                 executes functional test
 logs                 show services logs
 ```

## References

- Docker engine API specifications for building an image: https://docs.docker.com/engine/api/v1.39/#operation/ImageBuild
- The used taring files strategy is inspired by: https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07

## License

`go-docker-builder` is available under [MIT](https://github.com/apenella/go-docker-builder/blob/master/LICENSE) license.
