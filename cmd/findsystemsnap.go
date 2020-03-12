package cmd

import "andreyladmj/filessystemsnap/internal/scanners"

//
//type customError struct{}
//
//func (c *customError) Error() string {
//	return "Find the bug"
//}
//
//func fail() ([]byte, *customError) {
//	return nil, nil
//}

func main() {
	//var err error
	//_, err = fail()
	//
	//print(err)
	//fmt.Println(err)
	//
	//if err != nil {
	//	log.Fatal("why did this fail?")
	//}
	//
	//log.Println("No Error")

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

	//f := &utils.Filters{}
	//
	//scanner := scanners.NewDirsRecursiveScanner(8, f)
	//scanner.Scan("/home/ladyhin/Downloads")
	//scanner.Print(utils.DisplayLikeTreeFormatted)
	//scanner.Print(utils.DisplayInline)

	scanner := scanners.NewFastScanner()
	scanner.Scan("/home/ladyhin/Downloads")
}
