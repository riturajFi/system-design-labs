package metrics

import (
	"fmt"
	"net/http"
)

type Handler struct {
	reg *Registry
}

func NewHandler(reg *Registry) *Handler {
	return &Handler{
		reg: reg,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {

	snap := h.reg.Snapshot()

	for k, v := range snap {
		fmt.Fprintf(w, "%s %d\n", k, v)
	}
}
