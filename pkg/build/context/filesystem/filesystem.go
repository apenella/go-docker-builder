package filesystem

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/spf13/afero"
)

const DefaultRootPath = "/"

// ContextFilesystem is an afero fs extension
type ContextFilesystem struct {
	afero.Fs
	RootPath string
}

func NewContextFilesystem(fs afero.Fs) *ContextFilesystem {
	return &ContextFilesystem{fs, DefaultRootPath}
}

// Tar return an io.Reader with ContextFilesystem archived
func (f *ContextFilesystem) Tar() (io.Reader, error) {

	var err error
	var tarBuff bytes.Buffer
	var stat os.FileInfo
	var path afero.File

	errorContext := "(filesystem::Tar)"

	if f == nil {
		return nil, errors.New(errorContext, "ContextFilesystem is nil")
	}

	path, err = f.Fs.Open(f.RootPath)
	if err != nil {
		panic(err.Error())
	}

	// ensure the src actually exists before trying to tar it
	stat, err = f.Fs.Stat(path.Name())
	if err != nil {
		return nil, errors.New(errorContext, fmt.Sprintf("Stat error for '%s'", path.Name()), err)
	}

	// context to tar must be a directory
	if !stat.IsDir() {
		return nil, errors.New(errorContext, fmt.Sprintf("'%s' must be a directory", path.Name()))
	}

	tw := tar.NewWriter(&tarBuff)
	defer tw.Close()

	if f.RootPath == "" {
		f.RootPath = DefaultRootPath
	}

	err = afero.Walk(f.Fs, f.RootPath, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return errors.New(errorContext, "Error at the beginning of the walk", err)
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return errors.New(errorContext, fmt.Sprintf("Error creating '%s' header", file), err)
		}
		relativePath, err := filepath.Rel(path.Name(), file)
		if err != nil {
			return errors.New(errorContext, fmt.Sprintf("A relative path on '%s' could not be made from '%s'", file, path.Name()), err)
		}
		header.Name = relativePath

		// write the header
		if err := tw.WriteHeader(header); err != nil {
			return errors.New(errorContext, fmt.Sprintf("Error writing '%s' header", file), err)
		}

		// open files for taring
		f, err := f.Fs.Open(file)
		if err != nil {
			return errors.New(errorContext, fmt.Sprintf("Error opening '%s'", file), err)
		}

		if _, err := io.Copy(tw, f); err != nil {
			return errors.New(errorContext, fmt.Sprintf("Error copying '%s' into tar", file), err)
		}
		// manually close here after each file operation; defering would cause each file close
		// to wait until all operations have completed.
		f.Close()

		return nil
	})
	if err != nil {
		return nil, errors.New(errorContext, fmt.Sprintf("Error walking through '%s'", path.Name()), err)
	}

	return bytes.NewReader(tarBuff.Bytes()), nil
}

// Join create a new ContextFilesystem joining the input filesystems content
func Join(forceOverride bool, filesystem ...*ContextFilesystem) (*ContextFilesystem, error) {
	var err error
	errorContext := "(filesystem::Join)"

	memfs := NewContextFilesystem(afero.NewMemMapFs())

	// Iterate over each filesystem
	for _, fs := range filesystem {

		if fs == nil {
			return nil, errors.New(errorContext, "Error trying join a nil filesystem")
		}

		err = afero.Walk(fs.Fs, fs.RootPath, func(file string, fi os.FileInfo, err error) error {
			var src, dest afero.File

			if err != nil {
				return errors.New(errorContext, fmt.Sprintf("Error walking through file '%s", file), err)
			}

			// Extrat relative file relative path to fs.RootPath
			// example: when root path is /test/fs1 and file is /test/fs1/dir1/f1.txt, the relative path is /dir1/f1.txt
			relativePath, _ := filepath.Rel(fs.RootPath, file)
			if err != nil {
				return errors.New(errorContext, fmt.Sprintf("Unable to achive relative path for '%s'", file), err)
			}

			memfile := filepath.Join(memfs.RootPath, relativePath)

			if fi.IsDir() {
				_, err = memfs.Fs.Stat(memfile)
				if os.IsNotExist(err) {

					memfs.Fs.MkdirAll(memfile, fi.Mode())
				}
			} else {
				_, err = memfs.Fs.Stat(memfile)
				if os.IsNotExist(err) {
					dest, err = memfs.Fs.OpenFile(memfile, os.O_CREATE, fi.Mode())
					if err != nil {
						return errors.New(errorContext, fmt.Sprintf("Error creating file '%s' on destination filesystem", memfile), err)
					}
				} else {
					// When file already exists and the join is not force it fails
					if !forceOverride {
						return errors.New(errorContext, fmt.Sprintf("File '%s' already on destination filesystem", memfile))
					}

					dest, err = memfs.Fs.OpenFile(memfile, os.O_RDWR, fi.Mode())
					if err != nil {
						return errors.New(errorContext, fmt.Sprintf("Error opening '%s' on destination filesystem", memfile), err)
					}

					// truncate the existing file
					err = dest.Truncate(0)
					if err != nil {
						return errors.New(errorContext, fmt.Sprintf("Error truncating '%s' on destination filesystem", memfile), err)
					}
				}

				src, err = fs.Fs.Open(file)
				if err != nil {
					return errors.New(errorContext, fmt.Sprintf("Error openning source file '%s'", file), err)
				}

				if _, err := io.Copy(dest, src); err != nil {
					return errors.New(errorContext, fmt.Sprintf("Error copying '%s' into destination filesystem", memfile), err)
				}
				// manually close here after each file operation; defering would cause each file close
				// to wait until all operations have completed.
				src.Close()
				dest.Close()

			}
			return nil
		})
		if err != nil {
			return nil, errors.New(errorContext, "Error joinnig context filesystem", err)
		}
	}

	return memfs, nil
}
