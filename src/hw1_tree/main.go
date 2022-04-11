package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type node struct {
	value    os.FileInfo
	children []node
}

func getNodes(path string, withFiles bool) (tree []node, err error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var subTree []node
	for _, file := range files {
		if !withFiles && !file.IsDir() {
			continue
		}
		n := node{
			value: file,
		}

		if file.IsDir() {
			subNodes, err := getNodes(path+string(os.PathSeparator)+file.Name(), withFiles)
			if err != nil {
				return nil, err
			}

			n.children = subNodes
		}
		subTree = append(subTree, n)
	}
	return subTree, nil
}

func (t node) Name() (name string) {
	if t.value.IsDir() {
		return t.value.Name()
	} else {
		return fmt.Sprintf("%s (%s)", t.value.Name(), t.Size())
	}
}

func (t node) Size() (size string) {
	if t.value.Size() == 0 {
		return "empty"
	} else {
		return fmt.Sprintf("%db", t.value.Size())
	}
}

func printTree(out io.Writer, tree []node, parentPrefix string) (err error) {

	var (
		lastIdx     = len(tree) - 1
		prefix      = "├───"
		childPrefix = "│\t"
	)

	for idx, t := range tree {

		if idx == lastIdx {
			prefix = "└───"
			childPrefix = "\t"
		}

		_, err = fmt.Fprint(out, parentPrefix, prefix, t.Name(), "\n")
		if err != nil {
			return nil
		}

		if t.value.IsDir() {
			err := printTree(out, t.children, parentPrefix+childPrefix)
			if err != nil {
				return nil
			}
		}

	}
	return
}

func dirTree(out io.Writer, path string, withFiles bool) (err error) {
	fmt.Fprintln(out, path, withFiles)
	tree, err := getNodes(path, withFiles)
	if err != nil {
		return nil
	}
	printTree(out, tree, "")
	return
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
