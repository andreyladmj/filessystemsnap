package main

import (
	"./utils"
)

func main() {
	f := &utils.Filters{}
	f.FileSizeFilter(200, "mb")

	maxGoroutines := uint8(10)
	ds := utils.NewDirsRecursiveScanner(maxGoroutines, f)
	ds.Scan("/home/ladyhin/")
	ds.DropEmptyDirs()
	ds.Print(utils.DisplayLikeTree2)

}
