package utils

import (
	"os"
	"strings"
)

type Filters struct {
	filters []func(f os.FileInfo) bool
}

func (f *Filters) Append(fn func(f os.FileInfo) bool) {
	f.filters = append(f.filters, fn)
}

func (f *Filters) Filter(file os.FileInfo) bool {
	for _, fn := range f.filters {
		if !fn(file) {
			return false
		}
	}

	return true
}

func (f *Filters) FileSizeFilter(n float64, t string) {
	switch strings.ToLower(t) {
	case "kb":
		n = 1024 * n
	case "mb":
		n = 1024 * 1024 * n
	case "gb":
		n = 1024 * 1024 * 1024 * n
	case "tb":
		n = 1024 * 1024 * 1024 * 1024 * n
	}

	f.Append(func(file os.FileInfo) bool {
		if !file.IsDir() {
			return file.Size() > int64(n)
		}
		return true
	})
}
