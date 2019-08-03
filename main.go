package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"syscall"
)

var (
	l_flag, r_flag, R_flag, a_flag, t_flag bool
	col, row                               uint
	g_path                                 string
	getFileString                          func(spacing, os.FileInfo, string) string
	printDir                               func(spacing, []string, int64)
	compare                                func(file []os.FileInfo) compareFunc
)

type compareFunc func(i, j int) bool

type spacing struct {
	link  int
	user  int
	group int
	size  int
	name  int // only for reg?
}

func init() {
	flag.BoolVar(&l_flag, "l", false, "long format")
	flag.BoolVar(&a_flag, "a", false, "Hidden files")
	flag.BoolVar(&t_flag, "t", false, "sort time")
	flag.BoolVar(&r_flag, "r", false, "sort reverse")
	flag.BoolVar(&R_flag, "R", false, "Recursive")
	flag.StringVar(&g_path, "path", ".", "path to be used")
	flag.Parse()
	if l_flag {
		getFileString = getFileStringLong
		printDir = printLong
	} else {
		getFileString = getFileStringShort
		printDir = printShort
	}
	if r_flag {
		if t_flag {
			compare = byTimeReverse
		} else {
			compare = byNameReverse
		}
	} else if t_flag {
		compare = byTime
	} else {
		compare = byName
	}
	col, row = getTerminalSize()
	fmt.Println("col:", col, "row", row)
}
ß
/*
	TODO:
		* parse flags like ls
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
func intCheck(new, old int) int {
	tmp := countDigits(new)
	if tmp > old {
		return tmp
	}
	return old
}
func countDigits(i int) (count int) {
	for i != 0 {

		i /= 10
		count = count + 1
	}
	return count
}
func stringCheck(new string, old int) int {
	tmp := len(new)
	if tmp > old {
		return tmp
	}
	return old

}

//just used to get spacing
func getDirSpacing(files []os.FileInfo) spacing {
	var spacing spacing

	for _, file := range files {
		stat := file.Sys().(*syscall.Stat_t) //might not work
		group, user := getUserNGroup(stat)
		spacing.link = intCheck(int(stat.Nlink), spacing.link)
		spacing.user = stringCheck(user, spacing.user)
		spacing.group = stringCheck(group, spacing.group)
		spacing.size = intCheck(int(stat.Size), spacing.size)
		spacing.name = stringCheck(file.Name(), spacing.name)
	}
	return spacing
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

/*
	Reads a dir, prints, if recursive - calls function again on queue
*/
func walk(path string, info os.FileInfo) error {
	files, err := readDir(path)
	if err != nil {
		return err
	}

	queue := handleDir(files, path)

	if R_flag && len(queue) != 0 {
		for _, name := range queue {
			fileName := filepath.Join(path, name)
			fileInfo, err := os.Lstat(fileName)
			if err != nil {
				return err
			}
			fmt.Printf("\n%s:\n", fileName)
			walk(fileName, fileInfo)
		}
	}
	return nil
}

/*Will take a list of arguments to either print as dir, file or trash*/

func checkInput(root string) error {
	info, err := os.Lstat(root)
	if err != nil {
		fmt.Println("Bad path? or invalid name?")
	} else if info.IsDir() {
		return walk(root, info)
	} else {
		fmt.Println("it'not a dir - simple print", root)
	}
	return nil
}

func main() {
	// parse flags like in LS
	// pass in a list of inputs, print bads, regs, and then go into dirs

	err := checkInput(g_path)
	if err != nil {
		println(err)
	}
}
