package scanners

import (
	"andreyladmj/filessystemsnap/utils"
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
	rootDir           *utils.File
}

func NewDirsChanScanner(maxGoroutines uint8) DirsChanScanner {
	return DirsChanScanner{wg: &sync.WaitGroup{}, maxGoroutines: maxGoroutines}
}

func (ds *DirsChanScanner) Scan(path string) {
	ds.rootDir = &utils.File{Path: path}
	ds.ReadDir(ds.rootDir)
	ds.wg.Wait()
}

func (ds *DirsChanScanner) Print(f func(d *utils.File) string) error {
	if ds.rootDir == nil {
		return errors.New("Root Dir is nil")
	}
	fmt.Println(f(ds.rootDir))
	fmt.Println(ds.rootDir.CalcAllObjects() - 1)
	return nil
}

func (ds *DirsChanScanner) ReadDir(dir *utils.File) {
	ch := make(chan *utils.File)

	go func(dirname string, ch chan *utils.File) {
		ds.wg.Add(1)
		defer ds.wg.Done()

		files, err := ioutil.ReadDir(dir.Path)
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {
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
				ch <- d
			} else {
				ch <- utils.NewFile(f, fullpath)
			}
		}
		close(ch)

	}(dir.Path, ch)

	for file := range ch {
		dir.Append(file)
	}
}
