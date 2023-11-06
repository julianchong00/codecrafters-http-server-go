package http

import (
	"fmt"
	"io"
	"strings"
)

const (
	StatusDesriptionOK        = "OK"
	StatusDescriptionNotFound = "Not Found"
)

func WriteResponse(w io.Writer, statusCode int, statusDescription string) {
	var out strings.Builder
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s \r\n", statusCode, statusDescription)
	out.WriteString(statusLine)
	out.WriteString("\r\n")
	w.Write([]byte(out.String()))
}
