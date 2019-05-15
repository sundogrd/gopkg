package file_test

import (
	"github.com/sundogrd/gopkg/file"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTempDir(t *testing.T) {
	dir, err := file.CreateTempDir("prefix")
	if err != nil {
		t.Fatal(err)
	}

	// Delete the directory no matter what.
	defer dir.Delete()

	// Make sure we have created the dir, with the proper name
	_, err = os.Stat(dir.Name)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(filepath.Base(dir.Name), "prefix") {
		t.Fatalf(`Directory doesn't start with "prefix": %q`,
			dir.Name)
	}

	// Verify that the directory is empty
	entries, err := ioutil.ReadDir(dir.Name)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("Directory should be empty, has %d elements",
			len(entries))
	}

	// Create a couple of files
	_, err = dir.NewFile("ONE")
	if err != nil {
		t.Fatal(err)
	}
	_, err = dir.NewFile("TWO")
	if err != nil {
		t.Fatal(err)
	}
	// We can't create the same file twice
	_, err = dir.NewFile("TWO")
	if err == nil {
		t.Fatal("NewFile should fail to create the same file twice")
	}

	// We have created only two files
	entries, err = ioutil.ReadDir(dir.Name)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("ReadDir should have two elements, has %d elements",
			len(entries))
	}

	// Verify that deletion works
	err = dir.Delete()
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat(dir.Name)
	if err == nil {
		t.Fatal("Directory should be gone, still present.")
	}
}