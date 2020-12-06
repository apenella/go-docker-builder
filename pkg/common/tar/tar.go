package tar

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	errors "github.com/apenella/go-common-utils/error"
)

// Tar return an tar io.Reader from the gived directory. It returns an error when the file is not a directory.
func Tar(path *os.File) (io.Reader, error) {

	var err error
	var tarBuff bytes.Buffer
	var stat os.FileInfo

	// ensure the src actually exists before trying to tar it
	stat, err = os.Stat(path.Name())
	if err != nil {
		return nil, errors.New("(common::tar::Tar)", fmt.Sprintf("Stat error for '%s'", path.Name()), err)
	}

	// context to tar must be a directory
	if !stat.IsDir() {
		return nil, errors.New("(common::tar::Tar)", fmt.Sprintf("'%s' must be a directory", path.Name()))
	}

	tw := tar.NewWriter(&tarBuff)
	defer tw.Close()

	err = filepath.Walk(path.Name(), func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return errors.New("(common::tar::Tar::Walk)", "Error at the beginning of the walk", err)
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return errors.New("(common::tar::Tar::Walk)", fmt.Sprintf("Error creating '%s' header", file), err)
		}
		relativePath, err := filepath.Rel(path.Name(), file)
		if err != nil {
			return errors.New("(common::tar::Tar::Walk)", fmt.Sprintf("A relative path on '%s' could not be made from '%s'", file, path.Name()), err)
		}
		header.Name = relativePath

		// write the header
		if err := tw.WriteHeader(header); err != nil {
			return errors.New("(common::tar::Tar::Walk)", fmt.Sprintf("Error writing '%s' header", file), err)
		}

		// open files for taring
		f, err := os.Open(file)
		if err != nil {
			return errors.New("(common::tar::Tar::Walk)", fmt.Sprintf("Error opening '%s'", file), err)
		}

		if _, err := io.Copy(tw, f); err != nil {
			return errors.New("(common::tar::Tar::Walk)", fmt.Sprintf("Error copying '%s' into tar", file), err)
		}
		// manually close here after each file operation; defering would cause each file close
		// to wait until all operations have completed.
		f.Close()

		return nil
	})
	if err != nil {
		return nil, errors.New("(common::tar::Tar)", fmt.Sprintf("Error exploring '%s'", path.Name()), err)
	}

	return bytes.NewReader(tarBuff.Bytes()), nil
}
