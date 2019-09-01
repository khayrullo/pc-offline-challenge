package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"golang.org/x/text/language"
)

const (
	maxRetries  = 4
	retryPeriod = 100 * time.Millisecond
)

// Translator in an interface of the service which translates strings
// from one language to another.
type Translator interface {
	Translate(ctx context.Context, from, to language.Tag, data string) (string, error)
}

// randomTranslator in a Translator implementation which is used
// only for testing purposes
type randomTranslator struct {
	minDelay  time.Duration
	maxDelay  time.Duration
	errorProb float64
}

func newRandomTranslator(minDelay, maxDelay time.Duration, errorProbability float64) *randomTranslator {
	return &randomTranslator{
		minDelay:  minDelay,
		maxDelay:  maxDelay,
		errorProb: errorProbability,
	}
}

// Translate returns fake translation string or error. In any case it delays execution for some time
// to emulate remote service. Error is returned with probablity set by errorProb.
func (t randomTranslator) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	time.Sleep(t.randomDuration())

	if rand.Float64() < t.errorProb {
		return "", errors.New("translation failed")
	}

	res := fmt.Sprintf("%v -> %v : %v -> %v", from, to, data, strconv.FormatInt(rand.Int63(), 10))
	return res, nil
}

func (t randomTranslator) randomDuration() time.Duration {
	delta := t.maxDelay - t.minDelay
	var delay time.Duration = t.minDelay + time.Duration(rand.Int63n(int64(delta)))
	return delay
}

type smartTranslator struct {
	randomTranslator
}

func newSmartTranslator(minDelay, maxDelay time.Duration, errorProbability float64) *smartTranslator {
	st := &smartTranslator{}
	st.minDelay = minDelay
	st.maxDelay = maxDelay
	st.errorProb = errorProbability
	return st
}

func (st *smartTranslator) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	var success string
	var fail error
	i := 0
	for i < maxRetries {
		success, fail = st.randomTranslator.Translate(ctx, from, to, data)
		if fail == nil {
			return success, nil
		}
		timeout := time.Duration(int64(int64(math.Pow(2, float64(i))) * int64(retryPeriod)))
		time.Sleep(timeout)
		i++
	}
	return success, fail
}
