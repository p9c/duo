package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fileHandle, _ := os.Open("./main.go")
	defer fileHandle.Close()
	fileScanner := bufio.NewScanner(fileHandle)
	i := 1
	for fileScanner.Scan() {
		fmt.Printf("%02d: %s\n", i, fileScanner.Text())
		i++
	}
}
