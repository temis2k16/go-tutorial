package main

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	// "strings"
)

func dirTreeWalker(inPath string, printFiles bool,
	// [родитель] наследники (файлы и папки)
	resultMap map[string][]string) error {
	var dirs []string
	var fileAndSize []string
	var files []fs.FileInfo
	var e error
	files, e = ioutil.ReadDir(inPath)
	if e != nil {
		return e
	}

	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file.Name())
		} else if printFiles {
			if file.Size() == 0 {
				filename := fmt.Sprint(file.Name(), " (empty)")
				fileAndSize = append(fileAndSize, filename)
			} else {
				filename := fmt.Sprint(file.Name(), " (", file.Size(), "b)")
				fileAndSize = append(fileAndSize, filename)
			}
		}
	}
	inPathChildren := append(dirs, fileAndSize...)
	resultMap[inPath] = inPathChildren

	for _, p := range dirs {
		inDir := inPath + "/" + p
		e = dirTreeWalker(inDir, printFiles, resultMap)
		if e != nil {
			return e
		}
	}
	return nil
}

func dirTree(out io.Writer, inPath string, printFiles bool) error {

	/*var p []string
	var files []fs.FileInfo
	var e error
	files, e = ioutil.ReadDir(inPath)

	for _, file := range files {
		if file.IsDir() {
			p = append(p, file.Name())
		} else {
			if file.Size() == 0 {
				filename := fmt.Sprint(file.Name(), " (empty)")
				p = append(p, filename)
			} else {
				filename := fmt.Sprint(file.Name(), " (", file.Size(), "b)")
				p = append(p, filename)
			}
		}
	}
	outputString := strings.Join(p, "\n")
	outputString += "\n"
	fmt.Fprintf(out, outputString)*/
	var outputString string
	resultMap := make(map[string][]string)

	e := dirTreeWalker(inPath, printFiles, resultMap)
	for key := range resultMap {
		// p := strings.Join(key, "\n")
		outputString += key + "\n"
	}
	fmt.Fprintf(out, outputString)
	return e
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
