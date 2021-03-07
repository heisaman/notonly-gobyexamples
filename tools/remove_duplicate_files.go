package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func input(x []string, err error) []string {
	if err != nil {
		return x
	}
	var d string
	n, err := fmt.Scanf("%s", &d)
	if n == 1 {
		x = append(x, d)
	}
	return input(x, err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	fmt.Println("请输入需要检查的文件目录名（以空格区分）: ")
	folders := input([]string{}, nil)
	fmt.Println()
	folderFiles := make([][]string, len(folders))
	folderFilesModTime := make([]map[string]time.Time, len(folders))
	for i := range folders {
		folderFiles[i] = []string{}
		folderFilesModTime[i] = map[string]time.Time{}
		c, err := ioutil.ReadDir(folders[i])
		check(err)
		for _, entry := range c {
			folderFiles[i] = append(folderFiles[i], entry.Name())
			folderFilesModTime[i][entry.Name()] = entry.ModTime()
		}
	}
	//fmt.Println(folderFiles)
	//fmt.Println(folderFilesModTime)

	filesDeleted := []string{}
	for i := range folderFiles {
		for j := range folderFiles[i] {
			fileName := folderFiles[i][j]
			for k := range folderFilesModTime {
				if k != i {
					if modTime, ok := folderFilesModTime[k][fileName]; ok {
						filePath := fmt.Sprint(folders[k], "/", fileName)
						if _, err := os.Stat(filePath); err == nil {
							if modTime.Before(folderFilesModTime[i][fileName]) {
								filesDeleted = append(filesDeleted, filePath)
								err := os.Remove(filePath)
								check(err)
							}
						}
					}
				}
			}
		}
	}
	fmt.Println("已成功清理以下过时的文件：", filesDeleted)
	fmt.Println()

	fmt.Print("按任意键退出...")
	fmt.Scanln()
}
