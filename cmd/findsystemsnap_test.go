package cmd

import (
	"andreyladmj/filessystemsnap/scanners"
	"andreyladmj/filessystemsnap/utils"
	"testing"
)

func BenchmarkScanRecursiveDirs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f := &utils.Filters{}
		ds := scanners.NewDirsRecursiveScanner(8, f)
		ds.Scan("/home/ladyhin/Projects")
	}
}

func BenchmarkScanDirsWithChannels(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ds := scanners.NewDirsChanScanner(8)
		ds.Scan("/home/ladyhin/Projects")
	}
}

//func BenchmarkScanDirsWithChannels(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		ds := scanners.NewDirsChanScanner(8)
//		ds.Scan("/home/ladyhin")
//	}
//}
