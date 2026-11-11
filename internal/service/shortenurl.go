package service

import (
	"context"
	"time"

	"bookmark-management/internal/repository"
)

//go:generate mockery --name=ShortenUrl --filename shortenurl.go --outpkg svc_mocks

const DefaultCodeLength = 7

// ShortenUrl creates shortened links.
type ShortenUrl interface {
	// CreateShortenLink generates a code for url and stores it for the duration exp.
	CreateShortenLink(ctx context.Context, url string, exp time.Duration) (string, error)
}

type shortenUrl struct {
	r   repository.URLStorage
	gen GenPass
}

func NewShortenUrl(r repository.URLStorage, gen GenPass) ShortenUrl {
	return &shortenUrl{r: r, gen: gen}
}

func (s *shortenUrl) CreateShortenLink(ctx context.Context, url string, exp time.Duration) (string, error) {
	code, err := s.gen.GeneratePassword(DefaultCodeLength)
	if err != nil {
		return "", err
	}

	err = s.r.StoreURL(ctx, code, url, exp)
	if err != nil {
		return "", err
	}

	return code, nil
}
