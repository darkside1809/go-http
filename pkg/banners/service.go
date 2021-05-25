package banners

import (
	"context"
	"errors"
	"net/http"
	"io/ioutil"
	// "mime/multipart"
	"strconv"
	"strings"
	"sync"
	"log"
	"os"
)

type Service struct {
	mu    sync.RWMutex
	items []*Banner
}
type Banner struct {
	ID      int64
	Title   string
	Content string
	Button  string
	Link    string
	Image	  string
}
func NewService() *Service {
	return &Service{items: make([]*Banner, 0)}
}
func UploadFile(item *Banner, r *http.Request) (string, error) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println("Err:app:uploadImage(): ", err)
		return "", err
	}
	file, header, err := r.FormFile("image")
	if err != nil {
		log.Println("Err:app:uploadImage(): no such file founded")
		return "", nil
	}
	defer file.Close()
	imageName := string(strconv.Itoa(int(item.ID)) + "." + GetExtension(header.Filename))

	tempFile, err := os.Create("web/banners/" + imageName)
	if err != nil {
		return "", err
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("Err:app:uploadImage(): ", err)
		return "", err
	}

	tempFile.Write(fileBytes)
	log.Println("imageName:", imageName)
	return imageName, nil
}

func GetExtension(imageName string) string {
	return strings.Split(imageName, ".")[1]
}

var idx int64 = 0
func (s *Service) All(ctx context.Context) ([]*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.items, nil
}

func (s *Service) ByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, banner := range s.items {
		if banner.ID == id {
			return banner, nil
		}
	}
	return nil, errors.New("item not found")
}

func (s *Service) Save(ctx context.Context, item *Banner) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if item.ID == 0 {
		idx++
		item.ID = idx
		s.items = append(s.items, item)
		return item, nil
	}
	for i := 0; i < len(s.items); i++ {
		if s.items[i].ID == item.ID {
			if item.Image == "" {
				item.Image = s.items[i].Image
			}
			s.items[i] = item
			return item, nil
		}
	}
	return nil, errors.New("item not found")
}

func (s *Service) RemoveByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i, banner := range s.items {
		if banner.ID == id {
			s.items = append(s.items[:i], s.items[i + 1:]...)
			return banner, nil
		}
	}
	return nil, errors.New("item not found")
}