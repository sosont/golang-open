package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/scottkiss/gosshtool"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

var quit bool

const BUFFER_SIZE = 1024 * 4

var buffer = make([]byte, BUFFER_SIZE)

var (
	host   string
	passwd string
	user   string
)

func main() {
	flag.StringVar(&host, "h", "", "host")
	flag.StringVar(&passwd, "p", "", "password")
	flag.Parse()
	hostsp := strings.Split(host, "@")
	user = hostsp[0]
	host = hostsp[1]
	server := &server{
		Host: ":22",
	}
	go server.Start()
	time.Sleep(time.Second)
	conn, err := net.Dial("tcp", ":22")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	fmt.Println("connected to ssh server!")
	go onMessageRecived(conn)
	for {
		if quit {
			break
		}
		inputReader := bufio.NewReader(os.Stdin)
		input, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Println("There ware errors reading.")
			return
		}
		b := []byte(input)
		conn.Write(b)
	}
}

func onMessageRecived(conn net.Conn) {
	buffer = make([]byte, BUFFER_SIZE)
	for {
		n, err := conn.Read(buffer)
		if err == io.EOF {
			fmt.Printf("The RemoteAddr: %s is closed!\n", conn.RemoteAddr().String())
			return
		}
		if err != nil {
			break
		}
		if n > 0 {
			str := string(buffer[:n])
			fmt.Printf("%s", str)
			if strings.Contains(str, "logout") {
				quit = true
			}
		} else {
			break
		}
	}
}

type server struct {
	Host string
}

func handleConn(conn net.Conn) {
	config := &gosshtool.SSHClientConfig{
		User:     user,
		Password: passwd,
		Host:     host,
	}
	sshclient := gosshtool.NewSSHClient(config)
	_, err := sshclient.Connect()
	if err == nil {
		fmt.Println("ssh connect success")
	} else {
		fmt.Println("ssh connect failed")
	}
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	pty := &gosshtool.PtyInfo{
		Term:  "xterm-256color",
		H:     80,
		W:     40,
		Modes: modes,
	}
	session, err := sshclient.Pipe(conn, pty, nil, 30)
	if err != nil {
		fmt.Println(err)
	}
	defer session.Close()
}

func (s *server) Start() {
	ln, err := net.Listen("tcp", s.Host)
	if err != nil {
		fmt.Println(err)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go handleConn(conn)
	}

}