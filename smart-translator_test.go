package main

import (
	"context"
	"fmt"
	"golang.org/x/text/language"
	"math/rand"
	"testing"
	"time"
)

func TestTranslateFail(t *testing.T) {
	ctx := context.Background()
	rand.Seed(time.Now().UTC().UnixNano())
	st := newSmartTranslator(100*time.Millisecond, 500*time.Millisecond, 1, 1, 24*time.Hour)
	s := &Service{translator: st}
	_, err := s.translator.Translate(ctx, language.English, language.Japanese, "test")
	if err == nil {
		t.Errorf("Expected service failure but got %s", err)
	}
}

func TestTranslateSuccess(t *testing.T) {
	ctx := context.Background()
	rand.Seed(time.Now().UTC().UnixNano())
	st := newSmartTranslator(100*time.Millisecond, 500*time.Millisecond, 0, 1, 24*time.Hour)
	s := &Service{translator: st}
	_, err := s.translator.Translate(ctx, language.English, language.Japanese, "test")
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestTrasnalteFromCache(t *testing.T) {
	ctx := context.Background()
	rand.Seed(time.Now().UTC().UnixNano())
	//set errorProbability to 1 so it always fails.
	st := newSmartTranslator(100*time.Millisecond, 500*time.Millisecond, 1, 1, 24*time.Hour)
	s := &Service{translator: st}
	key := fmt.Sprintf("%s %s %s", language.English.String(), language.Japanese.String(), "test")
	translationCache[key] = translation{"cached", time.Now()}
	res, err := s.translator.Translate(ctx, language.English, language.Japanese, "test")
	if err != nil {
		t.Fail()
	}
	if res != "cached" {
		t.Errorf("Expected cached but got %s", res)
	}
}
