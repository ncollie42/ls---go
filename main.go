package main

import (
	"io/ioutil"
	"fmt"
	"flag"
	"os"
	// "path/filepath"
)
var	l_flag, r_flag, R_flag, a_flag bool
var g_path string

func init() {
	flag.BoolVar(&l_flag,"l", false, "long format")
	flag.BoolVar(&a_flag, "a", false, "long format")
	flag.StringVar(&g_path, "path", "/Users/Nico/Desktop/main", "long format")
	flag.Parse()
}

func test( path string, info os.FileInfo, err error) error {
	if a_flag || info.Name()[0] != '.' {
		if info.IsDir() {
			fmt.Println(path)
		}
		// fmt.Println(info.Name())
		// files = append(files, path)
	}
	return nil
}

func main() {
	// var files []string


	// filepath.Walk(g_path , test)

	fd , err := os.Open(".")
	if (err != nil) {
		os.Exit(0)
	}

	_, err= ioutil.ReadDir("main")
	if err != nil {
		fmt.Println(err)
	}

	// files. err := fd.Name/
	tmp , err := fd.Readdirnames(-1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fd.Name(), tmp)
	// for _, file := range files {
	// 	fmt.Println(file)
	// }
	// files, err := ioutil.ReadDir(".")
	// if err != nil {
	// 	fmt.Errorf("Error reading dir")
	// }
	// for _, file := range files {
	// 	fmt.Println(file.Mode(), file.Size(), file.Name(),"\n")
	// }
}