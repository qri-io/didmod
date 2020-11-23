package didmod

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPropsEqual(t *testing.T) {
	wd, err := ioutil.TempDir("", "didmod")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(wd)

	if _, err := NewProps(filepath.Join(wd, "does_not_exist")); err == nil {
		t.Errorf("expected NewProps on non-existant file to error. got nil")
	}

	filename := filepath.Join(wd, "eg.txt")

	if err := ioutil.WriteFile(filename, []byte("hello world"), 0644); err != nil {
		t.Fatal(err)
	}

	p, err := NewProps(filename)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := os.Stat(filename)
	if err != nil {
		t.Fatal(err)
	}
	b := NewPropsFileInfo(fi)

	if !p.Equal(b) {
		diff := cmp.Diff(p, b)
		t.Errorf("expected a to equal b. diff (-a +b):\n%s", diff)
	}

	t.Run("file_move", func(t *testing.T) {
		// If you mv a file to replace another, it will have a different inode number,
		// which we notice. It also probably has a different size and (even if not
		// newer than the target) mtime, any of which are sufficient.
		if err := os.Rename(filename, filepath.Join("eg_moved.txt")); err != nil {
			t.Fatal(err)
		}

		if err := ioutil.WriteFile(filename, []byte("hello world"), 0644); err != nil {
			t.Fatal(err)
		}

		b, err = NewProps(filename)
		if err != nil {
			t.Fatal(err)
		}

		if p.Equal(b) {
			diff := cmp.Diff(p, b)
			t.Errorf("expected a to NOT equal b. diff (-a +b):\n%s", diff)
		}
	})

}
