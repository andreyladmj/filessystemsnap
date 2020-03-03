package scanners

import (
	"andreyladmj/filessystemsnap/utils"
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
	rootDir           *utils.File
	filters           *utils.Filters
}

func NewDirsRecursiveScanner(maxGoroutines uint8, filters *utils.Filters) DirsRecursiveScanner {
	return DirsRecursiveScanner{wg: &sync.WaitGroup{}, maxGoroutines: maxGoroutines, filters: filters}
}

func (ds *DirsRecursiveScanner) Scan(path string) {
	ds.rootDir = &utils.File{Path: path}
	ds.ReadDir(ds.rootDir)
	ds.wg.Wait()
}

func (ds *DirsRecursiveScanner) DropEmptyDirs() {
	ds.rootDir.DropEmptyDirs()
}

func (ds *DirsRecursiveScanner) Print(f func(d *utils.File) string) error {
	if ds.rootDir == nil {
		return errors.New("Root Dir is nil")
	}
	fmt.Println(f(ds.rootDir))
	fmt.Println(ds.rootDir.CalcAllObjects() - 1)
	return nil
}

func (ds *DirsRecursiveScanner) ReadDir(dir *utils.File) {
	files, err := ioutil.ReadDir(dir.Path)
	if err != nil {
		//log.Fatal(err)
		return
	}
	for _, f := range files {
		if !ds.filters.Filter(f) {
			continue
		}

		fullpath := path.Join(dir.Path, f.Name())

		if f.IsDir() {
			d := utils.NewFile(f, fullpath)

			if ds.currentGoroutines < ds.maxGoroutines {
				ds.wg.Add(1)
				go func(d1 *utils.File) {
					ds.ReadDir(d1)
					defer ds.wg.Done()
				}(d)
			} else {
				ds.ReadDir(d)
			}
			dir.Append(d)
		} else {
			dir.Append(utils.NewFile(f, fullpath))
		}
	}
}
