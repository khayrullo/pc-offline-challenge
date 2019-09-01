package main

import (
	"context"
	"golang.org/x/text/language"
	"math"
	"time"
)

const (
	retryPeriod = 100 * time.Millisecond
)

type smartTranslator struct {
	randomTranslator
	maxRetry int
	cacheTTL time.Duration
}

func newSmartTranslator(minDelay, maxDelay time.Duration, errorProbability float64, maxRetry int, cacheTTL time.Duration) *smartTranslator {
	st := &smartTranslator{}
	st.minDelay = minDelay
	st.maxDelay = maxDelay
	st.errorProb = errorProbability
	st.maxRetry = maxRetry
	st.cacheTTL = cacheTTL
	return st
}

func (st *smartTranslator) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	var success string
	var fail error
	i := 0
	for i < st.maxRetry {
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
