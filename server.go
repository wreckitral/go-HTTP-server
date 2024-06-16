package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type HTTPRequest struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    string
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:8888")
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
		return
	}

	var res string

	switch {
	case req.Path == "/":
		res = getStatus(200, "OK") + "\r\n\r\n"
	case strings.HasPrefix(req.Path, "/echo/"):
		params := strings.TrimPrefix(req.Path, "/echo/")
		res = fmt.Sprintf("%s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", getStatus(200, "OK"), len(params), params)
	case req.Path == "/user-agent":
		res = fmt.Sprintf("%s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", getStatus(200, "OK"), len(req.Headers["User-Agent"]), req.Headers["User-Agent"])
	case strings.HasPrefix(req.Path, "/files/"):
		dir := os.Args[2]
		filename := strings.TrimPrefix(req.Path, "/files/")
		if req.Method == "GET" {
			data, err := os.ReadFile(dir + filename)
			if err != nil {
				res = getStatus(404, "Not Found") + "\r\n\r\n"
			} else {
				res = fmt.Sprintf("%s\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", getStatus(200, "OK"), len(data), data)
			}
		} else if req.Method == "POST" {
            if err := os.WriteFile(dir + filename, []byte(req.Body), 0644); err != nil {
                fmt.Errorf("failed writing a file")
            }
			res = getStatus(201, "Created") + "\r\n\r\n"
		}
	default:
		res = getStatus(404, "Not Found") + "\r\n\r\n"
	}

	// Write response to the connection
	_, err = conn.Write([]byte(res))
	if err != nil {
		fmt.Printf("error writing response: %s\n", err.Error())
	}
}

func scanHTTPRequest(scanner *bufio.Scanner) (*HTTPRequest, error) {
	var req HTTPRequest
	req.Headers = make(map[string]string)

	// Read request line
	if scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		if len(parts) < 3 { // Minimum parts should be 3 (GET / HTTP/1.1)
			return nil, fmt.Errorf("invalid http request")
		}

        fmt.Println(parts)

		req.Method = parts[0]
		req.Path = parts[1]
	}

	// Read headers
for scanner.Scan() {
        line := scanner.Text()
        if line == "" {
            // Empty line indicates end of headers
            break
        }
        headers := strings.SplitN(line, ": ", 2)
        if len(headers) < 2 {
            return nil, fmt.Errorf("invalid header line: %s", line)
        }
        req.Headers[headers[0]] = headers[1]
    }

    // Check if theres body
    if contentLengthStr, ok := req.Headers["Content-Length"]; ok {
        contentLength, err := strconv.Atoi(contentLengthStr)
        if err != nil {
            return nil, fmt.Errorf("invalid Content-Length header: %v", err)
        }

        // Read the body
        var bodyBuffer strings.Builder
        bytesRead := 0
        for scanner.Scan() {
            line := scanner.Text()
            if bytesRead < contentLength {
                if bytesRead+len(line) > contentLength {
                    line = line[:contentLength-bytesRead]
                }
                bodyBuffer.WriteString(line)
                bytesRead += len(line)

                if bytesRead >= contentLength {
                    break
                }
            }
        }

        if bytesRead != contentLength {
            return nil, fmt.Errorf("expected to read %d bytes, read %d", contentLength, bytesRead)
        }

        req.Body = bodyBuffer.String()
    }

    return &req, nil
}

func getStatus(statusCode int, statusText string) string {
	return fmt.Sprintf("HTTP/1.1 %d %s", statusCode, statusText)
}
