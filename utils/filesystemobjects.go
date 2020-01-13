package utils

import (
	"fmt"
	"os"
)

type File struct {
	name string
	path string
	size int64
}

func (f *File) String() string {
	return fmt.Sprintf("- %s", f.name)
}

type Dir struct {
	files []*File
	dirs  []*Dir
	name  string
	path  string
}

func (d *Dir) appendFile(f *File) {
	d.files = append(d.files, f)
}

func (d *Dir) appendDir(dir *Dir) {
	d.dirs = append(d.dirs, dir)
}

func (d *Dir) GetFilesCount() int {
	c := len(d.files)
	for _, dir := range d.dirs {
		c += dir.GetFilesCount()
	}
	return c
}

func (d *Dir) DropEmptyDirs() {
	for i := 0; i < len(d.dirs); i++ {
		if d.dirs[i].GetFilesCount() == 0 {
			d.dirs = append(d.dirs[:i], d.dirs[i+1:]...)
			i--
		} else {
			d.dirs[i].DropEmptyDirs()
		}
	}
}

func (d *Dir) String() string {
	s := ""
	for _, dir := range d.dirs {
		s += fmt.Sprintf("- %s\n", dir.String())
	}

	for _, f := range d.files {
		s += fmt.Sprintf("  %s\n", f.name)
	}

	return s
}

func (d *Dir) CalcAllObjects() int {
	n := len(d.files)

	for _, dir := range d.dirs {
		n += dir.CalcAllObjects() + 1
	}

	return n
}

func NewDir(f os.FileInfo, fullpath string) *Dir {
	return &Dir{name: f.Name(), path: fullpath}
}

func NewFile(f os.FileInfo, fullpath string) *File {
	return &File{name: f.Name(), path: fullpath, size: f.Size()}
}
