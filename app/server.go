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
		path := directory + filename

		// read file contents
		switch req.Method {
		case http.MethodGet:
			fileContent, err = readFile(path)
			if err != nil {
				fmt.Println("failed to read contents of file at directory: ", path)
				// Set status code to not found if error occurred while reading file
				statusCode = http.StatusNotFound
			}
		case http.MethodPost:
			err = writeFile(path, req.Body)
			if err != nil {
				fmt.Println("failed to write contents to file at directory: ", path)
			}
			statusCode = http.StatusCreated
		}
	}

	resp := utils.NewResponse(req, statusCode, fileContent)
	resp.WriteResponse(conn)
}

func readFile(path string) ([]byte, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		fmt.Println("file does not exist: ", path)
		return nil, err
	}

	file, err := os.Open(path)
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

func writeFile(path string, lines []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Println(line)
		fmt.Fprintln(writer, string(line))
	}

	return writer.Flush()
}

func validatePath(path string) int {
	validPathRegex := regexp.MustCompile(AllowedPaths)
	if validPathRegex.MatchString(path) {
		return http.StatusOK
	} else {
		return http.StatusNotFound
	}
}
