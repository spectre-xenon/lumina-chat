package middleware

import (
	"net/http"
	"os"
	"path/filepath"
)

// This funciton takes an http.FileServer
// if the route matches a static file it's passed to the FileServer
// otherwise the index.html is served to allow client side routing to work
func (m *Middleware) Root(static http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(m.StaticPath, r.URL.Path)

		file, err := os.Stat(path)
		if os.IsNotExist(err) || file.IsDir() {
			http.ServeFile(w, r, filepath.Join(m.StaticPath, m.IndexPath))
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		static.ServeHTTP(w, r)
	}
}
