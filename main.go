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

var l_flag, r_flag, R_flag, a_flag bool
var g_path string

func init() {
	flag.BoolVar(&l_flag, "l", false, "long format")
	flag.BoolVar(&a_flag, "a", false, "long format")
	flag.BoolVar(&R_flag, "R", false, "long format")
	flag.StringVar(&g_path, "path", ".", "long format")
	flag.Parse()
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
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() }) // sort how I'd like?
	return list, nil
}

func walk2(path string, info os.FileInfo, walkFn func(os.FileInfo)) error {
	if !info.IsDir() {
		walkFn(info)
		return nil
	}

	names, err := readDirNames(path)
	walkFn(info)
	// If err != nil, walk can't walk into this directory.
	// err1 != nil means walkFn want walk to skip this directory or stop walking.
	// Therefore, if one of err and err1 isn't nil, walk will return.
	if err != nil {
		// The caller's behavior is controlled by the return value, which is decided
		// by walkFn. walkFn may ignore err and return nil.
		// If walkFn returns SkipDir, it will be handled by the caller.
		// So walk should return whatever walkFn returns.
		return err
	}
	//	if -R dont go in if isDir
	//	limit names in readDirNames
	//	figure out return
	//	and how to sum up the totals for that file
	// print -- add fir que
	for _, name := range names {
		filename := filepath.Join(path, name)
		fileInfo, err := os.Lstat(filename)
		if err != nil {
			walkFn(info)
		} else {
			err = walk2(filename, fileInfo, walkFn)
			// print?
			if err != nil {
				if !fileInfo.IsDir() {
					return err
				}
			}
		}
	}
	return nil
}

func readDirNames(dirname string) ([]string, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Strings(names) // sort how ever i need			---- also remove the '.' here?
	return names, nil
}

func printFile(file os.FileInfo) {
	tmp := file.Sys().(*syscall.Stat_t)
	group, err := user.LookupGroupId(strconv.FormatUint(uint64(tmp.Gid), 10))
	if err != nil {
		fmt.Println(err)
	}
	user, err := user.LookupId(strconv.FormatUint(uint64(tmp.Uid), 10))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(file.Mode(),
		tmp.Nlink,
		user.Username,
		group.Name,
		tmp.Size,
		file.ModTime().Format("Jan _2 15:04"),
		file.Name())
}

func getFileString(file os.FileInfo) string {
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
		tmp.Size,
		file.ModTime().Format("Jan _2 15:04"),
		file.Name())
	return fileString

}

func printLong(lines []string, blocks int64) {
	fmt.Println("total:", blocks)
	for _, line := range lines {
		fmt.Println(line)
	}
}

/*
	Need a way to seperate reg print and long,
	Takes a list of fileInfos
	makes it into a string
	prints and keeps track of what are dirs.
*/

func handleDir(files []os.FileInfo) []string {
	// table to all function?
	lines := []string{}
	queue := []string{}
	var blocks int64

	if l_flag {
		for _, file := range files {
			if a_flag || file.Name()[0] != '.' {
				stat := file.Sys().(*syscall.Stat_t)
				tmp := getFileString(file)
				lines = append(lines, tmp) // returns a string, we append, and a a totoal block size? // then loop and print?
				blocks += stat.Blocks
				if file.IsDir() {
					queue = append(queue, file.Name())
				}
			}
		}
		printLong(lines, blocks)
	} else {
		for _, file := range files {
			if a_flag || file.Name()[0] != '.' {
				fmt.Println(file.Name())
			}
		}
	}
	return queue
}

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
		// return walk(root, info, printFile)
	} else {
		fmt.Println("it'not a dir - simple print", root)
	}
	return nil
}

func test(name string) {
	info, err := os.Lstat(name)
	// files, err := ioutil.ReadDir(name)
	if err != nil {
		fmt.Println(err)
	}
	// for _, info := range files {
	tmp := info.Sys().(*syscall.Stat_t)
	group, err := user.LookupGroupId(strconv.FormatUint(uint64(tmp.Gid), 10))
	if err != nil {
		fmt.Println(err)
	}
	user, err := user.LookupId(strconv.FormatUint(uint64(tmp.Uid), 10))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info.Mode(),
		tmp.Nlink,
		user.Username,
		group.Name,
		tmp.Size,
		info.ModTime().Format("Jan _2 15:04"),
		info.Name())

}

func main() {
	// parse flags like in LS
	// pass in a list of inputs, print bads, regs, and then go into dirs

	// test("/tmp/a")

	err := checkInput(g_path)
	if err != nil {
		println(err)
	}
}
