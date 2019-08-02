package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"syscall"
)

var (
	l_flag, r_flag, R_flag, a_flag, t_flag bool
	g_path                                 string
	getFileString                          func(os.FileInfo) string
	printDir                               func(lines []string, blocks int64)
	compare                                func(file []os.FileInfo) compareFunc
)

type compareFunc func(i, j int) bool

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
}

/*
	TODO:
		*simlinks
		* include the ./ on the long?
		* parse flags like ls
		//in c - loop, DIR and make a linked list of fileInfo - if ! -a skip '.' files then sort
		// in go is it worth removing here? it's looping again - it's a waste
*/

/*		Sorting functions		*/
func byTimeReverse(file []os.FileInfo) compareFunc {
	return func(i, j int) bool {
		return !file[i].ModTime().After(file[j].ModTime())
	}
}
func byTime(file []os.FileInfo) compareFunc {
	return func(i, j int) bool {
		return file[i].ModTime().After(file[j].ModTime())
	}
}
func byName(file []os.FileInfo) compareFunc {
	return func(i, j int) bool {
		return file[i].Name() < file[j].Name()
	}
}
func byNameReverse(file []os.FileInfo) compareFunc {
	return func(i, j int) bool {
		return file[i].Name() > file[j].Name()
	}
}

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
	sort.Slice(list, compare(list))
	return list, nil
}

func getFileStringLong(file os.FileInfo) string {
	tmp := file.Sys().(*syscall.Stat_t)
	group, err := user.LookupGroupId(strconv.FormatUint(uint64(tmp.Gid), 10))
	if err != nil {
		fmt.Println(err)
	}
	user, err := user.LookupId(strconv.FormatUint(uint64(tmp.Uid), 10))
	if err != nil {
		fmt.Println(err)
	}
	fileString := fmt.Sprintf("%s %d %s %s %d %s %s",
		file.Mode(),
		tmp.Nlink,
		user.Username,
		group.Name,
		tmp.Size, //file.size()
		file.ModTime().Format("Jan _2 15:04"),
		file.Name())
	return fileString

}
func getFileStringShort(file os.FileInfo) string {
	return fmt.Sprintf("%s", file.Name())
}
func printLong(lines []string, blocks int64) {
	fmt.Println("total", blocks)
	for _, line := range lines {
		fmt.Println(line)
	}
}
func printShort(lines []string, blocks int64) {
	//handle the spacing here with terminal size ??
	for _, line := range lines {
		fmt.Println(line)
	}
}

/*
	Takes a list of fileInfos from a DIR
	makes a [] of strings from each of the files inside
	prints and keeps a queue of Dirs in current dir if recursive
*/

func handleDir(files []os.FileInfo) []string {
	lines := []string{}
	queue := []string{}
	var blocks int64

	for _, file := range files {
		if a_flag || file.Name()[0] != '.' {
			stat := file.Sys().(*syscall.Stat_t)
			tmp := getFileString(file)
			lines = append(lines, tmp)
			blocks += stat.Blocks
			if R_flag && file.IsDir() && file.Name() != "." && file.Name() != ".." { // Dont go into . or ..
				queue = append(queue, file.Name())
			}
		}
	}
	printDir(lines, blocks)
	return queue
}

/*
	Reads a dir, prints, if recursive - calls function again on queue
*/
func walk(path string, info os.FileInfo) error {
	stats, err := readDir(path)
	if err != nil {
		return err
	}
	queue := handleDir(stats)

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
