package platform

import (
	"fmt"
)

func DisplayLikeTree(d *File) string {
	return PrintLikeTree(d, "")
}

func PrintLikeTree(d *File, pad string) string {
	s := ""
	for _, dir := range d.Files {
		if dir.GetAllFilesCount() == 0 {
			continue
		}

		l := len(dir.Files)
		s += fmt.Sprintf("%s%s: (%d files)\n%s", pad, dir.Name, l, PrintLikeTree(dir, pad+"  "))
	}

	for _, f := range d.Files {
		s += fmt.Sprintf("%s%s (%s)\n", pad, f.Name, ByteFormat(float64(f.Size), 2))
	}

	return s
}

func DisplayInline(d *File) string {
	s := ""
	for _, dir := range d.Files {
		s += fmt.Sprintf("%v\n", dir.Path)
		if dir.IsDir {
			s += DisplayInline(dir)
		}
	}
	return s
}

func DisplayLikeTreeFormatted(d *File) string {
	return PrintLikeTreeFormatted(d, 0)
}

var e = make(map[int]bool, 15)

func PrintLikeTreeFormatted(d *File, nesting int) string {
	if nesting < 1 {
		nesting = 1
	}
	s := ""

	for i, dir := range d.Files {
		if !dir.IsDir {
			continue
		}

		e[nesting] = true

		if i == d.GetDirsCount()-1 && d.GetFilesCount() == 0 {
			e[nesting] = false
		}

		l := dir.GetAllFilesCount()
		p := getPrintedSymbol(i, d.GetDirsCount()+d.GetFilesCount())
		gap := getNestedLine(nesting)
		s += fmt.Sprintf("%s%s%s: (%d files)\n%s", gap, p, dir.Name, l, PrintLikeTreeFormatted(dir, nesting+1))
	}

	for i, f := range d.Files {
		if f.IsDir {
			continue
		}

		p := getPrintedSymbol(i, d.GetFilesCount()+d.GetDirsCount())
		gap := getNestedLine(nesting)
		s += fmt.Sprintf("%s%s%s (%s)\n", gap, p, f.Name, ByteFormat(float64(f.Size), 2))
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
