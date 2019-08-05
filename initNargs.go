package main
 import (
	 "os"
	 "fmt"
 )
func parseArgs(args []string) []string {
	for index, arg := range args {
		if arg[0] == '-' && len(arg) > 1{
			parseFlag(arg)
		} else {
			return args[index:]
		}
	}
	return []string{"."}
}


func parseFlag(arg string)  {
	for _, char := range arg[1:] {
		switch char {
		case 'a':
			a_flag = true
		case 't':
			t_flag = true
		case 'l':
			l_flag = true
		case 'r':
			r_flag = true
		case 'R':
			R_flag = true
		default:
			fmt.Printf("ls: illegal option -- %c\n", char)
			fmt.Printf(flagUsage)
			os.Exit(-1)
		}
	}
}

func setFunctions() {
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
}

func separateArgs(args []string) (trash[]string, files[]string, dirs[]string) {
	for _, curent := range args {
		info, err := os.Lstat(curent)	
		if err != nil {
			trash = append(trash, curent)
		} else if info.IsDir() {
			dirs = append(dirs, curent)
		} else {
			files = append(files, curent)
		}
	}
	return
}

func handleDirs(args []string) error {
	for _, curent := range args {
		info, err := os.Lstat(curent)
		if err != nil {
			return err
		}
		return walk(curent, info)
	}
	return nil
}

func printTrash(trash []string) {
	for _, curent := range trash {
		fmt.Println("ls:", curent, "No such file or directory")
	}
}

func printFiles(trash []string) {
	for _, curent := range trash {
		fmt.Println(curent)
	}
}
