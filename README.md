go-docker-builder
=======

go-docker-builder's aim was to write an example which shows how to build a docker image from a dockerfile, using Golang's docker client SDK. This repository also provides a package with a set of helpers which makes easy to define how to build and image or how to read the docker's client responses.

## Packages

### Types
**buildResponseBodyStreamMessage**: This struct defines the docker's client build response body message.
```go
// BuildResponseBodyStreamMessage contains the ImageBuild's body data from buildResponse
type BuildResponseBodyStreamMessage struct {
	// Status represents the status value on response body stream message
	Status string `json:"status"`
	// Stream represents the stream value on response body stream message
	Stream string `json:"stream"`
	// Aux represents the aux value on response body stream message
	Aux *BuildResponseBodyStreamAuxMessage `json:"aux"`
}
```

**buildResponseBodyStreamAuxMessage**: This struct defines the docker's client build aux value on response body message.
```go
// BuildResponseBodyStreamAuxMessage contains the ImageBuild's aux data from buildResponse
type BuildResponseBodyStreamAuxMessage struct {
	// ID is response body stream aux's id
	ID string `json:"ID"`
}
```

### Builder
**dockerBuilderContext**: This struct defines the building context.
When a context is generated, it will try to generate the building context from path, in case there is no path define then it will try URL and finally Git.
```go
// DockerBuilderContext defines the building context
type DockerBuilderContext struct {
	// Path defines the path location when the contxext is placed on a local folder
	Path string `yaml:"path"`
	// URL defines the url when the contxext is located remotely and published via HTTP
	URL string `yaml:"status"`
	// Git defines git reposiotry url when the contxext is located remote repository
	Git string `yaml:"git"`
}
```
> Note: URL and Git context are not already available.

**dockerBuilderOptions**: This struct defines building options.
```go
// DockerBuilderOptions has an options set to build and image
type DockerBuilderOptions struct {
	// ImageName is the name of the image
	ImageName string
	// Tags is a list of the image tags
	Tags []string
	// BuildArgs ia a list of arguments to set during the building
	BuildArgs map[string]*string
	// Dockerfile is the file name for dockerfile file
	Dockerfile string
	// PushImage set to push or not an image
	PushImage bool
}
```

**dockerBuilderCmd**: This struct contains the docker client and runtime details.
```go
// DockerBuilderCmd defines the image building runtime
type DockerBuilderCmd struct {
	// Writer is where the response body will write and available
	Writer               io.Writer
	// Context is a golang context
	Context              context.Context
	// Cli is a docker client instance
	Cli                  *client.Client
	// DockerBuilderContext is the docker build context
	DockerBuilderContext *DockerBuilderContext
	// DockerBuilderOptions are a the options which define a build
	DockerBuilderOptions *DockerBuilderOptions
	// ExecPrefix is a prefix which will at the beginning of each response body line
	ExecPrefix           string
}
```

## Example
You could find an example on `examples` folder.

## References
- Here there is docker engine API specifications for building and image using it. https://docs.docker.com/engine/api/v1.39/#operation/ImageBuild
- Inspiration to how to tar files with Golang. https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07
