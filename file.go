// Package didmod implements a robust strategy for checking if a file may have
// changed, based on @apenwarr's blog post "mtime comparison considered harmful"
// https://apenwarr.ca/log/20181113
// From the article:
//    My implementation of redo remembers the following information about each
//    source and target:
//      * mtime
//      * size
//      * inode number
//      * file mode
//      * owner uid and gid
//      * (targets only) the sequence number of the last time it was built
//   redo considers a dependency dirty if any of those attributes changed since
//   the last time a target was built
package didmod

import (
	"os"
	"time"
)

// Props are tracked properties of a file used for checking for changes
type Props struct {
	Name  string      `json:"name,omitempty"`
	Mtime time.Time   `json:"mtime,omitempty"`
	Size  int64       `json:"size,omitempty"`
	Mode  os.FileMode `json:"mode,omitempty"`

	Inode   uint64 `json:"inode,omitempty"`
	Volume  uint32 `json:"volume",omitempty"`
	UserID  uint32 `json:"userid,omitempty"`
	GroupID uint32 `json:"groupid,omitempty"`
}

// NewProps constructs Props by 'stat'-ing a file
func NewProps(filename string) (p Props, err error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return p, err
	}
	return NewPropsFileInfo(filename, fi), nil
}

// NewPropsFileInfo constructs props from an os.FileInfo, filling in as many
// fields as possible
func NewPropsFileInfo(filename string, fi os.FileInfo) Props {
	p := Props{
		Name:  fi.Name(),
		Size:  fi.Size(),
		Mtime: fi.ModTime().In(time.UTC),
		Mode:  fi.Mode(),
	}
	p.Inode, p.Volume, p.UserID, p.GroupID = GetIds(filename, fi)
	return p
}

// Equal checks for strict equality between two props structs
func (p Props) Equal(b Props) bool {
	return p.Mtime == b.Mtime &&
		p.Size == b.Size &&
		p.Inode == b.Inode &&
		p.Volume == b.Volume &&
		p.Mode == b.Mode &&
		p.UserID == b.UserID &&
		p.GroupID == b.GroupID
}
