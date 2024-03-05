package storage

import (
	"fmt"

	"github.com/oktavarium/gophkeeper/internal/shared/models"
)

func (s *Storage) GetData() (
	map[string]models.SimpleCardData,
	map[string]models.SimpleData,
	map[string]models.SimpleBinaryData,
	error,
) {

	cd, err := s.getCards()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error on getting cards: %w", err)
	}

	sd, err := s.getSimple()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error on getting simple data: %w", err)
	}

	sb, err := s.getBinary()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error on getting binary data: %w", err)
	}

	return cd, sd, sb, nil
}
