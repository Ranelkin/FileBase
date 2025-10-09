package main

import (
	"bufio"
	"filebase/compare"
	"filebase/traverse"
	"filebase/util"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Hello welcome to FileBase!")
	fmt.Println("To traverse a directory please enter the command 'traverse' followed by the directorys path")
	fmt.Println("To get the difference of one directory to another please enter 'diff <path1> <path2>'")
	fmt.Println("Where path1 and path2 are both the compressed .txt files created by the traverse cli tool from FileBase")
	fmt.Print("FileBase: ")

	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	in := strings.TrimSpace(input.Text())
	arr := strings.Split(in, " ")
	cmd := strings.ToLower(arr[0])
	switch cmd {
	case "traverse":
		//Check for correct input
		if !(len(arr) != 2) {
			fmt.Println("Wrong ammount of path entries. Make sure to enter 1 path")
		}
		traverse.Traverse(arr[1])

	case "diff":
		if !(len(arr) == 3) {
			fmt.Println("Wrong ammount of path entries. Make sure to enter 2 paths")
		}

		diff, err := compare.Difference(arr[1], arr[2])
		if err != nil {
			panic(err)
		}
		path := arr[2]
		err = util.WriteToFile(*diff, path)
		if err != nil {
			panic(err)
		}
	default:
		fmt.Println("Mate learn how to read")
	}
	fmt.Println("Exiting...")

}
