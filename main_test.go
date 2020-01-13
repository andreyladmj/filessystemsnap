package main

import "testing"
import "./utils"

func BenchmarkScanDirs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		maxGoroutines := uint8(0)
		ds := utils.NewDirScanner(maxGoroutines)
		ds.Scan("/home/ladyhin/apt/GoLand-2019.1.1")
		ds.Wait()
	}
}
