package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
    l, err := net.Listen("tcp", "0.0.0.0:7777")
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    conn, err := l.Accept()
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
}



