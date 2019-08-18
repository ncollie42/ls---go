package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"syscall"
)

const flagUsage = "usage: ls [-atlrR] [file ...]"

var (
	l_flag, r_flag, R_flag, a_flag, t_flag bool
	col, row                               uint
	g_path                                 []string
	getFileString                          func(spacing, os.FileInfo, string) string
	printDir                               func(spacing, []string, int64)
	compare                                func(file []os.FileInfo) compareFunc
)

type compareFunc func(i, j int) bool

type spacing struct {
	link          int
	user          int
	group         int
	size          int
	name          int
	numberOfFiles uint
}

/*
	TODO:
		* -ls -t by nano
		* argument -r on params
		* formating on regular ls
		* color
		* major - minor?
		* include the ./ on the long?

		//in c - loop, DIR and make a linked list of fileInfo - if ! -a skip '.' files then sort
		// in go is it worth removing here? it's looping again - it's a waste
*/

/*Go's readdir skips . and .., Using this function to add them back*/

func DotDotDot(dirname string) []os.FileInfo {
	tmp := []os.FileInfo{}
	fileInfo, err := os.Lstat(dirname + "/" + ".")
	if err != nil {
		fmt.Println(err)
	}
	tmp = append(tmp, fileInfo)
	fileInfo, err = os.Lstat(dirname + "/" + "..")
	if err != nil {
		fmt.Println(err)
	}
	tmp = append(tmp, fileInfo)
	return tmp
}

func readDir(dirname string) ([]os.FileInfo, error) {

	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	if a_flag { // adding . and .. because f.Readdir doesn't return them
		list = append(list, DotDotDot(dirname)...)
	}
	sort.Slice(list, compare(list)) // add stat struct and  path

	return list, nil
}

/*
	Takes a list of fileInfos from a DIR
	makes a [] of strings from each of the files inside
	prints and keeps a queue of Dirs in current dir if recursive
*/

func handleDir(files []os.FileInfo, path string) []string {
	lines := []string{}
	queue := []string{}
	var blocks int64

	spacing := getDirSpacing(files) //Probably not efficient

	for _, file := range files {
		if a_flag || file.Name()[0] != '.' {
			stat := file.Sys().(*syscall.Stat_t)
			tmp := getFileString(spacing, file, path)
			lines = append(lines, tmp)
			blocks += stat.Blocks
			if R_flag && file.IsDir() && file.Name() != "." && file.Name() != ".." { // Dont go into . or ..
				queue = append(queue, file.Name())
			}
		}
	}
	printDir(spacing, lines, blocks)
	return queue
}

func walk(path string) error {
	files, err := readDir(path)
	if err != nil {
		return err
	}

	queue := handleDir(files, path)

	if R_flag && len(queue) != 0 {
		for _, name := range queue {
			fileName := filepath.Join(path, name)
			_, err := os.Lstat(fileName)
			if err != nil {
				return err
			}
			fmt.Printf("\n%s:\n", fileName)
			walk(fileName)
		}
	}
	return nil
}

// removed argument from walk??
func main() {
	args := os.Args
	args = parseArgs(args[1:])
	setFunctions()
	sort.Slice(args, func(i, j int) bool { return args[i] < args[j] }) // r?

	trash, files, dirs := separateArgs(args)
	printTrash(trash)
	printFiles(files)
	err := handleDirs(dirs)
	if err != nil {
		println(err)
	}
}
