package engine

import (
	"net/http"
)

type Engine struct {}

func New() *Engine {

	return &Engine{}
}

func (e *Engine) Health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

