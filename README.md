go-docker-builder
=======

`go-docker-builder` package simplifies the golang's docker client SDK basic operations usage. It lets to write your own golang code to pull, build and push docker images in a few lines. The package also manages the registy auth either on pull or push actions.

Another feature of the package is that it supports to build images from `path` either `git` contexts. 
- **Path**: To build docker images from a local directory context definition.
- **Git**: To build docker images from a git repository context definition.
	Regarding git context, go-docker-builder also supports the auth to a git server by username/password, ssh agent or private key file.

## Table of Contents
- [Packages](#packages)
- [Examples](#examples)
- [References](#references)
- [License](#license)

## Packages

### Auth
On this folder are located the pakcages which manages auth operations. These packages simplify the docker and git auth methods.

- files:
```
|-- docker
|   |-- auth.go
|   `-- auth_test.go
`-- git
    |-- auth.go
    |-- basic
    |   `-- basic_auth.go
    |-- key
    |   `-- key_auth.go
    `-- sshagent
        `-- sshagent_auth.go
```

### Build
On this folder is found the docker build logic.

- files:
```
.
|-- build.go
|-- buildOptions.go
|-- buildOptions_test.go
`-- context
    |-- context.go
    |-- git
    |   `-- git.go
    |-- path
        `-- path.go
```

Build details are defined on `DockerBuildOption` struct, which is, mainly, a `"github.com/docker/docker/api/types"` `ImageBuildOptions`options struct subset.
```go
type DockerBuildOptions struct {
	// ImageName is the name of the image
	ImageName string
	// Tags is a list of the image tags
	Tags []string
	// BuildArgs ia a list of arguments to set during the building
	BuildArgs map[string]*string
	// Dockerfile is the file name for dockerfile file
	Dockerfile string
	// PushAfterBuild push image to registry after building
	PushAfterBuild bool
	// Auth required to be authenticated to docker registry
	Auth map[string]dockertypes.AuthConfig
	// BuildContext
	DockerBuildContext context.DockerBuildContexter
}
```

The `DockerBuildCmd` element is the responsible to run the docker build. On this struct there is `DockerPushOptions` attribute that contains the options to push images once its build is finished.
```go
type DockerBuildCmd struct {
	// Writer to write the build output
	Writer             io.Writer
	// Context manages the build context
	Context            context.Context
	// Cli is the docker api client
	Cli                *client.Client
	// DockerBuildOptions are the options to build
	DockerBuildOptions *DockerBuildOptions
	// DockerPushOptions are the option to push
	DockerPushOptions  *push.DockerPushOptions
	// ExecPrefix defines a prefix to each output lines
	ExecPrefix         string
	// Response manages responses from docker client
	Response           types.Responser
}
```

`Response`attribute is `Responser` interface and manages docker cli responses.
```go
type Responser interface {
	Write(io.Writer, io.ReadCloser) error
}
```

#### Context
On this build's subfolder are lacated the docker build context packages. Any docker build context implements the interface 

```go
import "github.com/go-git/go-git/v5/plumbing/transport"

type GitAuther interface {
	Auth() (transport.AuthMethod, error)
}
```

### Common
On this folder are located common or shared packages.
`Tar` package manages the tar object required by docker API to build images.

- files:
```
.
`-- tar
    `-- tar.go
```

### Push
On this folder is found the docker push logic.

- files:
```
.
|-- push.go
|-- pushOptions.go
`-- pushOptions_test.go
```

When pushing a docker image to registry the struct below defines own to push the image.
```go
// DockerBuilderOptions has an options set to build and image
type DockerPushOptions struct {
	// ImageName is the name of the image
	ImageName string
	// Tags is a list of the images to push
	Tags []string
	// RegistryAuth is the base64 encoded credentials for the registry
	RegistryAuth *string
}
```

### Response
On this folder is defined `DefaultResponse`, a `Responser` interface interface implementation.

- files:
```
.
`-- response.go
```

### Types
On this folder there is a set of custom types for `go-docker-builder`.

- files:
```
.
|-- responseBodyStreamAuxMessage.go
|-- responseBodyStreamErrorDetailMessage.go
|-- responseBodyStreamMessage.go
`-- responser.go
```

## Examples
You could find examples about how to you build, push or pull docker images using `go-docker-build` on the [exmples](https://github.com/apenella/go-docker-builder/tree/master/examples) repository folder.

> Note: On the examples, are used unexisting values for users, passowrd or registry hosts. If you would like test any example, take care to modify these values.

```
.
|-- build-and-push
|   |-- build-and-push.go
|   `-- files
|       `-- Dockerfile
|-- build-git-context
|   `-- build.go
|-- build-git-context-auth
|   `-- build.go
|-- build-path-context
|   |-- build.go
|   `-- files
|       `-- Dockerfile
`-- push
    `-- push.go
```

### List of examples

- **build-and-push**: Build and push an image. [go to exmple](https://github.com/apenella/go-docker-builder/tree/master/examples/build-and-push)
- **build-git-context**: Build an images using a git repository as a context. [go to exmple](https://github.com/apenella/go-docker-builder/tree/master/examples/build-git-context)
On the snipped below, you could see how to run on of these examples.
```sh
apenella [go-docker-builder/examples/build-git-context] $ go run build.go
registry/namespace/ubuntu ──  Step 1/5 : FROM alpine:3.9
registry/namespace/ubuntu ──
registry/namespace/ubuntu ──   ‣ Pulling from library/alpine
registry/namespace/ubuntu ──   ‣ Digest: sha256:414e0518bb9228d35e4cd5165567fb91d26c6a214e9c95899e1e056fcd349011
registry/namespace/ubuntu ──   ‣ Status: Image is up to date for alpine:3.9
registry/namespace/ubuntu ──  ---> 78a2ce922f86
registry/namespace/ubuntu ──  Step 2/5 : RUN apk add --no-cache lua5.3 lua-filesystem lua-lyaml lua-http
registry/namespace/ubuntu ──
registry/namespace/ubuntu ──  ---> Using cache
registry/namespace/ubuntu ──  ---> 311d8b3d8000
registry/namespace/ubuntu ──  Step 3/5 : COPY fetch-latest-releases.lua /usr/local/bin
registry/namespace/ubuntu ──
registry/namespace/ubuntu ──  ---> Using cache
registry/namespace/ubuntu ──  ---> b320a9351508
registry/namespace/ubuntu ──  Step 4/5 : VOLUME /out
registry/namespace/ubuntu ──
registry/namespace/ubuntu ──  ---> Using cache
registry/namespace/ubuntu ──  ---> 435051cee5dd
registry/namespace/ubuntu ──  Step 5/5 : ENTRYPOINT [ "/usr/local/bin/fetch-latest-releases.lua" ]
registry/namespace/ubuntu ──
registry/namespace/ubuntu ──  ---> Using cache
registry/namespace/ubuntu ──  ---> 791db64c487e
registry/namespace/ubuntu ──   ‣ sha256:791db64c487e19b418367522376205460718e414687128752af4c5259c4b6d00
registry/namespace/ubuntu ──  Successfully built 791db64c487e
registry/namespace/ubuntu ──  Successfully tagged registry/namespace/ubuntu:tag1
registry/namespace/ubuntu ──  Successfully tagged registry/namespace/ubuntu:latest
apenella [go-docker-builder/examples/build-git-context] $

```
- **build-git-context-auth**: Build an images using a git repository as a context and which required git server authorization. [go to exmple](https://github.com/apenella/go-docker-builder/tree/master/examples/build-git-context-auth)
- **build-path-context**: Build an images using a location folder as a context. [go to exmple](https://github.com/apenella/go-docker-builder/tree/master/examples/build-path-context)
- **push**: Push and image. [go to exmple](https://github.com/apenella/go-docker-builder/tree/master/examples/push)

## References
- Here there is docker engine API specifications for building and image using it. https://docs.docker.com/engine/api/v1.39/#operation/ImageBuild
- Taring files strategy was inspired in: https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07

## License
go-docker-builder is available under [MIT](https://github.com/apenella/go-docker-builder/blob/master/LICENSE) license.
