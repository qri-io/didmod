// +build !windows

package didmod

import (
	"os"
	"syscall"
)

// GetIds returns ids for the file: the inode, the volume, user id, group id
func GetIds(filename string, fi os.FileInfo) (uint64, uint32, uint32, uint32) {
	if stat, ok := fi.Sys().(*syscall.Stat_t); ok {
		return stat.Ino, 0, stat.Uid, stat.Gid
	}
	return 0, 0, 0, 0
}
