package main

import (
	//"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/darkside1809/http/cmd/app"
	"github.com/darkside1809/http/pkg/banners"
)

type handler struct {
	mu 		*sync.RWMutex
	handlers	 map[string]http.HandlerFunc
}

func main() {
	host := "0.0.0.0"
	port := "9999"

	if err := execute(host, port); err != nil {
		os.Exit(1)
	}
}

func execute(host, port string) error {
	mux := http.NewServeMux()
	bannerSvc := banners.NewService()

	sr := app.NewServer(mux, bannerSvc)
	sr.Init()

	srv := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: sr,
	}
	return srv.ListenAndServe()
} 
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	handler, ok := h.handlers[r.URL.Path]
	if !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	handler(w, r)
}