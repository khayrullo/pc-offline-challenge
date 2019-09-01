package main

import (
	"context"
	"fmt"
	"golang.org/x/text/language"
	"math/rand"
	"time"
)

func main() {
	ctx := context.Background()
	rand.Seed(time.Now().UTC().UnixNano())
	s := NewService()
	fmt.Println(s.translator.Translate(ctx, language.English, language.Japanese, "test"))
}
