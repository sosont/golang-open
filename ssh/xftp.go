package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"io"
	"bufio"
	"encoding/csv"
	"strings"
	"container/list"
)

var (
	num int
)

func main() {
	if len(os.Args) == 1{
		fmt.Println("请输入文件名参数")
		return
	}
	list := listNode(os.Args[1])
	fmt.Println("请选择执行的语句")
	fmt.Scanln(&num)
	if num <= list.Len(){
		fmt.Println("您选择的是 ", num)
		ssh_to_do(list,num)
	}else {
		fmt.Println("您输入有误！ num:",num)
	}


}

///使用 container/list 包，实现匿名集合，少量开发便捷，待确定性能消耗
func ssh_to_do(list *list.List, num int) {
	if num != 0 {
		i := 1
		for node := list.Front(); node != nil; node = node.Next() {
			if i == num {
				switch value := node.Value.(type) {
				case BatchNode:
					SSH_do(value.User, value.Password, value.Ip_port, value.Cmd)
				}
			}
			i++
		}
	} else {
		for node := list.Front(); node != nil; node = node.Next() {

			switch value := node.Value.(type) {
			case BatchNode:
				SSH_do(value.User, value.Password, value.Ip_port, value.Cmd)
			}
		}
	}
}

//从文件读取命令，后期拓展
func listNode(fileName string) *list.List {
	list := readNode(fileName)
	fmt.Printf("共计 %d 条数据\n", list.Len())
	i := 1
	for node := list.Front(); node != nil; node = node.Next() {
		switch value := node.Value.(type) {
		case BatchNode:
			fmt.Println(i, "  ", value.String())
		}
		i++
	}
	return list
}

///SSH使用的封装
func SSH_do(user, password, ip_port string, cmd string) {
	PassWd := []ssh.AuthMethod{ssh.Password(password)}
	Conf := ssh.ClientConfig{User: user, Auth: PassWd}
	Client, _ := ssh.Dial("tcp", ip_port, &Conf)
	defer Client.Close()
	for {
		command := cmd
		if session, err := Client.NewSession(); err == nil {
			defer session.Close()
			session.Stdout = os.Stdout
			session.Stderr = os.Stderr
			session.Run(command)
			break
		}
	}
}

type BatchNode struct {
	User     string
	Password string
	Ip_port  string
	Cmd      string
}

func (batchNode *BatchNode) String() string {
	return "ssh " + batchNode.User + "@" + batchNode.Ip_port + "  with password: " + batchNode.Password + "  and run: " + batchNode.Cmd
}

func readNode(fileName string) *list.List {
	inputFile, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("在打开文件的时候出现错误\n文件存在吗?\n有权限吗?\n")
		return list.New()
	}
	defer inputFile.Close()

	batchNodeList := list.New()

	inputReader := bufio.NewReader(inputFile)
	for {
		inputString, err := inputReader.ReadString('\n')
		r := csv.NewReader(strings.NewReader(string(inputString)))
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("error !!! ", err)
				continue
			}
			batchNode := BatchNode{record[0], record[1], record[2], record[3]}
			batchNodeList.PushBack(batchNode)
		}
		if err == io.EOF {
			break
		}
	}
	return batchNodeList
}