package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Method string
type ParserState string

const bufferSize = 8
const (
	GET     Method = "GET"
	POST    Method = "POST"
	PUT     Method = "PUT"
	PATCH   Method = "PATCH"
	DELETE  Method = "DELETE"
	HEAD    Method = "HEAD"
	CONNECT Method = "CONNECT"
	TRACE   Method = "TRACE"
	OPTIONS Method = "OPTIONS"
)

const (
	INITIALIZED ParserState = "initialized"
	DONE        ParserState = "done"
)

const HTTP = "HTTP"
const HTTPVERSION = "1.1"

var validMethod = map[Method]struct{}{
	GET:     {},
	POST:    {},
	PUT:     {},
	PATCH:   {},
	DELETE:  {},
	HEAD:    {},
	CONNECT: {},
	TRACE:   {},
	OPTIONS: {},
}

func (method Method) IsValidMethod() bool {
	_, ok := validMethod[method]
	return ok
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        Method
}

type Request struct {
	RequestLine RequestLine
	ParserState ParserState
}

func NewRequest(rqLine RequestLine) *Request {
	rq := &Request{
		RequestLine: rqLine,
		ParserState: INITIALIZED,
	}
	return rq
}

func (r *Request) parse(data []byte) (int, error) {
	if r.ParserState != INITIALIZED {
		return 0, fmt.Errorf("Parser state not initialized!")
	}
	rqLine, numBytes, err := parseRequestLine(data)
	if err != nil {
		return 0, fmt.Errorf("T %w", err)
	}
	if numBytes == 0 {
		return 0, nil
	}
	r.ParserState = DONE
	r.RequestLine = rqLine
	return numBytes, nil
}

func NewRequestLine(HttpVersion string, RequestTarget string, Method Method) *RequestLine {
	rqLine := &RequestLine{
		HttpVersion:   HttpVersion,
		RequestTarget: RequestTarget,
		Method:        Method,
	}

	return rqLine
}

func parseRequestLine(rqBytes []byte) (RequestLine, int, error) {
	numBytes := len(rqBytes)
	if !bytes.Contains(rqBytes, []byte("\r\n")) {
		return RequestLine{}, 0, nil
	}

	rqParts := strings.Split(string(rqBytes), "\r\n")
	rqLine := rqParts[0]
	rqLineParts := strings.Split(rqLine, " ")
	if len(rqLineParts) != 3 {
		return RequestLine{}, numBytes, fmt.Errorf("Please provide a valid request line!")
	}

	rqMethod, rqTarget, rqHttpVersion := rqLineParts[0], rqLineParts[1], rqLineParts[2]
	method := Method(rqMethod)
	if !method.IsValidMethod() {
		return RequestLine{}, numBytes, fmt.Errorf("Invalid request method: %s\n", method)
	}
	if !strings.Contains(rqTarget, "/") {
		return RequestLine{}, numBytes, fmt.Errorf("Invalid request target: %s\n", rqTarget)
	}

	rqHttpVersionParts := strings.Split(rqHttpVersion, "/")

	if len(rqHttpVersionParts) != 2 {
		return RequestLine{}, numBytes, fmt.Errorf("Invalid http version: ")
	}

	httpName, httpVersion := rqHttpVersionParts[0], rqHttpVersionParts[1]
	if httpName != HTTP || httpVersion != HTTPVERSION {
		return RequestLine{}, numBytes, fmt.Errorf("Invalid http-version-name Output: %s/%s \n: Expected: %s/%s \n", httpName, httpVersion, HTTP, HTTPVERSION)
	}

	requestLine := NewRequestLine(httpVersion, rqTarget, Method(rqMethod))
	return *requestLine, numBytes, nil

}

func RequestFromReader(reader io.Reader) (*Request, error) {

	buffer := make([]byte, bufferSize)
	request := NewRequest(RequestLine{})
	request.ParserState = INITIALIZED
	indexStart := 0
	for {
		if request.ParserState == DONE {
			break
		}

		numBytes, err := reader.Read(buffer[indexStart:])
		fmt.Printf("NumBytes: %d", numBytes)
		if err != nil {
			if err != io.EOF {
				return nil, fmt.Errorf("Error reading")
			}

		}
		indexStart += numBytes
		if indexStart == len(buffer) {
			newBuffer := make([]byte, len(buffer)*2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}
		_, err = request.parse(buffer[:indexStart])
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

	}

	return request, nil

}
