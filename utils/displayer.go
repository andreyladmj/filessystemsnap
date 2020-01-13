package utils

import (
	"fmt"
)

func DisplayLikeTree(d *Dir) string {
	return PrintLikeTree(d, "")
}

func PrintLikeTree(d *Dir, pad string) string {
	s := ""
	for _, dir := range d.dirs {
		if dir.GetFilesCount() == 0 {
			continue
		}

		l := len(dir.files) + len(dir.dirs)
		s += fmt.Sprintf("%s%s: (%d files)\n%s", pad, dir.name, l, PrintLikeTree(dir, pad+"  "))
	}

	for _, f := range d.files {
		s += fmt.Sprintf("%s%s (%s)\n", pad, f.name, ByteFormat(float64(f.size), 2))
	}

	return s
}

func DisplayLikeTreeFormatted(d *Dir) string {
	return PrintLikeTreeFormatted(d, 0)
}

var e = make(map[int]bool, 15)

func PrintLikeTreeFormatted(d *Dir, nesting int) string {
	if nesting < 1 {
		nesting = 1
	}
	s := ""
	for i, dir := range d.dirs {
		e[nesting] = true

		if i == len(d.dirs)-1 && len(d.files) == 0 {
			e[nesting] = false
		}

		l := dir.GetFilesCount()
		p := getPrintedSymbol(i, len(d.dirs)+len(d.files))
		gap := getNestedLine(nesting)
		s += fmt.Sprintf("%s%s%s: (%d files)\n%s", gap, p, dir.name, l, PrintLikeTreeFormatted(dir, nesting+1))
	}

	for i, f := range d.files {
		p := getPrintedSymbol(i, len(d.files))
		gap := getNestedLine(nesting)
		s += fmt.Sprintf("%s%s%s (%s)\n", gap, p, f.name, ByteFormat(float64(f.size), 2))
	}

	return s
}

func getNestedLine(index int) string {
	s := ""
	for i := 0; i < index; i++ {
		if e[i] {
			s += "┃"
		} else {
			s += " "
		}
	}
	return s
}

func getPrintedSymbol(i, l int) string {
	if i < l-1 {
		return "┣"
	}
	return "┗"
}
