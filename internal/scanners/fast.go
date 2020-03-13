package scanners

import (
	"andreyladmj/filessystemsnap/internal"
	"andreyladmj/filessystemsnap/internal/platform"
	"fmt"
	"syscall"
	"unsafe"
)

type FastScanner struct {
	RootDir *platform.File
	filters *internal.Filters
}

func NewFastScanner() FastScanner {
	return FastScanner{}
}

var blockSize = 8192

func (fs *FastScanner) Scan(dir string) {
	fs.RootDir = &platform.File{
		Name:  dir,
		Size:  0,
		IsDir: true,
		Files: fs.ScanDir(dir),
	}
}

func (fs *FastScanner) ScanDir(dir string) []*platform.File {
	sysfd, err := syscall.Open(dir, syscall.O_RDONLY|syscall.O_NOCTTY|syscall.O_NONBLOCK|syscall.O_NOFOLLOW|syscall.O_CLOEXEC|syscall.O_DIRECTORY, 0)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer syscall.Close(sysfd)

	origbuf := make([]byte, blockSize)
	files := make([]*platform.File, 0)

	for {
		n, err := syscall.ReadDirent(sysfd, origbuf)

		if err != nil {
			fmt.Println(err)
			return nil
		}

		if n <= 0 {
			break // EOF
		}

		buf := origbuf[0:n]

		for len(buf) > 0 {
			dirent := (*syscall.Dirent)(unsafe.Pointer(&buf[0]))

			rec := buf[:dirent.Reclen]
			buf = buf[dirent.Reclen:]

			if dirent.Ino == 0 { // File absent in directory.
				continue
			}

			const namoff = uint64(unsafe.Offsetof(dirent.Name))
			namlen := uint64(dirent.Reclen) - namoff

			if namoff+namlen > uint64(len(rec)) {
				break
			}

			bname := rec[namoff : namoff+namlen]
			for i, c := range bname {
				if c == 0 {
					bname = bname[:i]
					break
				}
			}

			name := string(bname)

			if name == "." || name == ".." { // Useless names
				continue
			}

			var sfs syscall.Stat_t
			statErr := syscall.Lstat(name, &sfs)
			size := 0

			if statErr == nil {
				size = int(sfs.Size)
			}

			fullPath := dir + "/" + name
			file := new(platform.File)

			file = platform.NewFile(name, fullPath, size, dirent.Type == syscall.DT_DIR) // syscall.DT_REG - reg file

			if fs.filters != nil && !fs.filters.Filter(file) {
				continue
			}

			if file.IsDir {
				file.Files = fs.ScanDir(fullPath)
			}

			files = append(files, file)
		}
	}

	return files
}
