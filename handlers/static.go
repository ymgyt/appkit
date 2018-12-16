package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Static -
type Static struct {
	fs http.Handler
}

// MustStatic -
func MustStatic(root, prefix string) *Static {
	s, err := NewStatic(root, prefix)
	if err != nil {
		panic(err)
	}
	return s
}

// NewStatic -
func NewStatic(root, prefix string) (*Static, error) {
	s := &Static{
		fs: http.StripPrefix(prefix, http.FileServer(http.Dir(root))),
	}
	return s, nil
}

// StaticRoot -
func (s *Static) StaticRoot(root string, prefix string) *Static {
	s.fs = http.StripPrefix(prefix, http.FileServer(http.Dir(root)))

	return s
}

// ServeStatic -
func (s *Static) ServeStatic(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s.fs.ServeHTTP(w, r)
}
