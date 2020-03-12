package scanners

import (
	"andreyladmj/filessystemsnap/internal/platform"
	"fmt"
	//"sync"
	"syscall"
	"unsafe"
)

type FastScanner struct {
	RootDir *platform.File
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
		Files: ScanDir(dir),
	}
}

func ScanDir(dir string) []*platform.File {
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

			var fs syscall.Stat_t
			statErr := syscall.Lstat(name, &fs)
			size := 0

			if statErr == nil {
				size = int(fs.Size)
			}

			fullPath := dir + "/" + name
			file := new(platform.File)

			if dirent.Type == syscall.DT_REG {
				file = platform.NewFile(name, fullPath, size, false)
			} else if dirent.Type == syscall.DT_DIR {
				file = platform.NewFile(name, fullPath, size, true)
				file.Files = ScanDir(fullPath)
			}

			files = append(files, file)
		}
	}

	return files
}
