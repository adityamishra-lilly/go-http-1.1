package main

import (
	"fmt"
	"log"
	"net"

	"github.com/adityamishra-lilly/go-http-1.1/internal/request"
)

func main() {
	ln, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("Error: ", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("Error: ", err)
		}
		fmt.Println("Connection Accepted")
		rq, err := request.RequestFromReader(conn)
		if err != nil {
			log.Panicf("%s", err)
		}
		fmt.Printf("Request line: \n- Method: %s\n- Target: %s\n- Version: %s\n", string(rq.RequestLine.Method), rq.RequestLine.RequestTarget, rq.RequestLine.HttpVersion)

	}

}

// func getLinesChannel(con net.Conn) <-chan string {
// 	strCh := make(chan string, 15)
// 	go func() {
// 		defer con.Close()
// 		defer close(strCh)
// 		currentLine := ""
// 		for {
// 			buf := make([]byte, 8)
// 			n, err := con.Read(buf)
// 			if err != nil {
// 				if err == io.EOF {
// 					break
// 				}
// 				log.Fatal("Error reading file", err)
// 				break
// 			}
// 			buf = buf[:n]
// 			if i := bytes.IndexByte(buf, '\n'); i != -1 {
// 				currentLine += string(buf[:i])
// 				buf = buf[i+1:]
// 				strCh <- currentLine
// 				currentLine = ""
// 			}
// 			currentLine += string(buf)
// 		}
// 		if len(currentLine) != 0 {
// 			strCh <- currentLine
// 		}
// 	}()
// 	return strCh
// }
