package http

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

var StatusLineMap = map[int]string{
	http.StatusOK:       "HTTP/1.1 200 OK",
	http.StatusNotFound: "HTTP/1.1 404 Not Found",
}

const (
	CRLF = "\r\n"
)

type Response struct {
	StatusCode int
	StatusLine string
	Headers    map[string]string
	Body       []byte
}

func NewResponse(req Request, statusCode int) Response {
	var resp Response
	resp.StatusCode = statusCode
	resp.StatusLine = StatusLineMap[statusCode]
	resp.Headers = make(map[string]string)
	resp.Headers["Content-Type"] = "text/plain"

	var body string
	if strings.Contains(req.Path, "/echo/") {
		body = strings.Replace(req.Path, "/echo/", "", 1)
	} else if strings.Contains(req.Path, "/user-agent") {
		body = req.Headers["User-Agent"]
	}

	resp.Headers["Content-Length"] = fmt.Sprintf("%d", len(body))
	resp.Body = []byte(body)

	return resp
}

func (r Response) WriteResponse(w io.Writer) {
	var out strings.Builder
	out.WriteString(r.StatusLine + CRLF)
	for header, value := range r.Headers {
		out.WriteString(header + ": " + value + CRLF)
	}
	out.WriteString(CRLF)
	out.Write(r.Body)

	_, err := w.Write([]byte(out.String()))
	if err != nil {
		fmt.Println("failed to write to socket", err.Error())
		return
	}
}
