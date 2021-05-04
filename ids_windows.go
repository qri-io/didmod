// +build windows

package didmod

import (
	"os"
	"syscall"
)

// GetIds returns ids for the file: the inode, the volume, user id, group id
// Adapted from github.com/fsnotify/fsnotify/windows.go
func GetIds(filename string, _ os.FileInfo) (uint64, uint32, uint32, uint32) {
	h, e := syscall.CreateFile(
		syscall.StringToUTF16Ptr(filename),
		syscall.FILE_LIST_DIRECTORY,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_FLAG_BACKUP_SEMANTICS|syscall.FILE_FLAG_OVERLAPPED,
		0)
	if e != nil {
		return 0, 0, 0, 0
	}
	var fi syscall.ByHandleFileInformation
	if e = syscall.GetFileInformationByHandle(h, &fi); e != nil {
		syscall.CloseHandle(h)
		return 0, 0, 0, 0
	}
	volume := fi.VolumeSerialNumber
	inode := uint64(fi.FileIndexHigh)<<32 | uint64(fi.FileIndexLow)
	// https://stackoverflow.com/questions/58179647/getting-uid-and-gid-of-a-file
	return inode, volume, uint32(os.Getuid()), uint32(os.Getgid())
}
