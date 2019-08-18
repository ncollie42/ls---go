package main

import (
	"os"
	"syscall"
)

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
		spacing.numberOfFiles++
	}
	spacing.name++
	return spacing
}

func stringCheck(new string, old int) int {
	tmp := len(new)
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
func intCheck(new, old int) int {
	tmp := countDigits(new)
	if tmp > old {
		return tmp
	}
	return old
}
