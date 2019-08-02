package main

import (
	"os"
)

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
