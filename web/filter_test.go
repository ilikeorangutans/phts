package web

import (
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterChainOrder(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		log.Println("handler")
	}
	firstCalled := false
	secondCalled := false

	first := func(i http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			assert.False(t, secondCalled)
			firstCalled = true
			i(w, r)
		}
	}
	second := func(i http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			assert.True(t, firstCalled)
			secondCalled = true
			i(w, r)
		}
	}

	res := chain(h, second, first)

	res(nil, nil)

}

func filterA(i http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("filter a")
		i(w, r)
	}
}

func filterB(i http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("filter b")
		i(w, r)
	}
}
