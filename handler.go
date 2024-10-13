package expressgo

import (
	"errors"
	"net/http"
	"strings"
)

type Handler struct {
	mux *http.ServeMux
}

// For path registration

func (h *Handler) isHostIncluded(path string) bool {
	return path[0] != '/'
}

func (h *Handler) makePrecise(path string) string {
	return path + "/{$}"
}

func (h *Handler) pathToLower(path string) string {
	return strings.ToLower(path)
}

func (h *Handler) register(method string, path string, handler http.Handler) error {
	if h.isHostIncluded(path) {
		return errors.New("path cannot contain host")
	}

	p := h.makePrecise(h.pathToLower(path))
	if method != "" {
		p = method + " " + p
	}

	h.mux.Handle(p, handler)
	return nil
}

// For processing requests

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.ToLower(r.URL.Path)
	h.mux.ServeHTTP(w, r)
}
