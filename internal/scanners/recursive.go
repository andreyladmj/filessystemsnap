package scanners

import (
	"andreyladmj/filessystemsnap/internal"
	"andreyladmj/filessystemsnap/internal/platform"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"sync"
)

type DirsRecursiveScanner struct {
	wg                *sync.WaitGroup
	maxGoroutines     uint8
	currentGoroutines uint8
	rootDir           *platform.File
	filters           *internal.Filters
}

func NewDirsRecursiveScanner(maxGoroutines uint8, filters *internal.Filters) DirsRecursiveScanner {
	return DirsRecursiveScanner{wg: &sync.WaitGroup{}, maxGoroutines: maxGoroutines, filters: filters}
}

func (ds *DirsRecursiveScanner) Scan(path string) {
	ds.rootDir = &platform.File{Path: path}
	ds.ReadDir(ds.rootDir)
	ds.wg.Wait()
}

func (ds *DirsRecursiveScanner) DropEmptyDirs() {
	ds.rootDir.DropEmptyDirs()
}

func (ds *DirsRecursiveScanner) Print(f func(d *platform.File) string) error {
	if ds.rootDir == nil {
		return errors.New("Root Dir is nil")
	}
	fmt.Println(f(ds.rootDir))
	fmt.Println(ds.rootDir.CalcAllObjects() - 1)
	return nil
}

func (ds *DirsRecursiveScanner) ReadDir(dir *platform.File) {
	files, err := ioutil.ReadDir(dir.Path)
	if err != nil {
		fmt.Println("ERROR", err)
		return
	}
	for _, f := range files {
		if !ds.filters.Filter(f) {
			continue
		}

		fullpath := path.Join(dir.Path, f.Name())

		if f.IsDir() {
			d := platform.NewFile(f.Name(), fullpath, int(f.Size()), f.IsDir())

			if ds.currentGoroutines < ds.maxGoroutines {
				ds.wg.Add(1)
				go func(d1 *platform.File) {
					ds.ReadDir(d1)
					defer ds.wg.Done()
				}(d)
			} else {
				ds.ReadDir(d)
			}
			dir.Append(d)
		} else {
			dir.Append(platform.NewFile(f.Name(), fullpath, int(f.Size()), f.IsDir()))
		}
	}
}
