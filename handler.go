package expressgo

import (
	"errors"
	"net/http"
	"strings"
)

type Handler struct {
	mux *http.ServeMux
	app *App
}

// For path registration

func (h *Handler) isHostIncluded(path string) bool {
	return path[0] != '/'
}

func (h *Handler) makePrecise(path string) string {
	// remove the trailing "/"
	return strings.TrimSuffix(path, "/") + "/{$}"
}

func (h *Handler) pathToLower(path string) string {
	return strings.ToLower(path)
}

func (h *Handler) register(method string, path string, handler http.Handler) error {
	// apply config options
	if !h.app.allowHost && h.isHostIncluded(path) {
		return errors.New("path cannot contain host")
	}

	// apply config options
	p := path
	if !h.app.caseSensitive {
		p = h.pathToLower(path)
	}
	if !h.app.coarse {
		p = h.makePrecise(p)
	}

	if method != "" {
		p = method + " " + p
	}

	h.mux.Handle(p, handler)
	return nil
}

// For processing requests

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// apply config options
	if !h.app.caseSensitive {
		r.URL.Path = strings.ToLower(r.URL.Path)
	}

	h.mux.ServeHTTP(w, r)
}
