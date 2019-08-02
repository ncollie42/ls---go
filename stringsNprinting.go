package main

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

func checkLink(path string, file os.FileInfo) string {
	if file.Mode()&os.ModeSymlink != 0 {
		str, err := os.Readlink(path + "/" + file.Name())
		if err != nil {
			fmt.Println(err)
		}
		return "-> " + str
	}
	return ""
}

func getFileStringLong(path string, file os.FileInfo) string { //include path in c?
	tmp := file.Sys().(*syscall.Stat_t)
	group, err := user.LookupGroupId(strconv.FormatUint(uint64(tmp.Gid), 10))
	if err != nil {
		fmt.Println(err)
	}
	user, err := user.LookupId(strconv.FormatUint(uint64(tmp.Uid), 10))
	if err != nil {
		fmt.Println(err)
	}
	link := checkLink(path, file)

	fileString := fmt.Sprintf("%s  %2d %s %s %7d %s %s %s",
		file.Mode(),
		tmp.Nlink,
		user.Username,
		group.Name,
		tmp.Size, //file.size()
		file.ModTime().Format("Jan _2 15:04"),
		file.Name(),
		link)
	return fileString

}
func getFileStringShort(path string, file os.FileInfo) string {
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
