go-docker-builder
=======

`go-docker-builder` package simplifies the Golang's docker client SDK basic operations. It lets to write your own golang code to pull, build and push docker images in a few lines.
The package manages the authorization operation either on pull or push images to docker any docker registries.

The packages let you to build docker images from the docker build path and git context types.

## Packages

### Auth
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

On this folder are located authorization packages. These packages simplify the docker and git authentication methods.

### Build

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

On this folder is found the docker build logic. 
Build must be defined by `DockerBuildOption` struct, which is, mainly, a `"github.com/docker/docker/api/types"` `ImageBuildOptions`options struct subset.

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

`DockerBuildCmd` sturct contains all the details to run a build.
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
The `DockerBuildCmd` element is the responsible to run the docker build. On this struct there is `DockerPushOptions`, a struct containing the options to push the image once is built.

The response is managed by `Responser` interface, defined on types package.
```go
type Responser interface {
	Write(io.Writer, io.ReadCloser) error
}
```

#### Context
On this subfolder are also lacated the context packages definition. Any docker build context implements the interface 

```go
import "github.com/go-git/go-git/v5/plumbing/transport"

type GitAuther interface {
	Auth() (transport.AuthMethod, error)
}
```

- **Path**: To build docker images from a local directory context definition
- **Git**: To build docker images from a git repository context definition

### Common

- files:
```
.
`-- tar
    `-- tar.go
```

On this folder are located common or shared packages.
`Tar` package manages the tar object required by docker API to build images.

### Push

- files:
```
.
|-- push.go
|-- pushOptions.go
`-- pushOptions_test.go
```

On this folder is found the docker push logic. 

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

- files:
```
.
`-- response.go
```

On this folder there defined `DefaultResponse` which implements the `Responser` interface and manages the docker api responses.

### Types

- files:
```
.
|-- responseBodyStreamAuxMessage.go
|-- responseBodyStreamErrorDetailMessage.go
|-- responseBodyStreamMessage.go
`-- responser.go
```

On this folder there is a set of custom types for `go-docker-builder`.

## Example
You could find an example on `examples` folder.

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

- **build-and-push**: Build and push an image.
- **build-git-context**: Build an images using a git repository as a context.
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
- **build-git-context-auth**: Build an images using a git repository as a context and which required git server authorization.
- **build-path-context**: Build an images using a location folder as a context.
- **push**: Push and image.



## References
- Here there is docker engine API specifications for building and image using it. https://docs.docker.com/engine/api/v1.39/#operation/ImageBuild
- Taring files strategy was inspired in: https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07
