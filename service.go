package main

import "time"

// Service is a Translator user.
type Service struct {
	translator Translator
}

func NewService() *Service {
	t := newSmartTranslator(
		100*time.Millisecond,
		500*time.Millisecond,
		0.1,
		4,
		time.Duration(24*int64(time.Hour)),
	)

	return &Service{
		translator: t,
	}
}
