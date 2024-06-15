package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type HTTPRequest struct {
    Method string
    Path string
    Headers map[string]string
    Body string
}

func main() {
    l, err := net.Listen("tcp", "0.0.0.0:7777")
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    defer l.Close()

    for {
        conn, err := l.Accept()
        if err != nil {
            fmt.Println(err.Error())
            continue
        }

        go handleConn(conn)
    }
}

func handleConn(conn net.Conn) {
    defer conn.Close()

    scanner := bufio.NewScanner(conn)

    req, err := scanHTTPRequest(scanner)
    if err != nil {
        fmt.Printf("error on scanning http struct: %s\n", err.Error())
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("scanner failed:", err.Error())
        return
    }

    var res string

    if req.Path == "/" {
        res = getStatus(200, "OK") + "\r\n\r\n"
    } else if strings.Contains(req.Path, "/echo/") {
        params := strings.Split(req.Path, "/")[2]

        res = fmt.Sprintf("%s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", getStatus(200, "OK"), len(params), params)
    } else if req.Path == "/user-agent" {
        res = fmt.Sprintf("%s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", getStatus(200, "OK"), len(req.Headers["User-Agent"]), req.Headers["User-Agent"])
    } else if strings.Contains(req.Path, "/files/" ) {
        dir := os.Args[2]
        filename := strings.TrimPrefix(req.Path, "/files/")
        data, err := os.ReadFile(dir + filename)
        if err != nil {
            res = getStatus(404, "Not Found") + "\r\n\r\n"
        } else {
            res = fmt.Sprintf("%s\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", getStatus(200, "OK"), len(data), data)
        }
    } else {
        res = getStatus(404, "Not Found") + "\r\n\r\n"
    }

    conn.Write([]byte(res))
}

func scanHTTPRequest(scanner *bufio.Scanner) (*HTTPRequest, error) {
    var req HTTPRequest = HTTPRequest{}
    req.Headers = make(map[string]string)

    if scanner.Scan() {
        parts := strings.Split(scanner.Text(), " ")
        if len(parts) < 2 { // if theres no header
            return nil, fmt.Errorf("invalid http request")
        }

        req.Method = parts[0]
        req.Path = parts[1]
    }

    for scanner.Scan() {
        line := scanner.Text()
        if line == "" {
            break
        }
        headers := strings.SplitN(line, ": ", 2)
        if len(headers) < 2 {
            req.Body = headers[0]
        }

        req.Headers[headers[0]] = headers[1]
    }

    return &req, nil
}

func getStatus(statusCode int, statusText string) string {
   return fmt.Sprintf("HTTP/1.1 %d %s", statusCode, statusText)
}
