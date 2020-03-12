package scanners

import (
	"andreyladmj/filessystemsnap/internal/platform"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"sync"
)

type DirsChanScanner struct {
	wg                *sync.WaitGroup
	maxGoroutines     uint8
	currentGoroutines uint8
	rootDir           *platform.File
}

func NewDirsChanScanner(maxGoroutines uint8) DirsChanScanner {
	return DirsChanScanner{wg: &sync.WaitGroup{}, maxGoroutines: maxGoroutines}
}

func (ds *DirsChanScanner) Scan(path string) {
	ds.rootDir = &platform.File{Path: path}
	ds.ReadDir(ds.rootDir)
	ds.wg.Wait()
}

func (ds *DirsChanScanner) Print(f func(d *platform.File) string) error {
	if ds.rootDir == nil {
		return errors.New("Root Dir is nil")
	}
	fmt.Println(f(ds.rootDir))
	fmt.Println(ds.rootDir.CalcAllObjects() - 1)
	return nil
}

func (ds *DirsChanScanner) ReadDir(dir *platform.File) {
	ch := make(chan *platform.File)

	go func(dirname string, ch chan *platform.File) {
		ds.wg.Add(1)
		defer ds.wg.Done()

		files, err := ioutil.ReadDir(dir.Path)
		if err != nil {
			fmt.Println("ERROR", err)
			return
		}

		for _, f := range files {
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
				ch <- d
			} else {
				ch <- platform.NewFile(f.Name(), fullpath, int(f.Size()), f.IsDir())
			}
		}
		close(ch)

	}(dir.Path, ch)

	for file := range ch {
		dir.Append(file)
	}
}
