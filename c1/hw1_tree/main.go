package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func dirTreeWithOffset(
	out *bytes.Buffer,
	path string,
	printFiles bool,
	prefixes []string,
) error {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	var dirCount int
	for _, entryInfo := range entries {
		if entryInfo.IsDir() {
			dirCount++
		}
	}

	var dirCursor int
	for i, entryInfo := range entries {
		isDir := entryInfo.IsDir()
		if isDir {
			dirCursor++
		} else if !printFiles {
			continue
		}

		for _, prefix := range prefixes {
			fmt.Fprint(out, prefix)
		}

		var last bool
		if printFiles {
			last = len(entries) == i+1
		} else {
			last = dirCount == dirCursor
		}
		if last {
			fmt.Fprint(out, "└───")
		} else {
			fmt.Fprint(out, "├───")
		}

		if isDir {
			fmt.Fprintf(out, "%s\n", entryInfo.Name())
		} else {
			var size string
			if entryInfo.Size() == 0 {
				size = "empty"
			} else {
				size = fmt.Sprintf("%db", entryInfo.Size())
			}
			fmt.Fprintf(out, "%s (%s)\n", entryInfo.Name(), size)
		}

		if entryInfo.IsDir() {
			var nestedPrefixes []string
			if last {
				nestedPrefixes = append(prefixes, "\t")
			} else {
				nestedPrefixes = append(prefixes, "│\t")
			}
			nestedPath := filepath.Join(path, entryInfo.Name())
			err = dirTreeWithOffset(out, nestedPath, printFiles, nestedPrefixes)
		}
	}

	return nil
}

func dirTree(out *bytes.Buffer, path string, printFiles bool) error {
	prefixes := make([]string, 0)
	return dirTreeWithOffset(out, path, printFiles, prefixes)
}

func main() {
	out := os.Stdout
	buf := &bytes.Buffer{}
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(buf, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(out, buf.String())
}
