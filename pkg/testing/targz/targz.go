package targz

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

func Untar(repositoryName string) (string, func(), error) {
	_, filename, _, _ := runtime.Caller(0)
	directory := path.Dir(filename)
	destination := filepath.Join(directory, "fixtures", "generated-data", repositoryName)
	reader, err := os.Open(filepath.Join(directory, "fixtures", fmt.Sprintf("%s.tgz", repositoryName)))
	defer reader.Close()

	cleanup := func() {
		os.RemoveAll(destination)
	}

	if err != nil {
		return "", nil, err
	}

	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return "", nil, err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()

		switch {
		case err == io.EOF:
			return destination, cleanup, nil
		case err != nil:
			return "", nil, err
		}

		target := filepath.Join(destination, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return "", nil, err
				}
			}

		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, 0644)
			if err != nil {
				return "", nil, err
			}

			if _, err := io.Copy(f, tarReader); err != nil {
				return "", nil, err
			}

			f.Close()
		}
	}
}
