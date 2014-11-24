package util

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type myzip struct{}

// Zip utilities.
var Zip = myzip{}

type ZipFile struct {
	zipFile *os.File
	writer  *zip.Writer
}

func (*myzip) Create(filename string) (*ZipFile, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return &ZipFile{zipFile: file, writer: zip.NewWriter(file)}, nil
}

func (z *ZipFile) Close() error {
	err := z.writer.Close()
	if nil != err {
		return err
	}

	return z.zipFile.Close() // close the underlying writer
}

func (z *ZipFile) AddEntryN(path string, names ...string) error {
	for _, name := range names {
		zipPath := filepath.Join(path, name)
		err := z.AddEntry(zipPath, name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (z *ZipFile) AddEntry(path, name string) error {
	fi, err := os.Stat(name)
	if err != nil {
		return err
	}

	fh, err := zip.FileInfoHeader(fi)
	if err != nil {
		return err
	}

	fh.Name = filepath.ToSlash(filepath.Clean(path))
	fh.Method = zip.Deflate // data compression algorithm

	if fi.IsDir() {
		fh.Name = fh.Name + "/" // be care the ending separator
	}

	entry, err := z.writer.CreateHeader(fh)
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return nil
	}

	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(entry, file)

	return err
}

func (z *ZipFile) AddDirectoryN(path string, names ...string) error {
	for _, name := range names {
		err := z.AddDirectory(path, name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (z *ZipFile) AddDirectory(path, dirName string) error {
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		return err
	}

	for _, file := range files {
		localPath := filepath.Join(dirName, file.Name())
		zipPath := filepath.Join(path, file.Name())

		err = nil
		if file.IsDir() {
			z.AddEntry(path, dirName)

			err = z.AddDirectory(zipPath, localPath)
		} else {
			err = z.AddEntry(zipPath, localPath)
		}
		if err != nil {
			return err
		}
	}

	return nil
}
