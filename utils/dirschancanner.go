package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"sync"
)

type DirsChanScanner struct {
	wg                *sync.WaitGroup
	maxGoroutines     uint8
	currentGoroutines uint8
	rootDir           *Dir
}

func NewDirsChanScanner(maxGoroutines uint8) DirsChanScanner {
	return DirsChanScanner{wg: &sync.WaitGroup{}, maxGoroutines: maxGoroutines}
}

func (ds *DirsChanScanner) Scan(path string) {
	ds.rootDir = &Dir{path: path}
	ds.ReadDir(ds.rootDir)
	ds.wg.Wait()
}

func (ds *DirsChanScanner) Print(f func(d *Dir) string) error {
	if ds.rootDir == nil {
		return errors.New("Root Dir is nil")
	}
	fmt.Println(f(ds.rootDir))
	fmt.Println(ds.rootDir.CalcAllObjects() - 1)
	return nil
}

func (ds *DirsChanScanner) ReadDir(dir *Dir) {
	ch := make(chan interface{})

	go func(dirname string, ch chan interface{}) {
		ds.wg.Add(1)
		defer ds.wg.Done()

		files, err := ioutil.ReadDir(dir.path)
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {
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
				ch <- d
				//dir.appendDir(d)
			} else {
				ch <- NewFile(f, fullpath)
				//dir.appendFile()
			}
		}
		close(ch)

	}(dir.path, ch)

	for subd := range ch {
		fmt.Println(subd)
		if d, ok := subd.(*Dir); ok {
			dir.appendDir(d)
		} else if f, ok := subd.(*File); ok {
			dir.appendFile(f)
		} else {
			log.Fatal("cannot convert types")
		}
	}
}
