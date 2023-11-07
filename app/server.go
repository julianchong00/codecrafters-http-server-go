package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"

	utils "github.com/codecrafters-io/http-server-starter-go/app/http"
)

const (
	AllowedPaths = `^/echo/.*$|^/$|^/user-agent$|^/files/.*$`
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	directoryPtr := flag.String("directory", "", "a directory")
	flag.Parse()

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn, *directoryPtr)
	}
}

func handleConnection(conn net.Conn, directory string) {
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

	var fileContent []byte
	if directory != "" {
		filename := strings.Replace(req.Path, "/files/", "", 1)

		// read file contents
		fileContent, err = readFile(directory, filename)
		if err != nil {
			fmt.Println("failed to read contents of file at directory: ", directory)
			// Set status code to not found if error occurred while reading file
			statusCode = http.StatusNotFound
		}
	}

	resp := utils.NewResponse(req, statusCode, fileContent)
	resp.WriteResponse(conn)
}

func readFile(directory string, filename string) ([]byte, error) {
	fullPath := directory + filename
	if _, err := os.Stat(fullPath); errors.Is(err, os.ErrNotExist) {
		fmt.Println("file does not exist: ", directory+filename)
		return nil, err
	}

	// cwd, err := os.Getwd()
	// if err != nil {
	// 	fmt.Println("failed to get current working directory")
	// 	return nil, err
	// }

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var buffer []byte
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Bytes()
		fmt.Println(scanner.Text())
		buffer = append(buffer, line...)
	}

	return buffer, nil
}

func validatePath(path string) int {
	validPathRegex := regexp.MustCompile(AllowedPaths)
	if validPathRegex.MatchString(path) {
		return http.StatusOK
	} else {
		return http.StatusNotFound
	}
}
