package main

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type BuildResponseBodyStreamMessage struct {
	Status string                             `json:"status"`
	Stream string                             `json:"stream"`
	Aux    *BuildResponseBodyStreamAuxMessage `json:"aux"`
}
type BuildResponseBodyStreamAuxMessage struct {
	ID string `json:"ID"`
}

func (m *BuildResponseBodyStreamMessage) String() string {

	if m.Stream != "" {
		return strings.TrimSpace(m.Stream)
	}
	if m.Status != "" {
		return " \u2023 " + strings.TrimSpace(m.Status)
	}
	if m.Aux != nil && m.Aux.ID != "" {
		return " \u2023 " + m.Aux.ID
	}
	return ""
}

func Tar(path string) (io.Reader, error) {

	var err error
	var tarBuff bytes.Buffer

	// thanks to https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07

	// ensure the src actually exists before trying to tar it
	_, err = os.Stat(path)
	if err != nil {
		return nil, errors.New("(Tar) " + err.Error())
	}

	tw := tar.NewWriter(&tarBuff)
	defer tw.Close()

	err = filepath.Walk(path, func(file string, fi os.FileInfo, err error) error {
		// return on any error
		if err != nil {
			return errors.New("(Tar:Walk) Error at the beginning of the walk. " + err.Error())
		}

		// return on non-regular files (thanks to [kumo](https://medium.com/@komuw/just-like-you-did-fbdd7df829d3) for this suggested update)
		if !fi.Mode().IsRegular() {
			return nil
		}

		// create a new dir/file header
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return errors.New("(Tar:Walk) Error creating '" + file + "' header. " + err.Error())
		}
		// update the name to correctly reflect the desired destination when untaring
		header.Name = strings.TrimPrefix(strings.Replace(file, path, "", -1), string(filepath.Separator))

		// write the header
		if err := tw.WriteHeader(header); err != nil {
			return errors.New("(Tar:Walk) Error writing '" + file + "' header. " + err.Error())
		}

		// open files for taring
		f, err := os.Open(file)
		if err != nil {
			return err
		}

		// copy file data into tar writer
		if _, err := io.Copy(tw, f); err != nil {
			return errors.New("(Tar:Walk) Error copying '" + file + "' into tar. " + err.Error())
		}
		// manually close here after each file operation; defering would cause each file close
		// to wait until all operations have completed.
		f.Close()

		return nil
	})
	if err != nil {
		return nil, errors.New("(Tar) " + err.Error())
	}

	return bytes.NewReader(tarBuff.Bytes()), nil
}

func main() {
	var err error
	var dockerCli *client.Client

	imageDefinitionPath := filepath.Join("dockerfiles", "ubuntu")
	tarBuffReader, err := Tar(imageDefinitionPath)
	if err != nil {
		panic(err)
	}
	imageName := strings.Join([]string{"gobuild", filepath.Base(imageDefinitionPath)}, "-")

	dockerCli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic("(images::NewDockerBuilder) Error on docker client creation. " + err.Error())
	}

	options := types.ImageBuildOptions{
		Context:        tarBuffReader,
		SuppressOutput: false,
		Remove:         true,
		ForceRemove:    true,
		PullParent:     true,
		Dockerfile:     "Dockerfile",
		Tags:           []string{imageName},
		//   BuildArgs:      args,
	}
	buildResponse, err := dockerCli.ImageBuild(context.TODO(), tarBuffReader, options)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	defer buildResponse.Body.Close()

	scanner := bufio.NewScanner(buildResponse.Body)
	prefix := imageName

	for scanner.Scan() {
		streamMessage := &BuildResponseBodyStreamMessage{}
		line := scanner.Bytes()
		err = json.Unmarshal(line, &streamMessage)
		if err != nil {
			fmt.Println("Error unmarshalling line: ", string(line))
		}

		fmt.Fprintf(os.Stdout, "%s \u2500\u2500  %s\n", prefix, streamMessage.String())
	}
}
