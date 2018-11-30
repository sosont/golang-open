package main

import (
	"fmt"
	"net"
)

func main() {
	// args := os.Args //好像还有一个包是直接获取，回头看 ...todo
	// fmt.Println(args[0])
	// _s1 := args[1] //很有意思，0参数默认是执行文件全路径  - fmt.Println(args[0])
	// switch _s1 {
	// case "hi":
	// 	/* 这是我的第一个简单的程序 */
	// 	fmt.Println("Hello, World!")
	// default:
	// 	fmt.Printf("%s is comand", _s1)
	// }
	ipstr, err := LocalIP()
	if err == nil {
		fmt.Printf("ip: %s", ipstr)

	}

}

//获取本机IP的 NET包
func LocalIP() (net.IP, error) {
	tables, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, t := range tables {
		addrs, err := t.Addrs()
		if err != nil {
			return nil, err
		}
		for _, a := range addrs {
			ipnet, ok := a.(*net.IPNet)
			if !ok || ipnet.IP.IsLoopback() {
				continue
			}
			if v4 := ipnet.IP.To4(); v4 != nil {
				return v4, nil

			}
		}
	}
	return nil, fmt.Errorf("cannot find local IP address")
}
