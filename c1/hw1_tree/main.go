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
	dirs := make([]os.FileInfo, 0)
	files := make([]os.FileInfo, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry)
		} else {
			files = append(files, entry)
		}
	}

	if printFiles && len(prefixes) > 0 {
		for i, fileInfo := range files {
			for _, prefix := range prefixes {
				fmt.Fprint(out, prefix)
			}
			last := len(dirs) == 0 && len(files) == i+1
			if last {
				fmt.Fprint(out, "└───")
			} else {
				fmt.Fprint(out, "├───")
			}
			var size string
			if fileInfo.Size() == 0 {
				size = "empty"
			} else {
				size = fmt.Sprintf("%db", fileInfo.Size())
			}
			fmt.Fprintf(out, "%s (%s)\n", fileInfo.Name(), size)
		}
	}

	for i, dirInfo := range dirs {
		for _, prefix := range prefixes {
			fmt.Fprint(out, prefix)
		}
		var last bool
		if len(prefixes) > 0 {
			last = len(dirs) == i+1
		} else {
			if printFiles {
				last = len(files) == 0 && len(dirs) == i+1
			} else {
				last = len(dirs) == i+1
			}
		}
		if last {
			fmt.Fprint(out, "└───")
		} else {
			fmt.Fprint(out, "├───")
		}
		fmt.Fprintf(out, "%s\n", dirInfo.Name())

		var nestedPrefixes []string
		if last {
			nestedPrefixes = append(prefixes, "\t")
		} else {
			nestedPrefixes = append(prefixes, "│\t")
		}
		nestedPath := filepath.Join(path, dirInfo.Name())
		err = dirTreeWithOffset(out, nestedPath, printFiles, nestedPrefixes)
	}

	if printFiles && len(prefixes) == 0 {
		for i, fileInfo := range files {
			for _, prefix := range prefixes {
				fmt.Fprint(out, prefix)
			}
			if len(files) == i+1 {
				fmt.Fprint(out, "└───")
			} else {
				fmt.Fprint(out, "├───")
			}
			var size string
			if fileInfo.Size() == 0 {
				size = "empty"
			} else {
				size = fmt.Sprintf("%db", fileInfo.Size())
			}
			fmt.Fprintf(out, "%s (%s)\n", fileInfo.Name(), size)
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
