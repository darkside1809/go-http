package server

import (
	"net"
	"net/url"
	"sync"
	"strings"
	"strconv"
	"io"
	"log"
	"bytes"
	// "fmt"
	// "time"
	// "context"
	// "math/rand"
)

type HandleFunc func(req *Request)
type Server struct {
	addr		string
	mu			sync.RWMutex
	handlers map[string]HandleFunc
}
type Request struct {
	Conn        	net.Conn
	QueryParams 	url.Values
	PathParams  	map[string]string
	Headers     	map[string]string
	Body        	[]byte
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr, 
		handlers: make(map[string]HandleFunc)}
}
func (s *Server) Register(path string, handler HandleFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[path] = handler
}
func (s *Server) Response(body string) string {
	return "HTTP/1.1 200 OK\r\n" +
		"Content-Length: " + strconv.Itoa(len(body)) + "\r\n" +
		"Content-Type: text/html\r\n" +
		"Connection: close\r\n" +
		"\r\n" + body
} 
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Print(err)
		return err
	}

	defer func() {
		if cerr := listener.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(cerr)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue		
		}
		go s.handle(conn)
	}
}
func (s *Server) handle(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	buf := make([]byte, 8192)
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			log.Printf("%s", buf[:n])
		}
		if err != nil {
			log.Print(err)
			return
		}

		req := Request{}
		data := buf[:n]

		reqLineSep := []byte{'\r', '\n'}
		reqLineEnd := bytes.Index(data, reqLineSep)
		if reqLineEnd == -1 {
			return
		}
		lineEnders := []byte{'\r', '\n', '\r', '\n'}
		lineIndex := bytes.Index(data, lineEnders)
		if reqLineEnd == -1 {
			return
		}

		headersLine := string(data[reqLineEnd:lineIndex])
		headers := strings.Split(headersLine, "\r\n")[1:]

		headerTo := make(map[string]string)
		for _, v := range headers {
			headerLine := strings.Split(v, ": ")
			headerTo[headerLine[0]] = headerLine[1]
		}

		req.Headers = headerTo
		
		newData := string(data[lineIndex:])
		newData = strings.Trim(newData, "\r\n")

		req.Body = []byte(newData)
	
		requestLine := string(data[:reqLineEnd])
		parts := strings.Split(requestLine, " ")
		if len(parts) != 3 {
			return
		}

		path, version := parts[1], parts[2]

		if version != "HTTP/1.1" {
			return
		}

		decode, err := url.PathUnescape(path)
		if err != nil {
			log.Println(err)
			return
		}
		uri, err := url.ParseRequestURI(decode)
		if err != nil {
			log.Println(err)
			return
		}
		req.Conn = conn
		req.QueryParams = uri.Query()
		handler := func(req *Request) {
			conn.Close()
		}
		s.mu.RLock()
		pathParametres, hr := s.ratify(uri.Path)
		if hr != nil {
			handler = hr
			req.PathParams = pathParametres
		}
		s.mu.RUnlock()
		handler(&req)
	}
}
func (s *Server) ratify(path string) (map[string]string, HandleFunc) {
	routes:= make([]string, len(s.handlers))
	idx := 0
	for k := range s.handlers {
		routes[idx] = k
		idx++
	}

	header := make(map[string]string)

	for i := 0; i < len(routes); i++ {
		flag := false
		route := routes[i]
		partsRoute := strings.Split(route, "/")
		pRoutes := strings.Split(path, "/")

		for j, v := range partsRoute {
			if v != "" {
				f := v[0:1]
				l := v[len(v)-1:]
				if f == "{" && l == "}" {
					header[v[1:len(v)-1]] = pRoutes[j]
					flag = true
				} else if pRoutes[j] != v {

					strs := strings.Split(v, "{")
					if len(strs) > 0 {
						key := strs[1][:len(strs[1])-1]
						header[key] = pRoutes[j][len(strs[0]):]
						flag = true
					} else {
						flag = false
						break
					}
				}
				flag = true
			}
		}
		if flag {
			if hr, found := s.handlers[route]; found {
				return header, hr
			}
			break
		}
	}
	return nil, nil
}