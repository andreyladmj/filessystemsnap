package scanners

import (
	"andreyladmj/filessystemsnap/utils"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"
	"unsafe"
)

type FastScanner struct {
	RootDir *utils.File
}

var blockSize = 8192
var _AT_SYMLINK_NOFOLLOW = 0x100

func (fs *FastScanner) Scan(dir string) {
	fs.RootDir = &utils.File{
		Name:  dir,
		Size:  0,
		IsDir: true,
		Files: ScanDir(dir),
	}
}

func ScanDir(dir string) []*utils.File {
	fd, err := os.OpenFile(dir, syscall.O_RDONLY|syscall.O_NOCTTY|syscall.O_NONBLOCK|syscall.O_NOFOLLOW|syscall.O_CLOEXEC|syscall.O_DIRECTORY, 0)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer fd.Close()

	buf := make([]byte, blockSize)

	sysfd := fd.Fd()

	n, err := syscall.ReadDirent(int(sysfd), buf)

	names := make([]string, 0, 100)

	//dirent := (*syscall.Dirent)(unsafe.Pointer(&buf[0]))
	//
	//fmt.Println("dirent", dirent)

	_, _, names = syscall.ParseDirent(buf, n, names)

	files := make([]*utils.File, 0, len(names))

	for _, name := range names {
		path := fmt.Sprintf("%s/%s", dir, name)
		f_stat := getFStat(path, int(sysfd))

		file := &utils.File{
			Name:  f_stat.name,
			Size:  int(f_stat.size),
			IsDir: f_stat.mode.IsDir(),
		}

		//fmt.Printf("%s/%s\n", dir, f_stat.name)
		//outBuf.WriteString(fmt.Sprintf("%s/%s\n", dir, f_stat.name))
		if f_stat.mode.IsDir() {
			file.Files = ScanDir(path)
		}

		files = append(files, file)
	}

	return files
}

func getFStat(dir string, sysfd int) *fileStat {
	_p0, err := syscall.BytePtrFromString(dir)
	if err != nil {
		log.Fatal(err)
	}
	sys_stat := new(syscall.Stat_t)
	f_stat := new(fileStat)

	_, _, e1 := syscall.Syscall6(syscall.SYS_NEWFSTATAT, uintptr(sysfd), uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(sys_stat)), uintptr(_AT_SYMLINK_NOFOLLOW), 0, 0)

	if e1 != 0 {
		log.Fatal(e1)
	}

	fillFileStatFromSys(f_stat, sys_stat, dir)
	return f_stat
}

type fileStat struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func fillFileStatFromSys(fs *fileStat, sysfs *syscall.Stat_t, name string) {
	fs.name = basename(name)
	fs.size = sysfs.Size
	fs.modTime = timespecToTime(sysfs.Mtim)
	fs.mode = os.FileMode(sysfs.Mode & 0777)
	switch sysfs.Mode & syscall.S_IFMT {
	case syscall.S_IFBLK:
		fs.mode |= os.ModeDevice
	case syscall.S_IFCHR:
		fs.mode |= os.ModeDevice | os.ModeCharDevice
	case syscall.S_IFDIR:
		fs.mode |= os.ModeDir
	case syscall.S_IFIFO:
		fs.mode |= os.ModeNamedPipe
	case syscall.S_IFLNK:
		fs.mode |= os.ModeSymlink
	case syscall.S_IFREG:
		// nothing to do
	case syscall.S_IFSOCK:
		fs.mode |= os.ModeSocket
	}
	if sysfs.Mode&syscall.S_ISGID != 0 {
		fs.mode |= os.ModeSetgid
	}
	if sysfs.Mode&syscall.S_ISUID != 0 {
		fs.mode |= os.ModeSetuid
	}
	if sysfs.Mode&syscall.S_ISVTX != 0 {
		fs.mode |= os.ModeSticky
	}
}

func basename(name string) string {
	i := len(name) - 1
	// Remove trailing slashes
	for ; i > 0 && name[i] == '/'; i-- {
		name = name[:i]
	}
	// Remove leading directory name
	for i--; i >= 0; i-- {
		if name[i] == '/' {
			name = name[i+1:]
			break
		}
	}

	return name
}

func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}
