package tar

import (
	"archive/tar"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Tar return an tar io.Reader from the gived directory. It returns an error when the file is not a directory.
func Tar(path *os.File) (io.Reader, error) {

	var err error
	var tarBuff bytes.Buffer
	var stat os.FileInfo

	// ensure the src actually exists before trying to tar it
	stat, err = os.Stat(path.Name())
	if err != nil {
		return nil, errors.New("(common::tar::Tar) '" + path.Name() + "' stat. " + err.Error())
	}

	// context to tar must be a directory
	if !stat.IsDir() {
		return nil, errors.New("(common::tar::Tar) '" + path.Name() + "' must be a directory")
	}

	tw := tar.NewWriter(&tarBuff)
	defer tw.Close()

	err = filepath.Walk(path.Name(), func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return errors.New("(common::tar::Tar::Walk) Error at the beginning of the walk. " + err.Error())
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return errors.New("(common::tar::Tar::Walk) Error creating '" + file + "' header. " + err.Error())
		}
		// update the name to correctly reflect the desired destination when untaring
		header.Name = strings.TrimPrefix(strings.Replace(file, path.Name(), "", -1), string(filepath.Separator))

		// write the header
		if err := tw.WriteHeader(header); err != nil {
			return errors.New("(common::tar::Tar::Walk) Error writing '" + file + "' header. " + err.Error())
		}

		// open files for taring
		f, err := os.Open(file)
		if err != nil {
			return err
		}

		if _, err := io.Copy(tw, f); err != nil {
			return errors.New("(common::tar::Tar::Walk) Error copying '" + file + "' into tar. " + err.Error())
		}
		// manually close here after each file operation; defering would cause each file close
		// to wait until all operations have completed.
		f.Close()

		return nil
	})
	if err != nil {
		return nil, errors.New("(common::tar::Tar) Error explorint '" + path.Name() + "'. " + err.Error())
	}

	return bytes.NewReader(tarBuff.Bytes()), nil
}
