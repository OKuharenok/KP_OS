package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	go func(c net.Conn) {
		defer c.Close()
		sc := bufio.NewScanner(c)
		for sc.Scan() {
			text := sc.Text()
			if text == "file" {
					reader := bufio.NewReader(c)
					fileName, _ := reader.ReadString('\n')
					fileName = strings.TrimSpace(fileName)
					fmt.Println("its filename "+ fileName)
					file, _ := os.Create(fileName)
					data, _ := reader.ReadString('\r')
					data = data[:len(data)-1]
					fmt.Println("its data " + data)
					file.Write([]byte(data))
					file.Close()
					fmt.Println("file upload")

			}
			fmt.Println(text)
		}
	}(conn)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if scanner.Scan() {
			str := scanner.Text()
			fmt.Fprintf(conn, str + "\n")
		}
	}
}