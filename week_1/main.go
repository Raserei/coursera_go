package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"
)

const element string = "├───"
const lastElement string = "└───"
const emptySize = "empty"

type TreeNode struct {
	File      os.FileInfo
	PrintName string
	Children  []TreeNode
}

func dirTree(out io.Writer, path string, isPrintFiles bool) (err error) {
	nodes, err := getTree(path, isPrintFiles)
	if err != nil {
		return
	}
	printNodes(out, nodes, "")
	return
}

func PrintName(fileInfo os.FileInfo) string {
	if fileInfo.IsDir() {
		return fileInfo.Name()
	} else {
		return fmt.Sprintf("%s (%s)", fileInfo.Name(), Size(fileInfo))
	}
}

func Size(fileInfo os.FileInfo) string {
	if fileInfo.Size() == 0 {
		return emptySize
	} else {
		return fmt.Sprintf("%db", fileInfo.Size())
	}
}

func getTree(root string, isPrintFiles bool) ([]TreeNode, error) {
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
	var nodes []TreeNode
	for _, file := range files {
		if !isPrintFiles && !file.IsDir() {
			continue
		}
		node := TreeNode{
			File:      file,
			PrintName: PrintName(file),
		}

		if file.IsDir() {
			nodePath := path.Join(root, node.File.Name())
			children, err := getTree(nodePath, isPrintFiles)
			if err != nil {
				return nil, err
			}
			node.Children = children
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func printNodes(out io.Writer, nodes []TreeNode, parentPrefix string) {
	var (
		lastIdx     = len(nodes) - 1
		childPrefix = "│\t"
		prefix      = element
	)

	for idx, node := range nodes {
		if idx == lastIdx {
			prefix = lastElement
			childPrefix = "\t"
		}
		fmt.Fprint(out, parentPrefix, prefix, node.PrintName, "\n")
		if node.File.IsDir() {
			printNodes(out, node.Children, parentPrefix+childPrefix)
		}
	}
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
