package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var Acs map[string]string = make(map[string]string)

var ConnectionPool map[string]net.Conn = make(map[string]net.Conn)

var Home map[string]string = make(map[string]string)

func handleConnection(conn net.Conn) {
	name := conn.RemoteAddr().String()

	fmt.Printf("%+v connected\n", name)
	conn.Write([]byte("Hello, " + name + "\n"))

	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		command := strings.Split(text, " ")
		if text == "Exit" {
			for login, connect := range ConnectionPool {
				if connect == conn {
					delete(ConnectionPool, login)
				}
			}
			fmt.Fprintln(conn, "bye")
			fmt.Println(name, "disconnected")
			break
		} else if command[0] == "0" {
			if command[1] == "0" {
				if _, isExist := Acs[command[2]]; !isExist {
					Home[command[2]] = "/home/olya/" + command[2] +"/"
					Acs[command[2]] = command[3]
					fmt.Fprintln(conn, "successful registration")
				} else {
					fmt.Fprintln(conn, "user with entered login exists")
				}
			} else {
				fmt.Fprintln(conn, "wrong address")
			}
		} else if command[0] == "1" {
			if command[1] == "0" {
				if _, isExist := Acs[command[2]]; isExist {
					if Acs[command[2]] == command[3] {
						ConnectionPool[command[2]] = conn
						fmt.Fprintln(conn, "successful login")
					} else {
						fmt.Fprintln(conn, "failed to log in")
					}
				} else {
					fmt.Fprintln(conn, "failed to log in")
				}
			} else {
				fmt.Fprintln(conn, "wrong address")
			}
		} else if command[0] == "2" {
			if command[1] == "2" {
				for _, connect := range ConnectionPool {
					if connect != conn {
						fmt.Fprintln(connect, command[2])
					}
				}
				fmt.Fprintln(conn, "Message sent")
			} else if command[1] == "1" {
				if connect, isExist := ConnectionPool[command[2]]; isExist {
					fmt.Fprintln(connect, command[3])
					fmt.Fprintln(conn, "Message sent")
				} else {
					fmt.Fprintln(conn, "user not found")
				}
			} else if command[1] == "0" {
				fmt.Println(command[2])
			} else {
				fmt.Fprintln(conn, "wrong address")
			}
		} else if command[0] == "3" {
			if command[1] == "1" {
				var log string
				for login, connect := range ConnectionPool {
					if connect == conn {
						log = login
					}
				}
				fileName := command[3]
				file, err := os.Open(Home[log] + fileName)
				if err != nil {
					fmt.Fprintln(conn, "error")
				} else {
					fmt.Println(log)
					if connect, isExist := ConnectionPool[command[2]]; isExist {
						var home string
						for login, connectt := range ConnectionPool {
							if connectt == connect {
								home = Home[login]
							}
						}
						fileName = home + fileName
						fmt.Println(fileName)
						fmt.Fprintln(connect, "file")
						fmt.Println("send filename")
						fmt.Fprintln(connect, fileName)
						info, _ := file.Stat()
						size := info.Size()
						b := make([]byte, size)
						file.Read(b)
						data := string(b) + "\r"
						fmt.Println("its data " + data)
						file.Close()
						fmt.Fprint(connect, data)
						fmt.Fprintln(conn, "file sent")
					} else {
						fmt.Fprintln(conn, "wrong address")
					}
				}
			}
		} else if text != "" {
			fmt.Fprintln(conn, "you enter ", text, "\n")
		}
	}
}

func main() {
	listner, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listner.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn)
	}
}
