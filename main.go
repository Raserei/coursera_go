package main

import (
	"fmt"
	"path"

	"io"
	"os"
	"strings"
)

const element string = "├───"
const lastElement string = "└───"

func dirTree(out *os.File, inPath string, printfiles bool) error {
	return dirTreeHelp(out, inPath, printfiles, 0)
}

func dirTreeHelp(out *io.Reader, inPath string, printfiles bool, depth int) error {
	pathRoot, err := os.Open(inPath)
	defer pathRoot.Close()
	if err != nil {
		return err
	}
	files, err := pathRoot.ReadDir(-1)
	for key, file := range files {
		if !file.IsDir() && !printfiles {
			continue
		}
		fmt.Print(getFileLine(file.Name(), depth, key == len(files)-1))
		if file.IsDir() {
			dirTreeHelp(out, path.Join(inPath, file.Name()), printfiles, depth+1)
		}
	}
	return nil
}

func getFileLine(name string, depth int, isLast bool) string {
	builder := new(strings.Builder)
	if depth != 0 {
		builder.WriteString("│")
	}
	for i := 1; i <= depth; i++ {
		builder.WriteString("\t")
	}
	if isLast {
		builder.WriteString(lastElement)
	} else {
		builder.WriteString(element)
	}
	builder.WriteString(name)
	builder.WriteString("\n")
	return builder.String()
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	println(path)
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
