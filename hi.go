package main

import (
	"fmt"
	"os"
)

func main() {
    args := os.Args  //好像还有一个包是直接获取，回头看 ...todo
    fmt.Println(args[0])
	 _s1 := args[1] //很有意思，0参数默认是执行文件全路径  - fmt.Println(args[0])
	switch _s1 {
	case "hi":
		/* 这是我的第一个简单的程序 */
		fmt.Println("Hello, World!")
	default:
		fmt.Printf("%s is comand", _s1)
	}
}
