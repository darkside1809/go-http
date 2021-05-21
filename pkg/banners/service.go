package banners

import (
	"context"
	"errors"
	"sync"
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
	startID int64
}
func NewService() *Service {
	return &Service{items: make([]*Banner, 0)}
}

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
		item.startID++
		item.ID = item.startID

		s.items = append(s.items, item)
		return item, nil
	}
	for i, banner := range s.items {
		if banner.ID == item.ID {
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