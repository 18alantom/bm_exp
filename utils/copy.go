package utils

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
)

func CopyDir(source string, dest string) error {
	return copyFS(
		dest,
		os.DirFS(source),
	)
}

// Source: https://github.com/golang/go/commit/d9be60974b694a17e5c6c3e71fb7767e6bfe17e9
// TODO:
// - Vet the implementation (symlinks?), use another if required
// - Move this into a proper place, remove when Go 1.23
func copyFS(dir string, fsys fs.FS) error {
	return fs.WalkDir(fsys, ".", func(path_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		fpath, err := fromFS(path_)
		if err != nil {
			return err
		}
		newPath := path.Join(dir, fpath)
		if d.IsDir() {
			return os.MkdirAll(newPath, 0777)
		}

		// TODO(panjf2000): handle symlinks with the help of fs.ReadLinkFS
		// 		once https://go.dev/issue/49580 is done.
		//		we also need safefilepath.IsLocal from https://go.dev/cl/564295.
		if !d.Type().IsRegular() {
			return &fs.PathError{Op: "CopyFS", Path: path_, Err: fs.ErrInvalid}
		}

		r, err := fsys.Open(path_)
		if err != nil {
			return err
		}
		defer r.Close()
		info, err := r.Stat()
		if err != nil {
			return err
		}
		w, err := os.OpenFile(newPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666|info.Mode()&0777)
		if err != nil {
			return err
		}

		if _, err := io.Copy(w, r); err != nil {
			w.Close()
			return &fs.PathError{Op: "Copy", Path: newPath, Err: err}
		}
		return w.Close()
	})
}

var errInvalidPath = errors.New("invalid path")

// Source: https://cs.opensource.google/go/go/+/refs/tags/go1.22.5:src/internal/safefilepath/path.go
func fromFS(path string) (string, error) {
	for i := range path {
		if path[i] == 0 {
			return "", errInvalidPath
		}
	}
	return path, nil
}
