package cmd

import (
	"andreyladmj/filessystemsnap/internal/scanners"
	"testing"
)

var DIR = "/home/ladyhin/Projects"

func BenchmarkScanDirsWithChannels8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ds := scanners.NewDirsChanScanner(8)
		ds.Scan(DIR)
	}
}

func BenchmarkScanDirsWithChannels64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ds := scanners.NewDirsChanScanner(64)
		ds.Scan(DIR)
	}
}

func BenchmarkScanRecursiveDirs8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ds := scanners.NewDirsRecursiveScanner(8)
		ds.Scan(DIR)
	}
}

func BenchmarkScanRecursiveDirs64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ds := scanners.NewDirsRecursiveScanner(64)
		ds.Scan(DIR)
	}
}

func BenchmarkFastScanDirs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ds := scanners.NewFastScanner()
		ds.Scan(DIR)
	}
}
