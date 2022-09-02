package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var dir string

func init() {
	dir, _ = os.Getwd()
}

func exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func ls() {
	files, _ := ioutil.ReadDir(dir)
	filenames := make([]string, len(files))
	for _, file := range files {
		filenames = append(filenames, file.Name())
	}
	fmt.Println(filenames)
}

func cd(cmd string) {
	arr := strings.Split(cmd, " ")
	if len(arr) != 2 {
		fmt.Println("invalid cmd:", cmd)
	}

	path := arr[1]
	if strings.HasPrefix(path, "/") {
		if !exists(path) {
			fmt.Println("path not exist, path:", path)
		} else {
			dir = path
		}
	} else {
		if !exists(dir + "/" + path) {
			fmt.Println("path not exist, path:", dir+"/"+path)
		} else {
			dir = dir + "/" + path
		}
	}
}

func process(cmd string) {
	if cmd == "pwd" {
		fmt.Println(dir)
	} else if cmd == "ls" {
		ls()
	} else if strings.HasPrefix(cmd, "cd") {
		cd(cmd)
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cmd := scanner.Text()
		process(cmd)
	}
}
