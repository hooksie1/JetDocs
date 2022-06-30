/*
Copyright © 2022 John Hooks

*/

package sync

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/hooksie1/jetdocs/backend"
)

type Syncer interface {
	Sync() error
}

type FileSync struct {
	Directory string
	Files     []*File
}

type File struct {
	Name     string
	Contents []byte
}

func (f *FileSync) ReadFile(name string) error {
	contents, err := os.ReadFile(name)
	if err != nil {
		return err
	}

	file := File{
		Name:     strings.TrimSuffix(name, filepath.Ext(name)),
		Contents: contents,
	}

	f.Files = append(f.Files, &file)

	return nil
}

func (f *FileSync) ReadAllFiles() error {
	files, err := os.ReadDir(f.Directory)
	if err != nil {
		return err
	}

	for _, v := range files {
		if v.IsDir() {
			continue
		}

		if filepath.Ext(v.Name()) != ".md" {
			continue
		}

		content, err := os.ReadFile(v.Name())
		if err != nil {
			return err
		}

		file := File{
			Name:     strings.TrimSuffix(v.Name(), filepath.Ext(v.Name())),
			Contents: content,
		}

		f.Files = append(f.Files, &file)
	}

	return nil
}

func (f *FileSync) Sync(b backend.Backend) error {
	for _, v := range f.Files {
		b.Write(v.Name, v.Contents)
	}

	return nil
}
