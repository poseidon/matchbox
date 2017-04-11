package storage

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDir(t *testing.T) {
	cases := []struct {
		path     string
		expected string
	}{
		{"/a", "a"},
		{"../b", "b"},
		{"..c", "..c"},
		// creates parent directories as needed
		{"d/e/ff", "d/e/ff"},
		{"d/e/ff/../gg", "d/e/gg"},
	}
	tdir, err := ioutil.TempDir("", "matchbox")
	assert.Nil(t, err)
	defer os.RemoveAll(tdir)

	// create a Dir, restricted to the temp dir
	dir := Dir(tdir)
	// write files rooted in the dir
	for _, c := range cases {
		dir.writeFile(c.path, []byte(c.expected))
	}
	// ensure expected files were created
	for _, c := range cases {
		_, err := os.Stat(filepath.Join(tdir, c.expected))
		assert.Nil(t, err)
	}
	// ensure expected files can be read by rooted dir
	for _, c := range cases {
		b, err := dir.readFile(c.path)
		assert.Nil(t, err)
		assert.Equal(t, []byte(c.expected), b)
	}
	// delete the files that were written
	for _, c := range cases {
		err := dir.deleteFile(c.path)
		assert.Nil(t, err)
	}
	// ensure the expected files were removed
	for _, c := range cases {
		rpath := filepath.Join(tdir, c.expected)
		_, err := os.Stat(rpath)
		assert.True(t, os.IsNotExist(err), "expected path %s would not exist", rpath)
	}
}

func TestSanitizePath(t *testing.T) {
	cases := []struct {
		dir      Dir
		path     string
		expected string
		err      error
	}{
		{Dir(""), "", ".", nil},
		{Dir(""), "..", ".", nil},
		{Dir(""), "../../", ".", nil},
		{Dir("."), "", ".", nil},
		{Dir("."), "/../../", ".", nil},
		{Dir("/etc"), "/hosts", "/etc/hosts", nil},
		{Dir("/etc"), "hosts", "/etc/hosts", nil},
		{Dir("/etc"), "../../../hosts", "/etc/hosts", nil},
		{Dir("/etc/"), "/hosts", "/etc/hosts", nil},
		{Dir("/etc/"), "hosts", "/etc/hosts", nil},
		{Dir("/etc/"), "../../../hosts", "/etc/hosts", nil},
		// zero byte - don't repeat Go bug https://github.com/golang/go/issues/3842
		{Dir("/etc/"), "/..\x00", "", errInvalidFilePathCharacter},
	}
	for _, c := range cases {
		path, err := c.dir.sanitize(c.path)
		assert.Equal(t, c.expected, path)
		assert.Equal(t, c.err, err)
	}
}
