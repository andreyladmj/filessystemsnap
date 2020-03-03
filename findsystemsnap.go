package main

import (
	"andreyladmj/filessystemsnap/scanners"
	"andreyladmj/filessystemsnap/utils"
)

func main() {
	//f := &utils.Filters{}
	//f.FileSizeFilter(200, "mb")
	//
	//maxGoroutines := uint8(10)
	//ds := utils.NewDirsRecursiveScanner(maxGoroutines, f)
	//ds.Scan("/home/ladyhin/")
	//ds.DropEmptyDirs()
	//ds.Print(utils.DisplayLikeTree2)

	//scanner := scanners.FastScanner{}
	//scanner.Scan("/home/ladyhin/")
	//utils.DisplayLikeTreeFormatted(scanner.RootDir)

	f := &utils.Filters{}

	scanner := scanners.NewDirsRecursiveScanner(8, f)
	scanner.Scan("/home/ladyhin/apt")
	//scanner.Print(utils.DisplayLikeTreeFormatted)
	scanner.Print(utils.DisplayInline)
}
