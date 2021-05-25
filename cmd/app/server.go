package app

import (
	"github.com/darkside1809/http/pkg/banners"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	//"sort"
	//"fmt"
	// "strings"
)

type Server struct {
	mux       *http.ServeMux
	bannerSvc *banners.Service
}

func NewServer(mux *http.ServeMux, bannersSvc *banners.Service) *Server {
	return &Server{mux: mux, bannerSvc: bannersSvc}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) Init() {
	s.mux.HandleFunc("/banners.getAll",  s.handleGetAllBanners)
	s.mux.HandleFunc("/banners.getById", s.handleGetBannerById)
	s.mux.HandleFunc("/banners.save", 	 s.handleSaveBanner)
	s.mux.HandleFunc("/banners.removeById", s.handleRemoveById)
}

func (s *Server) handleGetAllBanners(w http.ResponseWriter, r *http.Request) {
	banners, err := s.bannerSvc.All(r.Context())
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(banners)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleGetBannerById(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.bannerSvc.ByID(r.Context(), id)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleSaveBanner(w http.ResponseWriter, r *http.Request) {
	idParam := r.PostFormValue("id")
	title := r.PostFormValue("title")
	content := r.PostFormValue("content")
	button := r.PostFormValue("button")
	link := r.PostFormValue("link")

	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item := &banners.Banner{
		ID:      id,
		Title:   title,
		Content: content,
		Button:  button,
		Link:    link,
	}

	updateBanner, err := s.bannerSvc.Save(r.Context(), item)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	imageName, err := banners.UploadFile(updateBanner, r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if imageName != "" {
		updateBanner.Image = imageName
		log.Println("update banner.Image: ", imageName)
	}


	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleRemoveById(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	banners, err := s.bannerSvc.RemoveByID(r.Context(), id)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(banners)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
	}
}