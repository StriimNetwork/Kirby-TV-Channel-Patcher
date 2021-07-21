package wad

// Copied from https://github.com/wii-tools/archon/blob/master/cmd/wad/directory.go
import (
	"errors"
	"fmt"
	"github.com/wii-tools/wadlib"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Type directory helps us avoid a ton of nested functions.
type directory struct {
	dir     string
	titleId string
}

// writeContents writes contents to the content's index within a directory.
func (d *directory) writeContents(content wadlib.WADFile, byId bool) error {
	// We default to reading by the index.
	name := ""
	if byId {
		name = fmt.Sprintf("%08x.app", content.ID)
	} else {
		name = fmt.Sprintf("%08x.app", content.Index)
	}
	return d.writeFile(name, content.RawData)
}

// writeSection writes contents to the content's index within a directory.
func (d *directory) writeSection(suffix string, content []byte) error {
	name := fmt.Sprintf("%s.%s", d.titleId, suffix)
	return d.writeFile(name, content)
}

// writeFile writes contents to the named file within a directory.
func (d *directory) writeFile(name string, contents []byte) error {
	path := filepath.Join(d.dir, name)
	return ioutil.WriteFile(path, contents, os.ModePerm)
}

// readFileWithSuffix reads one file from the directory with a given suffix and returns its contents.
func (d *directory) readFileWithSuffix(suffix string) ([]byte, error) {
	dir, err := ioutil.ReadDir(d.dir)
	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("%s.%s", d.titleId, suffix)

	var potentialFiles []os.FileInfo
	for _, file := range dir {
		// Only consider files with our prefix in this directory.
		if strings.HasSuffix(file.Name(), name) {
			potentialFiles = append(potentialFiles, file)
		}
	}

	// We can assume the file doesn't exist.
	if len(potentialFiles) == 0 {
		return nil, os.ErrNotExist
	} else if len(potentialFiles) != 1 {
		return nil, errors.New(fmt.Sprintf("more than one file found with suffix %s", suffix))
	}

	// We're now certain this file meets our needs.
	return d.readFile(potentialFiles[0].Name())
}

// readFile reads contents of the named file within a directory.
func (d *directory) readFile(name string) ([]byte, error) {
	path := filepath.Join(d.dir, name)
	return ioutil.ReadFile(path)
}

// readSection reads contents of the named file within a directory.
func (d *directory) readSection(suffix string) ([]byte, error) {
	name := fmt.Sprintf("%s.%s", d.titleId, suffix)
	path := filepath.Join(d.dir, name)
	return ioutil.ReadFile(path)
}
