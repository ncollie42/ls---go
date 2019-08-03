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

func getUserNGroup(stat *syscall.Stat_t) (string, string) {
	group, err := user.LookupGroupId(strconv.FormatUint(uint64(stat.Gid), 10))
	if err != nil {
		fmt.Println(err)
	}
	user, err := user.LookupId(strconv.FormatUint(uint64(stat.Uid), 10))
	if err != nil {
		fmt.Println(err)
	}
	return user.Username, group.Name
}
func getFileStringLong(spacing spacing, file os.FileInfo, path string) string { //include path in c?
	stat := file.Sys().(*syscall.Stat_t)
	group, user := getUserNGroup(stat)
	link := checkLink(path, file)

	fileString := fmt.Sprintf("%s  %*d %*s %*s %*d %s %s %s",
		file.Mode(),
		spacing.link,
		stat.Nlink,
		spacing.user,
		user,
		spacing.group,
		group,
		spacing.size+1,
		stat.Size,
		file.ModTime().Format("Jan _2 15:04"),
		file.Name(),
		link)
	return fileString

}
func getFileStringShort(spacing spacing, file os.FileInfo, path string) string {
	return fmt.Sprintf("%s", file.Name())
}
func printLong(spacing spacing, lines []string, blocks int64) {
	fmt.Println("total", blocks)
	for _, line := range lines {
		fmt.Println(line)
	}
}
func printShort(spacing spacing, lines []string, blocks int64) {

	fmt.Println(spacing.name)

	rows := col / uint(spacing.name)
	counter := rows
	for _, line := range lines {
		fmt.Printf("%-*s", spacing.name, line)
		counter--
		if counter == 0 {
			counter = rows
		}
	}
}
