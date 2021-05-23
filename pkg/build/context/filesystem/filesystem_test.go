package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestTar(t *testing.T) {

}

func TestJoin(t *testing.T) {
	var err error
	var fs1, fs2, fs3 *ContextFilesystem

	_ = fs3

	errorContext := "(filesystem::Join)"

	fs1 = NewContextFilesystem(afero.NewOsFs())
	fs1.RootPath = filepath.Join("test", "fs1")

	fs2 = NewContextFilesystem(afero.NewOsFs())
	fs2.RootPath = filepath.Join("test", "fs2")

	tests := []struct {
		desc        string
		filesystems []*ContextFilesystem
		override    bool
		assertFunc  func(*testing.T, *ContextFilesystem)
		err         error
	}{
		{
			desc:        "Testing error when merge nil filesystem",
			filesystems: []*ContextFilesystem{fs1, nil},
			err:         errors.New(errorContext, "Error trying join a nil filesystem"),
		},
		{
			desc:        "Testing error when joining and existing file without allowing file override",
			filesystems: []*ContextFilesystem{fs1, fs2},
			override:    false,
			err: errors.New(errorContext, "Error joinnig context filesystem",
				errors.New(errorContext, "File '/dir1/d1f1.txt' already on destination filesystem")),
		},
		{
			desc:        "Testing join filesystem",
			filesystems: []*ContextFilesystem{fs1, fs2},
			override:    true,
			err:         &errors.Error{},
			assertFunc: func(t *testing.T, fs *ContextFilesystem) {

				resultFiles := map[string]struct{}{}
				expectedFiles := map[string]struct{}{
					"/":              {},
					"/dir1":          {},
					"/dir1/d1f1.txt": {},
					"/dir2":          {},
					"/dir2/d1f1.txt": {},
					"/f1.txt":        {},
				}

				err = afero.Walk(fs, fs.RootPath, func(file string, fi os.FileInfo, err error) error {

					if err != nil {
						return errors.New("Testing join filesystem", fmt.Sprintf("Error during test. Walk trough '%s'", file))
					}

					resultFiles[file] = struct{}{}

					// if !fi.IsDir() {
					// 	f, _ := fs3.Fs.Open(file)
					// 	fmt.Println(file)
					// 	io.Copy(os.Stdout, f)
					// 	fmt.Println()
					// }

					return nil
				})
				if err != nil {
					t.Fatal(err.Error())
				}

				for efile, _ := range expectedFiles {
					_, exists := resultFiles[efile]
					assert.True(t, exists, fmt.Sprintf("Expected file '%s' does not exists", efile))
					delete(resultFiles, efile)
				}

				assert.Empty(t, resultFiles, fmt.Sprintf("Results contains unexpected file '%v'", resultFiles))

			},
		},
		{
			desc:        "Testing override file content",
			filesystems: []*ContextFilesystem{fs1, fs2},
			override:    true,
			err:         &errors.Error{},
			assertFunc: func(t *testing.T, fs *ContextFilesystem) {
				expectedContent := []byte("fs2 > dir1 > file1")
				content, _ := afero.ReadFile(fs, "/dir1/d1f1.txt")
				assert.Equal(t, expectedContent, content)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			fs3, err = Join(test.override, test.filesystems...)

			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				test.assertFunc(t, fs3)
			}

		})
	}

}