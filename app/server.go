package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"

	utils "github.com/codecrafters-io/http-server-starter-go/app/http"
)

const (
	AllowedPaths = `^/user-agent/.*$|^/$`
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	handleConnection(conn)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	input := make([]byte, 1024)
	_, err := conn.Read(input)
	if err != nil {
		fmt.Println("error reading connection: ", err.Error())
		os.Exit(1)
	}

	req, err := utils.ParseRequest(input)
	if err != nil {
		fmt.Println("error parsing request: ", err.Error())
	}

	statusCode := validatePath(req.Path)

	resp := utils.NewResponse(req, statusCode)
	resp.WriteResponse(conn)
}

func validatePath(path string) int {
	validPathRegex := regexp.MustCompile(AllowedPaths)
	if validPathRegex.MatchString(path) {
		return http.StatusOK
	} else {
		return http.StatusNotFound
	}
}
