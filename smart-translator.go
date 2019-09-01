package main

import (
	"context"
	"fmt"
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

var translationCache = make(map[string]translation)

type translation struct {
	result     string
	cachedTime time.Time
}

func (st *smartTranslator) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	var success string
	var err error
	i := 0
	key := fmt.Sprintf("%s %s %s", from.String(), to.String(), data)
	//check cache before resolving from outside
	if res, ok := translationCache[key]; ok {
		translation := res.result
		//delete if cache ttl expired
		// TODO : seperate cache ttl check to another goroutine
		expiration := time.Since(res.cachedTime).Hours()
		if expiration > st.cacheTTL.Hours() {
			delete(translationCache, key)
		}
		return translation, nil
	}
	for i < st.maxRetry {
		success, err = st.randomTranslator.Translate(ctx, from, to, data)
		if err == nil {
			translationCache[key] = translation{success, time.Now()}
			return success, nil
		}
		timeout := time.Duration(int64(int64(math.Pow(2, float64(i))) * int64(retryPeriod)))
		time.Sleep(timeout)
		i++
	}
	return success, err
}
