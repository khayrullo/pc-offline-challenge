package main

import (
	"context"
	"golang.org/x/text/language"
	"math/rand"
	"testing"
	"time"
)

func TestTranslateFail(t *testing.T) {
	ctx := context.Background()
	rand.Seed(time.Now().UTC().UnixNano())
	st := newSmartTranslator(100*time.Millisecond, 500*time.Millisecond, 1)
	s := &Service{translator: st}
	_, err := s.translator.Translate(ctx, language.English, language.Japanese, "test")
	if err == nil {
		t.Errorf("Expected service failure but got %s", err)
	}
}

func TestTranslateSuccess(t *testing.T) {
	ctx := context.Background()
	rand.Seed(time.Now().UTC().UnixNano())
	st := newSmartTranslator(100*time.Millisecond, 500*time.Millisecond, 0)
	s := &Service{translator: st}
	_, err := s.translator.Translate(ctx, language.English, language.Japanese, "test")
	if err != nil {
		t.Errorf("%s", err)
	}
}
