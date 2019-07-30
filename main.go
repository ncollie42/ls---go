package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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

// func readDirNames(dirname string) ([]string, error) {
// 	f, err := os.Open(dirname)
// 	if err != nil {
// 		return nil, err
// 	}
// 	names, err := f.Readdirnames(-1)
// 	f.Close()
// 	if err != nil {
// 		return nil, err
// 	}
// 	sort.Strings(names) // sort how ever i need			---- also remove the '.' here?
// 	return names, nil
// }

func printDir(files []os.FileInfo) {
	// table to fall function?
	if l_flag {
		for _, file := range files {
			if a_flag || file.Name()[0] != '.' {
				fmt.Println(file.Mode(), "links", "name1", "name2", file.Size(), file.ModTime().Format("Jan _2 15:04"), file.Name())
			}
		}
	} else {
		for _, file := range files {
			if a_flag || file.Name()[0] != '.' {
				fmt.Println(file.Name())
			}
		}
	}
}

func walk(path string, info os.FileInfo) error {
	stats, err := readDir(path)
	if err != nil {
		return err
	}
	printDir(stats)
	if R_flag {
		for _, file := range stats {
			if a_flag || file.Name()[0] != '.' {
				fileName := filepath.Join(path, file.Name())
				fileInfo, err := os.Lstat(fileName)
				if err != nil {
					return err
				}
				if fileInfo.IsDir() {
					fmt.Println(fileName, ":")
					walk(fileName, fileInfo)
				}
			}
		}
	}
	return nil
}

func checkInput(root string) error { //check if it's a valid file/ dir
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

func test(name string) {
	info, err := os.Lstat(name)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info)
	fmt.Printf("\n%v", info)
	fmt.Printf("\n\n%#v", info)
	fmt.Printf("\n\n%#v", info.Sys())
	tmp := info.Sys().(*syscall.Stat_t)
	fmt.Printf("\n\n%v", tmp.Nlink)
	fmt.Printf("\n\n%v", tmp.Blocks)
	fmt.Printf("\n\n%v", tmp.Blksize)
	//get group name from group id?
	// s, ok := info.Sys().(*syscall.Stat_t)

}

func main() {
	// parse flags like in LS
	// pass in a list of inputs, print bads, regs, and then go into dirs

	test("tmp")

	// err := checkInput(g_path)
	// if err != nil {
	// 	println(err)
	// }
}
