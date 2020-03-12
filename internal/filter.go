package internal

import (
	"andreyladmj/filessystemsnap/internal/platform"
	"strings"
)

type Filters struct {
	filters []func(f platform.File) bool
}

func (f *Filters) Append(fn func(f platform.File) bool) {
	f.filters = append(f.filters, fn)
}

func (f *Filters) Filter(file platform.File) bool {
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

	f.Append(func(file platform.File) bool {
		if !file.IsDir {
			return file.Size > int(n)
		}
		return true
	})
}
