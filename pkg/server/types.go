package server

import (
	"net"
	"net/url"
	"sync"
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