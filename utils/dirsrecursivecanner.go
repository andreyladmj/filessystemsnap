package utils

import (
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
	rootDir           *Dir
	filters           *Filters
}

func NewDirsRecursiveScanner(maxGoroutines uint8, filters *Filters) DirsRecursiveScanner {
	return DirsRecursiveScanner{wg: &sync.WaitGroup{}, maxGoroutines: maxGoroutines, filters: filters}
}

func (ds *DirsRecursiveScanner) Scan(path string) {
	ds.rootDir = &Dir{path: path}
	ds.ReadDir(ds.rootDir)
	ds.wg.Wait()
}
func (ds *DirsRecursiveScanner) DropEmptyDirs() {
	ds.rootDir.DropEmptyDirs()
}

func (ds *DirsRecursiveScanner) Print(f func(d *Dir) string) error {
	if ds.rootDir == nil {
		return errors.New("Root Dir is nil")
	}
	fmt.Println(f(ds.rootDir))
	fmt.Println(ds.rootDir.CalcAllObjects() - 1)
	return nil
}

func (ds *DirsRecursiveScanner) ReadDir(dir *Dir) {
	files, err := ioutil.ReadDir(dir.path)
	if err != nil {
		//log.Fatal(err)
		return
	}
	for _, f := range files {
		if !ds.filters.filter(f) {
			continue
		}

		fullpath := path.Join(dir.path, f.Name())

		if f.IsDir() {
			d := NewDir(f, fullpath)

			if ds.currentGoroutines < ds.maxGoroutines {
				ds.wg.Add(1)
				go func(d1 *Dir) {
					ds.ReadDir(d1)
					defer ds.wg.Done()
				}(d)
			} else {
				ds.ReadDir(d)
			}
			dir.appendDir(d)
		} else {
			dir.appendFile(NewFile(f, fullpath))
		}
	}
}
