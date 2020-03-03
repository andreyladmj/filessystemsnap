package utils

import (
	"fmt"
	"os"
)

type File struct {
	Name  string
	Path  string
	Size  int
	IsDir bool
	Files []*File
}

func (f *File) Append(dir *File) {
	f.Files = append(f.Files, dir)
}

func (d *File) GetAllFilesCount() int {
	c := 0
	for _, file := range d.Files {
		if file.IsDir {
			c += file.GetAllFilesCount()
		} else {
			c++
		}

	}
	return c
}

func (d *File) DropEmptyDirs() {
	for i := 0; i < len(d.Files); i++ {
		if d.Files[i].GetAllFilesCount() == 0 {
			d.Files = append(d.Files[:i], d.Files[i+1:]...)
			i--
		} else {
			d.Files[i].DropEmptyDirs()
		}
	}
}

func (d *File) String() string {
	sf := ""
	sd := ""
	for _, f := range d.Files {
		if f.IsDir {
			sd += fmt.Sprintf("- %s\n", f.String())
		} else {
			sf += fmt.Sprintf("- %s\n", f.String())
		}
	}

	return sd + sf
}

func (d *File) GetFilesCount() int {
	c := 0
	for _, file := range d.Files {
		if !file.IsDir {
			c++
		}
	}
	return c
}

func (d *File) GetDirsCount() int {
	c := 0
	for _, file := range d.Files {
		if file.IsDir {
			c++
		}
	}
	return c
}

func (d *File) CalcAllObjects() int {
	n := 0

	for _, f := range d.Files {
		if f.IsDir {
			n += f.CalcAllObjects() + 1
		} else {
			n++
		}
	}

	return n
}

func NewFile(f os.FileInfo, fullpath string) *File {
	return &File{Name: f.Name(), Path: fullpath, Size: int(f.Size()), IsDir: f.IsDir()}
}
