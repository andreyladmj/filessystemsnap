package cmd

import "andreyladmj/filessystemsnap/internal/scanners"

func main() {
	scanner := scanners.NewFastScanner()
	scanner.Scan("/home/user/Downloads")
}
