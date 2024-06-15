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

    _, err = l.Accept()
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
}



